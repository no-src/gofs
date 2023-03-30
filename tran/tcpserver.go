package tran

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/util/stringutil"
	"github.com/no-src/log"
)

type tcpServer struct {
	network   string
	ip        net.IP
	port      int
	listener  net.Listener
	conns     sync.Map
	closed    *cbool.CBool
	users     []*auth.HashUser
	certFile  string
	keyFile   string
	enableTLS bool
	reporter  report.Reporter
}

// NewServer create an instance of tcpServer
func NewServer(ip string, port int, enableTLS bool, certFile string, keyFile string, users []*auth.User, reporter report.Reporter) Server {
	srv := &tcpServer{
		network:   "tcp",
		ip:        net.ParseIP(ip),
		port:      port,
		closed:    cbool.New(true),
		enableTLS: enableTLS,
		certFile:  certFile,
		keyFile:   keyFile,
		reporter:  reporter,
	}
	if !enableTLS {
		log.Warn("the tcp server is not enable enableTLS, it is not a security connection")
	}
	hashUserList, err := auth.ToHashUserList(users)
	if err != nil {
		log.Error(err, "parse users to HashUser list error")
	} else {
		srv.users = hashUserList
	}
	if len(srv.users) == 0 {
		log.Warn("the tcp server allows anonymous access, you should set some server users by the -users or -rand_user_count flag for security reasons")
	}
	return srv
}

func (srv *tcpServer) Listen() (err error) {
	addr := &net.TCPAddr{
		IP:   srv.ip,
		Port: srv.port,
	}
	if srv.enableTLS {
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(srv.certFile, srv.keyFile)
		if err != nil {
			return err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			Time:         time.Now,
			Rand:         rand.Reader,
		}
		srv.listener, err = tls.Listen(srv.network, addr.String(), tlsConfig)
	} else {
		srv.listener, err = net.ListenTCP(srv.network, addr)
	}
	if err == nil {
		srv.closed.Set(false)
		log.Info("tcp server is listening at:%s:%d enableTLS=%v", srv.ip, srv.port, srv.enableTLS)
	}
	return err
}

func (srv *tcpServer) Accept(process func(client *Conn, data []byte)) (err error) {
	for {
		if srv.closed.Get() {
			return errors.New("tcp server is closed")
		}
		newConn, err := srv.listener.Accept()
		if err != nil {
			log.Error(err, "listener accept connection error")
			continue
		}
		clientConn, err := NewConn(newConn, srv.reporter)
		if err != nil {
			log.Error(err, "create connection error")
			continue
		}
		_, err = srv.addClient(clientConn)
		if err != nil {
			log.Error(err, "add client connection error")
			continue
		}

		go func() {
			reader := bufio.NewReader(clientConn)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					clientConn.Close()
					srv.removeClient(clientConn)
					log.Error(err, "client[%s]conn closed, current client connect count:%d", clientConn.RemoteAddrString(), srv.ClientCount())
					clientConn = nil
					return
				}
				process(clientConn, line)
			}
		}()
	}
}

func (srv *tcpServer) addClient(conn *Conn) (clientCount int, err error) {
	if conn == nil {
		return clientCount, errNilTranConn
	}
	conn.StartAuthCheck()
	addr := strings.ToLower(conn.RemoteAddrString())
	_, exist := srv.conns.Load(addr)
	srv.conns.Store(addr, conn)
	if exist {
		log.Debug("client[%s]conn is already exist, replace it now", conn.RemoteAddrString())
	}
	clientCount = srv.ClientCount()
	log.Info("client[%s]conn succeed, current client connect count:%d", conn.RemoteAddrString(), clientCount)
	srv.reporter.PutConnection(addr)
	return clientCount, err
}

// removeClient just remove client, not close conn
func (srv *tcpServer) removeClient(conn *Conn) (clientCount int, err error) {
	if conn == nil {
		return clientCount, errNilTranConn
	}
	conn.StopAuthCheck()
	addr := strings.ToLower(conn.RemoteAddrString())
	srv.conns.Delete(addr)
	clientCount = srv.ClientCount()
	log.Info("client[%s]conn removed, current client connect count:%d", conn.RemoteAddrString(), clientCount)
	srv.reporter.DeleteConnection(addr)
	return clientCount, err
}

func (srv *tcpServer) ClientCount() int {
	count := 0
	srv.conns.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

func (srv *tcpServer) Send(data []byte) error {
	srv.conns.Range(func(key, value any) bool {
		if value == nil {
			return true
		}
		c := value.(*Conn)
		if c == nil || !c.CheckPerm(auth.ReadPerm) {
			return true
		}
		writer := bufio.NewWriter(c)
		result := append(data, EndIdentity...)
		result = append(result, LFBytes...)
		_, err := writer.Write(result)
		if err != nil {
			log.Error(err, "tcp server:Send message error => Write")
		}
		log.ErrorIf(writer.Flush(), "tcp server:Send message error => Flush")
		return true
	})
	return nil
}

func (srv *tcpServer) Host() string {
	return srv.ip.String()
}

func (srv *tcpServer) Port() int {
	return srv.port
}

func (srv *tcpServer) Close() error {
	srv.closed.Set(true)
	if srv.listener == nil {
		return nil
	}
	return srv.listener.Close()
}

func (srv *tcpServer) Auth(user *auth.HashUser) (bool, auth.Perm) {
	var perm auth.Perm
	if len(srv.users) == 0 {
		return true, auth.FullPerm
	}
	if user == nil || len(user.UserNameHash) == 0 || len(user.PasswordHash) == 0 {
		return false, perm
	}
	if user.IsExpired() {
		log.Warn("user auth request info is expired, user => %s", stringutil.String(user))
		return false, perm
	}
	var loginUser *auth.HashUser
	for _, u := range srv.users {
		if u.UserNameHash == user.UserNameHash && u.PasswordHash == user.PasswordHash {
			loginUser = u
		}
	}

	if loginUser != nil {
		if !loginUser.Perm.IsValid() {
			log.Warn("the user has no permission, user => %s", stringutil.String(user))
			loginUser = nil
		} else {
			perm = loginUser.Perm
		}
	}

	return loginUser != nil, perm
}
