
package config

import (
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
		Format string `yaml:"format" env:"LOG_FORMAT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	config.overrideFromEnv()

	return &config, nil
}

func (c *Config) overrideFromEnv() {
	c.Server.Host = getEnvOrDefault("SERVER_HOST", c.Server.Host)
	c.Server.Port = getEnvIntOrDefault("SERVER_PORT", c.Server.Port)

	c.Database.Host = getEnvOrDefault("DB_HOST", c.Database.Host)
	c.Database.Port = getEnvIntOrDefault("DB_PORT", c.Database.Port)
	c.Database.Name = getEnvOrDefault("DB_NAME", c.Database.Name)
	c.Database.User = getEnvOrDefault("DB_USER", c.Database.User)
	c.Database.Password = getEnvOrDefault("DB_PASSWORD", c.Database.Password)
	c.Database.SSLMode = getEnvOrDefault("DB_SSL_MODE", c.Database.SSLMode)

	c.Logging.Level = getEnvOrDefault("LOG_LEVEL", c.Logging.Level)
	c.Logging.Format = getEnvOrDefault("LOG_FORMAT", c.Logging.Format)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}