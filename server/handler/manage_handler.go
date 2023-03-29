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
	addr   string
}

// NewManageHandlerFunc returns a gin.HandlerFunc that shows the application config
func NewManageHandlerFunc(logger log.Logger, addr string) gin.HandlerFunc {
	return (&manageHandler{
		logger: logger,
		addr:   addr,
	}).Handle
}

func (h *manageHandler) Handle(c *gin.Context) {
	format := strings.ToLower(c.Query(server.ParamFormat))
	var result server.ApiResult
	// copy the config and mask the user info for security
	cp := conf.GetGlobalConfig(h.addr)
	if cp == nil {
		result = server.NewErrorApiResult(contract.NotFound, contract.NotFoundDesc)
	} else {
		config := *cp
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
		result = server.NewApiResult(contract.Success, contract.SuccessDesc, config)
	}
	if format == conf.YamlFormat.Name() {
		c.YAML(http.StatusOK, result)
	} else {
		c.PureJSON(http.StatusOK, result)
	}
}
