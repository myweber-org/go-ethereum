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
	DBPassword string `env:"DB_PASSWORD"`
	CacheTTL   int    `env:"CACHE_TTL" default:"300"`
	EnableSSL  bool   `env:"ENABLE_SSL" default:"false"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		envTag := structField.Tag.Get("env")
		defaultTag := structField.Tag.Get("default")

		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			envValue = defaultTag
		}

		if envValue == "" && field.Kind() != reflect.String {
			return nil, fmt.Errorf("environment variable %s is required", envTag)
		}

		if err := setFieldValue(field, envValue); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	if value == "" && field.Kind() == reflect.String {
		field.SetString("")
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
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

	if cfg.DBHost == "" {
		return errors.New("database host cannot be empty")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.CacheTTL < 0 {
		return errors.New("cache TTL cannot be negative")
	}

	return nil
}

func (c *Config) String() string {
	safeConfig := *c
	safeConfig.DBPassword = "***REDACTED***"
	data, _ := json.MarshalIndent(safeConfig, "", "  ")
	return string(data)
}