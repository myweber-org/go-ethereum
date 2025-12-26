package config

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort int
    DebugMode  bool
    MaxWorkers int
    ApiKey     string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DebugMode:  getEnvAsBool("DEBUG_MODE", false),
        MaxWorkers: getEnvAsInt("MAX_WORKERS", 10),
        ApiKey:     getEnv("API_KEY", ""),
    }
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}