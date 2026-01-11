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
		Driver   string `yaml:"driver" env:"DB_DRIVER"`
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Username string `yaml:"username" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASS"`
		Name     string `yaml:"name" env:"DB_NAME"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	setFromEnv(&cfg.Server.Host, "SERVER_HOST")
	setFromEnvInt(&cfg.Server.Port, "SERVER_PORT")
	setFromEnv(&cfg.Database.Driver, "DB_DRIVER")
	setFromEnv(&cfg.Database.Host, "DB_HOST")
	setFromEnvInt(&cfg.Database.Port, "DB_PORT")
	setFromEnv(&cfg.Database.Username, "DB_USER")
	setFromEnv(&cfg.Database.Password, "DB_PASS")
	setFromEnv(&cfg.Database.Name, "DB_NAME")
	setFromEnv(&cfg.Logging.Level, "LOG_LEVEL")
	setFromEnv(&cfg.Logging.Output, "LOG_OUTPUT")
}

func setFromEnv(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func setFromEnvInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}