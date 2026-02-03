package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    if val := os.Getenv("DB_HOST"); val != "" {
        config.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Database.Port)
    }
    if val := os.Getenv("DB_USER"); val != "" {
        config.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        config.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        config.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.Port)
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.ReadTimeout)
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &config.Server.WriteTimeout)
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        config.Server.DebugMode = val == "true" || val == "1"
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        config.LogLevel = val
    }
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    return nil
}
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
    dbHost := getEnvOrDefault("DB_HOST", "localhost")
    dbPort := getEnvAsInt("DB_PORT", 5432)
    dbUser := getEnvOrDefault("DB_USER", "postgres")
    dbPass := getEnvOrDefault("DB_PASS", "")
    dbName := getEnvOrDefault("DB_NAME", "appdb")
    dbSSL := getEnvOrDefault("DB_SSL_MODE", "disable")

    serverPort := getEnvAsInt("SERVER_PORT", 8080)
    readTimeout := getEnvAsInt("READ_TIMEOUT", 30)
    writeTimeout := getEnvAsInt("WRITE_TIMEOUT", 30)
    debugMode := getEnvAsBool("DEBUG_MODE", false)
    logLevel := getEnvOrDefault("LOG_LEVEL", "info")

    if dbPort <= 0 || dbPort > 65535 {
        return nil, fmt.Errorf("invalid database port: %d", dbPort)
    }

    if serverPort <= 0 || serverPort > 65535 {
        return nil, fmt.Errorf("invalid server port: %d", serverPort)
    }

    if readTimeout <= 0 || writeTimeout <= 0 {
        return nil, fmt.Errorf("timeout values must be positive")
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLogLevels[strings.ToLower(logLevel)] {
        return nil, fmt.Errorf("invalid log level: %s", logLevel)
    }

    config := &Config{
        Database: DatabaseConfig{
            Host:     dbHost,
            Port:     dbPort,
            Username: dbUser,
            Password: dbPass,
            Database: dbName,
            SSLMode:  dbSSL,
        },
        Server: ServerConfig{
            Port:         serverPort,
            ReadTimeout:  readTimeout,
            WriteTimeout: writeTimeout,
            DebugMode:    debugMode,
        },
        LogLevel: strings.ToLower(logLevel),
    }

    return config, nil
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

    value, err := strconv.ParseBool(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}