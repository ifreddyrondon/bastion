package config

import (
	"io/ioutil"

	"errors"

	"fmt"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/envy"
	"github.com/markbates/going/defaults"
)

const DefaultPort = "8080"

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
	Server struct {
		// Address is the bind address provided to http.Server. Default is "127.0.0.1:8080"
		// Can be set using ENV vars "ADDR" and "PORT".
		Addr string `json:"address"`
	} `json:"server"`
}

// New, returns a new Config instance with sensible defaults.
func New() *Config {
	return configWithDefaults(new(Config))
}

func configWithDefaults(cfg *Config) *Config {
	cfg.API.BasePath = defaults.String(cfg.API.BasePath, "/")
	port := envy.Get("PORT", DefaultPort)
	host := envy.Get("ADDR", "")
	cfg.Server.Addr = defaults.String(cfg.Server.Addr, fmt.Sprintf("%s:%s", host, port))
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
