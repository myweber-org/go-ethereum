package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DebugMode  bool
	FeatureFlags []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	port, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port
	
	cfg.DBHost = getEnvString("DB_HOST", "localhost")
	
	dbPort, err := getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	cfg.DBPort = dbPort
	
	debug, err := getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}
	cfg.DebugMode = debug
	
	flags := getEnvString("FEATURE_FLAGS", "")
	if flags != "" {
		cfg.FeatureFlags = strings.Split(flags, ",")
	}
	
	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.Atoi(value)
	}
	return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.ParseBool(value)
	}
	return defaultValue, nil
}package config

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Name     string `yaml:"name"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(filename string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var config AppConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        log.Printf("Invalid server port: %d", config.Server.Port)
        return false
    }
    
    if config.Database.Host == "" {
        log.Print("Database host cannot be empty")
        return false
    }
    
    return true
}package config

import (
    "bufio"
    "os"
    "strings"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    config := make(Config)
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        
        if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
            envVar := strings.TrimPrefix(strings.TrimSuffix(value, "}"), "${")
            if envVal, exists := os.LookupEnv(envVar); exists {
                value = envVal
            }
        }
        
        config[key] = value
    }
    
    return config, scanner.Err()
}

func (c Config) Get(key, defaultValue string) string {
    if val, ok := c[key]; ok {
        return val
    }
    return defaultValue
}