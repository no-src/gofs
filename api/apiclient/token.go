package apiclient

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
)

// insecureTokenSource supplies PerRPCCredentials from an oauth2.TokenSource.
type insecureTokenSource struct {
	oauth2.TokenSource
}

// GetRequestMetadata gets the request metadata as a map from a TokenSource.
func (ts insecureTokenSource) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	ri, _ := credentials.RequestInfoFromContext(ctx)
	if err = credentials.CheckSecurityLevel(ri.AuthInfo, credentials.NoSecurity); err != nil {
		return nil, fmt.Errorf("unable to transfer TokenSource PerRPCCredentials: %v", err)
	}
	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (ts insecureTokenSource) RequireTransportSecurity() bool {
	return false
}
