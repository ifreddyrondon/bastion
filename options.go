package bastion

import (
	"io"
	"os"
)

const (
	developmentEnv       = "development"
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
	// API500ErrMessage message returned to the user when catch a 500 status error.
	API500ErrMessage string
	// Env "environment" in which the App is running. Default is "development".
	Env string
	// NoPrettyLogging don't output a colored human readable version on the out writer.
	NoPrettyLogging bool
	// LoggerLevel defines log levels. Default is DebugLevel defines an absent log level.
	LoggerLevel Level
	// LoggerOutput logger output writer. Default os.Stdout
	LoggerOutput io.Writer
	// DisablePingRouter

}

func (o *Options) isDEV() bool {
	return o.Env == "development"
}

func setDefaultsOpts(opts *Options) {
	opts.Env = defaultString(opts.Env, defaultString(os.Getenv("GO_ENV"), developmentEnv))
	opts.API500ErrMessage = defaultString(opts.API500ErrMessage, default500ErrMessage)
	if opts.LoggerOutput == nil {
		opts.LoggerOutput = os.Stdout
	}
}

func defaultString(s1, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
}

// Opt helper type to create functional options
type Opt func(*Bastion)

// API500ErrMessage set the message returned to the user when catch a 500 status error.
func API500ErrMessage(msg string) Opt {
	return func(app *Bastion) {
		app.API500ErrMessage = msg
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
