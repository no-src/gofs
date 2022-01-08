package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/log"
	"net/http"
)

type authHandler struct {
	logger log.Logger
}

func NewAuthHandler(logger log.Logger) handler.GinHandler {
	return &authHandler{
		logger: logger,
	}
}

func (h *authHandler) Handle(c *gin.Context) {
	ip, _ := c.RemoteIP()
	ip.IsPrivate()
	session := sessions.Default(c)
	if session == nil {
		h.logger.Error(errors.New("session is nil"), "auth handler => get session error, remote=%s", c.Request.RemoteAddr)
	}
	var user interface{}
	if session != nil {
		user = session.Get(server.SessionUser)
	}
	if user == nil {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(fmt.Sprintf("<html><head><script>window.location.href='%s';</script></head></html>", server.LoginIndexFullRoute)))
		c.Abort()
	}
}
