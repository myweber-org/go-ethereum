package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port        int
    DatabaseURL string
    Debug       bool
}

func Load() (*Config, error) {
    cfg := &Config{
        Port:        8080,
        DatabaseURL: "postgres://localhost:5432/app",
        Debug:       false,
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil {
            cfg.Port = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if debugStr := os.Getenv("APP_DEBUG"); debugStr != "" {
        if debug, err := strconv.ParseBool(debugStr); err == nil {
            cfg.Debug = debug
        }
    }

    return cfg, nil
}