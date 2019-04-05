package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/render"
)

func TestDataResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	render.Data.Response(rr, http.StatusOK, []byte("test"))
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).ContentType("application/octet-stream")
}
