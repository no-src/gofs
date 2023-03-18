//go:build integration_test

package integration

import (
	"testing"

	"github.com/no-src/fsctl/command"
)

func TestIntegration_RemoteDisk(t *testing.T) {
	testCases := []struct {
		name          string
		runServerConf string
		runClientConf string
		testConf      string
	}{
		{"gofs remote disk", "run-gofs-remote-server.yaml", "run-gofs-remote-client.yaml", "test-gofs-remote-disk.yaml"},
		{"gofs remote disk with HTTP3", "run-gofs-remote-server-with-http3.yaml", "run-gofs-remote-client-with-http3.yaml", "test-gofs-remote-disk.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntegrationClientServer(t, tc.runServerConf, tc.runClientConf, tc.testConf)
		})
	}
}

func testIntegrationClientServer(t *testing.T, runServerConf string, runClientConf string, testConf string) {
	runServerConf = getRunConf(runServerConf)
	runClientConf = getRunConf(runClientConf)
	testConf = getTestConf(testConf)

	commands, err := command.ParseConfigFile(testConf)
	if err != nil {
		t.Errorf("parse confile file error, err=%v", err)
		return
	}

	if err = commands.ExecInit(); err != nil {
		t.Errorf("execute init commands error, err=%v", err)
		return
	}

	sr := runWithConfigFile(runServerConf)
	if err = sr.WaitInit(); err != nil {
		t.Errorf("wait gofs server init error, err=%v", err)
		return
	}

	cr := runWithConfigFile(runClientConf)
	if err = cr.WaitInit(); err != nil {
		t.Errorf("wait gofs client init error, err=%v", err)
		// shutdown the server
		if err = sr.Shutdown(); err != nil {
			t.Errorf("gofs server shutdown error, %v", err)
		}
		return
	}

	if err = commands.ExecActions(); err != nil {
		t.Errorf("execute actions commands error, err=%v", err)
	}

	if err = cr.Shutdown(); err != nil {
		t.Errorf("gofs client shutdown error, %v", err)
	}

	if err = sr.Shutdown(); err != nil {
		t.Errorf("gofs server shutdown error, %v", err)
	}

	if err = cr.Wait(); err != nil {
		t.Errorf("wait for the gofs client exit error, %v", err)
	}

	if err = sr.Wait(); err != nil {
		t.Errorf("wait for the gofs server exit error, %v", err)
	}

	if err = commands.ExecClear(); err != nil {
		t.Errorf("execute clear commands error, err=%v", err)
	}
}
