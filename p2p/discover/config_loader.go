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
    AllowedIPs []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DBHost = getEnv("DB_HOST", "localhost")
    
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, err
    }
    cfg.DBPort = dbPort
    
    debug, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
    if err != nil {
        return nil, err
    }
    cfg.DebugMode = debug
    
    ips := getEnv("ALLOWED_IPS", "127.0.0.1,::1")
    cfg.AllowedIPs = strings.Split(ips, ",")
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
package config

import (
	"errors"
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
	Database string `yaml:"database" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	overrideString := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	overrideInt := func(field *int, envVar string) error {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err != nil {
				return fmt.Errorf("invalid integer value for %s: %s", envVar, val)
			}
			*field = intVal
		}
		return nil
	}

	overrideBool := func(field *bool, envVar string) error {
		if val := os.Getenv(envVar); val != "" {
			*field = val == "true" || val == "1" || val == "yes"
		}
		return nil
	}

	overrideString(&config.Database.Host, "DB_HOST")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")
	overrideString(&config.Database.SSLMode, "DB_SSL_MODE")
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")

	if err := overrideInt(&config.Database.Port, "DB_PORT"); err != nil {
		return err
	}
	if err := overrideInt(&config.Server.Port, "SERVER_PORT"); err != nil {
		return err
	}
	if err := overrideInt(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT"); err != nil {
		return err
	}
	if err := overrideInt(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT"); err != nil {
		return err
	}

	if err := overrideBool(&config.Server.DebugMode, "SERVER_DEBUG"); err != nil {
		return err
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Database.Database == "" {
		return errors.New("database name is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("server read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("server write timeout cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}

func (c *AppConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}