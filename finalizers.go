package gobastion

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Finalizer is a func that will be executed into the graceful shutdown
type Finalizer interface {
	Finalize() error
}

type serverFinalizer struct {
	server *http.Server
	ctx    context.Context
}

func (sf serverFinalizer) Finalize() error {
	log.Printf("[finalizer:server] stoping server")
	if err := sf.server.Shutdown(sf.ctx); err != nil {
		return fmt.Errorf("[finalizer:server] unable to finalize. Got %v", err)
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
