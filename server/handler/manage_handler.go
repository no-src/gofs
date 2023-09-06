package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/server"
)

type manageHandler struct {
	logger *logger.Logger
	conf   conf.Config
}

// NewManageHandlerFunc returns a gin.HandlerFunc that shows the application config
func NewManageHandlerFunc(logger *logger.Logger, conf conf.Config) gin.HandlerFunc {
	return (&manageHandler{
		logger: logger,
		conf:   conf,
	}).Handle
}

func (h *manageHandler) Handle(c *gin.Context) {
	format := strings.ToLower(c.Query(server.ParamFormat))
	config := h.conf
	mask := "******"
	if len(config.Users) > 0 {
		config.Users = mask
	}
	if len(config.SessionConnection) > 0 {
		config.SessionConnection = mask
	}
	if len(config.EncryptSecret) > 0 {
		config.EncryptSecret = mask
	}
	if len(config.DecryptSecret) > 0 {
		config.DecryptSecret = mask
	}
	result := server.NewApiResult(contract.Success, contract.SuccessDesc, config)
	if format == conf.YamlFormat.Name() {
		c.YAML(http.StatusOK, result)
	} else {
		c.PureJSON(http.StatusOK, result)
	}
}
