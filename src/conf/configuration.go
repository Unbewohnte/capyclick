package conf

import (
	"encoding/json"
	"io"
	"os"
)

const CurrentVersion uint8 = 1

type Configuration struct {
	ConfigurationVersion uint8  `json:"configurationVersion"`
	WindowSize           [2]int `json:"windowSize"`
	LastWindowPosition   [2]int `json:"lastWindowPosition"`
}

// Returns a reasonable default configuration
func Default() Configuration {
	return Configuration{
		ConfigurationVersion: CurrentVersion,
		WindowSize:           [2]int{640, 280},
		LastWindowPosition:   [2]int{0, 0},
	}
}

// Tries to retrieve configuration from given json file
func FromFile(path string) (*Configuration, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	confBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = json.Unmarshal(confBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Creates configuration file with given fields
func Create(path string, conf Configuration) error {
	configFile, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer configFile.Close()

	configJsonBytes, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		return nil
	}

	_, err = configFile.Write(configJsonBytes)
	if err != nil {
		return nil
	}

	return nil
}
