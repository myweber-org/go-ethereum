package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port" env:"SERVER_PORT"`
	DBHost     string `json:"db_host" env:"DB_HOST"`
	DBPort     int    `json:"db_port" env:"DB_PORT"`
	DebugMode  bool   `json:"debug_mode" env:"DEBUG_MODE"`
	APIKey     string `json:"api_key" env:"API_KEY"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	if configPath != "" {
		if err := loadFromFile(configPath, config); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}
	
	if err := loadFromEnv(config); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}
	
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return config, nil
}

func loadFromFile(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	return decoder.Decode(config)
}

func loadFromEnv(config *Config) error {
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()
	
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)
		
		envTag := structField.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		
		if err := setFieldValue(field, envValue); err != nil {
			return fmt.Errorf("failed to set field %s from env %s: %w", 
				structField.Name, envTag, err)
		}
	}
	
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func validateConfig(config *Config) error {
	if config.ServerPort <= 0 || config.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if config.DBHost == "" {
		return errors.New("database host is required")
	}
	
	if config.DBPort <= 0 || config.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	
	if config.APIKey == "" {
		return errors.New("API key is required")
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
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    Debug        bool           `yaml:"debug"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(config *AppConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }
    if config.LogLevel == "" {
        config.LogLevel = "info"
    }
    return nil
}