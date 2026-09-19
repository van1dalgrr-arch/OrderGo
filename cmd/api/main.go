package main

import (
	"fmt"
	"log/slog"

	"orderApi/internal/config"
)

func main() {
	cfg, err := config.Load("internal/config/config.yml")

	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	slog.Info("configuration has been loaded", "port", cfg.Server.Port)
}
