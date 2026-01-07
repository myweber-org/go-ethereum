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
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        EnableDebug:    false,
        MaxConnections: 100,
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
        }
        cfg.ServerPort = port
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if debugStr := os.Getenv("ENABLE_DEBUG"); debugStr != "" {
        debug, err := strconv.ParseBool(debugStr)
        if err != nil {
            return nil, fmt.Errorf("invalid ENABLE_DEBUG: %v", err)
        }
        cfg.EnableDebug = debug
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        maxConn, err := strconv.Atoi(maxConnStr)
        if err != nil {
            return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
        }
        if maxConn <= 0 {
            return nil, fmt.Errorf("MAX_CONNECTIONS must be positive")
        }
        cfg.MaxConnections = maxConn
    }

    if !strings.Contains(cfg.DatabaseURL, "://") {
        return nil, fmt.Errorf("invalid database URL format")
    }

    return cfg, nil
}