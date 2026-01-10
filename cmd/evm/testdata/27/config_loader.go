package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{}
	
	portStr := getEnvWithDefault("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	config.ServerPort = port
	
	debugStr := getEnvWithDefault("DEBUG_MODE", "false")
	config.DebugMode = strings.ToLower(debugStr) == "true"
	
	config.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
	
	hostsStr := getEnvWithDefault("ALLOWED_HOSTS", "localhost,127.0.0.1")
	config.AllowedHosts = strings.Split(hostsStr, ",")
	
	return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}