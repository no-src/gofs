//go:build integration_test

package integration

import (
	"testing"

	"github.com/no-src/fsctl/command"
)

func TestIntegration_LocalDisk(t *testing.T) {
	testCases := []struct {
		name     string
		runConf  string
		testConf string
	}{
		{"gofs local disk", "run-gofs-local-disk.yaml", "test-gofs-local-disk.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntegrationClient(t, tc.runConf, tc.testConf)
		})
	}
}

func testIntegrationClient(t *testing.T, runConf string, testConf string) {
	runConf = getRunConf(runConf)
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

	r := runWithConfigFile(runConf)
	if err = r.WaitInit(); err != nil {
		t.Errorf("wait gofs init error, err=%v", err)
		return
	}

	if err = commands.ExecActions(); err != nil {
		t.Errorf("execute actions commands error, err=%v", err)
	}

	if err = r.Shutdown(); err != nil {
		t.Errorf("gofs shutdown error, %v", err)
	}

	if err = r.Wait(); err != nil {
		t.Errorf("wait for the gofs exit error, %v", err)
	}

	if err = commands.ExecClear(); err != nil {
		t.Errorf("execute clear commands error, err=%v", err)
	}
}
