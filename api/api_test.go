package api

import (
	"errors"
	"testing"
	"time"

	"github.com/no-src/gofs/api/apiclient"
	"github.com/no-src/gofs/api/apiserver"
	"github.com/no-src/gofs/api/monitor"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/report"
	"github.com/no-src/log"
)

const (
	certFile      = "../util/httputil/testdata/cert.pem"
	keyFile       = "../util/httputil/testdata/key.pem"
	serverAddr    = "https://127.0.0.1"
	apiServerHost = "127.0.0.1"
	apiServerPort = 8128
	tokenSecret   = "123456abcdefghij"
)

func TestApiServerAndClient(t *testing.T) {
	server, err := runApiServer(t)
	if err != nil {
		t.Errorf("running api server error => %v", err)
		return
	}
	err = runApiClient()
	if err != nil {
		t.Errorf("running api client error => %v", err)
		return
	}
	server.Stop()
}

func runApiServer(t *testing.T) (apiserver.Server, error) {
	var users []*auth.User
	user, _ := auth.NewUser(1, "root", "123990", auth.FullPerm)
	users = append(users, user)
	srv, err := apiserver.New(apiServerHost, apiServerPort, true, certFile, keyFile, tokenSecret, users, report.NewReporter(), serverAddr, log.DefaultLogger())
	if err != nil {
		return nil, err
	}
	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("start api server error => %v", err)
		} else {
			t.Logf("start api server success => %s:%d", apiServerHost, apiServerPort)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond)
			srv.SendMonitorMessage(&monitor.MonitorMessage{
				BaseUrl: serverAddr,
			})
		}
	}()
	return srv, nil
}

func runApiClient() (err error) {
	user, _ := auth.NewUser(1, "root", "123990", auth.FullPerm)
	c := apiclient.New(apiServerHost, apiServerPort, true, certFile, user)
	for i := 0; i < 3; i++ {
		err = c.Start()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 3)
	}

	info, err := c.GetInfo()
	if err != nil {
		return err
	}
	if info.GetServerAddr() != serverAddr {
		return errors.New("invalid server addr")
	}

	ms, err := c.Monitor()
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		msg, err := ms.Recv()
		if err != nil {
			if c.IsClosed(err) {
				err = nil
			}
			return err
		}
		if msg.GetBaseUrl() != serverAddr {
			return errors.New("invalid baseurl")
		}
	}
	return c.Stop()
}
