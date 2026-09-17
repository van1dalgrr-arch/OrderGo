package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Database() *sql.DB {
	_ = godotenv.Load()
	config := LoadConfig()
	port := config.DB.Port
	if envPort := os.Getenv("DB_PORT"); envPort != "" {
		if parsedPort, err := strconv.Atoi(envPort); err == nil {
			port = parsedPort
		}
	}
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.User,
		password,
		config.DB.Host,
		port,
		config.DB.Name,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}
	return db
}
