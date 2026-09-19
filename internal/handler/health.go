package handler

import (
	"github.com/gin-gonic/gin"
)

// check server
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server is health",
		})
	}
}
