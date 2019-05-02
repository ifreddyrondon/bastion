package binder

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	errDefaultXMLDecodingMsg = "cannot unmarshal xml body"
)

// XML default XML binder.
var XML = NewXML()

// XMLDecodingErrMsg sets the error message output when xml.Decode fails to decode an object.
func XMLDecodingErrMsg(msg string) func(*XMLBinder) {
	return func(j *XMLBinder) {
		j.errDecodingMsg = msg
	}
}

// XMLBinder bind the xml encoded data present in the request.
// It implements the Binding and BindingBody interface.
type XMLBinder struct {
	errDecodingMsg string
}

// NewXML returns a new XMLBinder responder instance.
func NewXML(opts ...func(*XMLBinder)) *XMLBinder {
	j := &XMLBinder{errDecodingMsg: errDefaultXMLDecodingMsg}
	for _, o := range opts {
		o(j)
	}
	return j
}

func (b *XMLBinder) decode(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return errors.New(b.errDecodingMsg)
	}
	return valid(obj)
}

// Bind the XML request body to an object, if the object implements Validate the valid method will be called.
func (b *XMLBinder) FromReq(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New(errInvalidRequest)
	}
	return b.decode(req.Body, obj)
}

// Bind the XML body to an object, if the object implements Validate the valid method will be called.
func (b *XMLBinder) FromSrc(body []byte, obj interface{}) error {
	return b.decode(bytes.NewReader(body), obj)
}
