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
    Name     string `json:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Address string `json:"address" env:"SERVER_ADDR"`
    Port    int    `json:"port" env:"SERVER_PORT"`
    Debug   bool   `json:"debug" env:"SERVER_DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig
    
    if configPath == "" {
        configPath = getDefaultConfigPath()
    }

    fileData, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := json.Unmarshal(fileData, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config JSON: %w", err)
    }

    overrideFromEnv(&config)

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func getDefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.json"
    }
    return filepath.Join(homeDir, ".app", "config.json")
}

func overrideFromEnv(config *AppConfig) {
    overrideStruct(config)
}

func overrideStruct(v interface{}) {
    // Implementation would use reflection to check struct tags
    // and override values from environment variables
    // Simplified for brevity
}

func validateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    if config.LogLevel == "" {
        config.LogLevel = "info"
    }
    
    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    if !validLogLevels[strings.ToLower(config.LogLevel)] {
        return fmt.Errorf("invalid log level: %s", config.LogLevel)
    }
    
    return nil
}

func (c *AppConfig) GetDSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
        c.Database.Username,
        c.Database.Password,
        c.Database.Host,
        c.Database.Port,
        c.Database.Name)
}

func (c *AppConfig) GetServerAddr() string {
    return fmt.Sprintf("%s:%d", c.Server.Address, c.Server.Port)
}