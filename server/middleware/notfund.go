package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/server"
	"net/http"
)

// NoRoute the middleware of the 404 status
func NoRoute(context *gin.Context) {
	context.JSON(http.StatusNotFound, server.NewApiResult(contract.NotFound, contract.NotFoundDesc, nil))
}
