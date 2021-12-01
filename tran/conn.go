package tran

import (
	"github.com/no-src/log"
	"net"
)

type Conn struct {
	net.Conn
	authorized bool
	userName   string
	password   string
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:       conn,
		authorized: false,
	}
}

func (conn *Conn) MarkAuthorized(userName, password string) {
	conn.authorized = true
	conn.userName = userName
	conn.password = password
	log.Debug("tran client authorized [%s] => [username=%s password=%s]", conn.RemoteAddr().String(), userName, password)
}

func (conn *Conn) Authorized() bool {
	return conn.authorized
}
