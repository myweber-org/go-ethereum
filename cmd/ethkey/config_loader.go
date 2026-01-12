
package config

import (
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
	cfg := &AppConfig{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
		CacheTTL:   getEnvAsInt("CACHE_TTL", 300),
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.ToLower(valueStr) == "true" || valueStr == "1"
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return ErrInvalidPort
	}
	if cfg.DatabaseURL == "" {
		return ErrMissingDatabaseURL
	}
	if cfg.CacheTTL < 0 {
		return ErrInvalidCacheTTL
	}
	return nil
}

var (
	ErrInvalidPort        = ConfigError{Code: "INVALID_PORT", Message: "Port must be between 1 and 65535"}
	ErrMissingDatabaseURL = ConfigError{Code: "MISSING_DB_URL", Message: "Database URL is required"}
	ErrInvalidCacheTTL    = ConfigError{Code: "INVALID_CACHE_TTL", Message: "Cache TTL cannot be negative"}
)

type ConfigError struct {
	Code    string
	Message string
}

func (e ConfigError) Error() string {
	return e.Code + ": " + e.Message
}