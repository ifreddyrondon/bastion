package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/bastion/render"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestDataResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	render.NewData(rr).Response(http.StatusOK, []byte("test"))
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).ContentType("application/octet-stream")
}
