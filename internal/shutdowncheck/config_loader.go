package config

import (
    "fmt"
    "io/ioutil"
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
        Name     string `yaml:"name"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func (c *Config) Validate() error {
    if c.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if c.Server.Port <= 0 {
        return fmt.Errorf("server port must be positive")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    return nil
}package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	APIKeys    []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var err error

	cfg.ServerPort, err = getIntEnv("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}

	cfg.DBHost = getStringEnv("DB_HOST", "localhost")
	
	cfg.DBPort, err = getIntEnv("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg.DebugMode = getBoolEnv("DEBUG_MODE", false)
	
	apiKeysStr := getStringEnv("API_KEYS", "")
	if apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getStringEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid integer value for " + key)
		}
		return intVal, nil
	}
	return defaultValue, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolVal
	}
	return defaultValue
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("invalid server port")
	}
	
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return errors.New("invalid database port")
	}
	
	if cfg.DBHost == "" {
		return errors.New("database host cannot be empty")
	}
	
	return nil
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    AllowedIPs []string
}

func Load() (*Config, error) {
    cfg := &Config{}

    port, err := getInt("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    cfg.ServerPort = port

    cfg.DBHost = getString("DB_HOST", "localhost")

    dbPort, err := getInt("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }
    cfg.DBPort = dbPort

    debug, err := getBool("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }
    cfg.DebugMode = debug

    cfg.AllowedIPs = getStringSlice("ALLOWED_IPS", []string{"127.0.0.1"})

    return cfg, nil
}

func getString(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getInt(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.Atoi(value)
    }
    return defaultValue, nil
}

func getBool(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.ParseBool(value)
    }
    return defaultValue, nil
}

func getStringSlice(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")
    }
    return defaultValue
}