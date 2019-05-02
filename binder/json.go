package binder

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	errInvalidRequest         = "invalid request, body not present"
	errDefaultJSONDecodingMsg = "cannot unmarshal json body"
)

// JSON default JSON binder.
var JSON = NewJSON()

// EnableUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
func EnableUseNumber() func(*JSONBinder) {
	return func(j *JSONBinder) {
		j.enableDecoderUseNumber = true
	}
}

// DisallowUnknownFields causes the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func DisallowUnknownFields() func(*JSONBinder) {
	return func(j *JSONBinder) {
		j.disallowUnknownFields = true
	}
}

// JSONDecodingErrMsg sets the error message output when json.Decode fails to decode an object.
func JSONDecodingErrMsg(msg string) func(*JSONBinder) {
	return func(j *JSONBinder) {
		j.errDecodingMsg = msg
	}
}

// JSONBinder bind the json encoded data present in the request.
// It implements the Binding and BindingBody interface.
type JSONBinder struct {
	enableDecoderUseNumber bool
	disallowUnknownFields  bool
	errDecodingMsg         string
}

// NewJSON returns a new JSONRender responder instance.
func NewJSON(opts ...func(*JSONBinder)) *JSONBinder {
	j := &JSONBinder{errDecodingMsg: errDefaultJSONDecodingMsg}
	for _, o := range opts {
		o(j)
	}
	return j
}

func (b *JSONBinder) decode(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if b.disallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if b.enableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return errors.New(b.errDecodingMsg)
	}
	return valid(obj)
}

// Bind the JSON request body to an object, if the object implements Validate the valid method will be called.
func (b *JSONBinder) FromReq(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New(errInvalidRequest)
	}
	return b.decode(req.Body, obj)
}

// Bind the JSON body to an object, if the object implements Validate the valid method will be called.
func (b *JSONBinder) FromSrc(body []byte, obj interface{}) error {
	return b.decode(bytes.NewReader(body), obj)
}
