package bastion

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Finalizer is an interface to be implemented when is necessary
// to run something before a shutdown of the server or in graceful shutdown.
type Finalizer interface {
	// Finalize is a func that will be executed into the graceful shutdown
	Finalize() error
}

type serverFinalizer struct {
	server *http.Server
	ctx    context.Context
}

func (sf serverFinalizer) Finalize() error {
	log.Printf("[finalizer:server] stoping server")
	if err := sf.server.Shutdown(sf.ctx); err != nil {
		return errors.Wrap(err, "[finalizer:server] unable to finalize")
	}
	return nil
}

// Finalize run all app Finalizer on the same order they are defined.
// The finalizers chain will stop on first error.
func finalize(finalizers []Finalizer) error {
	for _, finalizer := range finalizers {
		if err := finalizer.Finalize(); err != nil {
			return err
		}
	}
	return nil
}
