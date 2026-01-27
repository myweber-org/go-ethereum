package config

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    Database struct {
        ConnectionString string `yaml:"connection_string"`
        MaxConnections   int    `yaml:"max_connections"`
    } `yaml:"database"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var config AppConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    log.Printf("Configuration loaded successfully from %s", filePath)
    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    
    if config.Database.MaxConnections < 1 {
        return fmt.Errorf("max connections must be positive: %d", config.Database.MaxConnections)
    }
    
    return nil
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort    int
    DatabaseURL   string
    LogLevel      string
    CacheEnabled  bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:    getEnvAsInt("SERVER_PORT", 8080),
        DatabaseURL:   getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:      getEnv("LOG_LEVEL", "info"),
        CacheEnabled:  getEnvAsBool("CACHE_ENABLED", true),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 100),
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
    return strings.ToLower(strValue) == "true"
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
    }
    if cfg.DatabaseURL == "" {
        return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
    }
    if cfg.MaxConnections < 1 {
        return &ConfigError{Field: "MaxConnections", Message: "must be at least 1"}
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