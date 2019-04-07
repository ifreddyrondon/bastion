package bastion

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func isServerClosed(server *http.Server, ch chan<- bool) {
	for {
		if err := server.ListenAndServe(); err != nil {
			ch <- true
			return
		}
	}
}

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()

	app := New()
	visited := false
	mu := sync.RWMutex{}
	f := func() {
		mu.Lock()
		visited = true
		mu.Unlock()
	}
	app.RegisterOnShutdown(f)
	ctx, cancel := context.WithCancel(context.Background())
	go graceful(ctx, app.server, app.logger)
	cancel()
	ch := make(chan bool, 1)
	isServerClosed(app.server, ch)
	<-ch

	mu.RLock()
	require.True(t, visited)
	mu.RUnlock()
}

func TestPrintRoutes(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}
	app := New(DisablePrettyLogging(), LoggerOutput(out))
	app.Get("/", handler)  // GET /todos - read a list of todos
	app.Post("/", handler) // POST /todos - create a new todo and persist it
	app.Route("/{id}", func(r chi.Router) {
		r.Get("/", handler)    // GET /todos/{id} - read a single todo by :id
		r.Put("/", handler)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", handler) // DELETE /todos/{id} - delete a single todo by :id
	})

	printRoutes(app.Mux, app.logger)
	assert.Contains(t, out.String(), `"message":"GET /"`)
	assert.Contains(t, out.String(), `"message":"POST /"`)
	assert.Contains(t, out.String(), `"message":"GET /ping"`)
	assert.Contains(t, out.String(), `"message":"DELETE /{id}/"`)
	assert.Contains(t, out.String(), `"message":"GET /{id}/"`)
	assert.Contains(t, out.String(), `"message":"PUT /{id}/"`)
}

func TestResolveAddress(t *testing.T) {
	tt := []struct {
		name         string
		givenAddr    []string
		expectedAddr string
		outputLog    string
	}{
		{
			name:         "default without env",
			givenAddr:    []string{},
			expectedAddr: ":8080",
			outputLog:    `Using addr :8080`,
		},
		{
			name:         "passing addr",
			givenAddr:    []string{":3000"},
			expectedAddr: ":3000",
			outputLog:    ``,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			app := New(DisablePrettyLogging(), LoggerOutput(out))
			addr := resolveAddress(tc.givenAddr, app.logger)
			assert.Equal(t, tc.expectedAddr, addr)
			assert.Contains(t, out.String(), tc.outputLog)
		})
	}
}

func TestResolveAddressWithEnv(t *testing.T) {
	t.Parallel()

	tempADDR := os.Getenv("ADDR")
	out := &bytes.Buffer{}
	app := New(DisablePrettyLogging(), LoggerOutput(out))
	os.Setenv("ADDR", ":3000")
	addr := resolveAddress(nil, app.logger)
	assert.Equal(t, ":3000", addr)
	assert.Contains(t, out.String(), `Environment variable ADDR=\":3000\"`)
	os.Setenv("ADDR", tempADDR)
}

func TestResolveAddressPanic(t *testing.T) {
	t.Parallel()

	f := func() {
		resolveAddress([]string{":3000", ":8080"}, nil)
	}
	assert.PanicsWithValue(t, "too much parameters", f)
}
