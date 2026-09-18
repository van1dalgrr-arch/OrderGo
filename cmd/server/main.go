package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"orderApi/internal/cache"
	"orderApi/internal/config"
	"orderApi/internal/database"
	"orderApi/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	redis := cache.NewRedis(cfg)

	if err := redis.Ping(); err != nil {
		log.Fatal(fmt.Errorf("failed to ping redis: %w", err))
	}

	db := database.Database(cfg)
	defer db.Close()

	r := gin.Default()

	r.GET("/health", handlers.Health())
	r.GET("/orders/status/:status", handlers.SortOrdersByStatus(db))
	r.POST("/orders", handlers.CreateOrder(db))
	r.GET("/orders/:id", handlers.GetOrder(db, redis))
	r.GET("/orders", handlers.GetOrders(db))
	r.DELETE("/orders/:id", handlers.DeleteOrder(db))
	r.PATCH("/orders/:id/status", handlers.UpdateOrderStatus(db))

	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
