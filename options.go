package bastion

import (
	"fmt"
	"io"
	"os"

	"github.com/gobuffalo/envy"
	"github.com/markbates/going/defaults"
)

const (
	developmentEnv = "development"
	defaultPort    = "8080"
	defaultAddrs   = "127.0.0.1"
	// api defaults
	defaultBasePath      = "/"
	default500ErrMessage = "looks like something went wrong"
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
	// APIBasepath path where the bastion api router is going to be mounted. Default `/`.
	APIBasepath string `yaml:"apiBasepath"`
	// API500ErrMessage message returned to the user when catch a 500 status error.
	API500ErrMessage string `yaml:"api500ErrMessage"`
	// Addr bind address provided to http.Server. Default is "127.0.0.1:8080"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string `yaml:"addr"`
	// Env "environment" in which the App is running. Default is "development".
	Env string `yaml:"env"`
	// NoPrettyLogging don't output a colored human readable version on the out writer.
	NoPrettyLogging bool `yaml:"prettyLogging"`
	// LoggerLevel defines log levels. Default is DebugLevel defines an absent log level.
	LoggerLevel Level `yaml:"loggerLevel"`
	// LoggerOutput logger output writer. Default os.Stdout
	LoggerOutput io.Writer
}

func (o *Options) isDEV() bool {
	return o.Env == "development"
}

func setDefaultsOpts(opts *Options) {
	port := envy.Get("PORT", defaultPort)
	opts.Env = defaults.String(opts.Env, envy.Get("GO_ENV", developmentEnv))
	addr := "0.0.0.0"
	if opts.Env == developmentEnv {
		addr = defaultAddrs
	}
	envAddr := envy.Get("ADDR", addr)
	opts.Addr = defaults.String(opts.Addr, fmt.Sprintf("%s:%s", envAddr, port))
	opts.APIBasepath = defaults.String(opts.APIBasepath, defaultBasePath)
	opts.API500ErrMessage = defaults.String(opts.API500ErrMessage, default500ErrMessage)
	if opts.LoggerOutput == nil {
		opts.LoggerOutput = os.Stdout
	}
}

// Opt helper type to create functional options
type Opt func(*Bastion)

// APIBasePath set path where the bastion api router is going to be mounted.
func APIBasePath(path string) Opt {
	return func(app *Bastion) {
		app.APIBasepath = path
	}
}

// API500ErrMessage set the message returned to the user when catch a 500 status error.
func API500ErrMessage(msg string) Opt {
	return func(app *Bastion) {
		app.API500ErrMessage = msg
	}
}

// Addr bind address provided to http.Server
func Addr(add string) Opt {
	return func(app *Bastion) {
		app.Addr = add
	}
}

// Env set the "environment" in which the App is running.
func Env(env string) Opt {
	return func(app *Bastion) {
		app.Env = env
	}
}

// NoPrettyLogging turn off the pretty logging.
func NoPrettyLogging() Opt {
	return func(app *Bastion) {
		app.NoPrettyLogging = true
	}
}

// LoggerLevel set the logger level.
func LoggerLevel(lvl Level) Opt {
	return func(app *Bastion) {
		app.LoggerLevel = lvl
	}
}

// LoggerOutput set the logger output writer
func LoggerOutput(w io.Writer) Opt {
	return func(app *Bastion) {
		app.LoggerOutput = w
	}
}
