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

func Load() (*Config, error) {
	cfg := &Config{}
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		envKey := structField.Tag.Get("env")
		defaultVal := structField.Tag.Get("default")

		envVal := os.Getenv(envKey)
		if envVal == "" {
			envVal = defaultVal
		}

		if err := setField(field, envVal); err != nil {
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
			return err
		}
		field.SetInt(int64(intVal))
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

	if cfg.CacheTTL < 0 {
		return errors.New("cache TTL cannot be negative")
	}

	return nil
}

func (c *Config) String() string {
	masked := *c
	masked.DBPassword = "***"
	data, _ := json.MarshalIndent(masked, "", "  ")
	return string(data)
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Port)
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    return nil
}