package auth

import (
	"github.com/gorilla/sessions"
	"github.com/no-src/gofs/server"
	"net/http"
)

type authHandler struct {
	store sessions.Store
}

func NewAuthHandler(store sessions.Store) http.Handler {
	return &authHandler{
		store: store,
	}
}

func (h *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			writer.Write([]byte("user auth error"))
		}
	}()
	session, _ := h.store.Get(request, server.SessionName)
	user := session.Values[server.SessionUser]
	if user == nil {
		http.Redirect(writer, request, "/login/index", http.StatusFound)
	}
}

func Auth(h http.Handler, store sessions.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		NewAuthHandler(store).ServeHTTP(writer, request)
		h.ServeHTTP(writer, request)
	}
}
