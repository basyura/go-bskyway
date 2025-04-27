package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var config_ *Config

type Config struct {
	Identifier string
	PassWord   string
}

func Initialize() (*Config, error) {
	config_ = &Config{}
	data, err := os.ReadFile("config.local.json")
	if err != nil {
		return nil, fmt.Errorf("error to load config.loca.json : %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error to unmarshal config : %w", err)
	}

	config_ = &cfg

	return config_, nil
}

func Instance() *Config {
	return config_
}
