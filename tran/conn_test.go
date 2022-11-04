//go:build tran_test

package tran

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/no-src/gofs/auth"
)

func TestNewConn_WithNilConn(t *testing.T) {
	_, err := NewConn(nil)
	if err == nil {
		t.Errorf("NewConn: expect get an error but get nil")
		return
	}
	if !errors.Is(err, errNilNetConn) {
		t.Errorf("NewConn: expect get error =>%v, but actual get error => %v", errNilNetConn, err)
	}
}

func TestConn_StartAuthCheck_Timeout(t *testing.T) {
	c, err := NewConn(&net.TCPConn{})
	if err != nil {
		t.Errorf("NewConn: create the instance of Conn error => %v", err)
		return
	}
	c.authCheckTimeout = time.Second
	c.MarkAuthorized(nil)
	c.StartAuthCheck()
	time.Sleep(c.authCheckTimeout * 2)
}

func TestConn_DataRace_MarkAuthorizedWhenCheckPerm(t *testing.T) {
	c, err := NewConn(&net.TCPConn{})
	if err != nil {
		t.Errorf("NewConn: create the instance of Conn error => %v", err)
		return
	}
	c.StartAuthCheck()
	user := getRandomUser(t)[0].ToHashUser()
	for i := 0; i < 10; i++ {
		go c.MarkAuthorized(user)
		c.Authorized()
		c.CheckPerm(auth.ReadPerm)
	}
}
