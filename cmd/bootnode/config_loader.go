package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    MaxWorkers int
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DebugMode:  getEnvAsBool("DEBUG_MODE", false),
        MaxWorkers: getEnvAsInt("MAX_WORKERS", 10),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
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
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    valueStr = strings.ToLower(valueStr)
    return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}

func validateConfig(cfg *AppConfig) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.ServerPort)
    }
    if cfg.DBPort < 1 || cfg.DBPort > 65535 {
        return fmt.Errorf("invalid database port: %d", cfg.DBPort)
    }
    if cfg.MaxWorkers < 1 {
        return fmt.Errorf("max workers must be positive: %d", cfg.MaxWorkers)
    }
    return nil
}package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.overrideFromEnv()

	return &config, nil
}

func (c *Config) overrideFromEnv() {
	c.Server.Host = getEnvOrDefault("SERVER_HOST", c.Server.Host)
	c.Server.Port = getEnvIntOrDefault("SERVER_PORT", c.Server.Port)
	c.Database.Host = getEnvOrDefault("DB_HOST", c.Database.Host)
	c.Database.Port = getEnvIntOrDefault("DB_PORT", c.Database.Port)
	c.Database.Name = getEnvOrDefault("DB_NAME", c.Database.Name)
	c.Database.User = getEnvOrDefault("DB_USER", c.Database.User)
	c.Database.Password = getEnvOrDefault("DB_PASSWORD", c.Database.Password)
	c.Logging.Level = getEnvOrDefault("LOG_LEVEL", c.Logging.Level)
	c.Logging.Output = getEnvOrDefault("LOG_OUTPUT", c.Logging.Output)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}