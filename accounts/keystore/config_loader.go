package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL  string
    MaxConnections int
    DebugMode    bool
    AllowedHosts []string
}

func LoadConfig(filename string) (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnvOrDefault("DB_URL", "postgres://localhost:5432/mydb"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
    }

    if filename != "" {
        err := loadFromFile(filename, cfg)
        if err != nil {
            return nil, err
        }
    }

    return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }
    return strings.ToLower(valueStr) == "true"
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, ",")
}

func loadFromFile(filename string, cfg *Config) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        switch key {
        case "DATABASE_URL":
            cfg.DatabaseURL = value
        case "MAX_CONNECTIONS":
            if v, err := strconv.Atoi(value); err == nil {
                cfg.MaxConnections = v
            }
        case "DEBUG_MODE":
            cfg.DebugMode = strings.ToLower(value) == "true"
        case "ALLOWED_HOSTS":
            cfg.AllowedHosts = strings.Split(value, ",")
        }
    }

    return nil
}