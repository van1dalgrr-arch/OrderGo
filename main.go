package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func validStatus(status string) bool {
	switch status {
	case "pending", "paid", "processing", "shipped", "completed", "cancelled":
		return true
	default:
		return false
	}
}

func UpdateOrderStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id = c.Param("id")

		var req UpdateOrderStatusRequest

		// Parse the JSON request body.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("failed to bind JSON: %w", err),
			})
			return
		}

		if !validStatus(req.Status) {
			c.JSON(400, gin.H{
				"error": "invalid status",
			})
			return
		}

		// Update the order status using the ID from the URL.
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

	// Build the PostgreSQL connection string from environment variables.
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

	// Check that the database is actually reachable.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetOrder returns an order by its ID.
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

func GetOrders(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		page := c.Query("page")
		limit := c.Query("limit")
		if page == "" {
			page = "1"
		}
		if limit == "" {
			limit = "10"
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			c.JSON(400, gin.H{
				"error": "invalid page",
			})
			return
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt <= 0 {
			c.JSON(400, gin.H{
				"error": "invalid limit",
			})
			return
		}
		offset := (pageInt - 1) * limitInt
		rows, err := db.Query(
			"SELECT id, user_id, status FROM orders ORDER BY id LIMIT $1 OFFSET $2",
			limitInt,
			offset,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to get orders: %w", err),
			})
			return
		}
		defer rows.Close()
		var orders []Order
		for rows.Next() {
			var order Order
			if err := rows.Scan(
				&order.ID,
				&order.UserID,
				&order.Status,
			); err != nil {
				c.JSON(500, gin.H{
					"error": fmt.Errorf("failed to scan order: %w", err),
				})
				return
			}
			orders = append(orders, order)
		}
		if err := rows.Err(); err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to iterate orders: %w", err),
			})
			return
		}
		c.JSON(200, gin.H{
			"orders": orders,
			"page":   pageInt,
			"limit":  limitInt,
		})
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

		if order.ID == "" || order.UserID == "" || order.Status == "" {
			c.JSON(400, gin.H{
				"error": "id, user_id, and status are required",
			})
			return
		}

		if !validStatus(order.Status) {
			c.JSON(400, gin.H{
				"error": "invalid status",
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
			return
		}

		c.JSON(201, gin.H{
			"message": "order has been created",
		})
	}
}

func main() {
	config := LoadConfig()

	// Open the PostgreSQL connection.
	db := Database()
	defer db.Close()

	r := gin.Default()

	r.POST("/orders", CreateOrder(db))
	r.GET("/orders/:id", GetOrder(db))
	r.GET("/orders", GetOrders(db))
	r.DELETE("/orders/:id", DeleteOrder(db))
	r.PATCH("/orders/:id/status", UpdateOrderStatus(db))

	if err := r.Run(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
