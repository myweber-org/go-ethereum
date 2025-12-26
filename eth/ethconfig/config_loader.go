package config

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

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
    Port         int    `yaml:"port"`
    ReadTimeout  int    `yaml:"read_timeout"`
    WriteTimeout int    `yaml:"write_timeout"`
}

type AppConfig struct {
    Environment string         `yaml:"environment"`
    Debug       bool           `yaml:"debug"`
    Database    DatabaseConfig `yaml:"database"`
    Server      ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        configPath = getDefaultConfigPath()
    }

    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func getDefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.yaml"
    }
    return filepath.Join(homeDir, ".app", "config.yaml")
}

func validateConfig(config *AppConfig) error {
    if config.Environment == "" {
        return fmt.Errorf("environment must be specified")
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host must be specified")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }

    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    return nil
}

func (c *AppConfig) GetDatabaseDSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        c.Database.Username,
        c.Database.Password,
        c.Database.Host,
        c.Database.Port,
        c.Database.Name)
}