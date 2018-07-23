package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// DefaultPrettyPrintJSONIndent defines the number of spaces to pretty print a json
const DefaultPrettyPrintJSONIndent = "  "

// PrettyPrintJSON set JSON encoding indent to DefaultPrettyPrintJSONIndent
func PrettyPrintJSON() func(*JSON) {
	return func(j *JSON) {
		j.indentValue = DefaultPrettyPrintJSONIndent
	}
}

// JSON encode the response as "application/json" content type
// and implement the Renderer and APIRenderer interface.
type JSON struct {
	w            http.ResponseWriter
	indentPrefix string
	indentValue  string
}

// NewJSON returns a new JSON responder instance.
func NewJSON(w http.ResponseWriter, opts ...func(*JSON)) *JSON {
	j := &JSON{w: w}
	for _, o := range opts {
		o(j)
	}
	return j
}

// Response sends a JSON-encoded v in the body of a request with the HTTP status code.
func (j *JSON) Response(code int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent(j.indentPrefix, j.indentValue)
	if err := enc.Encode(v); err != nil {
		http.Error(j.w, err.Error(), http.StatusInternalServerError)
		return
	}

	j.w.Header().Set("Content-Type", "application/json")
	j.w.WriteHeader(code)
	j.w.Write(buf.Bytes())
}

// Send sends a JSON-encoded v in the body of a request with the 200 status code.
func (j *JSON) Send(v interface{}) {
	j.Response(http.StatusOK, v)
}

// Created sends a JSON-encoded v in the body of a request with the 201 status code.
func (j *JSON) Created(v interface{}) {
	j.Response(http.StatusCreated, v)
}

// NoContent sends a v without no content with the 204 status code.
func (j *JSON) NoContent() {
	j.w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a JSON-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func (j *JSON) BadRequest(err error) {
	s := http.StatusBadRequest
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(http.StatusBadRequest, message)
}

// NotFound sends a JSON-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func (j *JSON) NotFound(err error) {
	s := http.StatusNotFound
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(http.StatusNotFound, message)
}

// MethodNotAllowed sends a JSON-encoded error response in the body of a request with the 405 status code.
// The response will contains the status 405 and error "Method Not Allowed".
func (j *JSON) MethodNotAllowed(err error) {
	s := http.StatusMethodNotAllowed
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(http.StatusMethodNotAllowed, message)
}

// InternalServerError sends a JSON-encoded error response in the body of a request with the 500 status code.
// The response will contains the status 500 and error "Internal Server Error".
func (j *JSON) InternalServerError(err error) {
	s := http.StatusInternalServerError
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(http.StatusInternalServerError, message)
}
