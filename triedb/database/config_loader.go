
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	port, err := getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port

	cfg.DBHost = getStringEnv("DB_HOST", "localhost")

	dbPort, err := getIntEnv("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort

	debug, err := getBoolEnv("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug

	apiKeysStr := getStringEnv("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intVal, nil
	}
	return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
	if value, exists := os.LookupEnv(key); exists {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.New("invalid boolean value for " + key)
		}
		return boolVal, nil
	}
	return defaultValue, nil
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.DBHost == "" {
		return errors.New("database host cannot be empty")
	}

	return nil
}