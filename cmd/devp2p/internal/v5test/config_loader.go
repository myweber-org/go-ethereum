package config

import (
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    Database struct {
        ConnectionString string `yaml:"connection_string"`
        MaxConnections   int    `yaml:"max_connections"`
    } `yaml:"database"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
}

func LoadConfig(filePath string) (*Config, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *Config) error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    if c.Database.MaxConnections < 1 {
        return fmt.Errorf("max connections must be positive: %d", c.Database.MaxConnections)
    }
    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[c.Logging.Level] {
        return fmt.Errorf("invalid log level: %s", c.Logging.Level)
    }
    return nil
}package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host" env:"SERVER_HOST"`
        Port int    `yaml:"port" env:"SERVER_PORT"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host" env:"DB_HOST"`
        Port     int    `yaml:"port" env:"DB_PORT"`
        Name     string `yaml:"name" env:"DB_NAME"`
        User     string `yaml:"user" env:"DB_USER"`
        Password string `yaml:"password" env:"DB_PASSWORD"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    config := &Config{}
    
    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    overrideFromEnv(config)
    
    return config, nil
}

func overrideFromEnv(config *Config) {
    overrideStruct(config, "")
}

func overrideStruct(s interface{}, prefix string) {
    val := reflect.ValueOf(s).Elem()
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        envTag := fieldType.Tag.Get("env")
        yamlTag := fieldType.Tag.Get("yaml")

        if field.Kind() == reflect.Struct {
            nestedPrefix := prefix
            if yamlTag != "" {
                nestedPrefix = strings.TrimSuffix(prefix+"_"+strings.ToUpper(yamlTag), "_")
            }
            overrideStruct(field.Addr().Interface(), nestedPrefix)
            continue
        }

        if envTag == "" {
            continue
        }

        envKey := envTag
        if prefix != "" {
            envKey = prefix + "_" + envKey
        }

        if envValue, exists := os.LookupEnv(envKey); exists {
            switch field.Kind() {
            case reflect.String:
                field.SetString(envValue)
            case reflect.Int:
                if intValue, err := strconv.Atoi(envValue); err == nil {
                    field.SetInt(int64(intValue))
                }
            case reflect.Bool:
                if boolValue, err := strconv.ParseBool(envValue); err == nil {
                    field.SetBool(boolValue)
                }
            }
        }
    }
}