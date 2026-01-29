package config

import (
	"io/ioutil"
	"log"

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
	Port        int            `yaml:"port"`
	Environment string         `yaml:"environment"`
	Database    DatabaseConfig `yaml:"database"`
}

func LoadConfig(filePath string) (*ServerConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded from %s", filePath)
	return &config, nil
}

func ValidateConfig(config *ServerConfig) error {
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Port)
	}
	if config.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	if config.Database.Port <= 0 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}
	return nil
}package config

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
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Port == 0 {
        config.Port = 8080
    }
    if config.ReadTimeout == 0 {
        config.ReadTimeout = 30
    }
    if config.WriteTimeout == 0 {
        config.WriteTimeout = 30
    }

    return &config, nil
}

func (c *ServerConfig) Validate() error {
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
}