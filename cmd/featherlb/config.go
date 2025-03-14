package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Host     string    `yaml:"host"`
	Port     int       `yaml:"port"`
	Backends []Backend `yaml:"backends"`
}

type Backend struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func readConfigFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}
