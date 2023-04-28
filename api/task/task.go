package task

import (
	"net/netip"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterServer register the monitor server
func RegisterServer(s grpc.ServiceRegistrar, taskConf string) error {
	d, err := newDispatcher(taskConf)
	if err != nil {
		return err
	}
	RegisterTaskServiceServer(s, &server{
		d: d,
	})
	return nil
}

type server struct {
	UnimplementedTaskServiceServer

	d Dispatcher
}

func (s *server) SubscribeTask(client *ClientInfo, rs TaskService_SubscribeTaskServer) error {
	p, ok := peer.FromContext(rs.Context())
	if !ok {
		return status.Errorf(codes.Unknown, "the peer information is not found")
	}
	ap, err := netip.ParseAddrPort(p.Addr.String())
	if err != nil {
		return status.Errorf(codes.Unknown, "parse client address error => %v", err)
	}
	for {
		select {
		case <-rs.Context().Done():
			return nil
		default:
		}
		task, err := s.d.Acquire(client, ap.Addr().String())
		if err != nil {
			return status.Errorf(codes.Unknown, "acquire task error => %v", err)
		}
		if task == nil {
			time.Sleep(time.Second)
			continue
		}
		err = rs.Send(task)
		if err != nil {
			return status.Errorf(codes.Unknown, "send task error => %v", err)
		}
		time.Sleep(time.Second * 3)
	}
}
