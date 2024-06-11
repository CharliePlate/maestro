package pb

import (
	"encoding/binary"
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

type ParsedTag struct {
	Data          []byte
	Version       int
	ContentLength int
}

type ProtobufDecoder struct{}

// Returns the message, the type of the message, and an error
func (pbd *ProtobufDecoder) ParseIncoming(data []byte) (maestro.Message, error) {
	d := parse(data)

	msg, err := unmarshalMessage(d.Data)
	if err != nil {
		return maestro.Message{}, err
	}

	m := maestro.Message{
		ConnID:     msg.GetConnId(),
		Content:    nil,
		ActionType: "",
	}

	if err = validProtoVersion(msg.GetProtoVersion()); err != nil {
		return m, err
	}

	c := msg.GetContent()
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

func parse(d []byte) ParsedTag {
	p := ParsedTag{}

	p.Version = int(binary.BigEndian.Uint16(d[:4]))
	p.ContentLength = int(binary.BigEndian.Uint16(d[4:8]))
	p.Data = d[8:]

	return p
}
