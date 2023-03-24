package server

import (
	"encoding/json"
	"os"
)

type databaseConfig struct {
	DBname   string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type serverConfig struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type config struct {
	databaseConfig `json:"db"`
	serverConfig   `json:"server"`
}

func NewConfig(configPath string) (*config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
