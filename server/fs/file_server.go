//go:build !no_server && !gin_server
// +build !no_server,!gin_server

package fs

import (
	"fmt"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/server/middleware/auth"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"html/template"
	"net/http"
)

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string, serverUsers string, serverTemplate string) error {
	enableFileApi := false

	t, err := template.ParseGlob(serverTemplate)
	if err != nil {
		return err
	}

	// init session
	store, err := server.DefaultSessionStore()
	if err != nil {
		return err
	}

	http.Handle("/", auth.Auth(handler.NewDefaultHandler(), store))

	http.HandleFunc("/login/index", func(writer http.ResponseWriter, request *http.Request) {
		t.ExecuteTemplate(writer, "login.html", nil)
	})

	http.Handle("/login/submit", auth.NewLoginHandler(store, serverUsers))

	if src.IsDisk() || src.Is(core.RemoteDisk) {
		http.Handle(server.SrcRoutePrefix, auth.Auth(http.StripPrefix(server.SrcRoutePrefix, http.FileServer(http.Dir(src.Path()))), store))
		enableFileApi = true
	}

	if target.IsDisk() {
		http.Handle(server.TargetRoutePrefix, auth.Auth(http.StripPrefix(server.TargetRoutePrefix, http.FileServer(http.Dir(target.Path()))), store))
		enableFileApi = true
	}

	if enableFileApi {
		http.Handle(server.QueryRoute, auth.Auth(handler.NewFileApiHandler(http.Dir(src.Path())), store))
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
