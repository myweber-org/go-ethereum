package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port    string `yaml:"port" env:"SERVER_PORT"`
		Timeout int    `yaml:"timeout" env:"SERVER_TIMEOUT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	config.overrideFromEnv()

	return config, nil
}

func (c *Config) overrideFromEnv() {
	overrideStruct(c, "")
}

func overrideStruct(s interface{}, prefix string) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envTag := fieldType.Tag.Get("env")
		yamlTag := fieldType.Tag.Get("yaml")

		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if yamlTag != "" {
				newPrefix = strings.ToUpper(yamlTag) + "_"
			}
			overrideStruct(field.Addr().Interface(), newPrefix)
			continue
		}

		if envTag == "" {
			continue
		}

		envKey := prefix + envTag
		if envValue := os.Getenv(envKey); envValue != "" {
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