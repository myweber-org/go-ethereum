package config

import (
    "errors"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    if path == "" {
        return nil, errors.New("config path cannot be empty")
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func validateConfig(config *ServerConfig) error {
    if config.Port <= 0 || config.Port > 65535 {
        return errors.New("invalid server port")
    }

    if config.Database.Host == "" {
        return errors.New("database host is required")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return errors.New("invalid database port")
    }

    if config.Database.Database == "" {
        return errors.New("database name is required")
    }

    return nil
}