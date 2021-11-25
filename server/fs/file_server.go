//go:build !no_server
// +build !no_server

package fs

import (
	"fmt"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"net/http"
)

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string) error {
	enableFileApi := false

	http.Handle("/", handler.NewDefaultHandler())

	if src.IsDisk() || src.Is(core.RemoteDisk) {
		http.Handle(server.SrcRoutePrefix, http.StripPrefix(server.SrcRoutePrefix, http.FileServer(http.Dir(src.Path()))))
		enableFileApi = true
	}

	if target.IsDisk() {
		http.Handle(server.TargetRoutePrefix, http.StripPrefix(server.TargetRoutePrefix, http.FileServer(http.Dir(target.Path()))))
		enableFileApi = true
	}

	if enableFileApi {
		http.Handle(server.QueryRoute, handler.NewFileApiHandler(http.Dir(src.Path())))
	}

	log.Log("file server [%s] starting...", addr)
	server.InitServerInfo(addr, enableTLS)
	init.Done()

	if enableTLS {
		exist, err := util.FileExist(certFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("cert file is not found for https => %s", certFile)
		}
		exist, err = util.FileExist(keyFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("key file is not found for https => %s", keyFile)
		}
		return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	} else {
		log.Warn("file server is not a security connection, you need the https replaced maybe!")
		return http.ListenAndServe(addr, nil)
	}
}
