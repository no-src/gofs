//go:build tran_test

package tran

import (
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
	client.Close()
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
	client.Close()

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
		client.Write(getServerResponseMockData())
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
	server.Close()

	// communication
	client.Write([]byte("hello server"))

	// close client
	client.Close()
}
