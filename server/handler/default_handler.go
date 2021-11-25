package handler

import (
	"fmt"
	"github.com/no-src/gofs/server"
	"net/http"
)

type defaultHandler struct {
}

func NewDefaultHandler() http.Handler {
	return &defaultHandler{}
}

func (h *defaultHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		e := recover()
		if e != nil {
			writer.Write([]byte("server internal error"))
		}
	}()
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Write([]byte(fmt.Sprintf("<html><head><title>gofs</title></head><body><div><p>welcome to gofs!</p><pre><a target='_blank' href='%s'>source</a></pre><pre><a target='_blank' href='%s'>target</a></pre></div></body></html>", server.SrcRoutePrefix, server.TargetRoutePrefix)))
}
