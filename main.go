package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func UpdateOrderStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id = c.Param("id")

		var req UpdateOrderStatusRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("failed to bind JSON: %w", err),
			})
			return
		}

		_, err := db.Exec(
			"UPDATE orders SET status = $1 WHERE id = $2",
			req.Status, id,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to update order status: %w", err),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "order status updated",
		})
	}
}

func Database() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var order Order

		err := db.QueryRow(
			"SELECT id, user_id, status FROM orders WHERE id = $1",
			id,
		).Scan(&order.ID, &order.UserID, &order.Status)

		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{
				"error": "order not found",
			})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to get order: %w", err),
			})
			return
		}

		c.JSON(200, order)
	}
}

func DeleteOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(
			"DELETE FROM orders WHERE id = $1",
			id,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to delete order: %w", err),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "order has been deleted",
		})
	}
}

func CreateOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order Order

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := db.Exec(
			"INSERT INTO orders (id, user_id, status) VALUES ($1, $2, $3)",
			order.ID,
			order.UserID,
			order.Status,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to create order: %w", err),
			})
		} else {
			c.JSON(201, gin.H{
				"message": "order has been created",
			})
			return
		}
	}
}

func main() {
	config := LoadConfig()

	db := Database()
	defer db.Close()

	r := gin.Default()

	r.POST("/orders", CreateOrder(db))
	r.GET("/orders/:id", GetOrder(db))
	r.DELETE("/orders/:id", DeleteOrder(db))
	r.PATCH("/orders/:id/status", UpdateOrderStatus(db))

	if err := r.Run(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
