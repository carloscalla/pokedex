package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port            string
	AllowedOrigin   string
	DatabaseURL     string
	JWTSecret       string
	AdminToken      string
	PokeAPIBaseURL  string
	SyncConcurrency int
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		AllowedOrigin:   getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AdminToken:      os.Getenv("ADMIN_TOKEN"),
		PokeAPIBaseURL:  getEnv("POKEAPI_BASE_URL", "https://pokeapi.co/api/v2"),
		SyncConcurrency: 15,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
