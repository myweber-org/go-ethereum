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
    EnableDebug bool
    MaxWorkers int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("DATABASE_URL must start with postgres://")
    }
    cfg.DatabaseURL = dbURL
    
    debugStr := getEnvWithDefault("ENABLE_DEBUG", "false")
    debug, err := strconv.ParseBool(debugStr)
    if err != nil {
        return nil, fmt.Errorf("invalid ENABLE_DEBUG: %v", err)
    }
    cfg.EnableDebug = debug
    
    workersStr := getEnvWithDefault("MAX_WORKERS", "10")
    workers, err := strconv.Atoi(workersStr)
    if err != nil || workers <= 0 {
        return nil, fmt.Errorf("invalid MAX_WORKERS: must be positive integer")
    }
    cfg.MaxWorkers = workers
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}