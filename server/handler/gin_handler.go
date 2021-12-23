package handler

import "github.com/gin-gonic/gin"

type GinHandler interface {
	Handle(*gin.Context)
}
