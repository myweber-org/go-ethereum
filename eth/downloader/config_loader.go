
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
}package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Format string `yaml:"format" env:"LOG_FORMAT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	overrideString(&cfg.Server.Host, "SERVER_HOST")
	overrideInt(&cfg.Server.Port, "SERVER_PORT")
	overrideString(&cfg.Database.Host, "DB_HOST")
	overrideInt(&cfg.Database.Port, "DB_PORT")
	overrideString(&cfg.Database.Name, "DB_NAME")
	overrideString(&cfg.Database.User, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASSWORD")
	overrideString(&cfg.Database.SSLMode, "DB_SSL_MODE")
	overrideString(&cfg.Logging.Level, "LOG_LEVEL")
	overrideString(&cfg.Logging.Format, "LOG_FORMAT")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}