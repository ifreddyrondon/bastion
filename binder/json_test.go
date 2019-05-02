package binder_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/binder"
)

type number struct {
	Foo interface{} `json:"foo"`
}

func TestNewJSONFromReq(t *testing.T) {
	t.Parallel()
	body := `{"address": "la comarca", "lat": 123}`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewJSON().FromReq(req, &a)
	assert.Nil(t, err)
	assert.Equal(t, "la comarca", a.Address)
	assert.Equal(t, 123.0, a.Lat)
	assert.Equal(t, 0.0, a.Lng)
}

func TestNewJSONFromReqMissingBody(t *testing.T) {
	t.Parallel()
	var a address
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := binder.NewJSON().FromReq(req, &a)
	assert.EqualError(t, err, "invalid request, body not present")
}

func TestJSONFromSrc(t *testing.T) {
	t.Parallel()
	src := []byte(`{"address": "la comarca", "lat": 123}`)
	var a address
	err := binder.JSON.FromSrc(src, &a)
	assert.Nil(t, err)
	assert.Equal(t, "la comarca", a.Address)
	assert.Equal(t, 123.0, a.Lat)
	assert.Equal(t, 0.0, a.Lng)
}

func TestJSONDecodeNumber(t *testing.T) {
	t.Parallel()
	body := `{"foo": 123}`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var n number
	err := binder.JSON.FromReq(req, &n)
	assert.Nil(t, err)
	assert.Equal(t, float64(123), n.Foo)
}

func TestJSONDecodeNumberWithEnableUseNumber(t *testing.T) {
	t.Parallel()
	body := `{"foo": 123}`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var n number
	err := binder.NewJSON(binder.EnableUseNumber()).FromReq(req, &n)
	assert.Nil(t, err)
	assert.Equal(t, json.Number("123"), n.Foo)
}

func TestJSONDecodeWithoutDisallowUnknownFields(t *testing.T) {
	t.Parallel()
	body := `{"foo": 123}`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.JSON.FromReq(req, &a)
	assert.Nil(t, err)
}

func TestJSONDecodeWithDisallowUnknownFields(t *testing.T) {
	t.Parallel()
	body := `{"foo": 123}`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewJSON(binder.DisallowUnknownFields()).FromReq(req, &a)
	assert.EqualError(t, err, "cannot unmarshal json body")
}

func TestJSONDecodeWithJSONDecodingErrMsg(t *testing.T) {
	t.Parallel()
	body := `{`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewJSON(binder.JSONDecodingErrMsg("test")).FromReq(req, &a)
	assert.EqualError(t, err, "test")
}

func TestJSONValidate(t *testing.T) {
	t.Parallel()
	body := []byte(`{"address": "la comarca", "lat": -1}`)
	var a address
	err := binder.JSON.FromSrc(body, &a)
	assert.EqualError(t, err, "address lat can't be lower than 0")
}
