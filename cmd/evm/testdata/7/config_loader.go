package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
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
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Database.Port)
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.Port)
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.ReadTimeout)
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.WriteTimeout)
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        config.Server.DebugMode = val == "true" || val == "1"
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port < 1 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    if config.Server.ReadTimeout < 0 {
        return fmt.Errorf("read timeout cannot be negative")
    }
    if config.Server.WriteTimeout < 0 {
        return fmt.Errorf("write timeout cannot be negative")
    }
    return nil
}