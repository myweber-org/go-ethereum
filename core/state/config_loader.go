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
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	var errs []string

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		errs = append(errs, "invalid SERVER_PORT value")
	} else {
		cfg.ServerPort = port
	}

	debugStr := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}
	cfg.DatabaseURL = dbURL

	cacheStr := os.Getenv("CACHE_TTL")
	if cacheStr == "" {
		cacheStr = "300"
	}
	cacheTTL, err := strconv.Atoi(cacheStr)
	if err != nil {
		errs = append(errs, "invalid CACHE_TTL value")
	} else {
		cfg.CacheTTL = cacheTTL
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	return cfg, nil
}