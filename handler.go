package gognar

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadJSON parses the JSON-encoded reader and stores the result in the model.
func ReadJSON(reader io.ReadCloser, model interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(model); err != nil {
		return err
	}
	reader.Close()
	return nil
}

// ResponseJson send a JSON-encoded response in the body of a request with the HTTP status code.
func ResponseJson(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// ResponseJson send a JSON-encoded response in the body of a request with the 200 code.
func Send(w http.ResponseWriter, response interface{}) {
	ResponseJson(w, http.StatusOK, response)
}

// ResponseJson send a JSON-encoded response in the body of a request with the 201 code.
func Created(w http.ResponseWriter, response interface{}) {
	ResponseJson(w, http.StatusCreated, response)
}

