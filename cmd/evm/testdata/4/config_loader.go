package config

import (
	"encoding/json"
	"os"
	"strings"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	overrideFromEnv(&config)
	return &config, nil
}

func overrideFromEnv(config *AppConfig) {
	if val := os.Getenv("SERVER_PORT"); val != "" {
		config.ServerPort = val
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		config.DBHost = val
	}
	if val := os.Getenv("DEBUG_MODE"); val != "" {
		config.DebugMode = strings.ToLower(val) == "true"
	}
}package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort   int
	DatabaseURL  string
	LogLevel     string
	CacheEnabled bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		CacheEnabled: getEnvAsBool("CACHE_ENABLED", true),
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	strValue := getEnv(key, "")
	if strings.ToLower(strValue) == "true" {
		return true
	}
	if strings.ToLower(strValue) == "false" {
		return false
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return &ConfigError{Field: "ServerPort", Message: "port must be between 1 and 65535"}
	}
	if cfg.DatabaseURL == "" {
		return &ConfigError{Field: "DatabaseURL", Message: "database URL cannot be empty"}
	}
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return &ConfigError{Field: "LogLevel", Message: "invalid log level"}
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
	"encoding/json"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"SERVER_DEBUG"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	LogLevel string         `json:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			SSLMode: "disable",
		},
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
			DebugMode:    false,
		},
		LogLevel: "info",
	}

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	overrideFromEnv(config)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func overrideFromEnv(config *AppConfig) {
	overrideStruct(config)
}

func overrideStruct(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Struct {
			overrideStruct(field.Addr().Interface())
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			if intVal, err := strconv.Atoi(envValue); err == nil {
				field.SetInt(int64(intVal))
			}
		case reflect.Bool:
			boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
			field.SetBool(boolVal)
		}
	}
}

func validateConfig(config *AppConfig) error {
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}

	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Server.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if config.Server.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[strings.ToLower(config.LogLevel)] {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	return nil
}