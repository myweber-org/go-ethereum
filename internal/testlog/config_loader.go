package config

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
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Port)
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }

    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    return nil
}package config

import (
	"errors"
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
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	Environment  string `yaml:"environment"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	if filePath == "" {
		return nil, errors.New("config file path cannot be empty")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.Environment != "development" && config.Server.Environment != "production" {
		return errors.New("environment must be either 'development' or 'production'")
	}
	return nil
}