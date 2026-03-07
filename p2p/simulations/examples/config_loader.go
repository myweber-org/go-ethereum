package config

import (
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    dbConfig := DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvAsInt("DB_PORT", 5432),
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASS", ""),
        Database: getEnv("DB_NAME", "appdb"),
        SSLMode:  getEnv("DB_SSL_MODE", "disable"),
    }

    serverConfig := ServerConfig{
        Port:         getEnvAsInt("SERVER_PORT", 8080),
        ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
        WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
    }

    config := &Config{
        Database: dbConfig,
        Server:   serverConfig,
        LogLevel: strings.ToUpper(getEnv("LOG_LEVEL", "INFO")),
    }

    if err := validateConfig(config); err != nil {
        return nil, err
    }

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func validateConfig(config *Config) error {
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
    }

    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
    }

    allowedLogLevels := map[string]bool{
        "DEBUG": true,
        "INFO":  true,
        "WARN":  true,
        "ERROR": true,
        "FATAL": true,
    }

    if !allowedLogLevels[config.LogLevel] {
        return &ConfigError{Field: "LOG_LEVEL", Message: "invalid log level"}
    }

    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    Debug        bool   `yaml:"debug" env:"DEBUG"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    overrideFromEnv(&config)
    return &config, nil
}

func overrideFromEnv(config *Config) {
    overrideString(&config.Database.Host, "DB_HOST")
    overrideInt(&config.Database.Port, "DB_PORT")
    overrideString(&config.Database.Username, "DB_USER")
    overrideString(&config.Database.Password, "DB_PASS")
    overrideString(&config.Database.Name, "DB_NAME")

    overrideInt(&config.Server.Port, "SERVER_PORT")
    overrideInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    overrideInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    overrideBool(&config.Server.Debug, "DEBUG")
}

func overrideString(field *string, envVar string) {
    if val, exists := os.LookupEnv(envVar); exists && val != "" {
        *field = val
    }
}

func overrideInt(field *int, envVar string) {
    if val, exists := os.LookupEnv(envVar); exists && val != "" {
        var intVal int
        if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
            *field = intVal
        }
    }
}

func overrideBool(field *bool, envVar string) {
    if val, exists := os.LookupEnv(envVar); exists {
        *field = val == "true" || val == "1" || val == "yes"
    }
}package config

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
		Name     string `yaml:"name" env:"DB_NAME"`
		Username string `yaml:"username" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
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
	cfg.Server.Host = getEnvOrDefault("SERVER_HOST", cfg.Server.Host)
	cfg.Server.Port = getEnvIntOrDefault("SERVER_PORT", cfg.Server.Port)

	cfg.Database.Driver = getEnvOrDefault("DB_DRIVER", cfg.Database.Driver)
	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvIntOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.Name = getEnvOrDefault("DB_NAME", cfg.Database.Name)
	cfg.Database.Username = getEnvOrDefault("DB_USERNAME", cfg.Database.Username)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)

	cfg.Logging.Level = getEnvOrDefault("LOG_LEVEL", cfg.Logging.Level)
	cfg.Logging.Output = getEnvOrDefault("LOG_OUTPUT", cfg.Logging.Output)
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