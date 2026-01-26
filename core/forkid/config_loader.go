package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	LogLevel string         `json:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	overrideFromEnv(&config)

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *Config) {
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Database.Port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.Username = user
	}
	if pass := os.Getenv("DB_PASS"); pass != "" {
		config.Database.Password = pass
	}
	if db := os.Getenv("DB_NAME"); db != "" {
		config.Database.Database = db
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.LogLevel = strings.ToUpper(level)
	}
}

func validateConfig(config *Config) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if config.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	validLogLevels := map[string]bool{
		"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true,
	}
	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}
	return nil
}