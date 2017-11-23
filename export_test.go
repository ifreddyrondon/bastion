package gobastion

import "github.com/go-chi/chi"

// GetInternalRouter is a helper function only exported for test.
// It returns the internal mux router to run the server.
func GetInternalRouter(app *Bastion) *chi.Mux {
	return app.r
}
