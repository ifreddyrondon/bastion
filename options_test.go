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
	assert.Equal(t, "127.0.0.1:8080", opts.Addr)
	assert.Equal(t, "development", opts.Env)
	assert.Equal(t, "/", opts.APIBasepath)
	assert.Equal(t, "looks like something went wrong", opts.API500ErrMessage)
	assert.False(t, opts.NoPrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
}

func TestOptionsEnvProduction(t *testing.T) {
	t.Parallel()

	opts := bastion.New(bastion.Env("production")).Options
	assert.Equal(t, "0.0.0.0:8080", opts.Addr)
	assert.Equal(t, "production", opts.Env)
	assert.Equal(t, "/", opts.APIBasepath)
	assert.Equal(t, "looks like something went wrong", opts.API500ErrMessage)
	assert.False(t, opts.NoPrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
}

func TestOptionsAddr(t *testing.T) {
	t.Parallel()

	opts := bastion.New(bastion.Addr("1.1.1.1:80")).Options
	assert.Equal(t, "1.1.1.1:80", opts.Addr)
	assert.Equal(t, "development", opts.Env)
	assert.Equal(t, "/", opts.APIBasepath)
	assert.Equal(t, "looks like something went wrong", opts.API500ErrMessage)
	assert.False(t, opts.NoPrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
}

func TestOptionsAPIBasePath(t *testing.T) {
	t.Parallel()

	opts := bastion.New(bastion.APIBasePath("/v1")).Options
	assert.Equal(t, "127.0.0.1:8080", opts.Addr)
	assert.Equal(t, "development", opts.Env)
	assert.Equal(t, "/v1", opts.APIBasepath)
	assert.Equal(t, "looks like something went wrong", opts.API500ErrMessage)
	assert.False(t, opts.NoPrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
}

func TestOptionsAPIErrMsg(t *testing.T) {
	t.Parallel()

	opts := bastion.New(bastion.API500ErrMessage("test")).Options
	assert.Equal(t, "127.0.0.1:8080", opts.Addr)
	assert.Equal(t, "development", opts.Env)
	assert.Equal(t, "/", opts.APIBasepath)
	assert.Equal(t, "test", opts.API500ErrMessage)
	assert.False(t, opts.NoPrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerOutput)
}
