package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"orderApi/internal/models"
	"orderApi/internal/validation"

	"github.com/gin-gonic/gin"
)

func SortOrdersByUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		rows, err := db.Query(
			"SELECT id, user_id, status, created_at FROM orders WHERE user_id = $1",
			userID,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to get orders by user: %w", err).Error(),
			})
			return
		}
		defer rows.Close()
		var orders []models.Order
		for rows.Next() {
			var order models.Order
			if err := rows.Scan(
				&order.ID,
				&order.UserID,
				&order.Status,
				&order.CreatedAt,
			); err != nil {
				c.JSON(500, gin.H{
					"error": fmt.Errorf("failed to scan order: %w", err).Error(),
				})
				return
			}
			orders = append(orders, order)
		}
		if err := rows.Err(); err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to iterate orders: %w", err).Error(),
			})
			return
		}
		c.JSON(200, orders)
	}
}

func SortOrdersByStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Param("status")
		rows, err := db.Query(
			"SELECT id, user_id, status, created_at FROM orders WHERE status = $1",
			status,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to get orders by status: %w", err).Error(),
			})
			return
		}
		defer rows.Close()
		var orders []models.Order
		for rows.Next() {
			var order models.Order
			if err := rows.Scan(
				&order.ID,
				&order.UserID,
				&order.Status,
				&order.CreatedAt,
			); err != nil {
				c.JSON(500, gin.H{
					"error": fmt.Errorf("failed to scan order: %w", err).Error(),
				})
				return
			}
			orders = append(orders, order)
		}
		if err := rows.Err(); err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to iterate orders: %w", err).Error(),
			})
			return
		}
		c.JSON(200, orders)
	}
}

// health is a handler that returns a 200 OK response to indicate the server is healthy.
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server is healthy",
		})
	}
}

func CreateOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order

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

		if !validation.ValidStatus(order.Status) {
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
				"error": fmt.Errorf("failed to create order: %w", err).Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "order has been created",
		})
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var order models.Order

		err := db.QueryRow(
			"SELECT id, user_id, status, created_at FROM orders WHERE id = $1",
			id,
		).Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.CreatedAt,
		)

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

		if err != nil || limitInt <= 0 || limitInt > 100 {
			c.JSON(400, gin.H{
				"error": "invalid limit",
			})
			return
		}

		offset := (pageInt - 1) * limitInt
		rows, err := db.Query(
			"SELECT id, user_id, status, created_at FROM orders ORDER BY id LIMIT $1 OFFSET $2",
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
		var orders []models.Order
		for rows.Next() {
			var order models.Order

			if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt); err != nil {
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

		result, err := db.Exec(
			"DELETE FROM orders WHERE id = $1",
			id,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to delete order: %w", err),
			})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to get rows affected: %w", err),
			})
			return
		}

		if rowsAffected == 0 {
			c.JSON(404, gin.H{
				"error": "order not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "order deleted successfully",
		})
	}
}

func UpdateOrderStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req models.UpdateOrderStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("failed to bind JSON: %w", err),
			})
			return
		}

		if !validation.ValidStatus(req.Status) {
			c.JSON(400, gin.H{
				"error": "invalid status",
			})
			return
		}

		result, err := db.Exec(
			"UPDATE orders SET status = $1 WHERE id = $2",
			req.Status, id,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to update order status %w", err),
			})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Errorf("failed to get rows affected %w", err),
			})
			return
		}

		if rowsAffected == 0 {
			c.JSON(404, gin.H{
				"error": "order not found",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "order status updated successfully",
			})
		}
	}
}
