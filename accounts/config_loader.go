
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type ServerConfig struct {
    Port         int    `json:"port"`
    ReadTimeout  int    `json:"read_timeout"`
    WriteTimeout int    `json:"write_timeout"`
}

type AppConfig struct {
    Environment string         `json:"environment"`
    Debug       bool           `json:"debug"`
    Database    DatabaseConfig `json:"database"`
    Server      ServerConfig   `json:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    var config AppConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("failed to decode config: %w", err)
    }

    overrideFromEnv(&config)

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    if env := os.Getenv("APP_ENV"); env != "" {
        config.Environment = env
    }

    if debug := os.Getenv("APP_DEBUG"); debug != "" {
        config.Debug = strings.ToLower(debug) == "true"
    }

    if host := os.Getenv("DB_HOST"); host != "" {
        config.Database.Host = host
    }

    if port := os.Getenv("DB_PORT"); port != "" {
        fmt.Sscanf(port, "%d", &config.Database.Port)
    }

    if port := os.Getenv("SERVER_PORT"); port != "" {
        fmt.Sscanf(port, "%d", &config.Server.Port)
    }
}

func validateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    if config.Server.Port <= 0 || config.Server.Port > 65535 {
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

func (c *AppConfig) String() string {
    maskedConfig := *c
    maskedConfig.Database.Password = "***"
    data, _ := json.MarshalIndent(maskedConfig, "", "  ")
    return string(data)
}