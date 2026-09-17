package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Name string `yaml:"name"`
}

func LoadConfig() Config {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal("failed to read config.yml:", err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("failed to unmarshal config.yml:", err)
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		config.DB.Host = host
	}

	if user := os.Getenv("DB_USER"); user != "" {
		config.DB.User = user
	}

	if name := os.Getenv("DB_NAME"); name != "" {
		config.DB.Name = name
	}
	return config
}
