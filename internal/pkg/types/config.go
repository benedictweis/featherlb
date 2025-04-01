package types

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Routes []Route `yaml:"routes" validate:"min=1,dive"`
}

type Route struct {
	Host      string        `yaml:"host" validate:"required"`
	Port      int           `yaml:"port" validate:"gt=0"`
	Endpoints []Endpoint    `yaml:"endpoints" validate:"min=1,dive"`
	Strategy  KnownStrategy `yaml:"strategy" validate:"oneof=roundrobin random iphash"`
}

type Endpoint struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"gt=0"`
}

type KnownStrategy string

const (
	StrategyRoundRobin KnownStrategy = "roundrobin"
	StrategyRandom     KnownStrategy = "random"
	StrategyIPHash     KnownStrategy = "iphash"
)

func ReadConfigFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return config, nil
}
