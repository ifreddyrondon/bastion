package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// DefaultPrettyPrintJSONIndent defines the default number of spaces to pretty print a json
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
	indentPrefix string
	indentValue  string
}

// NewJSON returns a new JSON responder instance.
func NewJSON(opts ...func(*JSON)) *JSON {
	j := &JSON{}
	for _, o := range opts {
		o(j)
	}
	return j
}

// Response sends a JSON-encoded v in the body of a request with the HTTP status code.
func (j *JSON) Response(w http.ResponseWriter, code int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent(j.indentPrefix, j.indentValue)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	write(w, code, buf.Bytes())
}

// Send sends a JSON-encoded v in the body of a request with the 200 status code.
func (j *JSON) Send(w http.ResponseWriter, v interface{}) {
	j.Response(w, http.StatusOK, v)
}

// Created sends a JSON-encoded v in the body of a request with the 201 status code.
func (j *JSON) Created(w http.ResponseWriter, v interface{}) {
	j.Response(w, http.StatusCreated, v)
}

// NoContent sends a v without no content with the 204 status code.
func (j *JSON) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a JSON-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func (j *JSON) BadRequest(w http.ResponseWriter, err error) {
	s := http.StatusBadRequest
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusBadRequest, message)
}

// NotFound sends a JSON-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func (j *JSON) NotFound(w http.ResponseWriter, err error) {
	s := http.StatusNotFound
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusNotFound, message)
}

// MethodNotAllowed sends a JSON-encoded error response in the body of a request with the 405 status code.
// The response will contains the status 405 and error "Method Not Allowed".
func (j *JSON) MethodNotAllowed(w http.ResponseWriter, err error) {
	s := http.StatusMethodNotAllowed
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusMethodNotAllowed, message)
}

// InternalServerError sends a JSON-encoded error response in the body of a request with the 500 status code.
// The response will contains the status 500 and error "Internal Server Error".
func (j *JSON) InternalServerError(w http.ResponseWriter, err error) {
	s := http.StatusInternalServerError
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusInternalServerError, message)
}
