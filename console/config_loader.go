package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Version  string         `json:"version"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		Version: "1.0.0",
	}

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	overrideWithEnvVars(config)

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func overrideWithEnvVars(config *Config) {
	overrideStruct(&config.Database)
	overrideStruct(&config.Server)
}

func overrideStruct(s interface{}) {
	// This would be implemented using reflection to read struct tags
	// and override values from environment variables
	// Simplified for this example
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if !isValidLogLevel(config.Server.LogLevel) {
		return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
	}
	return nil
}

func isValidLogLevel(level string) bool {
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	for _, valid := range validLevels {
		if strings.ToLower(level) == valid {
			return true
		}
	}
	return false
}