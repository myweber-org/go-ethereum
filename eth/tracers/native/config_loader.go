package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    ServerPort string `json:"server_port"`
    DatabaseURL string `json:"database_url"`
    LogLevel string `json:"log_level"`
    CacheEnabled bool `json:"cache_enabled"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config
    
    if configPath != "" {
        data, err := os.ReadFile(configPath)
        if err != nil {
            return nil, err
        }
        if err := json.Unmarshal(data, &cfg); err != nil {
            return nil, err
        }
    }
    
    if port := os.Getenv("SERVER_PORT"); port != "" {
        cfg.ServerPort = port
    }
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }
    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        cfg.LogLevel = logLevel
    }
    
    cfg.CacheEnabled = os.Getenv("CACHE_DISABLED") == ""
    
    if cfg.ServerPort == "" {
        cfg.ServerPort = "8080"
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }
    
    return &cfg, nil
}

func DefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return ""
    }
    return filepath.Join(homeDir, ".app", "config.json")
}