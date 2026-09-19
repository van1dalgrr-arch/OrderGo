package server

import (
	"github.com/gin-gonic/gin"
	"orderApi/internal/handler"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.Health())
	return r
}
