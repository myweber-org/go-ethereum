package config

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
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(path string) (*AppConfig, error) {
    if path == "" {
        path = "config.yaml"
    }

    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port <= 0 {
        return fmt.Errorf("database port must be positive")
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    return nil
}package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Name     string `yaml:"name"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    Logging struct {
        Level  string `yaml:"level"`
        Output string `yaml:"output"`
    } `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
    config := &Config{}

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }

    overrideWithEnv(config)

    return config, nil
}

func overrideWithEnv(config *Config) {
    if envHost := os.Getenv("SERVER_HOST"); envHost != "" {
        config.Server.Host = envHost
    }
    if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
        fmt.Sscanf(envPort, "%d", &config.Server.Port)
    }

    if envDBHost := os.Getenv("DB_HOST"); envDBHost != "" {
        config.Database.Host = envDBHost
    }
    if envDBPort := os.Getenv("DB_PORT"); envDBPort != "" {
        fmt.Sscanf(envDBPort, "%d", &config.Database.Port)
    }
    if envDBName := os.Getenv("DB_NAME"); envDBName != "" {
        config.Database.Name = envDBName
    }
    if envDBUser := os.Getenv("DB_USER"); envDBUser != "" {
        config.Database.Username = envDBUser
    }
    if envDBPass := os.Getenv("DB_PASS"); envDBPass != "" {
        config.Database.Password = envDBPass
    }

    if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
        config.Logging.Level = strings.ToUpper(envLogLevel)
    }
    if envLogOutput := os.Getenv("LOG_OUTPUT"); envLogOutput != "" {
        config.Logging.Output = envLogOutput
    }
}