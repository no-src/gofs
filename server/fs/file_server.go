//go:build http_server
// +build http_server

package fs

import (
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/server/middleware/auth"
	"github.com/no-src/log"
	"html/template"
	"net/http"
	"path/filepath"
)

// StartFileServer start a file server
func StartFileServer(opt server.Option) error {
	enableFileApi := false
	src := opt.Src
	target := opt.Target

	err := server.ReleaseTemplate(filepath.Dir(opt.ServerTemplate))
	if err != nil {
		log.Error(err, "release template resource error")
		return err
	}

	t, err := template.ParseGlob(opt.ServerTemplate)
	if err != nil {
		return err
	}

	// init session
	store, err := server.DefaultSessionStore()
	if err != nil {
		return err
	}

	authFunc := auth.Auth
	if len(opt.Users) == 0 {
		server.PrintAnonymousAccessWarning()
		authFunc = auth.NoAuth
	}

	if opt.EnableCompress {
		log.Warn("the file server doesn't support response compress yet")
	}

	http.Handle("/", authFunc(handler.NewDefaultHandler(opt.ServerTemplate), store))

	http.HandleFunc(server.LoginIndexFullRoute, func(writer http.ResponseWriter, request *http.Request) {
		t.ExecuteTemplate(writer, "login.html", nil)
	})

	http.Handle(server.LoginRoute+server.LoginSignInRoute, auth.NewLoginHandler(store, opt.Users))

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

	log.Log("file server [%s] starting...", opt.Addr)
	server.InitServerInfo(opt.Addr, opt.EnableTLS)
	opt.Init.Done()

	if opt.EnableTLS {
		return http.ListenAndServeTLS(opt.Addr, opt.CertFile, opt.KeyFile, nil)
	} else {
		log.Warn("file server is not a security connection, you need the https replaced maybe!")
		return http.ListenAndServe(opt.Addr, nil)
	}
}
