package auth

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
	"net/url"
)

type loginHandler struct {
	store sessions.Store
	users []*auth.User
}

func NewLoginHandler(store sessions.Store, users []*auth.User) http.Handler {
	return &loginHandler{
		store: store,
		users: users,
	}
}

func (h *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			log.Error(fmt.Errorf("%v", e), "user login error")
			writer.Write([]byte("user login error"))
		}
	}()

	request.ParseMultipartForm(32 << 20) // 32 MB
	userName := request.PostForm.Get(server.ServerParamUserName)
	password := request.PostForm.Get(server.ServerParamPassword)
	returnUrl := request.PostForm.Get(server.ServerParamReturnUrl)
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
		session, err := h.store.New(request, server.SessionName)
		if err != nil && session == nil {
			log.Error(err, "login handler => get session error, remote=%s", request.RemoteAddr)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("get session error"))
			return
		}
		session.Values[server.SessionUser] = loginUser
		err = session.Save(request, writer)
		if err != nil {
			log.Error(err, "save session error, remote=%s", request.RemoteAddr)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("save session error"))
			return
		}
		log.Debug("login success, userid=%d username=%s password=%s remote=%s", loginUser.UserId, loginUser.UserName, loginUser.Password, request.RemoteAddr)
		http.Redirect(writer, request, returnUrl, http.StatusFound)
	} else {
		log.Debug("login failed, username=%s password=%s remote=%s", userName, password, request.RemoteAddr)
		http.Redirect(writer, request, server.LoginIndexFullRoute, http.StatusFound)
	}
}
