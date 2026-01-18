
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
}