package tran

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/log"
	"io"
	"net"
	"os"
	"sync"
)

type tcpClient struct {
	network   string
	host      string
	port      int
	innerConn net.Conn
	closed    *cbool.CBool
	enableTLS bool
	mu        sync.Mutex
}

var (
	clientNotConnected = errors.New("client is not connected")
)

// NewClient create an instance of tcpClient
func NewClient(host string, port int, enableTLS bool) Client {
	client := &tcpClient{
		host:      host,
		port:      port,
		network:   "tcp",
		closed:    cbool.New(true),
		enableTLS: enableTLS,
		mu:        sync.Mutex{},
	}
	return client
}

func (client *tcpClient) Connect() (err error) {
	address := fmt.Sprintf("%s:%d", client.host, client.port)
	if client.enableTLS {
		client.innerConn, err = tls.Dial(client.network, address, &tls.Config{
			InsecureSkipVerify: true,
		})
		// innerConn is not nil whatever err is nil or not
		if err != nil {
			client.innerConn = nil
		}
	} else {
		client.innerConn, err = net.Dial(client.network, address)
	}
	if err != nil {
		client.checkAndTagState(err)
		log.Error(err, "client connect failed")
	} else {
		client.closed.Set(false)
	}
	return err
}

func (client *tcpClient) Write(data []byte) (err error) {
	if client.IsClosed() {
		return clientNotConnected
	}
	writer := bufio.NewWriter(client.innerConn)
	data = append(data, EndIdentity...)
	data = append(data, LFBytes...)
	_, err = writer.Write(data)
	if err != nil {
		client.checkAndTagState(err)
		log.Error(err, "client write failed")
		return err
	}
	err = writer.Flush()
	if err != nil {
		client.checkAndTagState(err)
		log.Error(err, "client flush failed")
	}
	return err
}

func (client *tcpClient) isClosedError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	err = errors.Unwrap(err)
	syscallErr := &os.SyscallError{}
	if errors.As(err, &syscallErr) {
		syscall := syscallErr.Syscall
		if syscall == "wsarecv" || syscall == "connectex" || syscall == "read" || syscall == "connect" {
			return true
		} else {
			log.Error(err, "get a unknown error")
		}
	}
	return false
}

func (client *tcpClient) checkAndTagState(err error) bool {
	if client.isClosedError(err) {
		client.Close()
		return true
	}
	return false
}

func (client *tcpClient) ReadAll() (result []byte, err error) {
	if client.IsClosed() {
		return nil, clientNotConnected
	}
	reader := bufio.NewReader(client.innerConn)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			client.checkAndTagState(err)
			return result, err
		}
		isEnd := false
		endIdentity := EndIdentity
		hasError := false
		if bytes.HasSuffix(line, endIdentity) {
			isEnd = true
			if bytes.HasSuffix(line, ErrorEndIdentity) {
				endIdentity = ErrorEndIdentity
				hasError = true
			}
			line = line[:len(line)-len(endIdentity)]
		}

		result = append(result, line...)
		result = append(result, LFBytes...)

		if isEnd {
			if hasError {
				err = ServerExecuteError
				log.Error(err, string(result))
			}
			return result, err
		}
	}
}

func (client *tcpClient) Host() string {
	return client.host
}

func (client *tcpClient) Port() int {
	return client.port
}

func (client *tcpClient) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.closed.Set(true)
	if client.innerConn != nil {
		return client.innerConn.Close()
	}
	client.innerConn = nil
	return nil
}

func (client *tcpClient) IsClosed() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closed.Get() || client.innerConn == nil {
		client.closed.Set(true)
		return true
	}
	return false
}
