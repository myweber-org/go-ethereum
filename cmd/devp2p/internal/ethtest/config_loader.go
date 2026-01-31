package config

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort: getEnvAsInt("SERVER_PORT", 8080),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DebugMode:  getEnvAsBool("DEBUG_MODE", false),
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
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    if intValue, err := strconv.Atoi(strValue); err == nil {
        return intValue
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    boolValue, err := strconv.ParseBool(strValue)
    if err != nil {
        return defaultValue
    }
    return boolValue
}