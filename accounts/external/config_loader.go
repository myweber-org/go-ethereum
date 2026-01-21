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
    DebugMode   bool
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    if port < 1 || port > 65535 {
        return nil, fmt.Errorf("SERVER_PORT out of range: %d", port)
    }
    cfg.ServerPort = port
    
    dbURL := getEnv("DATABASE_URL", "")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }
    cfg.DatabaseURL = dbURL
    
    debugStr := getEnv("DEBUG_MODE", "false")
    debug, err := strconv.ParseBool(debugStr)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %v", err)
    }
    cfg.DebugMode = debug
    
    originsStr := getEnv("ALLOWED_ORIGINS", "*")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}