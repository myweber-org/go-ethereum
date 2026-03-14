
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
	config := &AppConfig{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
		CacheTTL:   getEnvAsInt("CACHE_TTL", 300),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
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

func validateConfig(config *AppConfig) error {
	if config.ServerPort < 1 || config.ServerPort > 65535 {
		return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
	}
	if config.DatabaseURL == "" {
		return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
	}
	if config.CacheTTL < 0 {
		return &ConfigError{Field: "CacheTTL", Message: "cache TTL cannot be negative"}
	}
	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Field + " - " + e.Message
}