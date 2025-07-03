package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPServer   HTTPServer   `envPrefix:"HTTP_SERVER_"`
	DB           DB           `envPrefix:"DB_"`
	TokenService TokenService `envPrefix:"TOKEN_SERVICE_"`
}

type HTTPServer struct {
	ListenAddr string `env:"LISTEN_ADDR"`
}

type DB struct {
	Port     string `env:"PORT"`
	Host     string `env:"HOST"`
	DBName   string `env:"DB_NAME"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	SSLMode  string `env:"SSL_MODE"`
}

type TokenService struct {
	SecretKey string `env:"SECRET_KEY"`
}

func loadConfigFromEnv() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
