package pb

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/charlieplate/maestro"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const (
	MsgTypeSubscribe = "Subscribe"
)

type ProtobufParser struct{}

var registry = map[string]string{}

func NewProtobufParser() *ProtobufParser {
	return &ProtobufParser{}
}

// Returns the message, the type of the message, and an error
func (pbd *ProtobufParser) Parse(data any) (maestro.Message, error) {
	d, ok := data.([]byte)
	if !ok {
		return maestro.Message{}, errors.New("invalid data type")
	}

	msg, err := unmarshalMessage(d)
	if err != nil {
		return maestro.Message{}, err
	}

	m := maestro.Message{
		Content:    nil,
		ActionType: "",
	}

	if err = validProtoVersion(msg.GetProtoVersion()); err != nil {
		return m, err
	}
	c := msg.GetContent()

	var msgType string
	if mt, ok := registry[c.GetTypeUrl()]; !ok {
		t, readErr := protoregistry.GlobalTypes.FindMessageByURL(c.GetTypeUrl())
		if readErr != nil {
			return m, readErr
		}

		msgType = string(t.Descriptor().Name())
		registry[c.GetTypeUrl()] = msgType
	} else {
		msgType = mt
	}

	switch msgType {
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

func unmarshalMessage(d []byte) (*Message, error) {
	msg := &Message{}
	err := proto.Unmarshal(d, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
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

func IntToBytes(n int, byteCount int) []byte {
	b := make([]byte, byteCount)
	for i := range byteCount {
		b[byteCount-i-1] = byte(n >> (8 * i) & 0xFF)
	}
	return b
}
