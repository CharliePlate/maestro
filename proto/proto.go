package proto

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"reflect"
	"slices"
)

func Unmarshal(data io.Reader, v any) error {
	p := map[string]StructTag{}
	pos := map[int]bool{}
	m := 0

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("t must be a pointer")
	}

	fields := reflect.VisibleFields(reflect.ValueOf(v).Elem().Type())
	for _, field := range fields {
		tag := field.Tag.Get("maestro")
		if tag == "" {
			continue
		}

		d, err := parseTag(tag)
		if err != nil {
			return err
		}

		if _, ok := pos[d.Position]; ok {
			return errors.New("duplicate position")
		}
		pos[d.Position] = true

		p[field.Name] = d
		m = int(math.Max(float64(m), float64(d.Position)))
	}

	for i := range m {
		if _, ok := pos[i+1]; !ok {
			return errors.New("missing position")
		}
	}

	ts := sortedTagArr(p)
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
