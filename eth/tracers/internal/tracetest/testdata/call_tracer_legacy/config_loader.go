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
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    debugStr := getEnv("DEBUG_MODE", "false")
    cfg.DebugMode = strings.ToLower(debugStr) == "true"
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    hostsStr := getEnv("ALLOWED_HOSTS", "localhost,127.0.0.1")
    cfg.AllowedHosts = strings.Split(hostsStr, ",")
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}package config

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
    DebugMode    bool
    ReadTimeout  int
    WriteTimeout int
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    dbHost := getEnvWithDefault("DB_HOST", "localhost")
    dbPort := getEnvAsInt("DB_PORT", 5432)
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

    serverPort := getEnvAsInt("SERVER_PORT", 8080)
    debugMode := getEnvAsBool("DEBUG_MODE", false)
    readTimeout := getEnvAsInt("READ_TIMEOUT", 30)
    writeTimeout := getEnvAsInt("WRITE_TIMEOUT", 30)

    cfg.Server = ServerConfig{
        Port:         serverPort,
        DebugMode:    debugMode,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
    }

    cfg.LogLevel = getEnvWithDefault("LOG_LEVEL", "info")

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
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

    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }

    valueStr = strings.ToLower(valueStr)
    return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
    }

    return nil
}