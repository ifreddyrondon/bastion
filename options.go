package bastion

import (
	"fmt"
	"io"
	"os"

	"github.com/gobuffalo/envy"
	"github.com/markbates/going/defaults"
)

const (
	defaultPort     = "8080"
	defaultEnv      = "development"
	defaultDevAddrs = "127.0.0.1"
	// api defaults
	defaultBasePath      = "/"
	default500ErrMessage = "looks like something went wrong!"
)

// Level defines log levels.
type Level uint8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled
)

// Options are used to define how the application should run.
type Options struct {
	// APIBasepath is the path where the bastion api router is going to be mounted. Default `/`.
	APIBasepath string `yaml:"apiBasepath"`
	// API500ErrMessage is the default message returned to the user when catch a 500 status error.
	API500ErrMessage string `yaml:"api500ErrMessage"`
	// Addr is the bind address provided to http.Server. Default is "127.0.0.1:8080"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string `yaml:"addr"`
	// Env is the "environment" in which the App is running. Default is "development".
	Env string `yaml:"env"`
	// NoPrettyLogging don't output a colored human readable version on the out writer.
	NoPrettyLogging bool `yaml:"prettyLogging"`
	// LoggerLevel defines log levels. Default is DebugLevel defines an absent log level.
	LoggerLevel Level `yaml:"loggerLevel"`
	// LoggerWriter logger output writer. Default os.Stdout
	LoggerWriter io.Writer
}

func (o *Options) isDEV() bool {
	return o.Env == "development"
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
	opts.API500ErrMessage = defaults.String(opts.API500ErrMessage, default500ErrMessage)
	if opts.LoggerWriter == nil {
		opts.LoggerWriter = os.Stdout
	}

	return opts
}
