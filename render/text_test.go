package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/bastion/render"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestTextResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	render.NewText(rr).Response(http.StatusOK, "test")
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).
		Text().Equal("test")
}
