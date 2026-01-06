package config

import (
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

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.ServerPort = getEnvAsInt("SERVER_PORT", 8080)
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnvAsInt("DB_PORT", 5432)
	cfg.DebugMode = getEnvAsBool("DEBUG_MODE", false)
	cfg.APIKeys = getEnvAsSlice("API_KEYS", []string{}, ",")

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return NewConfigError("invalid server port")
	}
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return NewConfigError("invalid database port")
	}
	if len(cfg.APIKeys) == 0 {
		return NewConfigError("at least one API key is required")
	}
	return nil
}

type ConfigError struct {
	Message string
}

func NewConfigError(msg string) *ConfigError {
	return &ConfigError{Message: msg}
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}package config

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
	cfg.APIKeys = parseAPIKeys(apiKeysStr)

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	strValue := getStringEnv(key, "")
	if strValue == "" {
		return defaultValue, nil
	}
	
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, errors.New("invalid integer value for " + key)
	}
	return value, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	strValue := getStringEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	
	lowerValue := strings.ToLower(strValue)
	return lowerValue == "true" || lowerValue == "1" || lowerValue == "yes"
}

func parseAPIKeys(keysStr string) []string {
	if keysStr == "" {
		return []string{}
	}
	
	keys := strings.Split(keysStr, ",")
	var cleanedKeys []string
	for _, key := range keys {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey != "" {
			cleanedKeys = append(cleanedKeys, trimmedKey)
		}
	}
	return cleanedKeys
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