
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type UserPreferences struct {
    Theme      string `json:"theme"`
    Notifications bool `json:"notifications"`
    Language   string `json:"language"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open preferences file: %w", err)
    }
    defer file.Close()

    var prefs UserPreferences
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&prefs); err != nil {
        return nil, fmt.Errorf("failed to decode JSON: %w", err)
    }

    if prefs.Theme == "" {
        prefs.Theme = "light"
    }
    if prefs.Language == "" {
        prefs.Language = "en"
    }

    return &prefs, nil
}

func main() {
    prefs, err := LoadPreferences("preferences.json")
    if err != nil {
        fmt.Printf("Error loading preferences: %v\n", err)
        return
    }

    fmt.Printf("Loaded preferences: Theme=%s, Notifications=%v, Language=%s\n",
        prefs.Theme, prefs.Notifications, prefs.Language)
}