//go:build tran_test

package tran

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestTcpClient_Connect_WithInvalidCertFile(t *testing.T) {
	port := getServerPort()

	// client
	client := NewClient(serverHost, port, true, notExistCertFile, false)
	err := client.Connect()
	if err == nil {
		t.Errorf("Connect: expect get file not exist error but get nil")
		return
	}
	if !os.IsNotExist(err) {
		t.Errorf("Connect: expect get file not exist error but get =>%v", err)
		return
	}

	// close
	err = client.Close()
	if err != nil {
		t.Errorf("Close: close tcp client error => %v", err)
	}
}

func TestTcpClient_Connect_WithNotStartedServer(t *testing.T) {
	port := getServerPort()

	// client
	client := NewClient(serverHost, port, true, certFile, false)
	err := client.Connect()
	if err == nil {
		t.Errorf("Connect: expect to get an error but get nil")
		return
	}
	// close
	err = client.Close()
	if err != nil {
		t.Errorf("Close: close tcp client error => %v", err)
		return
	}

	err = client.Write([]byte("hello"))
	if err == nil {
		t.Errorf("Write: expect to get an error but get nil")
		return
	}
}

func TestTcpClient_Close(t *testing.T) {
	users := getRandomUser(t)
	port := getServerPort()

	// server
	server := NewServer(serverHost, port, true, certFile, certKey, users)
	err := server.Listen()
	if err != nil {
		t.Errorf("Listen: the tcp server listen error =>%v", err)
		return
	}
	go server.Accept(func(client *Conn, data []byte) {
		// return error data
		client.Write(getServerResponseMockErrorData())
	})
	t.Logf("tcp server started, host=%s port=%d", server.Host(), server.Port())

	// client
	client := NewClient(serverHost, port, true, certFile, true)
	err = client.Connect()
	if err != nil {
		server.Close()
		t.Errorf("tcp client connect to tcp server error => %v", err)
		return
	}
	t.Logf("tcp client conneted, host=%s port=%d", client.Host(), client.Port())

	// close server first
	err = server.Close()
	if err != nil {
		t.Errorf("Close: close tcp server error => %v", err)
		return
	}

	// communication
	err = client.Write([]byte("hello server"))
	if err != nil {
		t.Errorf("Write: send tcp client data error => %v", err)
		return
	}

	// read data before client closed
	_, err = client.ReadAll()
	if err == nil || !errors.Is(err, ErrServerExecute) {
		t.Errorf("ReadAll: tcp client expect get error %v, but get => %v", ErrServerExecute, err)
		return
	}

	// close client
	err = client.Close()
	if err != nil {
		t.Errorf("Close: close tcp client error => %v", err)
		return
	}

	// read data after client closed
	_, err = client.ReadAll()
	if err == nil || !errors.Is(err, errClientNotConnected) {
		t.Errorf("ReadAll: tcp client expect get error %v, but get => %v", errClientNotConnected, err)
		return
	}
}

func TestTcpClient_CheckAndTagState(t *testing.T) {
	port := getServerPort()

	// client
	client := NewClient(serverHost, port, true, certFile, false).(*tcpClient)
	err := errors.New("syscall error")
	testCases := []struct {
		name   string
		err    error
		expect bool
	}{
		{"nil error", nil, false},
		{"errClientNotConnected", errClientNotConnected, false},
		{"io EOF", io.EOF, true},
		{"SyscallError wsarecv", fmt.Errorf("%w", os.NewSyscallError("wsarecv", err)), true},
		{"SyscallError connectex", fmt.Errorf("%w", os.NewSyscallError("connectex", err)), true},
		{"SyscallError read", fmt.Errorf("%w", os.NewSyscallError("read", err)), true},
		{"SyscallError connect", fmt.Errorf("%w", os.NewSyscallError("connect", err)), true},
		{"SyscallError unknown error", fmt.Errorf("%w", os.NewSyscallError("unknown error", err)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := client.checkAndTagState(tc.err)
			if actual != tc.expect {
				t.Errorf("checkAndTagState: expect get %v, but actual get %v", tc.expect, actual)
			}
		})
	}
}
