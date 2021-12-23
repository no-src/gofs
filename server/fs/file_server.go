//go:build !no_server
// +build !no_server

package fs

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/server/middleware"
	"github.com/no-src/log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// StartFileServer start a file server by gin
func StartFileServer(opt server.Option) error {
	enableFileApi := false
	src := opt.Src
	target := opt.Target

	err := server.ReleaseTemplate(filepath.Dir(opt.ServerTemplate), opt.ServerTemplateOverride)
	if err != nil {
		log.Error(err, "release template resource error")
		return err
	}

	// change default mode is release
	mode := os.Getenv(gin.EnvGinMode)
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	gin.DefaultWriter = log.DefaultLogger()

	engine := gin.New()
	if opt.EnableCompress {
		// enable gzip compression
		engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: defaultLogFormatter,
		Output:    log.DefaultLogger(),
	}), gin.Recovery())
	engine.LoadHTMLGlob(opt.ServerTemplate)

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

	loginGroup.POST(server.LoginSignInRoute, middleware.NewLoginHandler(opt.Users).Handle)

	rootGroup := engine.Group("/")

	if len(opt.Users) > 0 {
		rootGroup.Use(middleware.NewAuthHandler().Handle)
	} else {
		server.PrintAnonymousAccessWarning()
	}

	rootGroup.GET("/", handler.NewDefaultHandler().Handle)

	if src.IsDisk() || src.Is(core.RemoteDisk) {
		rootGroup.StaticFS(server.SrcRoutePrefix, http.Dir(src.Path()))
		enableFileApi = true
	}

	if target.IsDisk() {
		rootGroup.StaticFS(server.TargetRoutePrefix, http.Dir(target.Path()))
		enableFileApi = true
	}

	if enableFileApi {
		rootGroup.GET(server.QueryRoute, handler.NewFileApiHandler(http.Dir(src.Path())).Handle)
	}

	log.Log("file server [%s] starting...", opt.Addr)
	server.InitServerInfo(opt.Addr, opt.EnableTLS)
	opt.Init.Done()

	if opt.EnableTLS {
		return engine.RunTLS(opt.Addr, opt.CertFile, opt.KeyFile)
	} else {
		log.Warn("file server is not a security connection, you need the https replaced maybe!")
		return engine.Run(opt.Addr)
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
