package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"

	"orderApi/internal/config"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg config.Config) *Redis {
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: password,
	})

	return &Redis{
		client: client,
	}
}

func (r *Redis) Ping() error {
	_, err := r.client.Ping(context.Background()).Result()
	return err
}
