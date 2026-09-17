package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config := LoadConfig()

	// Open the PostgreSQL connection.
	db := Database()
	defer db.Close()

	r := gin.Default()

	r.GET("/health", health())
	r.POST("/orders", CreateOrder(db))
	r.GET("/orders/:id", GetOrder(db))
	r.GET("/orders", GetOrders(db))
	r.DELETE("/orders/:id", DeleteOrder(db))
	r.PATCH("/orders/:id/status", UpdateOrderStatus(db))

	if err := r.Run(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
