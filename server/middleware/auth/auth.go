package auth

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
)

func Auth(h http.Handler, store sessions.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, err := store.Get(request, server.SessionName)
		if err != nil {
			log.Error(err, "auth handler => get session error, remote=%s", request.RemoteAddr)
		}
		var user interface{}
		if err == nil {
			user = session.Values[server.SessionUser]
		}
		if user == nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("<html><head><script>window.location.href='%s';</script></head></html>", server.LoginIndexFullRoute)))
		} else if h != nil {
			h.ServeHTTP(writer, request)
		}
	}
}

func NoAuth(h http.Handler, store sessions.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if h != nil {
			h.ServeHTTP(writer, request)
		}
	}
}
