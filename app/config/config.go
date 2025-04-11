package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPath       string `yaml:"db_path"`
	Port         int    `yaml:"port"`
	EmployeesURL string `yaml:"employees_url"`
}

func NewConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return cfg, nil
} 