
package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func ReadJSONFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("empty file")
	}

	return json.Unmarshal(data, v)
}

func WriteJSONFile(filename string, v interface{}, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}