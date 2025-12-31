package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Debug    bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    // Database configuration
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %v", err)
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASS", ""),
        Database: getEnv("DB_NAME", "appdb"),
    }

    // Server configuration
    svrPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }

    readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid READ_TIMEOUT: %v", err)
    }

    writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
    if err != nil {
        return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %v", err)
    }

    cfg.Server = ServerConfig{
        Port:         svrPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    // Debug mode
    debugStr := strings.ToLower(getEnv("DEBUG", "false"))
    cfg.Debug = debugStr == "true" || debugStr == "1"

    // Validate configuration
    if err := validateConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
        return fmt.Errorf("database port %d out of range", cfg.Database.Port)
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("server port %d out of range", cfg.Server.Port)
    }

    if cfg.Server.ReadTimeout < 1 {
        return fmt.Errorf("read timeout must be positive")
    }

    if cfg.Server.WriteTimeout < 1 {
        return fmt.Errorf("write timeout must be positive")
    }

    return nil
}