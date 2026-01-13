package config

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    Database struct {
        Name     string `yaml:"name"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
    } `yaml:"database"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var config AppConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    log.Printf("Configuration loaded from %s", filePath)
    return &config, nil
}