package yaml

import (
	"errors"
	"github.com/HunterPie/Longinus/core/configuration"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	ErrConfigurationNotFound = errors.New("longinus configuration not found")
	ErrInvalidConfiguration  = errors.New("longinus configuration has invalid data in it")
)

func LoadConfig(path string) (*configuration.LonginusConfiguration, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, ErrConfigurationNotFound
	}

	var config configuration.LonginusConfiguration
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, ErrInvalidConfiguration
	}

	return &config, nil
}
