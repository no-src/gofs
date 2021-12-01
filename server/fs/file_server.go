//go:build !no_server && !gin_server
// +build !no_server,!gin_server

package fs

import (
	auth2 "github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/server/middleware/auth"
	"github.com/no-src/log"
	"html/template"
	"net/http"
	"path/filepath"
)

// StartFileServer start a file server
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string, users []*auth2.User, serverTemplate string) error {
	enableFileApi := false

	err := server.ReleaseTemplate(filepath.Dir(serverTemplate))
	if err != nil {
		log.Error(err, "release template resource error")
		return err
	}

	t, err := template.ParseGlob(serverTemplate)
	if err != nil {
		return err
	}

	// init session
	store, err := server.DefaultSessionStore()
	if err != nil {
		return err
	}

	authFunc := auth.Auth
	if len(users) == 0 {
		server.PrintAnonymousAccessWarning()
		authFunc = auth.NoAuth
	}

	http.Handle("/", authFunc(handler.NewDefaultHandler(serverTemplate), store))

	http.HandleFunc(server.LoginIndexFullRoute, func(writer http.ResponseWriter, request *http.Request) {
		t.ExecuteTemplate(writer, "login.html", nil)
	})

	http.Handle(server.LoginRoute+server.LoginSignInRoute, auth.NewLoginHandler(store, users))

	if src.IsDisk() || src.Is(core.RemoteDisk) {
		http.Handle(server.SrcRoutePrefix, authFunc(http.StripPrefix(server.SrcRoutePrefix, http.FileServer(http.Dir(src.Path()))), store))
		enableFileApi = true
	}

	if target.IsDisk() {
		http.Handle(server.TargetRoutePrefix, authFunc(http.StripPrefix(server.TargetRoutePrefix, http.FileServer(http.Dir(target.Path()))), store))
		enableFileApi = true
	}

	if enableFileApi {
		http.Handle(server.QueryRoute, authFunc(handler.NewFileApiHandler(http.Dir(src.Path())), store))
	}

	log.Log("file server [%s] starting...", addr)
	server.InitServerInfo(addr, enableTLS)
	init.Done()

	if enableTLS {
		return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	} else {
		log.Warn("file server is not a security connection, you need the https replaced maybe!")
		return http.ListenAndServe(addr, nil)
	}
}
