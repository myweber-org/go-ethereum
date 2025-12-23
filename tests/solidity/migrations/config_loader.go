package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Server struct {
		Port    int    `yaml:"port"`
		Timeout int    `yaml:"timeout"`
		Host    string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded successfully from %s", filePath)
	return &config, nil
}

func ValidateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 {
		return ErrInvalidPort
	}
	if config.Database.Host == "" {
		return ErrMissingDatabaseHost
	}
	return nil
}

var (
	ErrInvalidPort         = errors.New("invalid server port configuration")
	ErrMissingDatabaseHost = errors.New("database host is required")
)