package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserPreferences struct {
	Theme       string  `json:"theme"`
	Language    string  `json:"language"`
	NotificationsEnabled bool `json:"notifications_enabled"`
	Volume      float64 `json:"volume"`
}

func (up *UserPreferences) Validate() error {
	if up.Theme == "" {
		up.Theme = "light"
	}
	if up.Language == "" {
		up.Language = "en"
	}
	if up.Volume < 0.0 || up.Volume > 1.0 {
		return fmt.Errorf("volume must be between 0.0 and 1.0")
	}
	return nil
}

func (up *UserPreferences) ApplyDefaults() {
	if up.Theme == "" {
		up.Theme = "light"
	}
	if up.Language == "" {
		up.Language = "en"
	}
	if up.Volume == 0.0 {
		up.Volume = 0.5
	}
}

func main() {
	prefsJSON := `{"theme":"dark","volume":0.8}`
	var prefs UserPreferences
	
	err := json.Unmarshal([]byte(prefsJSON), &prefs)
	if err != nil {
		log.Fatal(err)
	}
	
	prefs.ApplyDefaults()
	
	err = prefs.Validate()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Validated preferences: %+v\n", prefs)
}