package main

import (
	"database/sql"
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
	config := LoadConfig()

	db := Database()
	defer db.Close()

	r := gin.Default()

	r.GET("/health", health)
	r.POST("/orders", createOrder)
	r.GET("/orders/:id", getOrder)
	r.DELETE("/orders/:id", deleteOrder)

	if err := r.Run(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
