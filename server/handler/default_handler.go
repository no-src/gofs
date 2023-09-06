package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/server"
)

type defaultHandler struct {
	logger *logger.Logger
}

// NewDefaultHandlerFunc returns a gin.HandlerFunc that shows the default home page
func NewDefaultHandlerFunc(logger *logger.Logger) gin.HandlerFunc {
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
