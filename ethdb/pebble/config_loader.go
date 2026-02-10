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
    LogLevel string
    CacheEnabled bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("DATABASE_URL must start with postgres://")
    }
    cfg.DatabaseURL = dbURL
    
    cfg.LogLevel = getEnvWithDefault("LOG_LEVEL", "info")
    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
    }
    
    cacheEnabledStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cacheEnabled, err := strconv.ParseBool(cacheEnabledStr)
    if err != nil {
        return nil, fmt.Errorf("invalid CACHE_ENABLED value: %v", err)
    }
    cfg.CacheEnabled = cacheEnabled
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}