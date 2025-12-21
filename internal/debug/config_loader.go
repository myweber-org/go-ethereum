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

	return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
	if config.Database.Host == "" || config.Database.Port == 0 {
		log.Println("Invalid database configuration")
		return false
	}

	if config.Server.Port < 1 || config.Server.Port > 65535 {
		log.Println("Invalid server port configuration")
		return false
	}

	return true
}