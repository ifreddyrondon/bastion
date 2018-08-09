package render

import "net/http"

type ByteResponder string

const (
	// Data returns a ByteResponder who's Response writes a string into a ResponseWriter
	// with the Content-Type as application/octet-stream.
	Data ByteResponder = "application/octet-stream"
)

func (s ByteResponder) Response(w http.ResponseWriter, code int, response []byte) {
	w.Header().Set("Content-Type", string(s))
	write(w, code, response)
}
