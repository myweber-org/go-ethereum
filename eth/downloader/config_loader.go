
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
	CacheTTL   int    `env:"CACHE_TTL" default:"300"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		envKey := structField.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envValue := os.Getenv(envKey)
		defaultValue := structField.Tag.Get("default")

		var valueToSet string
		if envValue != "" {
			valueToSet = envValue
		} else {
			valueToSet = defaultValue
		}

		if err := setField(field, valueToSet); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", structField.Name, err)
		}
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func setField(field reflect.Value, value string) error {
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
}