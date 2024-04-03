package server

import (
	"testing"

	"github.com/no-src/gofs/contract"
)

func TestNewServerErrorResult(t *testing.T) {
	r := NewServerErrorResult()
	if r.Code != contract.ServerError || r.Message != contract.ServerErrorDesc {
		t.Errorf("expect: code=%d message=%s, but actual: code=%d message=%s", contract.ServerError, contract.ServerErrorDesc, r.Code, r.Message)
	}
}
