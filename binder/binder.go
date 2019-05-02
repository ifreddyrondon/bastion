package binder

import (
	"net/http"
)

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	FromReq(*http.Request, interface{}) error
}

// BindingSrc reads passes body from supplied source instead of a request.
type BindingSrc interface {
	Binding
	FromSrc([]byte, interface{}) error
}

// Validate represents types capable of validating themselves.
type Validate interface {
	Validate() error
}

func valid(obj interface{}) error {
	if val, ok := obj.(Validate); ok {
		return val.Validate()
	}
	return nil
}
