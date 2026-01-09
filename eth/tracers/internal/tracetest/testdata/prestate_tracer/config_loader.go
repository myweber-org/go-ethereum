package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Features []string       `json:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	}

	if err := loadEnvOverrides(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadEnvOverrides(config *AppConfig) error {
	overrideFromEnv(&config.Database.Host, "DB_HOST")
	overrideFromEnv(&config.Database.Port, "DB_PORT")
	overrideFromEnv(&config.Database.Username, "DB_USER")
	overrideFromEnv(&config.Database.Password, "DB_PASS")
	overrideFromEnv(&config.Database.Database, "DB_NAME")

	overrideFromEnv(&config.Server.Port, "SERVER_PORT")
	overrideFromEnv(&config.Server.ReadTimeout, "READ_TIMEOUT")
	overrideFromEnv(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
	overrideFromEnv(&config.Server.DebugMode, "DEBUG_MODE")
	overrideFromEnv(&config.Server.LogLevel, "LOG_LEVEL")

	return nil
}

func overrideFromEnv(target interface{}, envVar string) {
	val := os.Getenv(envVar)
	if val == "" {
		return
	}

	switch v := target.(type) {
	case *string:
		*v = val
	case *int:
		fmt.Sscanf(val, "%d", v)
	case *bool:
		*v = strings.ToLower(val) == "true" || val == "1"
	}
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if config.Server.LogLevel != "" {
		validLevels := map[string]bool{
			"debug": true, "info": true, "warn": true, "error": true,
		}
		if !validLevels[strings.ToLower(config.Server.LogLevel)] {
			return fmt.Errorf("invalid log level: %s", config.Server.LogLevel)
		}
	}
	return nil
}

func SaveConfig(config *AppConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}