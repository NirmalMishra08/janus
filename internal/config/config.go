package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Route struct {
	Path    string `yaml:"path"`
	Service string `yaml:"service"`
}

type ServiceConfig struct {
	Instances []string `yaml:"instances"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`

	// Routes & Services from config.yaml
	Routes   []Route                  `yaml:"routes"`
	Services map[string]ServiceConfig `yaml:"services"`


	// Database & other sensitive configs from .env
	PostgresConn string
	RedisURL     string
	JWTSECRET string
}

func LoadConfig() (*Config, error) {
	// Load .env file (ignores error if .env doesn't exist)
	_ = godotenv.Load()

	data, err := os.ReadFile("../../configs/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	cfg.PostgresConn = os.Getenv("POSTGRES_CONN")
	cfg.RedisURL = os.Getenv("REDIS_URL")
	cfg.JWTSECRET = os.Getenv("JWT_SECRET")

	// Basic validation
	if cfg.Server.Port == "" {
		cfg.Server.Port = os.Getenv("PORT")
		if cfg.Server.Port == "" {
			cfg.Server.Port = "8080"
		}
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = os.Getenv("PORT")
		if cfg.Server.Port == "" {
			cfg.Server.Port = "8080"
		}
	}

	return &cfg, err
}