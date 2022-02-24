package handler

import "github.com/gin-gonic/gin"

// GinHandler the handler interface of Gin
type GinHandler interface {
	Handle(*gin.Context)
}
