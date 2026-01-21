package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DBName     string `env:"DB_NAME" default:"appdb"`
	DBUser     string `env:"DB_USER" default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" default:""`
	CacheTTL   int    `env:"CACHE_TTL" default:"300"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		envKey := structField.Tag.Get("env")
		defaultVal := structField.Tag.Get("default")

		envValue := os.Getenv(envKey)
		if envValue == "" {
			envValue = defaultVal
		}

		if err := setFieldValue(field, envValue); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		field.SetInt(int64(intVal))
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return errors.New("invalid log level")
	}

	if cfg.CacheTTL < 0 {
		return errors.New("cache TTL cannot be negative")
	}

	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort    int
    DatabaseURL   string
    LogLevel      string
    CacheEnabled  bool
    MaxConnections int
}

func Load() (*Config, error) {
    cfg := &Config{
        ServerPort:    8080,
        DatabaseURL:   "localhost:5432",
        LogLevel:      "info",
        CacheEnabled:  true,
        MaxConnections: 100,
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
        }
        cfg.ServerPort = port
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        validLevels := map[string]bool{
            "debug": true,
            "info":  true,
            "warn":  true,
            "error": true,
        }
        if !validLevels[strings.ToLower(logLevel)] {
            return nil, fmt.Errorf("invalid LOG_LEVEL: %s", logLevel)
        }
        cfg.LogLevel = strings.ToLower(logLevel)
    }

    if cacheStr := os.Getenv("CACHE_ENABLED"); cacheStr != "" {
        enabled, err := strconv.ParseBool(cacheStr)
        if err != nil {
            return nil, fmt.Errorf("invalid CACHE_ENABLED: %v", err)
        }
        cfg.CacheEnabled = enabled
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        maxConn, err := strconv.Atoi(maxConnStr)
        if err != nil {
            return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
        }
        if maxConn <= 0 {
            return nil, fmt.Errorf("MAX_CONNECTIONS must be positive")
        }
        cfg.MaxConnections = maxConn
    }

    return cfg, nil
}

func (c *Config) Validate() error {
    if c.ServerPort < 1 || c.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }
    if c.MaxConnections < 1 {
        return fmt.Errorf("max connections must be at least 1")
    }
    return nil
}