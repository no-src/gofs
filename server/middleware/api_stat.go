package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/report"
)

type apiStatHandler struct {
}

// ApiStat returns a middleware that records user access statistics data
func ApiStat() gin.HandlerFunc {
	return (&apiStatHandler{}).Handle
}

func (h *apiStatHandler) Handle(c *gin.Context) {
	report.GlobalReporter.PutApiStat(c.ClientIP())
}
