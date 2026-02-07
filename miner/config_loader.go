package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
	DebugMode    bool   `yaml:"debug_mode"`
	LogLevel     string `yaml:"log_level"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features []string       `yaml:"features"`
}

func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.LogLevel == "" {
		config.Server.LogLevel = "info"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}
	return nil
}package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel   string `json:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{
		ServerPort: 8080,
		DBHost:     "localhost",
		DBPort:     5432,
		DebugMode:  false,
		LogLevel:   "info",
	}

	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func loadFromFile(path string, cfg *Config) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return nil
}

func loadFromEnv(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		fieldValue := v.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Int:
			var intVal int
			if _, err := fmt.Sscanf(envValue, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid integer value for %s: %s", envTag, envValue)
			}
			fieldValue.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
			fieldValue.SetBool(boolVal)
		}
	}

	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return fmt.Errorf("server_port must be between 1 and 65535")
	}
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return fmt.Errorf("db_port must be between 1 and 65535")
	}
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return fmt.Errorf("log_level must be one of: debug, info, warn, error, fatal")
	}
	return nil
}