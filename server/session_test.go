package server

import (
	"errors"
	"testing"
)

func TestNewSessionStore(t *testing.T) {
	testCases := []struct {
		conn string
	}{
		{"memory:"},
		//{"redis://127.0.0.1:6379"},
		//{"redis://127.0.0.1:6379?password=&db=10&max_idle=10&secret=redis_secret"},
	}

	for _, tc := range testCases {
		t.Run(tc.conn, func(t *testing.T) {
			store, err := NewSessionStore(tc.conn)
			if err != nil {
				t.Errorf("create session store error => %s", err)
				return
			}
			if store == nil {
				t.Errorf("get a nil session store")
			}
		})
	}
}

func TestNewSessionStore_InvalidParameter(t *testing.T) {
	testCases := []struct {
		conn      string
		expectErr error
	}{
		{"redis://127.0.0.1:6379?db=x", errInvalidRedisDB},
		{"redis://127.0.0.1:6379?db=100", errInvalidRedisDB},
		{"redis://127.0.0.1:6379?max_idle=x", errInvalidRedisMaxIdle},
		{"redis://127.0.0.1:6379?max_idle=0", errInvalidRedisMaxIdle},
	}

	for _, tc := range testCases {
		t.Run(tc.conn, func(t *testing.T) {
			_, err := NewSessionStore(tc.conn)
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("expect to get error [%s], but actual get error [%s]", tc.expectErr, err)
			}
		})
	}
}

func TestNewSessionStore_Unsupported(t *testing.T) {
	testCases := []struct {
		conn string
	}{
		{"hello://127.0.0.1:8888"},
	}

	for _, tc := range testCases {
		t.Run(tc.conn, func(t *testing.T) {
			_, err := NewSessionStore(tc.conn)
			if !errors.Is(err, errUnsupportedSession) {
				t.Errorf("expect to get error [%s], but actual get error [%s]", errUnsupportedSession, err)
			}
		})
	}
}

func TestNewSessionStore_Invalid(t *testing.T) {
	testCases := []struct {
		conn string
	}{
		{"hello_://"},
		{"\t"},
	}

	for _, tc := range testCases {
		t.Run(tc.conn, func(t *testing.T) {
			_, err := NewSessionStore(tc.conn)
			if !errors.Is(err, errInvalidSession) {
				t.Errorf("expect to get error [%s], but actual get error [%s]", errInvalidSession, err)
			}
		})
	}
}
