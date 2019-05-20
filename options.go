package bastion

import (
	"io"
	"os"
	"strings"
)

const (
	defaultInternalErrMsg = "looks like something went wrong"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
)

const (
	DebugMode      = "debug"
	ProductionMode = "production"
)

type codeMode uint8

const (
	debugCode codeMode = iota
	productionCode
)

func (mode codeMode) String() string {
	return [...]string{DebugMode, ProductionMode}[mode]
}

const defaultProfilerRoutePrefix = "/debug"

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
	// DisableLoggerMiddleware boolean flag to disable the logger middleware.
	DisableLoggerMiddleware bool
	// DisablePrettyLogging don't output a colored human readable version on the out writer.
	DisablePrettyLogging bool
	// LoggerOutput logger output writer. Default os.Stdout
	LoggerOutput io.Writer
	// LoggerLevel defines log levels. Default "debug".
	LoggerLevel string
	// Mode in which the App is running. Default is "debug".
	Mode string
	codeMode
	// ProfilerRoutePrefix is an optional path prefix for profiler subrouter. If left unspecified, `/debug/`
	// is used as the default path prefix.
	ProfilerRoutePrefix string
	// DisableProfiler boolean flag to disable the profiler router.
	DisableProfiler bool
}

// IsDebug check if app is running in debug mode
func (opts Options) IsDebug() bool {
	return opts.codeMode == debugCode
}

func resolveMode(mode string) codeMode {
	modeEnv := defaultString(os.Getenv("GO_ENV"), "")
	if modeEnv == "" {
		modeEnv = defaultString(os.Getenv("GO_ENVIRONMENT"), "")
	}
	v := defaultString(mode, modeEnv)
	var codeMode = debugCode
	switch v {
	case DebugMode, "":
		codeMode = debugCode
	case ProductionMode:
		codeMode = productionCode
	default:
		panic("bastion mode unknown: " + v)
	}
	return codeMode
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

func DisableLoggerMiddleware() Opt {
	return func(app *Bastion) {
		app.DisableLoggerMiddleware = true
	}
}

// DisablePrettyLogging turn off the pretty logging.
func DisablePrettyLogging() Opt {
	return func(app *Bastion) {
		app.DisablePrettyLogging = true
	}
}

// LoggerLevel set the logger level.
func LoggerLevel(lvl string) Opt {
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

// Mode set the mode in which the App is running.
func Mode(mode string) Opt {
	return func(app *Bastion) {
		app.Mode = mode
	}
}

// ProfilerRoutePrefix set the prefix path for the profile router.
func ProfilerRoutePrefix(prefix string) Opt {
	return func(app *Bastion) {
		if !strings.HasPrefix(prefix, "/") {
			app.ProfilerRoutePrefix = "/" + prefix
		} else {
			app.ProfilerRoutePrefix = prefix
		}
	}
}

// DisableProfiler turn on the profiler router.
func DisableProfiler() Opt {
	return func(app *Bastion) {
		app.DisableProfiler = true
	}
}
