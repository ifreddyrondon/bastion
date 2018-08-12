package bastion

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
)

// Ping endpoint is useful for load balancers or uptime testing
// external services can make a request before hitting any routes.
func pingHandler(w http.ResponseWriter, _ *http.Request) {
	render.Text.Response(w, http.StatusOK, "pong")
}
