package config

import (
    "errors"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Name     string `yaml:"name"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
    } `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    if config.Server.Host == "" {
        config.Server.Host = "localhost"
    }
    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    if config.Database.Host == "" || config.Database.Name == "" {
        return nil, errors.New("database configuration is incomplete")
    }

    return &config, nil
}