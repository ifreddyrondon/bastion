package bastion_test

import (
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	opts := bastion.NewOptions()
	assert.Equal(t, "127.0.0.1:8080", opts.Addr)
	assert.Equal(t, "development", opts.Env)
	assert.False(t, opts.Debug)
	assert.Equal(t, "/", opts.APIBasepath)
}
