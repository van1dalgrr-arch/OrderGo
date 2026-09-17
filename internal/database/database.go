package database

import (
	"database/sql"

	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"orderApi/internal/config"
	"os"
	"strconv"
)

func Database(cfg config.Config) *sql.DB {
	_ = godotenv.Load()
	port := cfg.DB.Port
	if envPort := os.Getenv("DB_PORT"); envPort != "" {
		if parsedPort, err := strconv.Atoi(envPort); err == nil {
			port = parsedPort
		}
	}
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User,
		password,
		cfg.DB.Host,
		port,
		cfg.DB.Name,
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
