
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	
	portStr := getEnvWithDefault("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port
	
	debugStr := getEnvWithDefault("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
	
	ttlStr := getEnvWithDefault("CACHE_TTL", "300")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return nil, errors.New("invalid CACHE_TTL value")
	}
	cfg.CacheTTL = ttl
	
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if cfg.DatabaseURL == "" {
		return errors.New("database URL cannot be empty")
	}
	
	if cfg.CacheTTL < 0 {
		return errors.New("cache TTL cannot be negative")
	}
	
	return nil
}package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func Load() (*Config, error) {
	cfg := &Config{}
	var err error

	cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getStringEnv("DB_HOST", "localhost")
	
	cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode = getBoolEnv("DEBUG_MODE", false)
	
	apiKeysStr := getStringEnv("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return nil, errors.New("invalid server port range")
	}

	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return nil, errors.New("invalid database port range")
	}

	return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}