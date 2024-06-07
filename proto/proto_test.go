package proto_test

import (
	"bytes"
	"testing"

	"github.com/charlieplate/maestro/proto"
	"github.com/stretchr/testify/require"
)

type byteString struct {
	input     any
	byteCount int
}

func stringToByteArray(s ...byteString) []byte {
	var b []byte

	for _, st := range s {
		var tb []byte
		switch v := st.input.(type) {
		case string:
			tb = []byte(v)
		case int:
			tb = intToBytes(v, st.byteCount)
		default:
			panic("Unknown type")
		}

		// Pad or truncate to the desired byteCount
		if len(tb) < st.byteCount {
			tb = append(make([]byte, st.byteCount-len(tb)), tb...)
		} else if len(tb) > st.byteCount {
			tb = tb[:st.byteCount]
		}

		b = append(b, tb...)
	}
	return b
}

func intToBytes(n int, byteCount int) []byte {
	b := make([]byte, byteCount)
	for i := range byteCount {
		b[byteCount-i-1] = byte(n >> (8 * i) & 0xFF)
	}
	return b
}

func TestStringToByteArray(t *testing.T) {
	testCases := []struct {
		name      string
		input     any
		want      []byte
		byteCount int
	}{
		{
			name:      "String Parsing",
			input:     "Hello",
			byteCount: 8,
			want:      []byte{0x00, 0x00, 0x00, 0x48, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			name:      "Parse Int",
			input:     1,
			byteCount: 4,
			want:      []byte{0x00, 0x00, 0x00, 0x01},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := stringToByteArray(byteString{input: tc.input, byteCount: tc.byteCount})
			require.Equal(t, tc.want, result)
		})
	}
}

func TestProto_Unmarshal(t *testing.T) {
	type TestStruct struct {
		Data          []byte `maestro:"position:3,bytecount:ContentLength"`
		Version       int    `maestro:"position:1,bytecount:4"`
		ContentLength int    `maestro:"position:2,bytecount:4"`
	}

	testCases := []struct {
		expectedError  error
		input          []byte
		name           string
		expectedOutput TestStruct
	}{
		{
			name: "Valid input",
			input: stringToByteArray(
				byteString{input: 1, byteCount: 4},
				byteString{input: 5, byteCount: 4},
				byteString{input: "hello", byteCount: 5},
			),
			expectedOutput: TestStruct{
				Version:       1,
				ContentLength: 5,
				Data:          []byte("hello"),
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ts TestStruct
			reader := bytes.NewReader(tc.input)
			err := proto.Unmarshal(reader, &ts)

			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedOutput, ts)
		})
	}
}
