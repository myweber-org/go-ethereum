package config

import (
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
    SSLMode  string
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
    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            Username: getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            Database: getEnv("DB_NAME", "appdb"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
        Server: ServerConfig{
            Port:         getEnvAsInt("SERVER_PORT", 8080),
            ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
            WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
            DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        },
        LogLevel: getEnv("LOG_LEVEL", "info"),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, err
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
    if value, err := strconv.Atoi(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    strValue = strings.ToLower(strValue)
    return strValue == "true" || strValue == "1" || strValue == "yes"
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
    }
    if cfg.Database.Password == "" {
        return &ConfigError{Field: "DB_PASSWORD", Message: "database password is required"}
    }
    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}