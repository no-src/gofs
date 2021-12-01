package server

import (
	"fmt"
	"github.com/no-src/gofs"
	"github.com/no-src/log"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var serverAddr *net.TCPAddr
var enableTLS bool

const (
	SrcRoutePrefix       = "/src/"
	TargetRoutePrefix    = "/target/"
	QueryRoute           = "/query"
	LoginRoute           = "/login"
	LoginIndexRoute      = "/index"
	LoginIndexFullRoute  = LoginRoute + LoginIndexRoute
	LoginSignInRoute     = "/signin"
	LoginSignInFullRoute = LoginRoute + LoginSignInRoute
)

const (
	DefaultAddrHttps = ":443"
	DefaultAddrHttp  = ":80"
	SchemeHttp       = "http"
	SchemeHttps      = "https"
	DefaultPortHttp  = 80
	DefaultPortHttps = 443
)

const (
	SessionName = "session_id"
	SessionUser = "user"
)

const (
	ResourceTemplatePath = "server/template"
)

func InitServerInfo(addr string, tls bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err == nil {
		serverAddr = tcpAddr
	} else {
		log.Error(err, "invalid server addr => %s", addr)
	}
	enableTLS = tls
}

// ServerAddr the addr of file server running
func ServerAddr() *net.TCPAddr {
	return serverAddr
}

// ServerPort the port of file server running
func ServerPort() int {
	if serverAddr != nil {
		return serverAddr.Port
	}
	return 0
}

// EnableTLS is using https on the file server
func EnableTLS() bool {
	return enableTLS
}

// GenerateAddr generate http or https address
func GenerateAddr(scheme, host string, port int) string {
	addr := ""
	scheme = strings.ToLower(scheme)
	if scheme == SchemeHttp && port == DefaultPortHttp {
		addr = fmt.Sprintf("%s://%s", SchemeHttp, host)
	} else if scheme == SchemeHttps && port == DefaultPortHttps {
		addr = fmt.Sprintf("%s://%s", SchemeHttps, host)
	} else {
		addr = fmt.Sprintf("%s://%s:%d", scheme, host, port)
	}
	return addr
}

func PrintAnonymousAccessWarning() {
	log.Warn("the file server allows anonymous access, you should set some server users by the -users or -rand_user_count flag for security reasons")
}

func ReleaseTemplate(releasePath string) error {
	files, err := gofs.Templates.ReadDir(ResourceTemplatePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.ToSlash(releasePath), os.ModePerm)
	if err != nil {
		return err
	}
	for _, f := range files {
		srcPath := filepath.ToSlash(filepath.Join(ResourceTemplatePath, f.Name()))
		targetPath := filepath.ToSlash(filepath.Join(releasePath, f.Name()))
		if f.IsDir() {
			os.MkdirAll(targetPath, os.ModePerm)
		} else {
			data, err := gofs.Templates.ReadFile(srcPath)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(targetPath, data, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
