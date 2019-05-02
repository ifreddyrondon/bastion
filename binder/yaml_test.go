package binder_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/binder"
)

func TestNewYAMLFromReq(t *testing.T) {
	t.Parallel()
	body := `
address: la comarca
lat: 123
`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewYAML().FromReq(req, &a)
	assert.Nil(t, err)
	assert.Equal(t, "la comarca", a.Address)
	assert.Equal(t, 123.0, a.Lat)
	assert.Equal(t, 0.0, a.Lng)
}

func TestYAMLFromReqMissingBody(t *testing.T) {
	t.Parallel()
	var a address
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := binder.YAML.FromReq(req, &a)
	assert.EqualError(t, err, "invalid request, body not present")
}

func TestYAMLFromSrc(t *testing.T) {
	t.Parallel()
	src := []byte(`
address: la comarca
lat: 123
`)
	var a address
	err := binder.YAML.FromSrc(src, &a)
	assert.Nil(t, err)
	assert.Equal(t, "la comarca", a.Address)
	assert.Equal(t, 123.0, a.Lat)
	assert.Equal(t, 0.0, a.Lng)
}

func TestYAMLDecodeWithYAMLDecodingErrMsg(t *testing.T) {
	t.Parallel()
	body := `address:\nla comarca`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewYAML(binder.YAMLDecodingErrMsg("test")).FromReq(req, &a)
	assert.EqualError(t, err, "test")
}

func TestYAMLValidate(t *testing.T) {
	t.Parallel()
	body := []byte(`
address: la comarca
lat: -1
`)
	var a address
	err := binder.YAML.FromSrc(body, &a)
	assert.EqualError(t, err, "address lat can't be lower than 0")
}
