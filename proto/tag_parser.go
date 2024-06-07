package proto

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var ErrUnknownTag = errors.New("unknown tag")

const (
	TagPosition  = "position"
	TagByteCount = "bytecount"
)

type StructTag struct {
	SizeKey  string
	Position int
	Size     int
}

func parseTag(tag string) (StructTag, error) {
	parsed := []string{}
	p := StructTag{
		SizeKey:  "",
		Position: 0,
		Size:     0,
	}

	tags := strings.Split(tag, ",")
	for _, t := range tags {
		kv := strings.Split(t, ":")
		if len(kv) != 2 {
			return p, fmt.Errorf("error parsing tag: %s", t)
		}

		switch strings.ToLower(kv[0]) {
		case TagPosition:
			v, err := strconv.Atoi(kv[1])
			if err != nil {
				return p, fmt.Errorf("error parsing tag: %w", err)
			}

			parsed = append(parsed, TagPosition)
			p.Position = v
		case TagByteCount:
			if s, err := strconv.Atoi(kv[1]); err == nil {
				p.Size = s
			} else {
				p.SizeKey = kv[1]
			}

			parsed = append(parsed, TagByteCount)
		default:
			return p, ErrUnknownTag
		}
	}

	for _, i := range []string{TagPosition, TagByteCount} {
		if !slices.Contains(parsed, i) {
			return p, fmt.Errorf("missing tag: %s", i)
		}
	}

	return p, nil
}
