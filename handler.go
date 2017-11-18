package gognar

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadJSON, parses the JSON-encoded reader and stores the result in the model.
func ReadJSON(reader io.ReadCloser, model interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(model); err != nil {
		return err
	}
	reader.Close()
	return nil
}

// ResponseJson, sends a JSON-encoded response in the body of a request with the HTTP status code.
func ResponseJson(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// ResponseJson, sends a JSON-encoded response in the body of a request with the 200 status code.
func Send(w http.ResponseWriter, response interface{}) {
	ResponseJson(w, http.StatusOK, response)
}

// Created, sends a JSON-encoded response in the body of a request with the 201 status code.
func Created(w http.ResponseWriter, response interface{}) {
	ResponseJson(w, http.StatusCreated, response)
}

// NoContent, sends a response without no content with the 204 status code.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Abort, sends a JSON-encoded error response in the body of a request with the HTTP status code.
// The error response contains:
//	* Message: (string) that contains message explaining the error.
//	* Errors: (string) identifier of error messages.
//	* Status: (int) HTTP response states. They can range from 400 (Client Errors) to 500 (Server Errors).
func Abort(w http.ResponseWriter, status int, err string, message string) {
	ResponseJson(w, status, responseError{
		Status:  status,
		Errors:  err,
		Message: message,
	})
}

// BadRequest, sends a JSON-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func BadRequest(w http.ResponseWriter, err error) {
	Abort(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
}

// Abort, sends a JSON-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func NotFound(w http.ResponseWriter, err error) {
	Abort(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
}

