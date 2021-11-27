//go:build gin_server
// +build gin_server

package fs

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/server/middleware/auth"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"io"
	"net/http"
	"os"
)

// StartFileServer start a file server by gin
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string, serverUsers string, serverTemplate string) error {
	enableFileApi := false

	// init log
	ginLog, err := os.OpenFile("gin.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	gin.DefaultWriter = io.MultiWriter(ginLog, os.Stdout)

	engine := gin.Default()
	engine.LoadHTMLGlob(serverTemplate)

	// init session
	store, err := server.DefaultSessionStore()
	if err != nil {
		return err
	}
	engine.Use(sessions.Sessions(server.SessionName, store))

	loginGroup := engine.Group("/login")

	loginGroup.GET("/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})

	loginGroup.POST("/submit", func(context *gin.Context) {
		auth.NewLoginHandler(store, serverUsers).ServeHTTP(context.Writer, context.Request)
	})

	rootGroup := engine.Group("/").Use(func(context *gin.Context) {
		auth.NewAuthHandler(store).ServeHTTP(context.Writer, context.Request)
	})

	rootGroup.GET("/", func(context *gin.Context) {
		handler.NewDefaultHandler().ServeHTTP(context.Writer, context.Request)
	})

	if src.IsDisk() || src.Is(core.RemoteDisk) {
		rootGroup.StaticFS(server.SrcRoutePrefix, http.Dir(src.Path()))
		enableFileApi = true
	}

	if target.IsDisk() {
		rootGroup.StaticFS(server.TargetRoutePrefix, http.Dir(target.Path()))
		enableFileApi = true
	}

	if enableFileApi {
		rootGroup.GET(server.QueryRoute, func(context *gin.Context) {
			handler.NewFileApiHandler(http.Dir(src.Path())).ServeHTTP(context.Writer, context.Request)
		})
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

		return engine.RunTLS(addr, certFile, keyFile)
	} else {
		log.Warn("file server is not a security connection, you need the https replaced maybe!")
		return engine.Run(addr)
	}
}
