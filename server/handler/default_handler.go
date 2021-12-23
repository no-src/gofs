package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/server"
	"net/http"
)

type defaultHandler struct {
}

func NewDefaultHandler() GinHandler {
	return &defaultHandler{}
}

func (h *defaultHandler) Handle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", struct {
		Src    string
		Target string
	}{
		server.SrcRoutePrefix,
		server.TargetRoutePrefix,
	})
}
