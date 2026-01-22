package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
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
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if config.Environment == "" {
        config.Environment = "development"
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    return &config, nil
}

func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port == 0 {
        return fmt.Errorf("database port is required")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}package config

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

func Load() (*Config, error) {
    cfg := &Config{
        DatabaseURL:  getEnv("DB_URL", "postgres://localhost:5432/app"),
        MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 10),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        AllowedHosts: getEnvAsSlice("ALLOWED_HOSTS", []string{"localhost"}),
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
    if value, err := strconv.ParseBool(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    return strings.Split(strValue, ",")
}