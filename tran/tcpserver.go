package tran

import (
	"bufio"
	"errors"
	"github.com/no-src/log"
	"net"
	"strings"
)

type tcpServer struct {
	network  string
	ip       net.IP
	port     int
	listener *net.TCPListener
	conns    map[string]net.Conn
	closed   bool
}

func NewServer(ip string, port int) Server {
	srv := &tcpServer{}
	srv.ip = net.ParseIP(ip)
	srv.port = port
	srv.network = "tcp"
	srv.conns = make(map[string]net.Conn)
	return srv
}

func (srv *tcpServer) Listen() (err error) {
	addr := &net.TCPAddr{
		IP:   srv.ip,
		Port: srv.port,
	}
	srv.listener, err = net.ListenTCP(srv.network, addr)
	if err == nil {
		log.Info("tcp server is listening at:%s:%d", srv.ip, srv.port)
	}
	return err
}

func (srv *tcpServer) Accept(process func(client net.Conn, data []byte)) (err error) {
	for {
		if srv.closed {
			return errors.New("tcp server is closed")
		}
		clientConn, err := srv.listener.Accept()
		if err != nil {
			continue
		}
		srv.addClient(clientConn)

		go func() {
			reader := bufio.NewReader(clientConn)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					clientConn.Close()
					srv.removeClient(clientConn)
					log.Error(err, "client[%s]conn closed ,current client connect count:%d", clientConn.RemoteAddr().String(), srv.ClientCount())
					clientConn = nil
					return
				} else {
					process(clientConn, line)
				}
			}

		}()
	}
	return err
}

func (srv *tcpServer) addClient(conn net.Conn) (clientCount int, err error) {
	if conn == nil {
		return clientCount, errors.New("conn is nil")
	}
	addr := strings.ToLower(conn.RemoteAddr().String())
	_, exist := srv.conns[addr]
	srv.conns[addr] = conn
	if exist {
		log.Debug("client[%s]conn is already exist ,replace it now", conn.RemoteAddr().String())
	}
	clientCount = srv.ClientCount()
	log.Debug("client[%s]conn succeed ,current client connect count:%d", conn.RemoteAddr().String(), clientCount)
	return clientCount, err
}

// removeClient just remove client ,not close conn
func (srv *tcpServer) removeClient(conn net.Conn) (clientCount int, err error) {
	if conn == nil {
		return clientCount, errors.New("conn is nil")
	}
	addr := strings.ToLower(conn.RemoteAddr().String())
	delete(srv.conns, addr)
	clientCount = srv.ClientCount()
	log.Debug("client[%s]conn removed ,current client connect count:%d", conn.RemoteAddr().String(), clientCount)
	return clientCount, err
}

func (srv *tcpServer) ClientCount() int {
	return len(srv.conns)
}

func (srv *tcpServer) Send(data []byte) error {
	for _, c := range srv.conns {
		writer := bufio.NewWriter(c)
		result := append(data, EndIdentity...)
		result = append(result, LFBytes...)
		_, err := writer.Write(result)
		if err != nil {
			log.Error(err, "tcp server:Send message error => Write")
		}
		err = writer.Flush()
		if err != nil {
			log.Error(err, "tcp server:Send message error => Flush")
		}
	}
	return nil
}

func (srv *tcpServer) Host() string {
	return srv.ip.String()
}

func (srv *tcpServer) Port() int {
	return srv.port
}

func (srv *tcpServer) Close() error {
	srv.closed = true
	return srv.listener.Close()
}
