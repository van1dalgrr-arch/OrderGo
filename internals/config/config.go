package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

func Load(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("failed to read path file", "error", err)
		return Config{}, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error("failed to unmarshal data", "error", err)
		return Config{}, err
	}
	return config, nil
}
