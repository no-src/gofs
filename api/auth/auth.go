package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterServer register the auth server
func RegisterServer(s grpc.ServiceRegistrar, token Token) {
	RegisterAuthServiceServer(s, &server{
		token: token,
	})
}

type server struct {
	UnimplementedAuthServiceServer

	token Token
}

func (s *server) Login(ctx context.Context, in *LoginUser) (*LoginReply, error) {
	token, err := s.token.GenerateToken(in)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &LoginReply{
		Token: token,
	}, nil
}
