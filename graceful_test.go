package bastion

import (
	"context"
	"net/http"
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
	app := New(nil)
	app.server = &http.Server{}
	visited := false
	f := func() {
		visited = true
	}
	app.RegisterOnShutdown(f)
	ctx, cancel := context.WithCancel(context.Background())
	go graceful(app.server, ctx)
	cancel()
	ch := make(chan bool, 1)
	isServerClosed(app.server, ch)
	<-ch

	require.True(t, visited)
}
