package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    dbHost := getEnvWithDefault("DB_HOST", "localhost")
    dbPort := getEnvIntWithDefault("DB_PORT", 5432)
    dbUser := getEnvWithDefault("DB_USER", "postgres")
    dbPass := getEnvWithDefault("DB_PASS", "")
    dbName := getEnvWithDefault("DB_NAME", "appdb")

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: dbUser,
        Password: dbPass,
        Database: dbName,
    }

    srvPort := getEnvIntWithDefault("SERVER_PORT", 8080)
    readTimeout := getEnvIntWithDefault("READ_TIMEOUT", 30)
    writeTimeout := getEnvIntWithDefault("WRITE_TIMEOUT", 30)
    debugMode := getEnvBoolWithDefault("DEBUG_MODE", false)

    cfg.Server = ServerConfig{
        Port:         srvPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
        DebugMode:    debugMode,
    }

    logLevel := getEnvWithDefault("LOG_LEVEL", "info")
    allowedLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !allowedLevels[strings.ToLower(logLevel)] {
        return nil, fmt.Errorf("invalid log level: %s", logLevel)
    }
    cfg.LogLevel = strings.ToLower(logLevel)

    if cfg.Database.Password == "" {
        return nil, fmt.Errorf("database password must be set")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return nil, fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getEnvBoolWithDefault(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolVal, err := strconv.ParseBool(value); err == nil {
            return boolVal
        }
    }
    return defaultValue
}