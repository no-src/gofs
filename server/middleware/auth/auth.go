package auth

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/no-src/gofs/server"
	"net/http"
)

func Auth(h http.Handler, store sessions.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, server.SessionName)
		user := session.Values[server.SessionUser]
		if user == nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("<html><head><script>window.location.href='%s';</script></head></html>", server.LoginIndexFullRoute)))
		} else if h != nil {
			h.ServeHTTP(writer, request)
		}
	}
}
