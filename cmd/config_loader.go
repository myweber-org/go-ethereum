
package config

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
	Database string `json:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `json:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `json:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	var config AppConfig

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, err
		}
	}

	loadFromEnv(&config)
	
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadFromEnv(config *AppConfig) {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	loadStructFromEnv(&config.Server, envMap)
	loadStructFromEnv(&config.Database, envMap)
}

func loadStructFromEnv(target interface{}, envMap map[string]string) {
	// Implementation would use reflection to read struct tags
	// and populate values from environment variables
	// Simplified for this example
}

func validateConfig(config *AppConfig) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return &ConfigError{Field: "server.port", Message: "port must be between 1 and 65535"}
	}

	if config.Database.Host == "" {
		return &ConfigError{Field: "database.host", Message: "database host is required"}
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return &ConfigError{Field: "database.port", Message: "database port must be between 1 and 65535"}
	}

	return nil
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}