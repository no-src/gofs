package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/log"
	"net/http"
)

type authHandler struct {
	logger log.Logger
	perm   auth.Perm
}

func NewAuthHandler(logger log.Logger, perm string) handler.GinHandler {
	p := auth.ToPermWithDefault(perm, auth.DefaultPerm)
	if !p.IsValid() {
		logger.Warn("the auth middleware get an invalid permission")
	}
	return &authHandler{
		logger: logger,
		perm:   p,
	}
}

func (h *authHandler) Handle(c *gin.Context) {
	session := sessions.Default(c)
	if session == nil {
		h.logger.Error(errors.New("session is nil"), "auth handler => get session error, remote=%s", c.Request.RemoteAddr)
	}
	var user *auth.SessionUser
	if session != nil {
		obj := session.Get(server.SessionUser)
		if obj != nil {
			tmp := obj.(auth.SessionUser)
			user = &tmp
		}
	}
	if user == nil {
		c.Abort()
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(fmt.Sprintf("<html><head><script>window.location.href='%s';</script></head></html>", server.LoginIndexFullRoute)))
	} else if !h.perm.CheckTo(user.Perm) {
		c.Abort()
		c.String(http.StatusUnauthorized, "you have no permission")
	}
}
