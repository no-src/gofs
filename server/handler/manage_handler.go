package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"net/http"
	"strings"
)

type manageHandler struct {
	logger log.Logger
}

// NewManageHandler create an instance of the manageHandler
func NewManageHandler(logger log.Logger) GinHandler {
	return &manageHandler{
		logger: logger,
	}
}

func (h *manageHandler) Handle(c *gin.Context) {
	format := strings.ToLower(c.Query(server.ParamFormat))
	// copy the config and mask the user info for security
	config := *conf.GlobalConfig
	if len(config.Users) > 0 {
		config.Users = "******"
	}
	if format == conf.YamlFormat.Name() {
		c.YAML(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, config))
	} else {
		c.PureJSON(http.StatusOK, server.NewApiResult(contract.Success, contract.SuccessDesc, config))
	}
}
