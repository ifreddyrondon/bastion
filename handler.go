package gognar

import (
	"encoding/json"
	"io"
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

