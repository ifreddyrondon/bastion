package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// DefaultPrettyPrintJSONIndent defines the default number of spaces to pretty print a json
const DefaultPrettyPrintJSONIndent = "  "

const jsonContentType = "application/json; charset=utf-8"

// JSON is the default JSON renderer
var JSON = NewJSON()

// PrettyPrintJSON set JSONRender encoding indent to DefaultPrettyPrintJSONIndent
func PrettyPrintJSON() func(*JSONRender) {
	return func(j *JSONRender) {
		j.indentValue = DefaultPrettyPrintJSONIndent
	}
}

// JSONRender encode the response as "application/json" content type
// and implement the Renderer and APIRenderer interface.
type JSONRender struct {
	indentPrefix string
	indentValue  string
}

// NewJSON returns a new JSONRender responder instance.
func NewJSON(opts ...func(*JSONRender)) *JSONRender {
	j := &JSONRender{}
	for _, o := range opts {
		o(j)
	}
	return j
}

// Response sends a JSONRender-encoded v in the body of a request with the HTTP status code.
func (j *JSONRender) Response(w http.ResponseWriter, code int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent(j.indentPrefix, j.indentValue)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeContentType(w, jsonContentType)
	write(w, code, buf.Bytes())
}

// Send sends a JSONRender-encoded v in the body of a request with the 200 status code.
func (j *JSONRender) Send(w http.ResponseWriter, v interface{}) {
	j.Response(w, http.StatusOK, v)
}

// Created sends a JSONRender-encoded v in the body of a request with the 201 status code.
func (j *JSONRender) Created(w http.ResponseWriter, v interface{}) {
	j.Response(w, http.StatusCreated, v)
}

// NoContent sends a v without no content with the 204 status code.
func (j *JSONRender) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a JSONRender-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func (j *JSONRender) BadRequest(w http.ResponseWriter, err error) {
	s := http.StatusBadRequest
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusBadRequest, message)
}

// NotFound sends a JSONRender-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func (j *JSONRender) NotFound(w http.ResponseWriter, err error) {
	s := http.StatusNotFound
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusNotFound, message)
}

// MethodNotAllowed sends a JSONRender-encoded error response in the body of a request with the 405 status code.
// The response will contains the status 405 and error "Method Not Allowed".
func (j *JSONRender) MethodNotAllowed(w http.ResponseWriter, err error) {
	s := http.StatusMethodNotAllowed
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusMethodNotAllowed, message)
}

// InternalServerError sends a JSONRender-encoded error response in the body of a request with the 500 status code.
// The response will contains the status 500 and error "Internal Server Error".
func (j *JSONRender) InternalServerError(w http.ResponseWriter, err error) {
	s := http.StatusInternalServerError
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	j.Response(w, http.StatusInternalServerError, message)
}
