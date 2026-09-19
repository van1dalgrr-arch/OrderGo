package main

import (
	"fmt"
	"log"
	"log/slog"

	"orderApi/internal/config"
	"orderApi/internal/server"
)

func main() {
	cfg, err := config.Load("internal/config/config.yml")

	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	slog.Info("configuration has been loaded", "port", cfg.Server.Port)

	r := server.New()

	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
