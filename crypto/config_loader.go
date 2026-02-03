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

func LoadConfig() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost", "127.0.0.1"}),
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

func getEnvAsSlice(key string, defaultValue []string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, ",")
}package config

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
    LogLevel string
    MaxConnections int
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
    
    cacheEnabledStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabledStr) == "true"
    
    cfg.LogLevel = strings.ToUpper(getEnvWithDefault("LOG_LEVEL", "INFO"))
    validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
    if !validLogLevels[cfg.LogLevel] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
    }
    
    maxConnStr := getEnvWithDefault("MAX_CONNECTIONS", "100")
    maxConn, err := strconv.Atoi(maxConnStr)
    if err != nil || maxConn <= 0 {
        return nil, fmt.Errorf("invalid MAX_CONNECTIONS value: %s", maxConnStr)
    }
    cfg.MaxConnections = maxConn
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func (c *AppConfig) Validate() error {
    if c.ServerPort < 1 || c.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if c.DatabaseURL == "" {
        return fmt.Errorf("database URL cannot be empty")
    }
    if c.MaxConnections < 1 {
        return fmt.Errorf("max connections must be positive")
    }
    return nil
}