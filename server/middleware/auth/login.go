package auth

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
	"net/url"
)

type loginHandler struct {
	store sessions.Store
	users []*contract.User
}

func NewLoginHandler(store sessions.Store, users []*contract.User) http.Handler {
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

	var loginUser *contract.SessionUser
	for _, user := range h.users {
		if user.UserName() == userName && user.Password() == password {
			loginUser = contract.MapperToSessionUser(user)
		}
	}
	if loginUser != nil {
		session, err := h.store.New(request, server.SessionName)
		session.Values[server.SessionUser] = loginUser
		err = session.Save(request, writer)
		if err != nil {
			log.Error(err, "save session error")
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("save session error"))
			return
		}
		log.Debug("login success, userid=%d username=%s password=%s", loginUser.UserId, loginUser.UserName, loginUser.Password)
		http.Redirect(writer, request, returnUrl, http.StatusFound)
	} else {
		log.Debug("login failed, username=%s password=%s", userName, password)
		http.Redirect(writer, request, server.LoginIndexFullRoute, http.StatusFound)
	}
}
