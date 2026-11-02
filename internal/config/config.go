// Package config загружает конфигурацию приложения из переменных окружения.
package config

import (
	"fmt"
	"os"
	"time"
)


type Config struct {
	HTTPAddr string
	DatabaseURL string
	ShutdownTimeout time.Duration
	TLSCertFile string
	TLSKeyFile string
}

// Load читает конфигурацию из переменных окружения и применяет значения по умолчанию.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ShutdownTimeout: 10 * time.Second,
		TLSCertFile: os.Getenv("TLS_CERT_FILE"), 
		TLSKeyFile: os.Getenv("TLS_KEY_FILE"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или defaultValue, если она пуста.
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
