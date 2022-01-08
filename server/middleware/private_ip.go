package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/server/handler"
	"github.com/no-src/log"
	"net"
	"net/http"
)

type privateIPHandler struct {
	logger log.Logger
}

func NewPrivateIPHandler(logger log.Logger) handler.GinHandler {
	return &privateIPHandler{
		logger: logger,
	}
}

func (h *privateIPHandler) Handle(c *gin.Context) {
	ip := net.ParseIP(c.ClientIP())
	if !ip.IsPrivate() {
		h.logger.Warn("access deny, client ip is [%s], path is [%s]", c.ClientIP(), c.FullPath())
		c.String(http.StatusUnauthorized, "access deny")
		c.Abort()
	}
}
