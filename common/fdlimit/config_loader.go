
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	
	portStr := getEnvWithDefault("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port
	
	debugStr := getEnvWithDefault("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
	
	ttlStr := getEnvWithDefault("CACHE_TTL", "300")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return nil, errors.New("invalid CACHE_TTL value")
	}
	cfg.CacheTTL = ttl
	
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if cfg.DatabaseURL == "" {
		return errors.New("database URL cannot be empty")
	}
	
	if cfg.CacheTTL < 0 {
		return errors.New("cache TTL cannot be negative")
	}
	
	return nil
}