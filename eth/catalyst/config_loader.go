package config

import (
	"errors"
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

type AppConfig struct {
	DB     DatabaseConfig
	Server ServerConfig
	APIKey string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "appdb")
	dbSSL := getEnv("DB_SSL_MODE", "disable")

	cfg.DB = DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUser,
		Password: dbPass,
		Database: dbName,
		SSLMode:  dbSSL,
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT value")
	}
	writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT value")
	}
	debugMode := strings.ToLower(getEnv("DEBUG_MODE", "false")) == "true"

	cfg.Server = ServerConfig{
		Port:         serverPort,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		DebugMode:    debugMode,
	}

	apiKey := getEnv("API_KEY", "")
	if apiKey == "" {
		return nil, errors.New("API_KEY is required")
	}
	cfg.APIKey = apiKey

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}