package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/server"
)

type authHandler struct {
	logger *logger.Logger
	perm   auth.Perm
}

// NewAuthHandlerFunc returns a middleware that checks whether the user is sign in
func NewAuthHandlerFunc(logger *logger.Logger, perm string) gin.HandlerFunc {
	p := auth.ToPermWithDefault(perm, auth.DefaultPerm)
	if !p.IsValid() {
		logger.Warn("the auth middleware get an invalid permission")
	}
	return (&authHandler{
		logger: logger,
		perm:   p,
	}).Handle
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
		c.JSON(http.StatusUnauthorized, server.NewApiResult(contract.NoPermission, contract.NoPermissionDesc, nil))
	}
}
