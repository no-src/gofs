package handler

import (
	"github.com/gin-gonic/gin"
)

type emptyHandler struct {
}

func NewEmptyHandler() GinHandler {
	return &emptyHandler{}
}

func (h *emptyHandler) Handle(c *gin.Context) {}
