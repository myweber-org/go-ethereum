
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
	var err error

	cfg.ServerPort, err = getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode, err = getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}

	cfg.DatabaseURL = getEnvString("DATABASE_URL", "postgres://localhost:5432/appdb")
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL cannot be empty")
	}

	cfg.CacheTTL, err = getEnvInt("CACHE_TTL", 300)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value := os.Getenv(key); value != "" {
		lowerValue := strings.ToLower(value)
		if lowerValue == "true" || lowerValue == "1" {
			return true, nil
		}
		if lowerValue == "false" || lowerValue == "0" {
			return false, nil
		}
		return false, errors.New("invalid boolean value for " + key)
	}
	return defaultValue, nil
}