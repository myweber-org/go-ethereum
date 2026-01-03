package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Server struct {
		Port    int    `yaml:"port"`
		Host    string `yaml:"host"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded successfully from %s", path)
	return &config, nil
}

func ValidateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 {
		return fmt.Errorf("server port must be positive")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	return nil
}