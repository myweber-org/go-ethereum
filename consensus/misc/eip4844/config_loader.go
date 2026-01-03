package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    ServerPort string `json:"server_port"`
    DBHost     string `json:"db_host"`
    DBPort     string `json:"db_port"`
    LogLevel   string `json:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config

    if configPath == "" {
        configPath = "config.json"
    }

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    if port := os.Getenv("SERVER_PORT"); port != "" {
        cfg.ServerPort = port
    }

    if host := os.Getenv("DB_HOST"); host != "" {
        cfg.DBHost = host
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        cfg.LogLevel = logLevel
    }

    return &cfg, nil
}