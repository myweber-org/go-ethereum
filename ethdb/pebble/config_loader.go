package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    LogLevel string
    CacheEnabled bool
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("DATABASE_URL must start with postgres://")
    }
    cfg.DatabaseURL = dbURL
    
    cfg.LogLevel = getEnvWithDefault("LOG_LEVEL", "info")
    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", cfg.LogLevel)
    }
    
    cacheEnabledStr := getEnvWithDefault("CACHE_ENABLED", "true")
    cacheEnabled, err := strconv.ParseBool(cacheEnabledStr)
    if err != nil {
        return nil, fmt.Errorf("invalid CACHE_ENABLED value: %v", err)
    }
    cfg.CacheEnabled = cacheEnabled
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
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

type AppConfig struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    LogLevel string
    MaxConnections int
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        CacheEnabled:   true,
        LogLevel:       "info",
        MaxConnections: 100,
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid APP_PORT: %v", err)
        }
        cfg.ServerPort = port
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if cacheStr := os.Getenv("CACHE_ENABLED"); cacheStr != "" {
        cacheEnabled, err := strconv.ParseBool(cacheStr)
        if err != nil {
            return nil, fmt.Errorf("invalid CACHE_ENABLED: %v", err)
        }
        cfg.CacheEnabled = cacheEnabled
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        validLevels := map[string]bool{
            "debug": true,
            "info":  true,
            "warn":  true,
            "error": true,
        }
        if !validLevels[strings.ToLower(logLevel)] {
            return nil, fmt.Errorf("invalid LOG_LEVEL: %s", logLevel)
        }
        cfg.LogLevel = strings.ToLower(logLevel)
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        maxConn, err := strconv.Atoi(maxConnStr)
        if err != nil {
            return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
        }
        if maxConn <= 0 {
            return nil, fmt.Errorf("MAX_CONNECTIONS must be positive")
        }
        cfg.MaxConnections = maxConn
    }

    return cfg, nil
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %v", err)
    }

    return &config, nil
}

func validateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Port)
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    return nil
}