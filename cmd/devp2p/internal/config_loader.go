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
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "")
    if cfg.DatabaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }
    
    cacheStr := getEnv("CACHE_ENABLED", "true")
    cacheEnabled, err := strconv.ParseBool(cacheStr)
    if err != nil {
        return nil, fmt.Errorf("invalid CACHE_ENABLED: %v", err)
    }
    cfg.CacheEnabled = cacheEnabled
    
    cfg.LogLevel = strings.ToUpper(getEnv("LOG_LEVEL", "INFO"))
    validLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
    if !validLevels[cfg.LogLevel] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}