package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var errs []string

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		errs = append(errs, "invalid SERVER_PORT value")
	} else {
		cfg.ServerPort = port
	}

	debugStr := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}
	cfg.DatabaseURL = dbURL

	cacheStr := os.Getenv("CACHE_TTL")
	if cacheStr == "" {
		cacheStr = "300"
	}
	cacheTTL, err := strconv.Atoi(cacheStr)
	if err != nil {
		errs = append(errs, "invalid CACHE_TTL value")
	} else {
		cfg.CacheTTL = cacheTTL
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	return cfg, nil
}package config

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty")
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(c *Config) error {
	if c.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}
	return nil
}