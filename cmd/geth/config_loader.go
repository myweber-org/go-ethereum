
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	CacheTTL   int
	AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}
	cfg.ServerPort = port
	
	debugStr := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	cfg.DatabaseURL = dbURL
	
	ttlStr := os.Getenv("CACHE_TTL")
	if ttlStr == "" {
		ttlStr = "300"
	}
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return nil, errors.New("invalid CACHE_TTL value")
	}
	cfg.CacheTTL = ttl
	
	hostsStr := os.Getenv("ALLOWED_HOSTS")
	if hostsStr != "" {
		cfg.AllowedHosts = strings.Split(hostsStr, ",")
	} else {
		cfg.AllowedHosts = []string{"localhost", "127.0.0.1"}
	}
	
	return cfg, nil
}

func ValidateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if cfg.CacheTTL < 0 {
		return errors.New("cache TTL cannot be negative")
	}
	
	if cfg.DatabaseURL == "" {
		return errors.New("database URL cannot be empty")
	}
	
	return nil
}package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	FeatureFlags map[string]bool
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DebugMode:  getEnvAsBool("DEBUG_MODE", false),
		FeatureFlags: parseFeatureFlags(getEnv("FEATURE_FLAGS", "")),
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
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
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func parseFeatureFlags(flagsStr string) map[string]bool {
	flags := make(map[string]bool)
	if flagsStr == "" {
		return flags
	}

	items := strings.Split(flagsStr, ",")
	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			flagName := strings.TrimSpace(parts[0])
			flagValue := strings.TrimSpace(parts[1])
			if boolValue, err := strconv.ParseBool(flagValue); err == nil {
				flags[flagName] = boolValue
			}
		}
	}
	return flags
}

func validateConfig(cfg *AppConfig) error {
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		return &ConfigError{Field: "SERVER_PORT", Value: cfg.ServerPort}
	}
	if cfg.DBPort < 1 || cfg.DBPort > 65535 {
		return &ConfigError{Field: "DB_PORT", Value: cfg.DBPort}
	}
	return nil
}

type ConfigError struct {
	Field string
	Value interface{}
}

func (e *ConfigError) Error() string {
	return "invalid configuration value for field: " + e.Field
}package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    AllowedOrigins []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    debugStr := getEnv("ENABLE_DEBUG", "false")
    cfg.EnableDebug = strings.ToLower(debugStr) == "true"
    
    originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
    cfg.AllowedOrigins = strings.Split(originsStr, ",")
    
    if err := validate(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func validate(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return strconv.ErrRange
    }
    
    if cfg.DatabaseURL == "" {
        return strconv.ErrSyntax
    }
    
    return nil
}