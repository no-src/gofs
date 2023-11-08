package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/logger"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/server"
)

type reportHandler struct {
	logger   *logger.Logger
	reporter report.Reporter
}

// NewReportHandlerFunc returns a gin.HandlerFunc that providers a report api to show the application status
func NewReportHandlerFunc(logger *logger.Logger, reporter report.Reporter) gin.HandlerFunc {
	return (&reportHandler{
		logger:   logger,
		reporter: reporter,
	}).Handle
}

func (h *reportHandler) Handle(c *gin.Context) {
	r := h.reporter.GetReport()
	c.JSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, r))
}
