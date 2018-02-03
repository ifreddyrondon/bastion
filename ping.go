package gobastion

import (
	"net/http"
)

// Ping endpoint is useful for load balancers or uptime testing
// external services can make a request before hitting any routes.
func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
