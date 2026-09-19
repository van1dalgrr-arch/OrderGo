package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"orderApi/internal/config"
	"orderApi/internal/server"
)

func main() {
	cfg, err := config.Load("internal/config/config.yml")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	slog.Info("configuration has been loaded", "port", cfg.Server.Port)

	// ! Graceful shutdown
	srv := server.New(cfg.Server.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go srv.ListenAndServe()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}
