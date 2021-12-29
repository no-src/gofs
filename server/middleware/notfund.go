package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRoute(context *gin.Context) {
	context.String(http.StatusNotFound, "404 page not found")
}