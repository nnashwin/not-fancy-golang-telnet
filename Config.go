package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port            string `json:"port"`
	Ip              string `json:"ip"`
	LogFileLocation string `json: "logFile"`
}

func loadConfigFile(fileLoc string) (Config, error) {
	file, err := os.Open(fileLoc)
	var config Config
	if err != nil {
		return config, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	log.Printf("%+v", config.Port)
	log.Printf("%+v", config.Ip)

	return config, nil
}
