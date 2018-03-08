package reader

import "io"

// Engine define method to read requests.
type Engine interface {
	Read(model interface{}) error
}

// Reader returns a Engine to read data from a encoded request.
type Reader func(io.ReadCloser) Engine
