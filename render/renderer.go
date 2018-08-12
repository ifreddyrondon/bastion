package render

import "net/http"

// StringRenderer interface manage string responses.
type StringRenderer interface {
	// Response encoded string into ResponseWriter with the HTTP status code.
	Response(w http.ResponseWriter, code int, response string)
}

// ByteRenderer interface manage []byte responses.
type ByteRenderer interface {
	// Response encoded []byte into ResponseWriter with the HTTP status code.
	Response(w http.ResponseWriter, code int, response []byte)
}

// Renderer interface for managing response payloads.
type Renderer interface {
	// Response encoded responses in the ResponseWriter with the HTTP status code.
	Response(w http.ResponseWriter, code int, response interface{})
}

// APIRenderer interface for managing API response payloads.
type APIRenderer interface {
	Renderer
	OKRenderer
	ClientErrRenderer
	ServerErrRenderer
}

// OKRenderer interface for managing success API response payloads.
type OKRenderer interface {
	Send(w http.ResponseWriter, response interface{})
	Created(w http.ResponseWriter, response interface{})
	NoContent(w http.ResponseWriter)
}

// ClientErrRenderer interface for managing API responses when client error.
type ClientErrRenderer interface {
	BadRequest(w http.ResponseWriter, err error)
	NotFound(w http.ResponseWriter, err error)
	MethodNotAllowed(w http.ResponseWriter, err error)
}

// ServerErrRenderer interface for managing API responses when server error.
type ServerErrRenderer interface {
	InternalServerError(w http.ResponseWriter, err error)
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Message string `json:"message,omitempty" xml:"message,attr,omitempty"`
	Error   string `json:"error,omitempty" xml:"error,attr,omitempty"`
	Status  int    `json:"status,omitempty" xml:"status,attr,omitempty"`
}

// NewHTTPError returns a new HTTPError instance.
func NewHTTPError(message, err string, status int) *HTTPError {
	return &HTTPError{
		Message: message,
		Error:   err,
		Status:  status,
	}
}

func write(w http.ResponseWriter, code int, v []byte) {
	w.WriteHeader(code)
	w.Write(v)
}
