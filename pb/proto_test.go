package pb_test

import (
	"testing"

	"github.com/charlieplate/maestro"
	"github.com/charlieplate/maestro/pb"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
)

type testMsg struct {
	Content proto.Message
	Version string
}

func mustUnmarshalMessage(tm testMsg, content proto.Message) []byte {
	a := &anypb.Any{}
	err := a.MarshalFrom(content)
	if err != nil {
		panic(err)
	}

	w := &pb.Message{
		Content:      a,
		ProtoVersion: tm.Version,
	}

	data, err := proto.Marshal(w)
	if err != nil {
		panic(err)
	}

	return data
}

func TestNewProtobufParser(t *testing.T) {
	tests := []struct {
		want *pb.ProtobufParser
		name string
	}{
		{
			name: "Test NewProtobufParser",
			want: &pb.ProtobufParser{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, pb.NewProtobufParser())
			require.Implements(t, (*maestro.Parser)(nil), tt.want)
		})
	}
}

func TestProtobuf_ParseIncoming(t *testing.T) {
	testCases := []struct {
		ExpectedContent     proto.Message
		ExpectedContentType any
		ExpectedError       error
		Name                string
		Incoming            []byte
	}{
		{
			Name: "Valid Input",
			Incoming: mustUnmarshalMessage(testMsg{
				Version: "3.0.0",
			}, &pb.Subscribe{
				Queue: "test123",
			}),
			ExpectedContent: &pb.Subscribe{
				Queue: "test123",
			},
			ExpectedError: nil,
		},
		{
			Name: "Invalid Version",
			Incoming: mustUnmarshalMessage(testMsg{
				Version: "2.0.0",
			}, &pb.Subscribe{}),
			ExpectedContent: nil,
			ExpectedError:   pb.ErrInvalidProtoVersion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			pbc := pb.ProtobufParser{}
			msg, err := pbc.Parse(tc.Incoming)
			if tc.ExpectedError != nil {
				require.Error(t, err, "Error expected")
			} else {
				require.NoError(t, err, "Error not expected")
				require.Zero(t, msg.ConnID, "ConnID should be nil")
				require.True(t, cmp.Equal(tc.ExpectedContent, msg.Content, protocmp.Transform()), cmp.Diff(tc.ExpectedContent, msg.Content, protocmp.Transform()))
			}
		})
	}
}
