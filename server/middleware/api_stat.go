package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/report"
)

type apiStatHandler struct {
	reporter report.Reporter
}

// ApiStat returns a middleware that records user access statistics data
func ApiStat(reporter report.Reporter) gin.HandlerFunc {
	return (&apiStatHandler{
		reporter: reporter,
	}).Handle
}

func (h *apiStatHandler) Handle(c *gin.Context) {
	h.reporter.PutApiStat(c.ClientIP())
}
