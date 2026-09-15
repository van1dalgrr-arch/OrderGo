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
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("failed to unmarshal config:", err)
	}

	return config
}
