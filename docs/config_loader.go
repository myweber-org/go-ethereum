package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		URL      string `yaml:"url" env:"DB_URL"`
		MaxConns int    `yaml:"max_connections" env:"DB_MAX_CONNS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	overrideWithEnv(config)

	return config, nil
}

func overrideWithEnv(c *Config) {
	if val := os.Getenv("SERVER_HOST"); val != "" {
		c.Server.Host = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		if port, err := parseInt(val); err == nil {
			c.Server.Port = port
		}
	}
	if val := os.Getenv("DB_URL"); val != "" {
		c.Database.URL = val
	}
	if val := os.Getenv("DB_MAX_CONNS"); val != "" {
		if maxConns, err := parseInt(val); err == nil {
			c.Database.MaxConns = maxConns
		}
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.LogLevel = strings.ToUpper(val)
	}
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}