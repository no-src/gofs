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

// NewDefaultHandlerFunc returns a gin.HandlerFunc that shows the default home page
func NewDefaultHandlerFunc(logger log.Logger) gin.HandlerFunc {
	return (&defaultHandler{
		logger: logger,
	}).Handle
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
