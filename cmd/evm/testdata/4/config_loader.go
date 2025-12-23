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
}