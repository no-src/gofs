package monitor

import (
	"sync"

	authapi "github.com/no-src/gofs/api/auth"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RegisterServer register the monitor server
func RegisterServer(s grpc.ServiceRegistrar, monitors *sync.Map, reporter report.Reporter, token authapi.Token) {
	RegisterMonitorServiceServer(s, &server{
		monitors: monitors,
		reporter: reporter,
		token:    token,
	})
}

type server struct {
	UnimplementedMonitorServiceServer

	monitors *sync.Map
	reporter report.Reporter
	token    authapi.Token
}

func (s *server) Monitor(in *emptypb.Empty, m MonitorService_MonitorServer) error {
	p, ok := peer.FromContext(m.Context())
	if !ok {
		return status.Errorf(codes.Aborted, "the peer information is not found")
	}
	k := p.Addr.String()
	var msgChan chan *MonitorMessage
	v, ok := s.monitors.Load(k)
	if ok {
		msgChan = v.(chan *MonitorMessage)
	} else {
		msgChan = make(chan *MonitorMessage)
		s.monitors.Store(k, msgChan)
		user, _ := s.token.IsLogin(m.Context())
		s.reporter.PutConnection(k, auth.MapperToSessionUser(user))
	}
	for {
		select {
		case msg := <-msgChan:
			m.Send(msg)
		case <-m.Context().Done():
			s.monitors.Delete(k)
			s.reporter.DeleteConnection(k)
			return nil
		}
	}
	return nil
}
