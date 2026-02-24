package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v3"
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
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config.Database)
    overrideFromEnv(&config.Server)

    return &config, nil
}

func overrideFromEnv(config interface{}) {
    val := reflect.ValueOf(config).Elem()
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        envTag := fieldType.Tag.Get("env")
        if envTag == "" {
            continue
        }

        envValue := os.Getenv(envTag)
        if envValue == "" {
            continue
        }

        switch field.Kind() {
        case reflect.String:
            field.SetString(envValue)
        case reflect.Int:
            if intVal, err := strconv.Atoi(envValue); err == nil {
                field.SetInt(int64(intVal))
            }
        case reflect.Bool:
            if boolVal, err := strconv.ParseBool(envValue); err == nil {
                field.SetBool(boolVal)
            }
        }
    }
}

func (c *AppConfig) Validate() error {
    var errors []string

    if c.Database.Host == "" {
        errors = append(errors, "database host is required")
    }
    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        errors = append(errors, "database port must be between 1 and 65535")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        errors = append(errors, "server port must be between 1 and 65535")
    }
    if c.Server.LogLevel != "" {
        validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
        if !validLevels[strings.ToLower(c.Server.LogLevel)] {
            errors = append(errors, "invalid log level")
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("config validation failed: %s", strings.Join(errors, ", "))
    }
    return nil
}