package gobastion

import (
	"encoding/json"
	"io"
)

type Reader interface {
	Read(reader io.ReadCloser, model interface{}) error
}

// JsonReader, parses the JSON-encoded reader and stores the result in the model.
type JsonReader struct{}

func (_ *JsonReader) Read(reader io.ReadCloser, model interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(model); err != nil {
		return err
	}
	reader.Close()
	return nil
}
