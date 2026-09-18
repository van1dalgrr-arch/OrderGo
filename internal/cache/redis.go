package cache

import (
	"context"

	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

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

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	err := r.client.Set(context.Background(), key, value, expiration).Err()

	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string) (string, error) {
	value, err := r.client.Get(context.Background(), key).Result()
	return value, err

}

func (r *Redis) Delete(key string) error {
	err := r.client.Del(context.Background(), key).Err()

	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) SetJSON(key string, value any, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(context.Background(), key, jsonValue, expiration).Err()
	if err != nil {
		slog.Error("failed to set json", "key", key, "err", err)
		return err
	}
	return nil
}

func (r *Redis) GetJSON(key string, dest any) error {
	jsonValue, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		slog.Error("failed to get json", "key", key, "err", err)
		return err
	}
	err = json.Unmarshal([]byte(jsonValue), dest)
	if err != nil {
		slog.Error("failed to unmarshal json", "key", key, "err", err)
		return err
	}
	return nil
}
