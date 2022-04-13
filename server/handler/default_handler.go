package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
)

type defaultHandler struct {
	logger log.Logger
}

// NewDefaultHandler create an instance of the defaultHandler
func NewDefaultHandler(logger log.Logger) GinHandler {
	return &defaultHandler{
		logger: logger,
	}
}

func (h *defaultHandler) Handle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", struct {
		Source string
		Dest   string
	}{
		server.SourceRoutePrefix,
		server.DestRoutePrefix,
	})
}
