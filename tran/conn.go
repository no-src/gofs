package tran

import (
	"net"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/gofs/report"
	"github.com/no-src/log"
)

// Conn the component of network connection
type Conn struct {
	net.Conn

	authorized     *cbool.CBool
	user           *auth.HashUser
	connTime       *time.Time
	authTime       *time.Time
	startAuthCheck *cbool.CBool
}

// NewConn create a Conn instance
func NewConn(conn net.Conn) *Conn {
	now := time.Now()
	c := &Conn{
		Conn:           conn,
		authorized:     cbool.New(false),
		connTime:       &now,
		authTime:       nil,
		startAuthCheck: cbool.New(false),
	}
	return c
}

// MarkAuthorized mark the current connection is authorized with the user info
func (conn *Conn) MarkAuthorized(user *auth.HashUser) {
	if user == nil {
		return
	}
	conn.authorized.Set(true)
	conn.user = user
	now := time.Now()
	conn.authTime = &now
	addr := conn.RemoteAddr().String()
	log.Info("the conn authorized [local=%s][remote=%s] => [username=%s password=%s perm=%s]", conn.LocalAddr().String(), addr, user.UserNameHash, user.PasswordHash, user.Perm.String())
	report.GlobalReporter.PutAuth(addr, user)
}

// Authorized check the current connection is authorized or not
func (conn *Conn) Authorized() bool {
	return conn.authorized.Get()
}

// CheckPerm check the current connection's permission whether accord with the specified permission
func (conn *Conn) CheckPerm(perm auth.Perm) bool {
	if !conn.Authorized() || conn.user == nil {
		return false
	}
	return perm.CheckTo(conn.user.Perm)
}

// StartAuthCheck auto check auth state per second, close the connection if unauthorized after one minute
func (conn *Conn) StartAuthCheck() {
	if !conn.startAuthCheck.Get() {
		conn.startAuthCheck.Set(true)
		conn.authCheck()
	}
}

// StopAuthCheck stop auto auth check
func (conn *Conn) StopAuthCheck() {
	conn.startAuthCheck.Set(false)
}

func (conn *Conn) authCheck() {
	go func() {
		for {
			if !conn.startAuthCheck.Get() {
				break
			}
			authorized := conn.Authorized()
			if authorized {
				conn.startAuthCheck.Set(false)
				break
			} else if !authorized && time.Now().After(conn.connTime.Add(time.Minute)) {
				log.Info("conn auth check ==> [%s] is unauthorized for more than one minute since connected ", conn.Conn.RemoteAddr().String())
				if conn.Conn != nil {
					conn.Close()
				}
				conn.startAuthCheck.Set(false)
				break
			}
			<-time.After(time.Second)
		}
	}()
}
