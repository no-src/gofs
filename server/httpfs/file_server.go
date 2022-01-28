//go:build !no_server
// +build !no_server

package httpfs

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/gofs/server/middleware"
	"github.com/no-src/log"
	"html/template"
	"net/http"
	"os"
	"time"
)

// StartFileServer start a file server by gin
func StartFileServer(opt server.Option) error {
	enableFileApi := false
	source := opt.Source
	dest := opt.Dest
	logger := opt.Logger

	// change default mode is release
	mode := os.Getenv(gin.EnvGinMode)
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	gin.DefaultWriter = logger

	engine := gin.New()

	engine.NoRoute(middleware.NoRoute)

	if opt.EnableCompress {
		// enable gzip compression
		engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: defaultLogFormatter,
		Output:    logger,
	}), gin.Recovery())

	tmpl, err := template.ParseFS(gofs.Templates, server.ResourceTemplatePath)
	if err != nil {
		log.Error(err, "parse template fs error")
		return err
	}
	engine.SetHTMLTemplate(tmpl)

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

	loginGroup.POST(server.LoginSignInRoute, handler.NewLoginHandler(opt.Users, logger).Handle)

	rootGroup := engine.Group("/")

	if len(opt.Users) > 0 {
		rootGroup.Use(middleware.NewAuthHandler(logger).Handle)
	} else {
		server.PrintAnonymousAccessWarning()
	}

	rootGroup.GET("/", handler.NewDefaultHandler(logger).Handle)

	if opt.EnablePprof {
		debugGroup := rootGroup.Group("/debug")
		if opt.PprofPrivate {
			debugGroup.Use(middleware.NewPrivateAccessHandler(logger).Handle)
		}
		pprof.RouteRegister(debugGroup, "pprof")
	}

	if source.IsDisk() || source.Is(core.RemoteDisk) {
		rootGroup.StaticFS(server.SourceRoutePrefix, http.Dir(source.Path()))
		enableFileApi = true
	}

	if dest.IsDisk() {
		rootGroup.StaticFS(server.DestRoutePrefix, http.Dir(dest.Path()))
		enableFileApi = true
	}

	if enableFileApi {
		rootGroup.GET(server.QueryRoute, handler.NewFileApiHandler(http.Dir(source.Path()), logger).Handle)
	}

	log.Info("file server [%s] starting...", opt.Addr)
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
