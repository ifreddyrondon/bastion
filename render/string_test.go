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
	render.Text.Response(rr, http.StatusOK, "test")
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).
		Text().Equal("test")
}

func TestHTMLResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	render.HTML.Response(rr, http.StatusOK, "<h1>Hello World</h1>")
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).ContentType("text/html", "utf-8")
}
