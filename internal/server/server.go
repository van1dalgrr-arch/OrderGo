package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"orderApi/internal/handler"
)

func New(port int) *http.Server {
	r := gin.Default()

	r.GET("/health", handler.Health())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
}
