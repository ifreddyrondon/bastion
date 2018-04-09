package bastion

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

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

	app := New(Options{})
	app.server = &http.Server{}
	visited := false
	mu := sync.RWMutex{}
	f := func() {
		mu.Lock()
		visited = true
		mu.Unlock()
	}
	app.RegisterOnShutdown(f)
	ctx, cancel := context.WithCancel(context.Background())
	go graceful(ctx, app)
	cancel()
	ch := make(chan bool, 1)
	isServerClosed(app.server, ch)
	<-ch

	mu.RLock()
	require.True(t, visited)
	mu.RUnlock()
}
