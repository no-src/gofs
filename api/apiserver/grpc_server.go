package apiserver

import (
	"net"
	"sync"
	"time"

	authapi "github.com/no-src/gofs/api/auth"
	"github.com/no-src/gofs/api/info"
	"github.com/no-src/gofs/api/monitor"
	"github.com/no-src/gofs/api/task"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/clist"
	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcServer struct {
	network         string
	ip              net.IP
	port            int
	users           []*auth.User
	token           authapi.Token
	certFile        string
	keyFile         string
	enableTLS       bool
	reporter        report.Reporter
	httpServerAddr  string
	server          *grpc.Server
	monitors        *sync.Map
	monitorMessages *clist.CList
	logger          *logger.Logger
	taskConf        string
}

// New create the instance of the Server
func New(ip string, port int, enableTLS bool, certFile string, keyFile string, tokenSecret string, users []*auth.User, reporter report.Reporter, httpServerAddr string, logger *logger.Logger, taskConf string) (Server, error) {
	if len(users) == 0 {
		logger.Warn("the grpc server allows anonymous access, you should set some server users by the -users or -rand_user_count flag for security reasons")
		users = append(users, auth.GetAnonymousUser())
	}
	token, err := authapi.NewToken(users, tokenSecret)
	if err != nil {
		return nil, err
	}
	srv := &grpcServer{
		network:         "tcp",
		ip:              net.ParseIP(ip),
		port:            port,
		users:           users,
		token:           token,
		enableTLS:       enableTLS,
		certFile:        certFile,
		keyFile:         keyFile,
		reporter:        reporter,
		httpServerAddr:  httpServerAddr,
		monitors:        &sync.Map{},
		monitorMessages: clist.New(),
		logger:          logger,
		taskConf:        taskConf,
	}
	creds := insecure.NewCredentials()
	if enableTLS {
		if creds, err = credentials.NewServerTLSFromFile(srv.certFile, srv.keyFile); err != nil {
			return nil, err
		}
	} else {
		logger.Warn("the grpc server is not enable enableTLS, it is not a security connection")
	}
	srv.server = grpc.NewServer(grpc.Creds(creds), grpc.StreamInterceptor(srv.StreamServerInterceptor), grpc.UnaryInterceptor(srv.UnaryServerInterceptor))
	if err = srv.initRoute(srv.server); err != nil {
		return nil, err
	}
	return srv, nil
}

func (gs *grpcServer) Start() error {
	addr := &net.TCPAddr{
		IP:   gs.ip,
		Port: gs.port,
	}
	listener, err := net.ListenTCP(gs.network, addr)
	if err != nil {
		return err
	}
	gs.logger.Info("grpc server is listening at:%s:%d enableTLS=%v", gs.ip, gs.port, gs.enableTLS)
	go gs.processMonitorMessage()
	return gs.server.Serve(listener)
}

func (gs *grpcServer) Stop() {
	gs.server.GracefulStop()
}

func (gs *grpcServer) SendMonitorMessage(message *monitor.MonitorMessage) {
	gs.monitorMessages.PushBack(message)
}

func (gs *grpcServer) initRoute(s *grpc.Server) (err error) {
	info.RegisterServer(s, gs.httpServerAddr)
	monitor.RegisterServer(s, gs.monitors, gs.reporter, gs.token)
	authapi.RegisterServer(s, gs.token)
	err = task.RegisterServer(s, gs.taskConf)
	return err
}

func (gs *grpcServer) processMonitorMessage() {
	for {
		e := gs.monitorMessages.Front()
		if e != nil {
			msg := e.Value.(*monitor.MonitorMessage)
			gs.monitors.Range(func(key, value any) bool {
				msgChan := value.(chan *monitor.MonitorMessage)
				msgChan <- msg
				return true
			})
			gs.monitorMessages.Remove(e)
		} else {
			time.Sleep(time.Millisecond)
		}
	}
}
