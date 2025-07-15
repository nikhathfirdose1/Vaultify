package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config mirrors the structure of vaultify.yml
type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		LogPath string `yaml:"log_path"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Encryption struct {
		KeyPath    string `yaml:"key_path"`
		RotateDays int    `yaml:"rotate_days"`
	} `yaml:"encryption"`

	Auth struct {
		Tokens []string `yaml:"tokens"`
	} `yaml:"auth"`
}

// LoadConfig opens and parses the YAML config file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return &cfg, nil
}
