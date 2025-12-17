
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
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
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    var config AppConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("failed to decode config: %w", err)
    }

    overrideFromEnv(&config)

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    overrideStruct(config)
}

func overrideStruct(s interface{}) {
    // Implementation would use reflection to check struct tags
    // and override values from environment variables
    // Simplified for this example
}

func validateConfig(config *AppConfig) error {
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

func (c *AppConfig) String() string {
    var sb strings.Builder
    sb.WriteString("Application Configuration:\n")
    sb.WriteString(fmt.Sprintf("  Debug Mode: %v\n", c.Debug))
    sb.WriteString(fmt.Sprintf("  Server Port: %d\n", c.Server.Port))
    sb.WriteString(fmt.Sprintf("  Database Host: %s\n", c.Database.Host))
    sb.WriteString(fmt.Sprintf("  Database Name: %s\n", c.Database.Database))
    return sb.String()
}