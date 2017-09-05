package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port    string `json:"port"`
	Ip      string `json:"ip"`
	LogFile string `json: "logFile"`
}

func loadConfigFile(fileLoc string) (Config, error) {
	var config Config
	file, err := os.Open(fileLoc)
	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
