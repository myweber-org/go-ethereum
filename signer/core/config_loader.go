package config

import (
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    Port        int
    DatabaseURL string
    DebugMode   bool
    AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
    config := &AppConfig{}
    
    portStr := getEnvWithDefault("APP_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    config.Port = port
    
    config.DatabaseURL = getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/appdb")
    
    debugStr := getEnvWithDefault("DEBUG_MODE", "false")
    config.DebugMode = strings.ToLower(debugStr) == "true"
    
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
}package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	LogLevel   string `json:"log_level"`
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			LogLevel:   getEnv("LOG_LEVEL", "info"),
		}

		configFile := os.Getenv("CONFIG_FILE")
		if configFile != "" {
			loadFromFile(configFile)
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadFromFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return
	}

	if fileConfig.ServerPort != "" {
		instance.ServerPort = fileConfig.ServerPort
	}
	if fileConfig.DBHost != "" {
		instance.DBHost = fileConfig.DBHost
	}
	if fileConfig.DBPort != "" {
		instance.DBPort = fileConfig.DBPort
	}
	if fileConfig.LogLevel != "" {
		instance.LogLevel = fileConfig.LogLevel
	}
}