package pb_test

import (
	"testing"

	"github.com/charlieplate/maestro/pb"
	"github.com/charlieplate/maestro/protocol"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
)

type testWrapper struct {
	Content proto.Message
	ConnID  string
	Version string
}

func unsafeMarshalWrapper(tw testWrapper, content proto.Message) []byte {
	a := &anypb.Any{}
	err := a.MarshalFrom(content)
	if err != nil {
		panic(err)
	}

	d := []byte{0x00, 0x00, 0x00, 0x01}

	w := &pb.Wrapper{
		Content:      a,
		ConnId:       tw.ConnID,
		ProtoVersion: tw.Version,
	}

	data, err := proto.Marshal(w)
	if err != nil {
		panic(err)
	}

	d = append(d, protocol.IntToBytes(len(data), 4)...)
	d = append(d, data...)

	return d
}

func TestProtobuf_ParseIncoming(t *testing.T) {
	testCases := []struct {
		ExpectedContent     proto.Message
		ExpectedContentType any
		ExpectedError       error
		Name                string
		ExpectedConnID      string
		Incoming            []byte
	}{
		{
			Name: "Valid Input",
			Incoming: unsafeMarshalWrapper(testWrapper{
				Version: "3.0.0",
				ConnID:  "123",
			}, &pb.Subscribe{
				Queue: "test123",
			}),
			ExpectedContent: &pb.Subscribe{
				Queue: "test123",
			},
			ExpectedConnID: "123",
			ExpectedError:  nil,
		},
		{
			Name: "Invalid Version",
			Incoming: unsafeMarshalWrapper(testWrapper{
				Version: "2.0.0",
			}, &pb.Subscribe{}),
			ExpectedContent: nil,
			ExpectedConnID:  "",
			ExpectedError:   pb.ErrInvalidProtoVersion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			pbc := pb.ProtobufDecoder{}
			msg, err := pbc.ParseIncoming(tc.Incoming)
			if tc.ExpectedError != nil {
				require.Error(t, err, "Error expected")
			} else {
				require.NoError(t, err, "Error not expected")
				require.NotNil(t, msg, "ConnID IS nil")
				require.Equal(t, tc.ExpectedConnID, msg.ConnID, "ConnID mismatch")
				require.True(t, cmp.Equal(tc.ExpectedContent, msg.Content, protocmp.Transform()), cmp.Diff(tc.ExpectedContent, msg.Content, protocmp.Transform()))
			}
		})
	}
}
