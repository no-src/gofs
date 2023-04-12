package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/util/jsonutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Token an authentication and token component
type Token interface {
	// GenerateToken generate a new token by user info
	GenerateToken(in *LoginUser) (token string, err error)
	// IsLogin resolve the token in the context.Context
	IsLogin(ctx context.Context) (user *auth.User, err error)
}

type token struct {
	users               []*auth.User
	timeoutSeconds      int64
	tokenExpiresSeconds int64
	secret              string
	iv                  string
}

// NewToken create a default implementation of the Token
func NewToken(users []*auth.User, secret string) (Token, error) {
	err := checkTokenSecret(secret)
	if err != nil {
		return nil, err
	}
	return &token{
		users:               users,
		timeoutSeconds:      60,
		tokenExpiresSeconds: 60 * 30,
		secret:              secret,
		iv:                  "nosrc-gofs-token",
	}, nil
}

func (t *token) GenerateToken(in *LoginUser) (token string, err error) {
	var user *auth.User
	for _, u := range t.users {
		if u.UserName() == in.GetUsername() && u.Password() == in.GetPassword() && in.GetTimestamp()+t.timeoutSeconds > time.Now().Unix() {
			user = u
		}
	}
	if user != nil {
		return t.encodeToken(user)
	}
	return token, errors.New("login failed")
}

func (t *token) IsLogin(ctx context.Context) (user *auth.User, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "login failed")
	}
	var authorization string
	mdv := md.Get("authorization")
	if len(mdv) > 0 {
		authorization = strings.TrimPrefix(mdv[0], "Bearer ")
	}
	tu, err := t.decodeToken(authorization)
	if err != nil {
		return nil, err
	}
	for _, u := range t.users {
		if u.UserName() == tu.GetUsername() && tu.Expires > time.Now().Unix() {
			user = u
		}
	}
	return user, nil
}

func (t *token) encodeToken(u *auth.User) (token string, err error) {
	tu := TokenUser{
		UserId:   int32(u.UserId()),
		Username: u.UserName(),
		Expires:  time.Now().Unix() + t.tokenExpiresSeconds,
	}
	data, err := jsonutil.Marshal(tu)
	if err != nil {
		return token, err
	}
	block, err := aes.NewCipher([]byte(t.secret))
	if err != nil {
		return token, err
	}
	stream := cipher.NewCFBEncrypter(block, []byte(t.iv))
	dst := make([]byte, len(data))
	stream.XORKeyStream(dst, data)
	token = base64.StdEncoding.EncodeToString(dst)
	return token, nil
}

func (t *token) decodeToken(token string) (tu *TokenUser, err error) {
	if len(token) == 0 {
		return nil, errors.New("token can't be empty")
	}
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(t.secret))
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, []byte(t.iv))
	dst := make([]byte, len(data))
	stream.XORKeyStream(dst, data)
	err = jsonutil.Unmarshal(dst, &tu)
	return tu, err
}

func checkTokenSecret(secret string) error {
	length := len(secret)
	if length == 16 || length == 24 || length == 32 {
		return nil
	}
	return fmt.Errorf("invalid token secret size => %d, current must be either 16, 24, or 32 bytes, please check the -token_secret flag", length)
}
