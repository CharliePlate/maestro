package maestro_test

import (
	"testing"

	"github.com/charlieplate/maestro"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

type TestAuthenticator struct {
	Error error
}

func (ta *TestAuthenticator) Authenticate(data any) (map[string]any, error) {
	if ta.Error != nil {
		return nil, ta.Error
	}

	return map[string]any{
		"auth": data,
	}, nil
}

func TestNewNilAuthenticator(t *testing.T) {
	tests := []struct {
		want *maestro.NilAuthenticator
		name string
	}{
		{
			name: "Test NewNilAuthenticator",
			want: &maestro.NilAuthenticator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, maestro.NewNilAuthenticator())
		})
	}
}

func TestNilAuthenticator_Authenticate(t *testing.T) {
	type fields struct {
		na *maestro.NilAuthenticator
	}
	type args struct {
		any
	}
	tests := []struct {
		args    args
		wantErr error
		fields  fields
		name    string
	}{
		{
			name: "Test NilAuthenticator Authenticate",
			fields: fields{
				na: &maestro.NilAuthenticator{},
			},
			args: args{
				any: nil,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			na := &maestro.NilAuthenticator{}
			if err := na.Authenticate(tt.args.any); err != nil {
				require.Equal(t, tt.wantErr, err)
			}
		})
	}
}

func TestNewJWTAuthenticator(t *testing.T) {
	type args struct {
		opts maestro.JWTAuthenticatorOpts
	}
	tests := []struct {
		want *maestro.JWTAuthenticator
		args args
		name string
	}{
		{
			name: "Test NewJWTAuthenticator",
			args: args{
				opts: maestro.JWTAuthenticatorOpts{
					SigningMethod: "HS256",
					Secret:        "secret",
				},
			},
			want: &maestro.JWTAuthenticator{
				Opts: maestro.JWTAuthenticatorOpts{
					SigningMethod: "HS256",
					Secret:        "secret",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, maestro.NewJWTAuthenticator(tt.args.opts))
		})
	}
}

func TestJWTAuthenticator_Authenticate(t *testing.T) {
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
	})
	validTokenString, _ := validToken.SignedString([]byte("secret"))

	type fields struct {
		jwt *maestro.JWTAuthenticator
	}
	type args struct {
		data any
	}
	tests := []struct {
		args    args
		wantErr error
		fields  fields
		want    map[string]any
		name    string
	}{
		{
			name: "Test JWTAuthenticator Authenticate",
			fields: fields{
				jwt: &maestro.JWTAuthenticator{
					Opts: maestro.JWTAuthenticatorOpts{
						SigningMethod: "HS256",
						Secret:        "secret",
					},
				},
			},
			args: args{
				data: validTokenString,
			},
			want: map[string]any{
				"sub":  "1234567890",
				"name": "John Doe",
			},
			wantErr: nil,
		},
		{
			name: "Test JWTAuthenticator Authenticate with invalid secret",
			fields: fields{
				jwt: &maestro.JWTAuthenticator{
					Opts: maestro.JWTAuthenticatorOpts{
						SigningMethod: "HS256",
						Secret:        "not the right secret",
					},
				},
			},
			args: args{
				data: validTokenString,
			},
			want:    map[string]any{},
			wantErr: maestro.ErrUnauthorized,
		},
		{
			name: "Test JWTAuthenticator Authenticate with invalid token",
			fields: fields{
				jwt: &maestro.JWTAuthenticator{
					Opts: maestro.JWTAuthenticatorOpts{
						SigningMethod: "HS256",
						Secret:        "secret",
					},
				},
			},
			args: args{
				data: "invalid token",
			},
			want:    map[string]any{},
			wantErr: maestro.ErrUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt := &maestro.JWTAuthenticator{
				Opts: tt.fields.jwt.Opts,
			}
			got, err := jwt.Authenticate(tt.args.data)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func makeBinaryAuthStream(m maestro.BinaryAuthContentMessage) []byte {
	b := append([]byte{}, maestro.IntToBytes(m.Version, 4)...)
	b = append(b, 0x1E)
	b = append(b, maestro.IntToBytes(m.AuthSize, 4)...)
	b = append(b, 0x1E)
	b = append(b, m.Auth...)
	b = append(b, 0x1E)
	b = append(b, maestro.IntToBytes(m.ContentSize, 4)...)
	b = append(b, 0x1E)
	b = append(b, m.Content...)
	b = append(b, []byte{0x1E, 0x1E, 0x1E}...)
	return b
}

type BinaryAuthTestAuthenticator struct {
	Error error
	Valid bool
}

func (ta BinaryAuthTestAuthenticator) Authenticate(data any) (map[string]any, error) {
	if ta.Error != nil {
		return nil, ta.Error
	}

	if !ta.Valid {
		return nil, maestro.ErrUnauthorized
	}

	return map[string]any{
		"auth": data,
	}, nil
}

type BinaryAuthTestParser struct {
	ActionType maestro.ActionType
	ConnID     string
	Error      error
}

func (tp BinaryAuthTestParser) Parse(data any) (maestro.Message, error) {
	if tp.Error != nil {
		return maestro.Message{}, tp.Error
	}

	m := maestro.Message{
		Content:    data,
		ConnID:     tp.ConnID,
		ActionType: tp.ActionType,
	}

	return m, nil
}

func TestBinaryAuthContentProtocol_ParseIncoming(t *testing.T) {
	inputResult := maestro.Message{
		ConnID:     "conn_id",
		Content:    []byte("content"),
		ActionType: maestro.ActionTypeSubscribe,
	}

	type fields struct {
		Authenticator maestro.Authenticator
		Parser        maestro.Parser
	}
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr error
	}{
		{
			name: "Valid Input",
			fields: fields{
				Authenticator: BinaryAuthTestAuthenticator{
					Error: nil,
					Valid: true,
				},
				Parser: BinaryAuthTestParser{
					ActionType: maestro.ActionTypeSubscribe,
					ConnID:     "conn_id",
					Error:      nil,
				},
			},
			args: args{
				data: makeBinaryAuthStream(maestro.BinaryAuthContentMessage{
					Version:     1,
					AuthSize:    4,
					Auth:        []byte("auth"),
					ContentSize: len([]byte("content")),
					Content:     []byte("content"),
				}),
			},
			want:    inputResult,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := &maestro.BinaryAuthContentProtocol{
				Authenticator: tt.fields.Authenticator,
				Parser:        tt.fields.Parser,
			}
			got, err := au.ParseIncoming(tt.args.data)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}
