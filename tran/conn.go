package tran

import (
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/log"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	authorized     *cbool.CBool
	userName       string
	password       string
	connTime       *time.Time
	authTime       *time.Time
	startAuthCheck bool
}

func NewConn(conn net.Conn) *Conn {
	now := time.Now()
	c := &Conn{
		Conn:           conn,
		authorized:     cbool.New(false),
		connTime:       &now,
		authTime:       nil,
		startAuthCheck: false,
	}
	return c
}

func (conn *Conn) MarkAuthorized(userName, password string) {
	conn.authorized.Set(true)
	conn.userName = userName
	conn.password = password
	now := time.Now()
	conn.authTime = &now
	log.Info("the conn authorized [local=%s][remote=%s] => [username=%s password=%s]", conn.LocalAddr().String(), conn.RemoteAddr().String(), userName, password)
}

func (conn *Conn) Authorized() bool {
	return conn.authorized.Get()
}

// StartAuthCheck auto check auth state per second, close the connection if unauthorized after one minute
func (conn *Conn) StartAuthCheck() {
	if !conn.startAuthCheck {
		conn.startAuthCheck = true
		conn.authCheck()
	}
}

// StopAuthCheck stop auto auth check
func (conn *Conn) StopAuthCheck() {
	conn.startAuthCheck = false
}

func (conn *Conn) authCheck() {
	go func() {
		for {
			if !conn.startAuthCheck {
				break
			}
			authorized := conn.authorized.Get()
			if authorized {
				conn.startAuthCheck = false
				break
			} else if !authorized && time.Now().After(conn.connTime.Add(time.Minute)) {
				log.Info("conn auth check ==> [%s] is unauthorized for more than one minute since connected ", conn.Conn.RemoteAddr().String())
				if conn.Conn != nil {
					conn.Close()
				}
				conn.startAuthCheck = false
				break
			}
			<-time.After(time.Second)
		}
	}()
}
