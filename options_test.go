package bastion_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

	opts := bastion.New().Options
	assert.Equal(t, "looks like something went wrong", opts.InternalErrMsg)
	assert.False(t, opts.DisableInternalErrorMiddleware)
	assert.False(t, opts.DisableRecoveryMiddleware)
	assert.False(t, opts.DisablePingRouter)
	assert.False(t, opts.DisablePrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
	assert.Equal(t, "debug", opts.Mode)
	assert.Equal(t, "/debug", opts.ProfilerRoutePrefix)
	assert.True(t, opts.EnableProfiler)
}

func TestOptionsLoggerLevel(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.LoggerLevel(bastion.ErrorLevel)).Options
	assert.Equal(t, opts.LoggerLevel, bastion.ErrorLevel)
}

func TestOptionsDefaultLoggerLevelWhenProd(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.Mode(bastion.ProductionMode)).Options
	assert.Equal(t, opts.LoggerLevel, bastion.ErrorLevel)
}

func TestOptionsLoggerLevelBadArg(t *testing.T) {
	t.Parallel()
	f := func() {
		bastion.New(bastion.LoggerLevel("bad"))
	}
	assert.Panics(t, f)
}

func TestOptionsInternalErrMsg(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.InternalErrMsg("test")).Options
	assert.Equal(t, "test", opts.InternalErrMsg)
}

func TestBooleanFunctionalOptions(t *testing.T) {
	t.Parallel()

	assert.True(t, bastion.New(bastion.DisableInternalErrorMiddleware()).Options.DisableInternalErrorMiddleware)
	assert.True(t, bastion.New(bastion.DisableRecoveryMiddleware()).Options.DisableRecoveryMiddleware)
	assert.True(t, bastion.New(bastion.DisablePingRouter()).Options.DisablePingRouter)
	assert.True(t, bastion.New(bastion.DisableLoggerMiddleware()).Options.DisableLoggerMiddleware)
	assert.True(t, bastion.New(bastion.DisablePrettyLogging()).Options.DisablePrettyLogging)
	assert.True(t, bastion.New(bastion.EnableProfiler()).Options.EnableProfiler)
}

func TestOptionsDisablePrettyLoggingWhenProd(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.Mode(bastion.ProductionMode)).Options
	assert.True(t, opts.DisablePrettyLogging)
}

func TestModeWithOption(t *testing.T) {
	t.Parallel()
	app := bastion.New(bastion.Mode("production"))
	assert.Equal(t, "production", app.Options.Mode)
	assert.False(t, app.IsDebug())
}

func TestModeWithGO_ENV(t *testing.T) {
	tempADDR := os.Getenv("GO_ENV")
	os.Setenv("GO_ENV", "production")
	app := bastion.New()
	assert.Equal(t, "production", app.Options.Mode)
	assert.False(t, app.IsDebug())
	os.Setenv("GO_ENV", tempADDR)
}

func TestModeWithGO_ENVIRONMENT(t *testing.T) {
	tempADDR := os.Getenv("GO_ENVIRONMENT")
	os.Setenv("GO_ENVIRONMENT", "production")
	app := bastion.New()
	assert.Equal(t, "production", app.Options.Mode)
	assert.False(t, app.IsDebug())
	os.Setenv("GO_ENVIRONMENT", tempADDR)
}

func TestModeOptionPreferenceOverEnv(t *testing.T) {
	tempADDR := os.Getenv("GO_ENVIRONMENT")
	os.Setenv("GO_ENVIRONMENT", "production")
	app := bastion.New(bastion.Mode("debug"))
	assert.Equal(t, "debug", app.Options.Mode)
	assert.True(t, app.IsDebug())
	os.Setenv("GO_ENVIRONMENT", tempADDR)
}

func TestModeWithOptionBadArg(t *testing.T) {
	t.Parallel()
	f := func() {
		bastion.New(bastion.Mode("bad"))
	}
	assert.PanicsWithValue(t, "bastion mode unknown: bad", f)
}

func TestOptionsProfilerRoutePrefixWhenMissingTrailingSlash(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.ProfilerRoutePrefix("dbg")).Options
	assert.Equal(t, "/dbg", opts.ProfilerRoutePrefix)
}

func TestOptionsProfilerRoutePrefix(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.ProfilerRoutePrefix("/abc")).Options
	assert.Equal(t, "/abc", opts.ProfilerRoutePrefix)
}

func TestEnableProfilerShouldBeFalseForProd(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.Mode(bastion.ProductionMode)).Options
	assert.False(t, opts.EnableProfiler)
}

func TestEnableProfilerShouldBeTrueWhenOptionForProd(t *testing.T) {
	t.Parallel()
	opts := bastion.New(bastion.Mode(bastion.ProductionMode), bastion.EnableProfiler()).Options
	assert.True(t, opts.EnableProfiler)
}
