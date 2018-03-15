package json_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/render/json"
	"gopkg.in/gavv/httpexpect.v1"
)

type address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func TestResponseJson(t *testing.T) {
	tc := struct {
		name       string
		toResponse interface{}
		result     map[string]interface{}
		status     int
	}{
		"response json",
		address{"test address", 1, 1},
		map[string]interface{}{"address": "test address", "lat": 1, "lng": 1},
		http.StatusOK,
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).Response(tc.status, tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(tc.status).
		JSON().Object().Equal(tc.result)
}

func TestSend(t *testing.T) {
	tc := struct {
		name       string
		toResponse interface{}
		result     map[string]interface{}
	}{
		"send",
		address{"test address", 1, 1},
		map[string]interface{}{"address": "test address", "lat": 1, "lng": 1},
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).Send(tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).
		JSON().Object().Equal(tc.result)
}

func TestCreated(t *testing.T) {
	tc := struct {
		name       string
		toResponse interface{}
		result     map[string]interface{}
	}{
		"send",
		address{"test address", 1, 1},
		map[string]interface{}{"address": "test address", "lat": 1, "lng": 1},
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).Created(tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusCreated).
		JSON().Object().Equal(tc.result)
}

func TestNoContent(t *testing.T) {
	rr := httptest.NewRecorder()
	json.NewRender(rr).NoContent()
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNoContent).NoContent()
}

func TestBadRequest(t *testing.T) {
	tc := struct {
		name       string
		toResponse error
		result     map[string]interface{}
	}{
		"Bad Request",
		errors.New("test"),
		map[string]interface{}{"message": "test", "error": "Bad Request", "status": 400},
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).BadRequest(tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusBadRequest).
		JSON().Object().Equal(tc.result)
}

func TestNotFound(t *testing.T) {
	tc := struct {
		name       string
		toResponse error
		result     map[string]interface{}
	}{
		"Not Found",
		errors.New("test"),
		map[string]interface{}{"message": "test", "error": "Not Found", "status": 404},
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).NotFound(tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNotFound).
		JSON().Object().Equal(tc.result)
}

func TestMethodNotAllowed(t *testing.T) {
	tc := struct {
		name       string
		toResponse error
		result     map[string]interface{}
	}{
		"Method Not Allowed",
		errors.New("test"),
		map[string]interface{}{"message": "test", "error": "Method Not Allowed", "status": 405},
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).MethodNotAllowed(tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusMethodNotAllowed).
		JSON().Object().Equal(tc.result)
}

func TestInternalServerError(t *testing.T) {
	tc := struct {
		name       string
		toResponse error
		result     map[string]interface{}
	}{
		"Internal Server Error",
		errors.New("test"),
		map[string]interface{}{"message": "test", "error": "Internal Server Error", "status": 500},
	}

	rr := httptest.NewRecorder()
	err := json.NewRender(rr).InternalServerError(tc.toResponse)
	assert.Nil(t, err)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(tc.result)
}
