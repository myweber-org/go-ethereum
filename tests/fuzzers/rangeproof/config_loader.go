package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	MaxWorkers int
	APIKey     string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
		MaxWorkers: getEnvAsInt("MAX_WORKERS", 10),
		APIKey:     getEnv("API_KEY", ""),
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
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value
	}
	return defaultValue
}