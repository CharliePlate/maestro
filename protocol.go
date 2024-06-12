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

func (au *BinaryAuthContentProtocol) Authenticate(auth any) (map[string]any, error) {
	return au.Authenticator.Authenticate(auth)
}

func (au *BinaryAuthContentProtocol) Parse(data any) (Message, error) {
	return au.Parser.Parse(data)
}

func (au *BinaryAuthContentProtocol) ParseIncoming(data any) (Message, error) {
	d, ok := data.([]byte)
	if !ok {
		return Message{}, fmt.Errorf("ParseIncoming: %w", errors.New("invalid data"))
	}

	acm := BinaryAuthContentMessage{}
	acm.Version = int(binary.BigEndian.Uint32(d[:4]))
	acm.AuthSize = int(binary.BigEndian.Uint32(d[4:8]))
	authOffset := 8 + acm.AuthSize
	acm.Auth = d[8:authOffset]
	acm.ContentSize = int(binary.BigEndian.Uint32(d[authOffset : authOffset+4]))
	acm.Content = d[authOffset+4 : acm.ContentSize+authOffset+4]

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

	if connID, ok := auth["conn_id"]; ok {
		if connIDStr, ok := connID.(string); ok {
			msg.ConnID = connIDStr
		} else {
			return Message{}, fmt.Errorf("ParseIncoming: %w", errors.New("invalid conn_id"))
		}
	}

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
	Authenticate(auth any) (map[string]any, error)
}

type NilAuthenticator struct{}

func NewNilAuthenticator() *NilAuthenticator {
	return &NilAuthenticator{}
}

func (na *NilAuthenticator) Authenticate(any) error {
	return nil
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
func (j *JWTAuthenticator) Authenticate(data any) (map[string]any, error) {
	ts, ok := data.(string)
	if !ok {
		return map[string]any{}, fmt.Errorf("authenticate: %w: %s", ErrUnauthorized, "invalid token")
	}

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != j.Opts.SigningMethod {
			return nil, fmt.Errorf("%w: %s", ErrUnauthorized, "invalid signing method")
		}

		return []byte(j.Opts.Secret), nil
	})
	if err != nil {
		return map[string]any{}, fmt.Errorf("%w: %s", ErrUnauthorized, err.Error())
	}

	if !token.Valid {
		return map[string]any{}, ErrUnauthorized
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return map[string]any{}, errors.New("invalid claims")
	}

	for k, v := range tokenClaims {
		tokenClaims[k] = v
	}

	return tokenClaims, nil
}

func IntToBytes(n int, byteCount int) []byte {
	b := make([]byte, byteCount)
	for i := range byteCount {
		b[byteCount-i-1] = byte(n >> (8 * i) & 0xFF)
	}
	return b
}
