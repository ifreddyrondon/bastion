package bastion_test

import (
	"os"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

	opts := bastion.NewOptions()
	assert.Equal(t, "127.0.0.1:8080", opts.Addr)
	assert.Equal(t, "development", opts.Env)
	assert.Equal(t, "/", opts.APIBasepath)
	assert.Equal(t, "looks like something went wrong", opts.API500ErrMessage)
	assert.False(t, opts.NoPrettyLogging)
	assert.Equal(t, opts.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, opts.LoggerWriter)
}

func TestOptionsEnvProduction(t *testing.T) {
	t.Parallel()

	app := bastion.New(bastion.Options{Env: "production"})
	assert.Equal(t, "0.0.0.0:8080", app.Options.Addr)
	assert.Equal(t, "production", app.Options.Env)
	assert.Equal(t, "/", app.Options.APIBasepath)
	assert.Equal(t, "looks like something went wrong", app.Options.API500ErrMessage)
	assert.False(t, app.Options.NoPrettyLogging)
	assert.Equal(t, app.Options.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, app.Options.LoggerWriter)
}

func TestOptionsAddr(t *testing.T) {
	t.Parallel()

	app := bastion.New(bastion.Options{Addr: "1.1.1.1:80"})
	assert.Equal(t, "1.1.1.1:80", app.Options.Addr)
	assert.Equal(t, "development", app.Options.Env)
	assert.Equal(t, "/", app.Options.APIBasepath)
	assert.False(t, app.Options.NoPrettyLogging)
	assert.Equal(t, app.Options.LoggerLevel, bastion.DebugLevel)
	assert.Equal(t, os.Stdout, app.Options.LoggerWriter)
}
