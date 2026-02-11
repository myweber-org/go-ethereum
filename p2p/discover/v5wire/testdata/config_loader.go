package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			return nil, fmt.Errorf("failed to decode config: %w", err)
		}
	}

	if err := loadEnvVars(cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadEnvVars(cfg *Config) error {
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
	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.ServerPort)
	}

	if cfg.DBHost == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if cfg.DBPort <= 0 || cfg.DBPort > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.DBPort)
	}

	return nil
}package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := make(Config)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format on line %d", lineNumber)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = os.ExpandEnv(value)
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

func (c Config) Get(key string) (string, error) {
	value, exists := c[key]
	if !exists {
		return "", fmt.Errorf("key '%s' not found in configuration", key)
	}
	return value, nil
}