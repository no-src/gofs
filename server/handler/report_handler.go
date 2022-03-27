package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
)

type reportHandler struct {
	logger log.Logger
}

// NewReportHandler create an instance of the reportHandler
func NewReportHandler(logger log.Logger) GinHandler {
	return &reportHandler{
		logger: logger,
	}
}

func (h *reportHandler) Handle(c *gin.Context) {
	r := report.GlobalReporter.GetReport()
	c.JSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, r))
}
