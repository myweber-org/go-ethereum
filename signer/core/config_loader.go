package config

import (
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    Port        int
    DatabaseURL string
    DebugMode   bool
    AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
    config := &AppConfig{}
    
    portStr := getEnvWithDefault("APP_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    config.Port = port
    
    config.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
    
    debugStr := getEnvWithDefault("DEBUG_MODE", "false")
    config.DebugMode = strings.ToLower(debugStr) == "true"
    
    hostsStr := getEnvWithDefault("ALLOWED_HOSTS", "localhost,127.0.0.1")
    config.AllowedHosts = strings.Split(hostsStr, ",")
    
    return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}