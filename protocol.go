package maestro

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Protocol interface {
	Authenticator
	Parser
	ParseIncoming(data any) (Message, error)
}

type BinaryAuthContentProtocol struct {
	Authenticator Authenticator
	Parser        Parser
}

type BinaryAuthContentMessage struct {
	Auth        []byte
	Content     []byte
	Version     int
	AuthSize    int
	ContentSize int
}

type AuthInfo struct {
	Claims map[string]any
	ConnID string
}

func (au *BinaryAuthContentProtocol) Authenticate(m AuthInfo) (AuthInfo, error) {
	return au.Authenticator.Authenticate(m)
}

func (au *BinaryAuthContentProtocol) Parse(data any) (Message, error) {
	return au.Parser.Parse(data)
}

func checkSeparator(d []byte, offset int) (int, error) {
	if d[offset] != 0x1E {
		return 0, fmt.Errorf("checkNextByteIsSeperator: %w", errors.New("invalid seperator"))
	}
	return offset + 1, nil
}

var ErrInvalidTerminator = errors.New("invalid terminator")

func (au *BinaryAuthContentProtocol) parseToMessage(d []byte) (BinaryAuthContentMessage, error) {
	offset := 0
	acm := BinaryAuthContentMessage{}
	var err error

	ver, err := safeByteRange(d, offset, 4)
	if err != nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: failed to read version: %w", err)
	}
	acm.Version = int(binary.BigEndian.Uint32(ver))
	offset += 4

	if offset, err = checkSeparator(d, offset); err != nil {
		return BinaryAuthContentMessage{}, err
	}

	as, err := safeByteRange(d, offset, offset+4)
	if err != nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: failed to read auth size: %w", err)
	}
	acm.AuthSize = int(binary.BigEndian.Uint32(as))
	offset += 4

	if offset, err = checkSeparator(d, offset); err != nil {
		return BinaryAuthContentMessage{}, err
	}

	auth, err := safeByteRange(d, offset, offset+acm.AuthSize)
	if err != nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: failed to read auth: %w", err)
	}
	acm.Auth = auth
	offset += acm.AuthSize

	if offset, err = checkSeparator(d, offset); err != nil {
		return BinaryAuthContentMessage{}, err
	}

	cs, err := safeByteRange(d, offset, offset+4)
	if err != nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: failed to read content size: %w", err)
	}
	acm.ContentSize = int(binary.BigEndian.Uint32(cs))
	offset += 4

	if offset, err = checkSeparator(d, offset); err != nil {
		return BinaryAuthContentMessage{}, err
	}

	content, err := safeByteRange(d, offset, offset+acm.ContentSize)
	if err != nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: failed to read content: %w", err)
	}
	acm.Content = content
	offset += acm.ContentSize

	term, err := safeByteRange(d, offset, offset+3)
	if err != nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: could not read terminator: %w", ErrInvalidTerminator)
	}
	offset += 3

	if _, isNotEOF := safeByteRange(d, offset, offset+1); isNotEOF == nil {
		return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: unexpected data after terminator: %w", ErrInvalidTerminator)
	}

	for _, b := range term {
		if b != 0x1E {
			return BinaryAuthContentMessage{}, fmt.Errorf("parseToMessage: invalid character in terminator: %w", ErrInvalidTerminator)
		}
	}

	return acm, nil
}

func (au *BinaryAuthContentProtocol) ParseIncoming(data any) (Message, error) {
	d, ok := data.([]byte)
	if !ok {
		return Message{}, fmt.Errorf("ParseIncoming: %w", errors.New("invalid data"))
	}

	acm, err := au.parseToMessage(d)
	if err != nil {
		return Message{}, err
	}

	// Validate Version at some point
	_ = acm.Version

	auth, err := au.Authenticator.Authenticate(acm.Auth)
	if err != nil {
		return Message{}, err
	}

	msg, err := au.Parser.Parse(acm.Content)
	if err != nil {
		return Message{}, err
	}

	msg.ConnID = auth.ConnID

	return msg, nil
}

type AuthParserProtocol struct {
	Authenticator Authenticator
	Parser        Parser
}

type Parser interface {
	Parse(data any) (Message, error)
}

type Authenticator interface {
	Authenticate(auth any) (AuthInfo, error)
}

type NilAuthenticator struct{}

func NewNilAuthenticator() *NilAuthenticator {
	return &NilAuthenticator{}
}

func (na *NilAuthenticator) Authenticate(any) (AuthInfo, error) {
	return AuthInfo{}, nil
}

type JWTAuthenticator struct {
	Opts JWTAuthenticatorOpts
}
type JWTAuthenticatorOpts struct {
	SigningMethod string
	Secret        string
}

func NewJWTAuthenticator(opts JWTAuthenticatorOpts) *JWTAuthenticator {
	return &JWTAuthenticator{
		Opts: opts,
	}
}

var ErrUnauthorized = errors.New("unauthorized")

// Authenticate takes a token string and returns the claims if the token is valid
func (j *JWTAuthenticator) Authenticate(data any) (AuthInfo, error) {
	ts, ok := data.(string)
	if !ok {
		return AuthInfo{}, fmt.Errorf("authenticate: %w: %s", ErrUnauthorized, "invalid token")
	}

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != j.Opts.SigningMethod {
			return nil, fmt.Errorf("%w: %s", ErrUnauthorized, "invalid signing method")
		}

		return []byte(j.Opts.Secret), nil
	})
	if err != nil {
		return AuthInfo{}, fmt.Errorf("%w: %s", ErrUnauthorized, err.Error())
	}

	if !token.Valid {
		return AuthInfo{}, ErrUnauthorized
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthInfo{}, errors.New("invalid claims")
	}

	connID, ok := tokenClaims["conn_id"]
	if !ok {
		return AuthInfo{}, errors.New("missing conn_id")
	}

	return AuthInfo{
		ConnID: connID.(string),
		Claims: tokenClaims,
	}, nil
}

func IntToBytes(n int, byteCount int) []byte {
	b := make([]byte, byteCount)
	for i := range byteCount {
		b[byteCount-i-1] = byte(n >> (8 * i) & 0xFF)
	}
	return b
}

func safeByteRange(b []byte, start, end int) ([]byte, error) {
	if start < 0 || start > len(b) {
		return nil, fmt.Errorf("safeByteRange: %w", errors.New("start index out of range"))
	}
	if end < 0 || end > len(b) {
		return nil, fmt.Errorf("safeByteRange: %w", errors.New("end index out of range"))
	}
	if start > end {
		return nil, fmt.Errorf("safeByteRange: %w", errors.New("start index greater than end index"))
	}
	return b[start:end], nil
}
