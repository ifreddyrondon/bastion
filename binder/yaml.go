package binder

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	errDefaultYAMLDecodingMsg = "cannot unmarshal yaml body"
)

// YAML default YAML binder.
var YAML = NewYAML()

// YAMLLDecodingErrMsg sets the error message output when yaml.Decode fails to decode an object.
func YAMLDecodingErrMsg(msg string) func(*YAMLBinder) {
	return func(j *YAMLBinder) {
		j.errDecodingMsg = msg
	}
}

// YAMLBinder bind the yaml encoded data present in the request.
// It implements the Binding and BindingBody interface.
type YAMLBinder struct {
	errDecodingMsg string
}

// NewYAML returns a new YAMLBinder responder instance.
func NewYAML(opts ...func(*YAMLBinder)) *YAMLBinder {
	j := &YAMLBinder{errDecodingMsg: errDefaultYAMLDecodingMsg}
	for _, o := range opts {
		o(j)
	}
	return j
}

func (b *YAMLBinder) decode(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return errors.New(b.errDecodingMsg)
	}
	return valid(obj)
}

// Bind the YAML request body to an object, if the object implements Validate the valid method will be called.
func (b *YAMLBinder) FromReq(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New(errInvalidRequest)
	}
	return b.decode(req.Body, obj)
}

// Bind the YAML body to an object, if the object implements Validate the valid method will be called.
func (b *YAMLBinder) FromSrc(body []byte, obj interface{}) error {
	return b.decode(bytes.NewReader(body), obj)
}
