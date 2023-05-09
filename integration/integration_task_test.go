//go:build integration_test_task

package integration

import (
	"testing"
)

func TestIntegration_Task(t *testing.T) {
	testCases := []struct {
		name          string
		runServerConf string
		testConf      string
	}{
		{"gofs local disk", "run-gofs-local-disk-task-server.yaml", "test-gofs-local-disk.yaml"},
		{"gofs remote disk", "run-gofs-remote-disk-server.yaml", "test-gofs-remote-disk.yaml"},
		{"gofs remote disk with HTTP3", "run-gofs-remote-disk-server-with-http3.yaml", "test-gofs-remote-disk.yaml"},
		{"gofs remote push", "run-gofs-remote-push-server.yaml", "test-gofs-remote-push.yaml"},
		{"gofs remote push with HTTP3", "run-gofs-remote-push-server-with-http3.yaml", "test-gofs-remote-push.yaml"}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntegrationClientServer(t, tc.runServerConf, "run-gofs-task-client.yaml", tc.testConf)
		})
	}
}
