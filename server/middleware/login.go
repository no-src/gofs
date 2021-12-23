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
	"net/url"
)

type loginHandler struct {
	users []*auth.User
}

func NewLoginHandler(users []*auth.User) handler.GinHandler {
	return &loginHandler{
		users: users,
	}
}

func (h *loginHandler) Handle(c *gin.Context) {
	defer func() {
		e := recover()
		if e != nil {
			log.Error(fmt.Errorf("%v", e), "user login error")
			c.String(http.StatusOK, "user login error")
		}
	}()

	userName := c.PostForm(server.ServerParamUserName)
	password := c.PostForm(server.ServerParamPassword)
	returnUrl := c.PostForm(server.ServerParamReturnUrl)
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
			log.Error(errors.New("session is nil"), "login handler => get session error, remote=%s", c.Request.RemoteAddr)
			c.String(http.StatusInternalServerError, "get session error")
			return
		}
		session.Set(server.SessionUser, loginUser)
		err := session.Save()
		if err != nil {
			log.Error(err, "save session error, remote=%s", c.Request.RemoteAddr)
			c.String(http.StatusInternalServerError, "save session error")
			return
		}
		log.Info("login success, userid=%d username=%s password=%s remote=%s", loginUser.UserId, loginUser.UserName, loginUser.Password, c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, returnUrl)
	} else {
		log.Info("login failed, username=%s password=%s remote=%s", userName, password, c.Request.RemoteAddr)
		c.Redirect(http.StatusFound, server.LoginIndexFullRoute)
	}
}
