package config

import (
	"encoding/json"
	"os"
	"strings"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)
	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	if val := os.Getenv("SERVER_PORT"); val != "" {
		config.ServerPort = val
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		config.DBHost = val
	}
	if val := os.Getenv("DEBUG_MODE"); val != "" {
		config.DebugMode = strings.ToLower(val) == "true"
	}
}package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort   int
	DatabaseURL  string
	LogLevel     string
	CacheEnabled bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		CacheEnabled: getEnvAsBool("CACHE_ENABLED", true),
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
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	strValue := getEnv(key, "")
	if strings.ToLower(strValue) == "true" {
		return true
	}
	if strings.ToLower(strValue) == "false" {
		return false
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
	}
	if cfg.DatabaseURL == "" {
		return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
	}
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return &ConfigError{Field: "LogLevel", Message: "invalid log level"}
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