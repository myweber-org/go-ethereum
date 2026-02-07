package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DebugMode  bool
    DatabaseURL string
    AllowedHosts []string
}

func Load() (*Config, error) {
    cfg := &Config{
        ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
        DebugMode:   getEnvAsBool("DEBUG_MODE", false),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost"}, ","),
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
    if val, err := strconv.ParseBool(valueStr); err == nil {
        return val
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
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Version  string         `json:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath == "" {
		configPath = "config.json"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return loadFromEnvironment()
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	overrideFromEnvironment(&config)

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func loadFromEnvironment() (*AppConfig, error) {
	config := &AppConfig{
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASS", ""),
			Database: getEnvOrDefault("DB_NAME", "appdb"),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30),
			DebugMode:    getEnvAsBool("SERVER_DEBUG", false),
			LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
		},
		Version: getEnvOrDefault("APP_VERSION", "1.0.0"),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func overrideFromEnvironment(config *AppConfig) {
	overrideStruct(&config.Database)
	overrideStruct(&config.Server)
	config.Version = getEnvOrDefault("APP_VERSION", config.Version)
}

func overrideStruct(target interface{}) {
	// Implementation would use reflection to check struct tags
	// and override values from environment variables
	// Simplified for this example
}

func validateConfig(config *AppConfig) error {
	var errors []string

	if config.Database.Host == "" {
		errors = append(errors, "database host cannot be empty")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		errors = append(errors, "database port must be between 1 and 65535")
	}
	if config.Database.Username == "" {
		errors = append(errors, "database username cannot be empty")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		errors = append(errors, "server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout <= 0 {
		errors = append(errors, "server read timeout must be positive")
	}
	if config.Server.WriteTimeout <= 0 {
		errors = append(errors, "server write timeout must be positive")
	}
	if !isValidLogLevel(config.Server.LogLevel) {
		errors = append(errors, "invalid log level")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func isValidLogLevel(level string) bool {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	return validLevels[strings.ToLower(level)]
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}