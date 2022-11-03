//go:build tran_test

package tran

import (
	"bytes"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/no-src/gofs/auth"
)

var (
	certFile               = "../util/httputil/testdata/cert.pem"
	certKey                = "../util/httputil/testdata/key.pem"
	notExistCertFile       = "./not_exist_cert.pem"
	serverHost             = "127.0.0.1"
	serverPort       int32 = 9630
	serverResponse         = "world"
)

func getRandomUser(t *testing.T) []*auth.User {
	users, err := auth.RandomUser(1, 8, 8, "rwx")
	if err != nil {
		t.Fatalf("generate random user error =>%v", err)
	}
	return users
}

func getServerPort() int {
	return int(atomic.AddInt32(&serverPort, 1))
}

func getServerResponseMockData() []byte {
	data := []byte(serverResponse)
	data = append(data, EndIdentity...)
	data = append(data, LFBytes...)
	return data
}

func TestTcpServer_Listen_WithNilUser(t *testing.T) {
	users := getRandomUser(t)
	port := getServerPort()
	users = append(users, nil)

	// server
	server := NewServer(serverHost, port, true, certFile, certKey, users)
	err := server.Listen()
	if err != nil {
		t.Errorf("Listen: the tcp server listen error =>%v", err)
		return
	}
	server.Close()
}

func TestTcpServer_Listen_WithInvalidCertFile(t *testing.T) {
	users := getRandomUser(t)
	port := getServerPort()

	// server
	server := NewServer(serverHost, port, true, notExistCertFile, certKey, users)
	err := server.Listen()
	if err == nil {
		t.Errorf("Listen: expect get file not exist error but get nil")
		return
	}
	if !os.IsNotExist(err) {
		t.Errorf("Listen: expect get file not exist error but get =>%v", err)
		return
	}
	server.Close()
}

func TestTcpServer_Send(t *testing.T) {
	users := getRandomUser(t)
	port := getServerPort()

	// server
	server := NewServer(serverHost, port, true, certFile, certKey, users)
	err := server.Listen()
	if err != nil {
		t.Errorf("Listen: the tcp server listen error =>%v", err)
		return
	}
	defer server.Close()

	go server.Accept(func(client *Conn, data []byte) {
		hashUser := users[0].ToHashUser()
		hashUser.RefreshExpires()
		authed, perm := server.Auth(hashUser)
		if authed {
			hashUser.Perm = perm
			client.MarkAuthorized(hashUser)
		}
		client.Write(getServerResponseMockData())
	})
	t.Logf("tcp server started, host=%s port=%d", server.Host(), server.Port())

	// client
	client := NewClient(serverHost, port, true, certFile, true)
	err = client.Connect()
	if err != nil {
		t.Errorf("tcp client connect to tcp server error => %v", err)
		return
	}
	defer client.Close()

	// communication
	client.Write([]byte("hello server"))
	// wait for authorized
	time.Sleep(time.Second * 2)
	server.Send([]byte("hello client"))
}

func TestTcpServer_Listen_DisableTLS(t *testing.T) {
	port := getServerPort()

	// server
	server := NewServer(serverHost, port, false, "", "", nil)
	err := server.Listen()
	if err != nil {
		t.Errorf("Listen: the tcp server listen error =>%v", err)
		return
	}
	defer server.Close()

	go server.Accept(func(client *Conn, data []byte) {
		client.Write(getServerResponseMockData())
	})

	// client
	client := NewClient(serverHost, port, false, "", true)
	err = client.Connect()
	if err != nil {
		t.Errorf("tcp client connect to tcp server error => %v", err)
		return
	}
	defer client.Close()

	// communication
	client.Write([]byte("hello"))
	data, err := client.ReadAll()
	if err != nil {
		t.Errorf("tcp client read server data error => %v", err)
		return
	}
	if !bytes.Equal(bytes.TrimSpace(data), []byte(serverResponse)) {
		t.Errorf("tcp client read server data expect => %s, but get => %v", serverResponse, string(data))
		return
	}
}
