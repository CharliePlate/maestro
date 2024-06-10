package pb

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/charlieplate/maestro"
	"github.com/charlieplate/maestro/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	MsgTypeSubscribe = "Subscribe"
)

type ParsedTag struct {
	Data          []byte `maestro:"position:3,bytecount:ContentLength"`
	Version       int    `maestro:"position:1,bytecount:4"`
	ContentLength int    `maestro:"position:2,bytecount:4"`
}

type ProtobufDecoder struct{}

// Returns the message, the type of the message, and an error
func (pbd *ProtobufDecoder) ParseIncoming(data []byte) (maestro.Message, error) {
	var d ParsedTag
	err := protocol.Unmarshal(bytes.NewReader(data), &d)
	if err != nil {
		return maestro.Message{}, err
	}

	wrapper, err := unmarshalWrapper(d.Data)
	if err != nil {
		return maestro.Message{}, err
	}

	m := maestro.Message{
		ConnID:     wrapper.GetConnId(),
		Content:    nil,
		ActionType: "",
	}

	if err = validProtoVersion(wrapper.GetProtoVersion()); err != nil {
		return m, err
	}

	c := wrapper.GetContent()
	msgType, err := protoregistry.GlobalTypes.FindMessageByURL(c.GetTypeUrl())
	if err != nil {
		return m, err
	}

	t := msgType.Descriptor().Name()

	switch t {
	case MsgTypeSubscribe:
		m.ActionType = maestro.ActionTypeSubscribe
		sub := &Subscribe{}
		err = c.UnmarshalTo(sub)
		if err != nil {
			return m, err
		}
		m.Content = sub
	default:
		return m, errors.New("unknown message type")
	}

	return m, nil
}

func unmarshalWrapper(d []byte) (*Wrapper, error) {
	wrapper := &Wrapper{}
	err := proto.Unmarshal(d, wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper, nil
}

var (
	ErrInvalidProtoVersion = errors.New("invalid protobuf version")
	minimumProtobufVersion = []int64{3, 0, 0}
)

func validProtoVersion(v string) error {
	match, err := regexp.MatchString(`\d+\.\d+\.\d+`, v)
	if err != nil {
		return err
	} else if !match {
		return ErrInvalidProtoVersion
	}

	for idx, sub := range strings.Split(v, ".") {
		subV, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			return err
		}

		if subV < minimumProtobufVersion[idx] {
			return ErrInvalidProtoVersion
		}
	}
	return nil
}
