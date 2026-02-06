package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
	Database   struct {
		Host     string `env:"DB_HOST" default:"localhost"`
		Port     int    `env:"DB_PORT" default:"5432"`
		Username string `env:"DB_USER"`
		Password string `env:"DB_PASS"`
	}
	FeatureFlags map[string]bool `env:"FEATURE_FLAGS" default:"{}"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	
	v := reflect.ValueOf(cfg).Elem()
	if err := parseStruct(v, ""); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func parseStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		envTag := field.Tag.Get("env")
		defaultTag := field.Tag.Get("default")
		
		if field.Type.Kind() == reflect.Struct {
			if err := parseStruct(fieldValue, prefix+field.Name+"_"); err != nil {
				return err
			}
			continue
		}
		
		if envTag == "" {
			continue
		}
		
		envKey := prefix + envTag
		envValue := os.Getenv(envKey)
		
		if envValue == "" {
			envValue = defaultTag
		}
		
		if err := setValue(fieldValue, envValue); err != nil {
			return fmt.Errorf("failed to set %s: %w", envKey, err)
		}
	}
	
	return nil
}

func setValue(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}
	
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Map:
		if field.Type().Key().Kind() == reflect.String && 
		   field.Type().Elem().Kind() == reflect.Bool {
			var m map[string]bool
			if err := json.Unmarshal([]byte(value), &m); err != nil {
				return err
			}
			field.Set(reflect.ValueOf(m))
		}
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
	
	return nil
}

func (c *Config) String() string {
	var sb strings.Builder
	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("  ServerPort: %d\n", c.ServerPort))
	sb.WriteString(fmt.Sprintf("  LogLevel: %s\n", c.LogLevel))
	sb.WriteString("  Database:\n")
	sb.WriteString(fmt.Sprintf("    Host: %s\n", c.Database.Host))
	sb.WriteString(fmt.Sprintf("    Port: %d\n", c.Database.Port))
	sb.WriteString(fmt.Sprintf("    Username: %s\n", c.Database.Username))
	sb.WriteString("    Password: [REDACTED]\n")
	sb.WriteString(fmt.Sprintf("  FeatureFlags: %v\n", c.FeatureFlags))
	return sb.String()
}