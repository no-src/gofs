package handler

import (
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"html/template"
	"net/http"
)

type defaultHandler struct {
	serverTemplate string
}

func NewDefaultHandler(serverTemplate string) http.Handler {
	return &defaultHandler{
		serverTemplate: serverTemplate,
	}
}

func (h *defaultHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			writer.Write([]byte("server internal error"))
		}
	}()
	t, err := template.ParseGlob(h.serverTemplate)
	if err != nil {
		log.Error(err, "parse template error")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("parse template error"))
		return
	}
	t.ExecuteTemplate(writer, "index.html", struct {
		Src    string
		Target string
	}{
		server.SrcRoutePrefix,
		server.TargetRoutePrefix,
	})
}
