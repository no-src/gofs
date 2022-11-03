package tran

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/gofs/report"
	"github.com/no-src/log"
)

var (
	errNilConn = errors.New("the instance of net.Conn is nil")
)

// Conn the component of network connection
type Conn struct {
	nc             net.Conn
	authorized     *cbool.CBool
	user           *auth.HashUser
	connTime       *time.Time
	authTime       *time.Time
	startAuthCheck *cbool.CBool
	mu             sync.RWMutex
}

// NewConn create a Conn instance
func NewConn(nc net.Conn) (*Conn, error) {
	if nc == nil {
		return nil, errNilConn
	}
	now := time.Now()
	c := &Conn{
		nc:             nc,
		authorized:     cbool.New(false),
		connTime:       &now,
		authTime:       nil,
		startAuthCheck: cbool.New(false),
	}
	return c, nil
}

// MarkAuthorized mark the current connection is authorized with the user info
func (conn *Conn) MarkAuthorized(user *auth.HashUser) {
	if user == nil {
		return
	}
	conn.authorized.Set(true)
	conn.mu.Lock()
	conn.user = user
	conn.mu.Unlock()
	now := time.Now()
	conn.authTime = &now
	addr := conn.RemoteAddrString()
	log.Info("the conn authorized [local=%s][remote=%s] => [username=%s password=%s perm=%s]", conn.LocalAddrString(), addr, user.UserNameHash, user.PasswordHash, user.Perm.String())
	report.GlobalReporter.PutAuth(addr, user)
}

// Authorized check the current connection is authorized or not
func (conn *Conn) Authorized() bool {
	return conn.authorized.Get()
}

// CheckPerm check the current connection's permission whether accord with the specified permission
func (conn *Conn) CheckPerm(perm auth.Perm) bool {
	conn.mu.RLock()
	defer conn.mu.RUnlock()
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

// RemoteAddrString returns the remote network address, if known, or else return empty string
func (conn *Conn) RemoteAddrString() string {
	if conn.nc.RemoteAddr() == nil {
		return ""
	}
	return conn.nc.RemoteAddr().String()
}

// LocalAddrString returns the local network address, if known, or else return empty string
func (conn *Conn) LocalAddrString() string {
	if conn.nc.LocalAddr() == nil {
		return ""
	}
	return conn.nc.LocalAddr().String()
}

// Write writes data to the connection
func (conn *Conn) Write(b []byte) (n int, err error) {
	return conn.nc.Write(b)
}

// Read reads data from the connection
func (conn *Conn) Read(b []byte) (n int, err error) {
	return conn.nc.Read(b)
}

// Close closes the connection
func (conn *Conn) Close() error {
	return conn.nc.Close()
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
				log.Info("conn auth check ==> [%s] is unauthorized for more than one minute since connected ", conn.RemoteAddrString())
				conn.Close()
				conn.startAuthCheck.Set(false)
				break
			}
			<-time.After(time.Second)
		}
	}()
}
