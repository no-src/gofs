//go:build integration_test

package integration

import (
	"testing"
)

func TestIntegration_RemoteDisk(t *testing.T) {
	testCases := []struct {
		name          string
		runServerConf string
		runClientConf string
		testConf      string
	}{
		{"gofs remote disk", "run-gofs-remote-disk-server.yaml", "run-gofs-remote-disk-client.yaml", "test-gofs-remote-disk.yaml"},
		{"gofs remote disk with HTTP3", "run-gofs-remote-disk-server-with-http3.yaml", "run-gofs-remote-disk-client-with-http3.yaml", "test-gofs-remote-disk.yaml"},
		{"gofs remote push", "run-gofs-remote-push-server.yaml", "run-gofs-remote-push-client.yaml", "test-gofs-remote-push.yaml"},
		{"gofs remote push with HTTP3", "run-gofs-remote-push-server-with-http3.yaml", "run-gofs-remote-push-client-with-http3.yaml", "test-gofs-remote-push.yaml"}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntegrationClientServer(t, tc.runServerConf, tc.runClientConf, tc.testConf)
		})
	}
}
