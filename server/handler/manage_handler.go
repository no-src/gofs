package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
)

type manageHandler struct {
	logger log.Logger
}

// NewManageHandlerFunc returns a gin.HandlerFunc that shows the application config
func NewManageHandlerFunc(logger log.Logger) gin.HandlerFunc {
	return (&manageHandler{
		logger: logger,
	}).Handle
}

func (h *manageHandler) Handle(c *gin.Context) {
	format := strings.ToLower(c.Query(server.ParamFormat))
	// copy the config and mask the user info for security
	config := *conf.GlobalConfig
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
	if format == conf.YamlFormat.Name() {
		c.YAML(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, config))
	} else {
		c.PureJSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, config))
	}
}
