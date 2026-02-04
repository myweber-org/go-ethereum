
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
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Debug    bool
}

func Load() (*Config, error) {
	cfg := &Config{}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}

	cfg.Database = DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", ""),
		Database: getEnv("DB_NAME", "appdb"),
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}

	cfg.Server = ServerConfig{
		Port:         serverPort,
		ReadTimeout:  parseInt(getEnv("READ_TIMEOUT", "30")),
		WriteTimeout: parseInt(getEnv("WRITE_TIMEOUT", "30")),
	}

	debugStr := strings.ToLower(getEnv("DEBUG", "false"))
	cfg.Debug = debugStr == "true" || debugStr == "1"

	if cfg.Database.Password == "" {
		return nil, errors.New("database password is required")
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return nil, errors.New("server port must be between 1 and 65535")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}