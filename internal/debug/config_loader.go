package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    MaxConnections int
    FeatureFlags map[string]bool
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        FeatureFlags: make(map[string]bool),
    }

    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port

    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("DATABASE_URL must start with postgres://")
    }
    cfg.DatabaseURL = dbURL

    cacheEnabled := getEnvWithDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabled) == "true"

    maxConnStr := getEnvWithDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil || maxConn <= 0 {
        return nil, fmt.Errorf("MAX_CONNECTIONS must be positive integer")
    }
    cfg.MaxConnections = maxConn

    featureFlags := getEnvWithDefault("FEATURE_FLAGS", "")
    if featureFlags != "" {
        flags := strings.Split(featureFlags, ",")
        for _, flag := range flags {
            parts := strings.Split(strings.TrimSpace(flag), "=")
            if len(parts) == 2 {
                cfg.FeatureFlags[parts[0]] = strings.ToLower(parts[1]) == "true"
            }
        }
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}