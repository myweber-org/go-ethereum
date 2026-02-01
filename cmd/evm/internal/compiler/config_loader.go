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
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DBName     string `env:"DB_NAME" default:"appdb"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		envKey := structField.Tag.Get("env")
		defaultValue := structField.Tag.Get("default")

		envValue := os.Getenv(envKey)
		if envValue == "" {
			envValue = defaultValue
		}

		if envValue == "" {
			return nil, fmt.Errorf("missing value for %s", envKey)
		}

		if err := setFieldValue(field, envValue); err != nil {
			return nil, fmt.Errorf("invalid value for %s: %w", envKey, err)
		}
	}

	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
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
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func (c *Config) Validate() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.DBPort < 1 || c.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if strings.TrimSpace(c.DBHost) == "" {
		return errors.New("database host cannot be empty")
	}
	if strings.TrimSpace(c.DBName) == "" {
		return errors.New("database name cannot be empty")
	}
	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}