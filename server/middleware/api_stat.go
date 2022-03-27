package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/server/handler"
)

type apiStatHandler struct {
}

// ApiStat returns a middleware that records user access statistics data
func ApiStat() gin.HandlerFunc {
	return NewApiStatHandler().Handle
}

// NewApiStatHandler create an instance of the apiStatHandler middleware
func NewApiStatHandler() handler.GinHandler {
	return &apiStatHandler{}
}

func (h *apiStatHandler) Handle(c *gin.Context) {
	report.GlobalReporter.PutApiStat(c.ClientIP())
}
