package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    APIKeys    []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    var err error
    cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    
    cfg.DBHost = getStringEnv("DB_HOST", "localhost")
    
    cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }
    
    cfg.DebugMode = getBoolEnv("DEBUG_MODE", false)
    
    apiKeysStr := getStringEnv("API_KEYS", "")
    if apiKeysStr != "" {
        cfg.APIKeys = strings.Split(apiKeysStr, ",")
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

func getBoolEnv(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return defaultValue
        }
        return boolValue
    }
    return defaultValue
}