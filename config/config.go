package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/envy"
	"github.com/markbates/going/defaults"
)

const (
	DefaultAddress  = ""
	DefaultBasePath = "/"
	DefaultPort     = "8080"
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
	Server struct {
		// Address is the bind address provided to http.Server. Default is "127.0.0.1:8080"
		// Can be set using ENV vars "ADDR" and "PORT".
		Addr string `json:"address"`
	} `json:"server"`
	// Debug flag if Bastion should enable debugging features.
	Debug bool `json:"debug"`
}

// New, returns a new Config instance with sensible defaults.
func DefaultConfig() *Config {
	cfg := new(Config)
	cfg.API.BasePath = DefaultBasePath
	cfg.Server.Addr = fmt.Sprintf("%s:%s", DefaultAddress, DefaultPort)
	return cfg
}

// FromFile is an util function to load the bastion configuration from a config file.
// The config file could it be in YAML or JSON format. Is some attributes are missing
// from the config file it'll be set with the default.
// FromFile takes a special consideration for `server.address` default.
// When it's not provided it'll search the ADDR and PORT environment variables
// first before set the default.
func FromFile(path string) (*Config, error) {
	cfg := new(Config)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrorMissingConfigFile
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, ErrorUnmarshalConfig
	}
	return setDefaults(cfg), nil
}

func setDefaults(cfg *Config) *Config {
	cfg.API.BasePath = defaults.String(cfg.API.BasePath, "/")
	port := envy.Get("PORT", DefaultPort)
	host := envy.Get("ADDR", "")
	cfg.Server.Addr = defaults.String(cfg.Server.Addr, fmt.Sprintf("%s:%s", host, port))
	return cfg
}
