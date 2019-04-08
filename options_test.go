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
	assert.Equal(t, "development", opts.Env)
}

func TestOptionsEnvProduction(t *testing.T) {
	t.Parallel()

	opts := bastion.New(bastion.Env("production")).Options
	assert.Equal(t, "production", opts.Env)
	assert.Equal(t, "looks like something went wrong", opts.InternalErrMsg)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
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
	assert.True(t, bastion.New(bastion.DisablePrettyLogging()).Options.DisablePrettyLogging)
}
