package bastion

import (
	"fmt"

	"github.com/gobuffalo/envy"
	"github.com/markbates/going/defaults"
)

const (
	defaultPort     = "8080"
	defaultEnv      = "development"
	defaultDevAddrs = "127.0.0.1"
	// api defaults
	defaultBasePath = "/"
)

// Options are used to define how the application should run.
type Options struct {
	// APIBasepath is the path where the bastion api router is going to be mounted. Default `/`.
	APIBasepath string `yaml:"apiBasepath"`
	// Addr is the bind address provided to http.Server. Default is "127.0.0.1:8080"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string `yaml:"addr"`
	// Env is the "environment" in which the App is running. Default is "development".
	Env string `yaml:"env"`
	// Debug flag if Bastion should enable debugging features.
	Debug bool `yaml:"debug"`
}

// NewOptions returns a new Options instance with sensible defaults
func NewOptions() *Options {
	return optionsWithDefaults(&Options{})
}

func optionsWithDefaults(opts *Options) *Options {
	port := envy.Get("PORT", defaultPort)
	opts.Env = defaults.String(opts.Env, envy.Get("GO_ENV", defaultEnv))
	addr := "0.0.0.0"
	if opts.Env == defaultEnv {
		addr = defaultDevAddrs
	}
	envAddr := envy.Get("ADDR", addr)
	opts.Addr = defaults.String(opts.Addr, fmt.Sprintf("%s:%s", envAddr, port))
	opts.APIBasepath = defaults.String(opts.APIBasepath, defaultBasePath)
	return opts
}
