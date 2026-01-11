package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
		APIKeys:    getEnvAsSlice("API_KEYS", []string{}, ","),
	}

	if config.ServerPort < 1 || config.ServerPort > 65535 {
		return nil, &ConfigError{Field: "SERVER_PORT", Value: strconv.Itoa(config.ServerPort)}
	}

	if config.DBPort < 1 || config.DBPort > 65535 {
		return nil, &ConfigError{Field: "DB_PORT", Value: strconv.Itoa(config.DBPort)}
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
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}

type ConfigError struct {
	Field string
	Value string
}

func (e *ConfigError) Error() string {
	return "invalid configuration value for " + e.Field + ": " + e.Value
}