package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
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

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./config.yaml"
	}
	return filepath.Join(homeDir, ".app", "config.yaml")
}

func overrideWithEnvVars(config *Config) {
	if host := os.Getenv("APP_SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("APP_SERVER_PORT"); port != "" {
		config.Server.Port = parsePort(port)
	}
	if host := os.Getenv("APP_DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("APP_DB_PORT"); port != "" {
		config.Database.Port = parsePort(port)
	}
	if user := os.Getenv("APP_DB_USERNAME"); user != "" {
		config.Database.Username = user
	}
	if pass := os.Getenv("APP_DB_PASSWORD"); pass != "" {
		config.Database.Password = pass
	}
	if name := os.Getenv("APP_DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if level := os.Getenv("APP_LOG_LEVEL"); level != "" {
		config.Logging.Level = strings.ToUpper(level)
	}
	if output := os.Getenv("APP_LOG_OUTPUT"); output != "" {
		config.Logging.Output = output
	}
}

func parsePort(portStr string) int {
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return 0
	}
	return port
}

func validateConfig(config *Config) error {
	if config.Server.Host == "" {
		return errors.New("server host is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("invalid server port")
	}
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("invalid database port")
	}
	if config.Database.Name == "" {
		return errors.New("database name is required")
	}
	if config.Logging.Level != "" {
		validLevels := map[string]bool{
			"DEBUG": true,
			"INFO":  true,
			"WARN":  true,
			"ERROR": true,
			"FATAL": true,
		}
		if !validLevels[config.Logging.Level] {
			return errors.New("invalid log level")
		}
	}
	return nil
}package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.overrideFromEnv()

	return &config, nil
}

func (c *Config) overrideFromEnv() {
	overrideString(&c.Server.Host, "SERVER_HOST")
	overrideInt(&c.Server.Port, "SERVER_PORT")
	overrideString(&c.Database.Host, "DB_HOST")
	overrideInt(&c.Database.Port, "DB_PORT")
	overrideString(&c.Database.Name, "DB_NAME")
	overrideString(&c.Database.User, "DB_USER")
	overrideString(&c.Database.Password, "DB_PASSWORD")
	overrideString(&c.Logging.Level, "LOG_LEVEL")
	overrideString(&c.Logging.Format, "LOG_FORMAT")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}