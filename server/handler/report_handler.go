package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
)

type reportHandler struct {
	logger log.Logger
}

// NewReportHandlerFunc returns a gin.HandlerFunc that providers a report api to show the application status
func NewReportHandlerFunc(logger log.Logger) gin.HandlerFunc {
	return (&reportHandler{
		logger: logger,
	}).Handle
}

func (h *reportHandler) Handle(c *gin.Context) {
	r := report.GlobalReporter.GetReport()
	c.JSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, r))
}
