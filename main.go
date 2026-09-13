package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

func getOrder(c *gin.Context) {
	id := c.Param("id")

	// TODO: get order from database
	_ = id

	c.JSON(404, gin.H{
		"error": "order not found",
	})
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	// TODO: delete order from database
	_ = id

	c.JSON(200, gin.H{
		"message": "order has been deleted",
	})
}

func createOrder(c *gin.Context) {
	var order Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: save order to database

	c.JSON(201, order)
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func main() {
	// TODO: connect to database

	r := gin.Default()

	r.GET("/health", health)
	r.POST("/orders", createOrder)
	r.GET("/orders/:id", getOrder)
	r.DELETE("/orders/:id", deleteOrder)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
