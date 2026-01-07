package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    loadFromEnv(&config)

    return &config, nil
}

func loadFromEnv(config *AppConfig) {
    setFromEnv(&config.Database.Host, "DB_HOST")
    setFromEnvInt(&config.Database.Port, "DB_PORT")
    setFromEnv(&config.Database.Username, "DB_USER")
    setFromEnv(&config.Database.Password, "DB_PASS")
    setFromEnv(&config.Database.Name, "DB_NAME")

    setFromEnvInt(&config.Server.Port, "SERVER_PORT")
    setFromEnvInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    setFromEnvInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    setFromEnvBool(&config.Server.DebugMode, "DEBUG_MODE")

    setFromEnv(&config.LogLevel, "LOG_LEVEL")
}

func setFromEnv(field *string, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val
    }
}

func setFromEnvInt(field *int, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        var intVal int
        if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
            *field = intVal
        }
    }
}

func setFromEnvBool(field *bool, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val == "true" || val == "1" || val == "yes"
    }
}

func DefaultConfigPath() string {
    if path := os.Getenv("CONFIG_PATH"); path != "" {
        return path
    }
    return filepath.Join("config", "app.yaml")
}package config

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST"`
    Port     int    `json:"port" env:"DB_PORT"`
    Username string `json:"username" env:"DB_USER"`
    Password string `json:"password" env:"DB_PASS"`
    SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
    Port         int    `json:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `json:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `json:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `json:"debug_mode" env:"DEBUG_MODE"`
}

type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Env      string         `json:"env" env:"APP_ENV"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config
    
    if configPath != "" {
        file, err := os.Open(configPath)
        if err != nil {
            return nil, fmt.Errorf("failed to open config file: %w", err)
        }
        defer file.Close()
        
        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&cfg); err != nil {
            return nil, fmt.Errorf("failed to decode config: %w", err)
        }
    }
    
    overrideFromEnv(&cfg)
    
    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
    overrideStruct(cfg)
}

func overrideStruct(v interface{}) {
    val := reflect.ValueOf(v).Elem()
    typ := val.Type()
    
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        
        if field.Kind() == reflect.Struct {
            overrideStruct(field.Addr().Interface())
            continue
        }
        
        envTag := fieldType.Tag.Get("env")
        if envTag == "" {
            continue
        }
        
        envValue := os.Getenv(envTag)
        if envValue == "" {
            continue
        }
        
        switch field.Kind() {
        case reflect.String:
            field.SetString(envValue)
        case reflect.Int, reflect.Int64:
            if intVal, err := strconv.Atoi(envValue); err == nil {
                field.SetInt(int64(intVal))
            }
        case reflect.Bool:
            boolVal := strings.ToLower(envValue) == "true" || envValue == "1"
            field.SetBool(boolVal)
        }
    }
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if cfg.Server.ReadTimeout < 0 {
        return fmt.Errorf("read timeout cannot be negative")
    }
    if cfg.Server.WriteTimeout < 0 {
        return fmt.Errorf("write timeout cannot be negative")
    }
    return nil
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
        c.Database.Host,
        c.Database.Port,
        c.Database.Username,
        c.Database.Password,
        c.Database.SSLMode,
    )
}