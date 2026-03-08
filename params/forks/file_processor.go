package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJSONFile(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func WriteJSONFile(filename string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func main() {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	config := Config{Host: "localhost", Port: 8080}

	err := WriteJSONFile("config.json", config)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	var loadedConfig Config
	err = ReadJSONFile("config.json", &loadedConfig)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Loaded config: %+v\n", loadedConfig)
}