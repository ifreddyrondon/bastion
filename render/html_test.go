package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/bastion/render"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestHTMLResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	render.NewHTML(rr).Response(http.StatusOK, "<h1>Hello World</h1>")
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).ContentType("text/html", "utf-8")
}
