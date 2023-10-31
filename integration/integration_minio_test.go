//go:build integration_test_minio

package integration

import (
	"testing"
)

func TestIntegration_MinIO(t *testing.T) {
	testCases := []struct {
		name          string
		runServerConf string
		runClientConf string
		testConf      string
	}{
		{"gofs MinIO push", "", "run-gofs-minio-push-client.yaml", "test-gofs-minio-push.yaml"},
		{"gofs MinIO pull", "", "run-gofs-minio-pull-client.yaml", "test-gofs-minio-pull.yaml"},
		{"gofs MinIO push partial", "", "run-gofs-minio-push-client-partial.yaml", "test-gofs-minio-push-partial.yaml"},
		{"gofs MinIO pull partial", "", "run-gofs-minio-pull-client-partial.yaml", "test-gofs-minio-pull-partial.yaml"},
		{"gofs MinIO push once", "", "run-gofs-minio-push-client-once.yaml", "test-gofs-minio-push-once.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testIntegrationClientServer(t, tc.runServerConf, tc.runClientConf, tc.testConf)
		})
	}
}
