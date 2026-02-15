package config

import (
	"encoding/json"
	"os"
	"sync"
)

type AppConfig struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
}

var (
	config     *AppConfig
	configOnce sync.Once
)

func LoadConfig() (*AppConfig, error) {
	var err error
	configOnce.Do(func() {
		config = &AppConfig{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     5432,
			DebugMode:  getEnv("DEBUG", "false") == "true",
		}

		configFile := getEnv("CONFIG_FILE", "")
		if configFile != "" {
			fileConfig, fileErr := loadConfigFromFile(configFile)
			if fileErr == nil {
				mergeConfigs(config, fileConfig)
			} else {
				err = fileErr
			}
		}
	})

	return config, err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadConfigFromFile(filename string) (*AppConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileConfig AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fileConfig); err != nil {
		return nil, err
	}

	return &fileConfig, nil
}

func mergeConfigs(base, override *AppConfig) {
	if override.ServerPort != "" {
		base.ServerPort = override.ServerPort
	}
	if override.DBHost != "" {
		base.DBHost = override.DBHost
	}
	if override.DBPort != 0 {
		base.DBPort = override.DBPort
	}
	base.DebugMode = override.DebugMode
}

func GetConfig() *AppConfig {
	if config == nil {
		LoadConfig()
	}
	return config
}