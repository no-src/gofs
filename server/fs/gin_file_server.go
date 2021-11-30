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
	"net/http"
	"os"
	"time"
)

// StartFileServer start a file server by gin
func StartFileServer(src core.VFS, target core.VFS, addr string, init retry.WaitDone, enableTLS bool, certFile string, keyFile string, users []*auth.User, serverTemplate string) error {
	enableFileApi := false

	// change default mode is release
	mode := os.Getenv(gin.EnvGinMode)
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	gin.DefaultWriter = log.DefaultLogger()

	engine := gin.New()
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: defaultLogFormatter,
		Output:    log.DefaultLogger(),
	}), gin.Recovery())
	engine.LoadHTMLGlob(serverTemplate)

	// init session
	store, err := server.DefaultSessionStore()
	if err != nil {
		return err
	}
	engine.Use(sessions.Sessions(server.SessionName, store))

	loginGroup := engine.Group(server.LoginRoute)

	loginGroup.GET(server.LoginIndexRoute, func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})

	loginGroup.POST(server.LoginSignInRoute, func(context *gin.Context) {
		auth.NewLoginHandler(store, users).ServeHTTP(context.Writer, context.Request)
	})

	rootGroup := engine.Group("/").Use(func(context *gin.Context) {
		auth.Auth(nil, store).ServeHTTP(context.Writer, context.Request)
	})

	rootGroup.GET("/", func(context *gin.Context) {
		handler.NewDefaultHandler(serverTemplate).ServeHTTP(context.Writer, context.Request)
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

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[%v] [GIN] |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006-01-02 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
