package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST"`
    Port     int    `json:"port" env:"DB_PORT"`
    Username string `json:"username" env:"DB_USER"`
    Password string `json:"password" env:"DB_PASS"`
    Name     string `json:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Address string `json:"address" env:"SERVER_ADDR"`
    Port    int    `json:"port" env:"SERVER_PORT"`
    Debug   bool   `json:"debug" env:"SERVER_DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    var config AppConfig

    if configPath == "" {
        configPath = getDefaultConfigPath()
    }

    fileData, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := json.Unmarshal(fileData, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config JSON: %w", err)
    }

    overrideFromEnv(&config)

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func getDefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.json"
    }
    return filepath.Join(homeDir, ".app", "config.json")
}

func overrideFromEnv(config *AppConfig) {
    overrideStruct(config)
}

func overrideStruct(s interface{}) {
    v := reflect.ValueOf(s).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)

        if field.Kind() == reflect.Struct {
            overrideStruct(field.Addr().Interface())
            continue
        }

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

func validateConfig(config *AppConfig) error {
    var errors []string

    if config.Database.Host == "" {
        errors = append(errors, "database host is required")
    }
    if config.Database.Port < 1 || config.Database.Port > 65535 {
        errors = append(errors, "database port must be between 1 and 65535")
    }
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        errors = append(errors, "server port must be between 1 and 65535")
    }
    if config.LogLevel == "" {
        config.LogLevel = "info"
    }

    validLogLevels := map[string]bool{
        "debug": true, "info": true, "warn": true, "error": true, "fatal": true,
    }
    if !validLogLevels[strings.ToLower(config.LogLevel)] {
        errors = append(errors, "invalid log level")
    }

    if len(errors) > 0 {
        return fmt.Errorf(strings.Join(errors, "; "))
    }

    return nil
}package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(filePath string) (*ServerConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func ValidateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Port)
    }
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    return nil
}