//go:build integration_test_sftp

package integration

import (
	"testing"
)

func TestIntegration_SFTP(t *testing.T) {
	testCases := []struct {
		name          string
		runServerConf string
		runClientConf string
		testConf      string
	}{
		{"gofs SFTP push", "", "run-gofs-sftp-push-client.yaml", "test-gofs-sftp-push.yaml"},
		{"gofs SFTP pull", "", "run-gofs-sftp-pull-client.yaml", "test-gofs-sftp-pull.yaml"},
		{"gofs SFTP push once", "", "run-gofs-sftp-push-client-once.yaml", "test-gofs-sftp-push-once.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntegrationClientServer(t, tc.runServerConf, tc.runClientConf, tc.testConf)
		})
	}
}
