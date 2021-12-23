//go:build !no_server && !http_server
// +build !no_server,!http_server

package auth

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
)

// GinAuth auth middleware for gin server
func GinAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		if session == nil {
			log.Error(errors.New("session is nil"), "auth handler => get session error, remote=%s", context.Request.RemoteAddr)
		}
		var user interface{}
		if session != nil {
			user = session.Get(server.SessionUser)
		}
		if user == nil {
			context.Writer.Write([]byte(fmt.Sprintf("<html><head><script>window.location.href='%s';</script></head></html>", server.LoginIndexFullRoute)))
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
