package info

import (
	"context"

	srv "github.com/no-src/gofs/server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RegisterServer register the info server
func RegisterServer(s grpc.ServiceRegistrar, serverAddr string) {
	RegisterInfoServiceServer(s, &server{
		serverAddr: serverAddr,
	})
}

type server struct {
	UnimplementedInfoServiceServer

	serverAddr string
}

func (s *server) GetInfo(context.Context, *emptypb.Empty) (*FileServerInfo, error) {
	info := &FileServerInfo{
		ServerAddr: s.serverAddr,
		SourcePath: srv.SourceRoutePrefix,
		DestPath:   srv.DestRoutePrefix,
		QueryAddr:  srv.QueryRoute,
		PushAddr:   srv.PushFullRoute,
	}
	return info, nil
}
