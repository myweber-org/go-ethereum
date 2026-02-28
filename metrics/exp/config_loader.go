package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type DatabaseConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	Username string `json:"username" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Database string `json:"database" env:"DB_NAME"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Features struct {
		CacheEnabled bool `json:"cache_enabled" env:"CACHE_ENABLED"`
		RateLimit    int  `json:"rate_limit" env:"RATE_LIMIT"`
	} `json:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	fileData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fileData, &config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(&config)

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getDefaultConfigPath() string {
	possiblePaths := []string{
		"config.json",
		"config/config.json",
		"/etc/app/config.json",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "config.json"
}

func overrideWithEnvVars(config *AppConfig) {
	overrideStruct(config, "")
}

func overrideStruct(v interface{}, prefix string) {
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag != "" {
			envValue := os.Getenv(envTag)
			if envValue != "" {
				setFieldValue(field, fieldType.Type, envValue)
			}
		}

		if field.Kind() == reflect.Struct {
			overrideStruct(field.Addr().Interface(), prefix+fieldType.Name+"_")
		}
	}
}

func setFieldValue(field reflect.Value, fieldType reflect.Type, value string) {
	switch fieldType.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intVal)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolVal)
		}
	}
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("write timeout cannot be negative")
	}
	if config.Features.RateLimit < 0 {
		return errors.New("rate limit cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[strings.ToLower(config.Server.LogLevel)] {
		return errors.New("invalid log level")
	}

	return nil
}

func SaveConfig(config *AppConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}