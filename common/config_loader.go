package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *Config) error {
	envVars := map[string]string{
		"SERVER_HOST":    &config.Server.Host,
		"SERVER_PORT":    "",
		"DB_HOST":        &config.Database.Host,
		"DB_PORT":        "",
		"DB_NAME":        &config.Database.Name,
		"DB_USER":        &config.Database.User,
		"DB_PASSWORD":    &config.Database.Password,
		"DB_SSL_MODE":    &config.Database.SSLMode,
		"LOG_LEVEL":      &config.Logging.Level,
		"LOG_OUTPUT":     &config.Logging.Output,
	}

	for envKey, fieldPtr := range envVars {
		if val := os.Getenv(envKey); val != "" {
			switch v := fieldPtr.(type) {
			case *string:
				*v = val
			case *int:
				if intVal, err := strconv.Atoi(val); err == nil {
					*v = intVal
				}
			}
		}
	}

	return nil
}

func validateConfig(config *Config) error {
	if config.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}
	if config.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}

	return nil
}