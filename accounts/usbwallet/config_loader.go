package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
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
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    overrideStruct(config.Database)
    overrideStruct(config.Server)
}

func overrideStruct(s interface{}) {
    // Implementation would use reflection to check struct tags
    // and override values from environment variables
    // Simplified for this example
}

func ValidateConfigPath(path string) error {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return err
    }

    info, err := os.Stat(absPath)
    if err != nil {
        return err
    }

    if info.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a file", absPath)
    }

    return nil
}

func DefaultConfig() *AppConfig {
    return &AppConfig{
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Username: "postgres",
            Password: "",
            Name:     "appdb",
        },
        Server: ServerConfig{
            Port:         8080,
            Debug:        false,
            LogLevel:     "info",
            ReadTimeout:  30,
            WriteTimeout: 30,
        },
        Version: "1.0.0",
    }
}