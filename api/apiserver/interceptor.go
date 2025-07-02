package apiserver

import (
	"context"

	"github.com/no-src/gofs/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (gs *grpcServer) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if info.FullMethod == auth.AuthService_Login_FullMethodName {
		return handler(ctx, req)
	}
	if len(gs.users) == 0 {
		return handler(ctx, req)
	}
	loginUser, err := gs.token.IsLogin(ctx)
	if err != nil || loginUser == nil {
		gs.logger.ErrorIf(err, "login failed")
		return nil, status.Errorf(codes.Unauthenticated, "login failed")
	}
	return handler(ctx, req)
}

func (gs *grpcServer) StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if len(gs.users) == 0 {
		return nil
	}
	loginUser, err := gs.token.IsLogin(ss.Context())
	if err != nil || loginUser == nil {
		gs.logger.ErrorIf(err, "login failed")
		return status.Errorf(codes.Unauthenticated, "login failed")
	}
	return handler(srv, ss)
}
