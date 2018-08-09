package render_test

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/bastion/render"
	"gopkg.in/gavv/httpexpect.v1"
)

type address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func TestJSONResponse(t *testing.T) {
	t.Parallel()

	a := address{"test address", 1, 1}
	expected := map[string]interface{}{"address": "test address", "lat": 1, "lng": 1}

	rr := httptest.NewRecorder()
	render.NewJSON().Response(rr, http.StatusOK, &a)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestJSONOptions(t *testing.T) {
	t.Parallel()

	a := &address{"test address", 1, 1}

	tt := []struct {
		name     string
		opts     []func(*render.JSON)
		a        *address
		expected string
	}{
		{
			"marshal without indent",
			[]func(*render.JSON){},
			a,
			"{\"address\":\"test address\",\"lat\":1,\"lng\":1}\n",
		},
		{
			"marshal with indent (pretty print)",
			[]func(*render.JSON){render.PrettyPrintJSON()},
			a,
			"{\n  \"address\": \"test address\",\n  \"lat\": 1,\n  \"lng\": 1\n}\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			render.NewJSON(tc.opts...).Response(rr, http.StatusOK, tc.a)
			httpexpect.NewResponse(t, rr.Result()).
				Status(http.StatusOK).
				Body().
				Equal(tc.expected)
		})
	}
}

func TestJSONResponseError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	render.NewJSON().Response(rr, http.StatusOK, math.Inf(1))
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusInternalServerError).
		Text().
		Equal("json: unsupported value: +Inf\n")
}

func TestJSONSend(t *testing.T) {
	t.Parallel()

	a := address{"test address", 1, 1}
	expected := map[string]interface{}{"address": "test address", "lat": 1, "lng": 1}

	rr := httptest.NewRecorder()
	render.NewJSON().Send(rr, &a)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestJSONCreated(t *testing.T) {
	t.Parallel()

	a := address{"test address", 1, 1}
	expected := map[string]interface{}{"address": "test address", "lat": 1, "lng": 1}

	rr := httptest.NewRecorder()
	render.NewJSON().Created(rr, &a)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusCreated).
		JSON().Object().Equal(expected)
}

func TestJSONNoContent(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	render.NewJSON().NoContent(rr)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNoContent).NoContent()
}

func TestJSONBadRequest(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := map[string]interface{}{"message": "test", "error": "Bad Request", "status": 400}

	rr := httptest.NewRecorder()
	render.NewJSON().BadRequest(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusBadRequest).
		JSON().Object().Equal(expected)
}

func TestJSONNotFound(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := map[string]interface{}{"message": "test", "error": "Not Found", "status": 404}

	rr := httptest.NewRecorder()
	render.NewJSON().NotFound(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNotFound).
		JSON().Object().Equal(expected)
}

func TestJSONMethodNotAllowed(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := map[string]interface{}{"message": "test", "error": "Method Not Allowed", "status": 405}

	rr := httptest.NewRecorder()
	render.NewJSON().MethodNotAllowed(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusMethodNotAllowed).
		JSON().Object().Equal(expected)
}

func TestJSONInternalServerError(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := map[string]interface{}{"message": "test", "error": "Internal Server Error", "status": 500}

	rr := httptest.NewRecorder()
	render.NewJSON().InternalServerError(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(expected)
}
