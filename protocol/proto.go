package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"reflect"
	"slices"
)

var (
	ErrDuplicatePosition = errors.New("duplicate position")
	ErrInvalidPosition   = errors.New("invalid position")
	ErrRequiresPointer   = errors.New("t must be a pointer")
)

const SizeKeyDynamic = -1

func Unmarshal(data io.Reader, v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return ErrRequiresPointer
	}

	fields := reflect.VisibleFields(reflect.ValueOf(v).Elem().Type())
	t, err := parseTags(fields)
	if err != nil {
		return err
	}

	ts := sortedTagArr(t)
	for _, nt := range ts {
		if nt.Size == 0 && nt.SizeKey != "" {
			refSize := reflect.ValueOf(v).Elem().FieldByName(nt.SizeKey)
			nt.Size = int(refSize.Int())
		}

		lr := io.LimitReader(data, int64(nt.Size))
		b, err := io.ReadAll(lr)
		if err != nil {
			return err
		}

		rf := reflect.ValueOf(v).Elem().FieldByName(nt.Name)
		//nolint:exhaustive // throw in the default
		switch rf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if len(b) != 8 {
				b = append(make([]byte, 8-len(b)), b...)
			}
			i := int64(binary.BigEndian.Uint64(b))
			rf.SetInt(i)
		case reflect.Slice:
			rf.SetBytes(b)
		case reflect.String:
			rf.SetString(string(b))
		case reflect.Bool:
			if len(b) != 1 {
				return errors.New("invalid bool")
			}
			rf.SetBool(b[0] != 0)
		default:
			return errors.New("unknown type")
		}
	}

	return nil
}

type NamedStructTag struct {
	Name string
	StructTag
}

func sortedTagArr(p map[string]StructTag) []NamedStructTag {
	tags := []NamedStructTag{}
	for k, v := range p {
		tags = append(tags, NamedStructTag{StructTag: v, Name: k})
	}

	slices.SortFunc(tags, func(i, j NamedStructTag) int {
		return i.Position - j.Position
	})

	return tags
}

func Marshal(v any) ([]byte, error) {
	fields := reflect.VisibleFields(reflect.TypeOf(v))
	t, err := parseTags(fields)
	if err != nil {
		return nil, err
	}

	ts := sortedTagArr(t)
	b := make([][]byte, len(ts))

	for _, nt := range ts {
		size := nt.Size
		if size == 0 && nt.SizeKey != "" {
			refSize := reflect.ValueOf(v).FieldByName(nt.SizeKey)
			size = int(refSize.Int())
		}

		rf := reflect.ValueOf(v).FieldByName(nt.Name)
		var fieldBytes []byte
		//nolint:exhaustive // throw in the default
		switch rf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldBytes = IntToBytes(int(rf.Int()), size)
		case reflect.Slice:
			if rf.Type().Elem().Kind() == reflect.Uint8 {
				fieldBytes = rf.Bytes()
			}
		case reflect.String:
			fieldBytes = []byte(rf.String())
		case reflect.Bool:
			if rf.Bool() {
				fieldBytes = []byte{1}
			} else {
				fieldBytes = []byte{0}
			}
		default:
			return nil, errors.New("unknown type")
		}

		fieldBytes = padOrTrucBytes(fieldBytes, size)
		b[nt.Position-1] = fieldBytes
	}

	return bytes.Join(b, nil), nil
}

func IntToBytes(n int, byteCount int) []byte {
	b := make([]byte, byteCount)
	for i := range byteCount {
		b[byteCount-i-1] = byte(n >> (8 * i) & 0xFF)
	}
	return b
}

func parseTags(fields []reflect.StructField) (map[string]StructTag, error) {
	t := map[string]StructTag{}
	pos := map[int]struct{}{}
	m := 0

	for _, field := range fields {
		tag := field.Tag.Get("maestro")
		if tag == "" {
			continue
		}

		d, err := parseTag(tag)
		if err != nil {
			return map[string]StructTag{}, err
		}

		if _, ok := pos[d.Position]; ok {
			return map[string]StructTag{}, ErrDuplicatePosition
		}
		pos[d.Position] = struct{}{}

		t[field.Name] = d
		m = int(math.Max(float64(m), float64(d.Position)))
	}

	for i := range m {
		if _, ok := pos[i+1]; !ok {
			return map[string]StructTag{}, ErrInvalidPosition
		}
	}

	return t, nil
}

func padOrTrucBytes(b []byte, s int) []byte {
	if len(b) < s {
		return append(make([]byte, s-len(b)), b...)
	}

	return b[:s]
}
