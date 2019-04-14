package render

import "net/http"

// ByteResponder response []byte with application/octet-stream Content-Type
type ByteResponder string

const (
	// Data returns a ByteResponder who's Response writes a string into a ResponseWriter
	// with the Content-Type as application/octet-stream.
	Data ByteResponder = "application/octet-stream"
)

// Response encoded []byte into ResponseWriter with the HTTP status code.
func (s ByteResponder) Response(w http.ResponseWriter, code int, response []byte) {
	writeContentType(w, string(s))
	write(w, code, response)
}
