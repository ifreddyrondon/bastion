package config

import (
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
)

// Config, contains the configuration of the bastion.
type Config struct {
	API struct {
		BasePath string `json:"base_path"`
	} `json:"api"`
}

func getDefault() *Config {
	cfg := new(Config)
	cfg.API.BasePath = "/"
	return cfg
}

// NewConfig, returns a configuration with the default values.
func NewConfig() *Config {
	return getDefault()
}

func (cfg *Config) FromFile(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Panic(err)
	}
}
