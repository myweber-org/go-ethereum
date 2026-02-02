package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    if configPath == "" {
        return nil, fmt.Errorf("config path cannot be empty")
    }

    file, err := os.Open(configPath)
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
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Server.Host == "" {
        config.Server.Host = "localhost"
    }
    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }
    if config.LogLevel == "" {
        config.LogLevel = "info"
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}package config

import (
    "encoding/json"
    "os"
    "strings"
)

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST"`
    Port     int    `json:"port" env:"DB_PORT"`
    Username string `json:"username" env:"DB_USER"`
    Password string `json:"password" env:"DB_PASS"`
    SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
    Address      string `json:"address" env:"SERVER_ADDRESS"`
    Port         int    `json:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
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

    overrideFromEnv(&config)
    
    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
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
    if config.Database.Host == "" {
        return errors.New("database host is required")
    }
    if config.Database.Port < 1 || config.Database.Port > 65535 {
        return errors.New("database port must be between 1 and 65535")
    }
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return errors.New("server port must be between 1 and 65535")
    }
    if config.Server.ReadTimeout < 0 {
        return errors.New("server read timeout cannot be negative")
    }
    if config.Server.WriteTimeout < 0 {
        return errors.New("server write timeout cannot be negative")
    }
    return nil
}