package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `json:"host" yaml:"host"`
		Port int    `json:"port" yaml:"port"`
	} `json:"server" yaml:"server"`
	Database struct {
		Driver   string `json:"driver" yaml:"driver"`
		Host     string `json:"host" yaml:"host"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	} `json:"database" yaml:"database"`
	LogLevel string `json:"log_level" yaml:"log_level"`
}

func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("config file path cannot be empty")
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		err = json.Unmarshal(fileData, config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(fileData, config)
	default:
		return nil, errors.New("unsupported config file format")
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Database.Driver == "" {
		return errors.New("database driver cannot be empty")
	}
	return nil
}