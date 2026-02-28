package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST"`
    Port     int    `json:"port" env:"DB_PORT"`
    Username string `json:"username" env:"DB_USER"`
    Password string `json:"password" env:"DB_PASS"`
    Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `json:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
}

type AppConfig struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Debug    bool           `json:"debug" env:"APP_DEBUG"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config AppConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, err
    }

    if err := loadEnvOverrides(&config); err != nil {
        return nil, err
    }

    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func loadEnvOverrides(config *AppConfig) error {
    // Implementation would read environment variables
    // and override config values based on struct tags
    return nil
}

func validateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return &ConfigError{Field: "database.host", Reason: "cannot be empty"}
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return &ConfigError{Field: "database.port", Reason: "must be between 1 and 65535"}
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return &ConfigError{Field: "server.port", Reason: "must be between 1 and 65535"}
    }
    return nil
}

type ConfigError struct {
    Field  string
    Reason string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " " + e.Reason
}

func GetDefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.json"
    }
    return filepath.Join(homeDir, ".app", "config.json")
}