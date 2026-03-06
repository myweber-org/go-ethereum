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
}package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	DebugMode    bool   `json:"debug_mode"`
	LogLevel     string `json:"log_level"`
}

type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
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
	if val := os.Getenv("DB_HOST"); val != "" {
		config.Database.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		if port, err := parseInt(val); err == nil {
			config.Database.Port = port
		}
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		if port, err := parseInt(val); err == nil {
			config.Server.Port = port
		}
	}
	if val := os.Getenv("DEBUG_MODE"); val != "" {
		config.Server.DebugMode = val == "true" || val == "1"
	}
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		config.Server.LogLevel = val
	}
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}