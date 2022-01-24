package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
)

type defaultHandler struct {
	logger log.Logger
}

func NewDefaultHandler(logger log.Logger) GinHandler {
	return &defaultHandler{
		logger: logger,
	}
}

func (h *defaultHandler) Handle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", struct {
		Src  string
		Dest string
	}{
		server.SrcRoutePrefix,
		server.DestRoutePrefix,
	})
}
