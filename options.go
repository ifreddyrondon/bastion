package bastion

import (
	"io"
	"os"
)

const (
	developmentEnv        = "development"
	defaultInternalErrMsg = "looks like something went wrong"
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
	// InternalErrMsg message returned to the user when catch a 500 status error.
	InternalErrMsg string
	// DisableInternalErrorMiddleware boolean flag to disable the internal error middleware.
	DisableInternalErrorMiddleware bool
	// DisableRecoveryMiddleware boolean flag to disable the recovery middleware.
	DisableRecoveryMiddleware bool
	// DisablePingRouter boolean flag to disable the ping router.
	DisablePingRouter bool
	// DisablePrettyLogging don't output a colored human readable version on the out writer.
	DisablePrettyLogging bool
	// LoggerLevel defines log levels. Default is DebugLevel defines an absent log level.
	LoggerLevel Level
	// LoggerOutput logger output writer. Default os.Stdout
	LoggerOutput io.Writer
	// Env "environment" in which the App is running. Default is "development".
	Env string
}

func (o *Options) isDEV() bool {
	return o.Env == "development"
}

func setDefaultsOpts(opts *Options) {
	opts.Env = defaultString(opts.Env, defaultString(os.Getenv("GO_ENV"), developmentEnv))
	opts.InternalErrMsg = defaultString(opts.InternalErrMsg, defaultInternalErrMsg)
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

// InternalErrMsg set the message returned to the user when catch a 500 status error.
func InternalErrMsg(msg string) Opt {
	return func(app *Bastion) {
		app.InternalErrMsg = msg
	}
}

// DisableInternalErrorMiddleware turn off internal error middleware.
func DisableInternalErrorMiddleware() Opt {
	return func(app *Bastion) {
		app.DisableInternalErrorMiddleware = true
	}
}

// DisableRecoveryMiddleware turn off recovery middleware.
func DisableRecoveryMiddleware() Opt {
	return func(app *Bastion) {
		app.DisableRecoveryMiddleware = true
	}
}

// DisablePingRouter turn off ping route.
func DisablePingRouter() Opt {
	return func(app *Bastion) {
		app.DisablePingRouter = true
	}
}

// DisablePrettyLogging turn off the pretty logging.
func DisablePrettyLogging() Opt {
	return func(app *Bastion) {
		app.DisablePrettyLogging = true
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

// Env set the "environment" in which the App is running.
func Env(env string) Opt {
	return func(app *Bastion) {
		app.Env = env
	}
}
