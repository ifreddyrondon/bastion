package json

import (
	"encoding/json"
	"io"

	"github.com/ifreddyrondon/bastion/reader"
)

// NewReader returns the Engine that unmarshal the json.
func NewReader(body io.ReadCloser) reader.Engine {
	return &Reader{body}
}

// Reader parses the JSON-encoded body and stores the result in the model.
type Reader struct {
	body io.ReadCloser
}

func (r *Reader) Read(model interface{}) error {
	decoder := json.NewDecoder(r.body)
	if err := decoder.Decode(model); err != nil {
		return err
	}
	r.body.Close()
	return nil
}
