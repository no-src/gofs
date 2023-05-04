package api

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/no-src/gofs/api/apiclient"
	"github.com/no-src/gofs/api/apiserver"
	"github.com/no-src/gofs/api/monitor"
	"github.com/no-src/gofs/api/task"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
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
	taskConfFile  = "file://./testdata/tasks.yaml"
	taskLabels    = "local-disk-sync-once-test,local-disk-sync-test"
)

func TestApiServerAndClient(t *testing.T) {
	user, _ := auth.NewUser(1, "root", "123990", auth.FullPerm)
	runApiServerAndClient(t, user)
}

func TestApiServerAndClient_WithAnonymous(t *testing.T) {
	runApiServerAndClient(t, nil)
}

func runApiServerAndClient(t *testing.T, user *auth.User) {
	server, err := runApiServer(t, user)
	if err != nil {
		t.Errorf("running api server error => %v", err)
		return
	}
	err = runApiClient(user)
	if err != nil {
		t.Errorf("running api client error => %v", err)
		return
	}
	server.Stop()
}

func runApiServer(t *testing.T, user *auth.User) (apiserver.Server, error) {
	var users []*auth.User
	if user != nil {
		users = append(users, user)
	}
	srv, err := apiserver.New(apiServerHost, apiServerPort, true, certFile, keyFile, tokenSecret, users, report.NewReporter(), serverAddr, log.DefaultLogger(), taskConfFile)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("start api server error => %v", err)
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

func runApiClient(user *auth.User) (err error) {
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

	rc, err := c.SubscribeTask(&task.ClientInfo{
		Labels: strings.Split(taskLabels, ","),
	})
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		msg, err := rc.Recv()
		if err != nil {
			if c.IsClosed(err) {
				err = nil
			}
			return err
		}
		var c conf.Config
		err = conf.ParseContent([]byte(msg.Content), msg.Ext, &c)
		if err != nil {
			return err
		}
		if i == 0 && (!c.SyncOnce || c.Source.Path() != "source" || c.Dest.Path() != "dest") {
			return errors.New("unexpect arguments")
		}
		if i == 1 && (c.SyncOnce || c.Source.Path() != "source" || c.Dest.Path() != "dest") {
			return errors.New("unexpect arguments")
		}
	}
	return c.Stop()
}
