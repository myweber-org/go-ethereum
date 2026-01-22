package config

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
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        configPath = "config.yaml"
    }

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Database.Host == "" {
        config.Database.Host = "localhost"
    }

    if config.Database.Port == 0 {
        config.Database.Port = 5432
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    if config.Server.Mode == "" {
        config.Server.Mode = "development"
    }

    if config.LogLevel == "" {
        config.LogLevel = "info"
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }

    if config.Database.Username == "" {
        return fmt.Errorf("database username is required")
    }

    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    validModes := map[string]bool{
        "development": true,
        "staging":     true,
        "production":  true,
    }

    if !validModes[config.Server.Mode] {
        return fmt.Errorf("invalid server mode: %s", config.Server.Mode)
    }

    validLogLevels := map[string]bool{
        "debug":   true,
        "info":    true,
        "warning": true,
        "error":   true,
    }

    if !validLogLevels[config.LogLevel] {
        return fmt.Errorf("invalid log level: %s", config.LogLevel)
    }

    return nil
}