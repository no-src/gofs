package auth

import (
	"github.com/gorilla/sessions"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
	"strings"
)

type loginHandler struct {
	store sessions.Store
	users []*User
}

func NewLoginHandler(store sessions.Store, serverUsers string) http.Handler {
	return &loginHandler{
		store: store,
		users: parseUsers(serverUsers),
	}
}

func (h *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			writer.Write([]byte("user login error"))
		}
	}()
	request.ParseMultipartForm(32 << 20) // 32 MB
	userName := request.PostForm.Get("username")
	password := request.PostForm.Get("password")

	var loginUser *User
	for _, user := range h.users {
		if user.UserName == userName && user.Password == password {
			loginUser = user
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
		http.Redirect(writer, request, "/", http.StatusFound)
	} else {
		log.Debug("login failed, username=%s password=%s", userName, password)
		http.Redirect(writer, request, "/login/index", http.StatusFound)
	}
}

func parseUsers(serverUsers string) []*User {
	var users []*User
	if len(serverUsers) == 0 {
		return users
	}
	all := strings.Split(serverUsers, ",")
	userCount := 0
	for _, user := range all {
		userInfo := strings.Split(user, "|")
		if len(userInfo) == 2 {
			userName := strings.TrimSpace(userInfo[0])
			password := strings.TrimSpace(userInfo[1])
			if len(userName) > 0 && len(password) > 0 {
				userCount++
				users = append(users, &User{
					UserId:   userCount,
					UserName: userName,
					Password: password,
				})
			}
		}
	}
	return users
}
