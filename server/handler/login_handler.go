package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/server"
)

type loginHandler struct {
	users  []*auth.User
	logger *logger.Logger
}

// NewLoginHandlerFunc returns a gin.HandlerFunc that providers a login api
func NewLoginHandlerFunc(users []*auth.User, logger *logger.Logger) gin.HandlerFunc {
	return (&loginHandler{
		users:  users,
		logger: logger,
	}).Handle
}

func (h *loginHandler) Handle(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			h.logger.Error(fmt.Errorf("%v", e), "user login error")
			c.String(http.StatusOK, "user login error")
		}
	}()

	userName := c.PostForm(server.ParamUserName)
	password := c.PostForm(server.ParamPassword)
	returnUrl := c.PostForm(server.ParamReturnUrl)
	if len(returnUrl) == 0 {
		returnUrl = "/"
	} else {
		_, parseErr := url.Parse(returnUrl)
		if parseErr != nil {
			returnUrl = "/"
		}
	}

	var loginUser *auth.SessionUser
	for _, user := range h.users {
		if user.UserName() == userName && user.Password() == password {
			loginUser = auth.MapperToSessionUser(user)
		}
	}
	if loginUser != nil {
		session := sessions.Default(c)
		if session == nil {
			h.logger.Error(errors.New("session is nil"), "login handler => get session error, remote=%s", c.Request.RemoteAddr)
			c.String(http.StatusInternalServerError, "get session error")
			return
		}
		session.Set(server.SessionUser, loginUser)
		err := session.Save()
		if err != nil {
			h.logger.Error(err, "save session error, remote=%s", c.Request.RemoteAddr)
			c.String(http.StatusInternalServerError, "save session error")
			return
		}
		h.logger.Info("login success, userid=%d username=%s remote=%s", loginUser.UserId, loginUser.UserName, c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, returnUrl)
	} else {
		h.logger.Info("login failed, username=%s password=%s remote=%s", userName, password, c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, server.LoginIndexFullRoute)
	}
}
