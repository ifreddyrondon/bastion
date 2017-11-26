package config

import (
	"io/ioutil"

	"errors"

	"github.com/ghodss/yaml"
	"github.com/markbates/going/defaults"
)

var (
	ErrorMissingConfigFile = errors.New("missing configuration file at path")
	ErrorUnmarshalConfig   = errors.New("cannot unmarshal configuration file")
)

// Config represents the configuration for bastion. Config are used to define how the application should run.
type Config struct {
	API struct {
		// BasePath is the path where the application is going to be mounted. Default `/`.
		BasePath string `json:"base_path"`
	} `json:"api"`
}

// New, returns a new Config instance with sensible defaults.
func New() *Config {
	return configWithDefaults(new(Config))
}

func configWithDefaults(cfg *Config) *Config {
	cfg.API.BasePath = defaults.String(cfg.API.BasePath, "/")
	return cfg
}

func (cfg *Config) FromFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ErrorMissingConfigFile
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return ErrorUnmarshalConfig
	}
	return nil
}
