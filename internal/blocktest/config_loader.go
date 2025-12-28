package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
    AllowedOrigins []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    var err error
    cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    
    cfg.DatabaseURL = getStringEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    cfg.CacheEnabled, err = getBoolEnv("CACHE_ENABLED", true)
    if err != nil {
        return nil, fmt.Errorf("invalid CACHE_ENABLED: %w", err)
    }
    
    cfg.MaxConnections, err = getIntEnv("MAX_CONNECTIONS", 100)
    if err != nil {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %w", err)
    }
    
    cfg.AllowedOrigins = getStringSliceEnv("ALLOWED_ORIGINS", []string{"*"})
    
    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, fmt.Errorf("cannot parse %s as integer: %w", key, err)
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return false, fmt.Errorf("cannot parse %s as boolean: %w", key, err)
        }
        return boolValue, nil
    }
    return defaultValue, nil
}

func getStringSliceEnv(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    
    if cfg.MaxConnections < 1 {
        return fmt.Errorf("max connections must be positive")
    }
    
    if cfg.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }
    
    return nil
}