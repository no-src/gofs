package apiclient

import (
	"context"
	"fmt"
	"time"

	authapi "github.com/no-src/gofs/api/auth"
	"github.com/no-src/gofs/api/info"
	"github.com/no-src/gofs/api/monitor"
	"github.com/no-src/gofs/api/task"
	"github.com/no-src/gofs/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type client struct {
	info.InfoServiceClient
	monitor.MonitorServiceClient
	authapi.AuthServiceClient
	task.TaskServiceClient

	host       string
	port       int
	enableTLS  bool
	certFile   string
	user       *auth.User
	clientConn *grpc.ClientConn
	creds      credentials.PerRPCCredentials
}

// New create the instance of the Client
func New(host string, port int, enableTLS bool, certFile string, user *auth.User) Client {
	if user == nil {
		user = auth.GetAnonymousUser()
	}
	return &client{
		host:      host,
		port:      port,
		enableTLS: enableTLS,
		certFile:  certFile,
		user:      user,
	}
}

func (c *client) Start() (err error) {
	if err = c.connect(); err != nil {
		return err
	}
	return c.login()
}

func (c *client) connect() (err error) {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	tranCreds := insecure.NewCredentials()
	if c.enableTLS {
		if tranCreds, err = credentials.NewClientTLSFromFile(c.certFile, c.host); err != nil {
			return err
		}
	}
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(tranCreds))
	if err != nil {
		return err
	}
	c.InfoServiceClient = info.NewInfoServiceClient(clientConn)
	c.MonitorServiceClient = monitor.NewMonitorServiceClient(clientConn)
	c.AuthServiceClient = authapi.NewAuthServiceClient(clientConn)
	c.TaskServiceClient = task.NewTaskServiceClient(clientConn)
	c.clientConn = clientConn
	return nil
}

func (c *client) Stop() error {
	return c.clientConn.Close()
}

func (c *client) getInfo() (*info.FileServerInfo, error) {
	return c.InfoServiceClient.GetInfo(context.Background(), &emptypb.Empty{}, grpc.PerRPCCredentials(c.creds))
}

func (c *client) GetInfo() (*info.FileServerInfo, error) {
	fsi, err := c.getInfo()
	if !c.needLogin(err) {
		return fsi, err
	}
	if err = c.login(); err != nil {
		return nil, err
	}
	return c.getInfo()
}

func (c *client) monitor() (monitor.MonitorService_MonitorClient, error) {
	return c.MonitorServiceClient.Monitor(context.Background(), &emptypb.Empty{}, grpc.PerRPCCredentials(c.creds))
}

func (c *client) Monitor() (monitor.MonitorService_MonitorClient, error) {
	fsi, err := c.monitor()
	if !c.needLogin(err) {
		return fsi, err
	}
	if err = c.login(); err != nil {
		return nil, err
	}
	return c.monitor()
}

func (c *client) IsClosed(err error) bool {
	return status.Code(err) == codes.Unavailable
}

func (c *client) subscribeTask(clientInfo *task.ClientInfo) (task.TaskService_SubscribeTaskClient, error) {
	return c.TaskServiceClient.SubscribeTask(context.Background(), clientInfo, grpc.PerRPCCredentials(c.creds))
}

func (c *client) SubscribeTask(clientInfo *task.ClientInfo) (task.TaskService_SubscribeTaskClient, error) {
	rc, err := c.subscribeTask(clientInfo)
	if !c.needLogin(err) {
		return rc, err
	}
	if err = c.login(); err != nil {
		return nil, err
	}
	return c.subscribeTask(clientInfo)
}

func (c *client) getToken() (token string, err error) {
	reply, err := c.AuthServiceClient.Login(context.Background(), &authapi.LoginUser{
		Username:  c.user.UserName(),
		Password:  c.user.Password(),
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return token, err
	}
	return reply.Token, nil
}

func (c *client) needLogin(err error) bool {
	return status.Code(err) == codes.Unauthenticated
}

func (c *client) login() (err error) {
	token, err := c.getToken()
	if err == nil {
		c.creds = &oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	}
	return err
}
