package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var config_ *Config

type Config struct {
	Identifier   string
	PassWord     string
	IconCacheDir string
}

func Initialize() (*Config, error) {
	// load local settings
	data, err := os.ReadFile("config.local.json")
	if err != nil {
		return nil, fmt.Errorf("error to load config.local.json : %w", err)
	}

	// convert to instance
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error to unmarshal config : %w", err)
	}

	config_ = &cfg

	// icon cache
	iconCacheDir, err := getIconCacheDir()
	if err != nil {
		return nil, fmt.Errorf("error to create icon cache dir : %w", err)
	}
	config_.IconCacheDir = iconCacheDir

	return config_, nil
}

func Instance() *Config {
	return config_
}

func getIconCacheDir() (string, error) {
	path := filepath.Join(os.TempDir(), "bskyway")
	err := os.MkdirAll(path, 0755)
	return path, err
}
