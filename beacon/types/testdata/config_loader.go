package config

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
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
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

func ValidateConfig(config *ServerConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port == 0 {
        return fmt.Errorf("database port is required")
    }
    if config.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	APIKey      string
	Debug       bool
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	content = os.ExpandEnv(content)

	lines := strings.Split(content, "\n")
	cfg := &Config{}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DATABASE_URL":
			cfg.DatabaseURL = value
		case "API_KEY":
			cfg.APIKey = value
		case "DEBUG":
			cfg.Debug = value == "true"
		}
	}

	return cfg, nil
}