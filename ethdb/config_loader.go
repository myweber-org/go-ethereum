package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DebugMode  bool
    DatabaseURL string
    AllowedHosts []string
}

func Load() (*Config, error) {
    cfg := &Config{
        ServerPort:  getInt("SERVER_PORT", 8080),
        DebugMode:   getBool("DEBUG_MODE", false),
        DatabaseURL: getString("DATABASE_URL", "postgres://localhost:5432/app"),
        AllowedHosts: getStringSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
    }
    return cfg, nil
}

func getString(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getStringSlice(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}