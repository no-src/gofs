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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	certFile      = "../util/httputil/testdata/cert.pem"
	keyFile       = "../util/httputil/testdata/key.pem"
	serverAddr    = "https://127.0.0.1"
	apiServerHost = "127.0.0.1"
	apiServerPort = 52172
	tokenSecret   = "123456abcdefghij"
)

func TestApiServerAndClient(t *testing.T) {
	server, err := runApiServer()
	if err != nil {
		t.Errorf("running api server error => %v", err)
		return
	}
	time.Sleep(time.Second * 3)
	err = runApiClient()
	if err != nil {
		t.Errorf("running api client error => %v", err)
		return
	}
	server.Stop()
}

func runApiServer() (apiserver.Server, error) {
	var users []*auth.User
	user, _ := auth.NewUser(1, "root", "123990", auth.FullPerm)
	users = append(users, user)
	srv, err := apiserver.New(apiServerHost, apiServerPort, true, certFile, keyFile, tokenSecret, users, report.NewReporter(), serverAddr, log.DefaultLogger())
	if err != nil {
		return nil, err
	}
	go srv.Start()
	go func() {
		for {
			time.Sleep(time.Second)
			srv.SendMonitorMessage(&monitor.MonitorMessage{
				BaseUrl: serverAddr,
			})
		}
	}()
	return srv, nil
}

func runApiClient() error {
	user, _ := auth.NewUser(1, "root", "123990", auth.FullPerm)
	c := apiclient.New(apiServerHost, apiServerPort, true, certFile, user)
	err := c.Start()
	if err != nil {
		return err
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
	go func() {
		time.Sleep(time.Second * 3)
		c.Stop()
	}()
	for {
		msg, err := ms.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled {
				err = nil
			}
			return err
		}
		if msg.GetBaseUrl() != serverAddr {
			return errors.New("invalid baseurl")
		}
	}
	return nil
}
