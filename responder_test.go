package gobastion_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/gobastion"
)

var responder gobastion.JsonResponder

type address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func TestResponseJson(t *testing.T) {
	expected := struct {
		body, contentType string
		status            int
	}{
		"{\"address\":\"test address\",\"lat\":1,\"lng\":1}\n",
		"application/json",
		http.StatusOK,
	}

	a := address{"test address", 1, 1}
	rr := httptest.NewRecorder()
	responder.Response(rr, http.StatusOK, a)

	if expected.status != rr.Code {
		t.Errorf("Expected response code to be '%v'. Got '%v'", expected.status, rr.Code)
	}
	if expected.contentType != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be '%v'. Got '%v'",
			expected.contentType, rr.Header().Get("Content-type"))
	}
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expected.body != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expected.body, string(resBody))
	}
}

func TestSend(t *testing.T) {
	expected := struct {
		body, contentType string
	}{
		"{\"address\":\"test address\",\"lat\":1,\"lng\":1}\n",
		"application/json",
	}

	a := address{"test address", 1, 1}
	rr := httptest.NewRecorder()
	responder.Send(rr, a)

	if 200 != rr.Code {
		t.Errorf("Expected response code to be 200. Got '%v'", rr.Code)
	}
	if expected.contentType != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be '%v'. Got '%v'",
			expected.contentType, rr.Header().Get("Content-type"))
	}
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expected.body != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expected.body, string(resBody))
	}
}

func TestCreated(t *testing.T) {
	expected := struct {
		body, contentType string
	}{
		"{\"address\":\"test address\",\"lat\":1,\"lng\":1}\n",
		"application/json",
	}

	a := address{"test address", 1, 1}
	rr := httptest.NewRecorder()
	responder.Created(rr, a)

	if 201 != rr.Code {
		t.Errorf("Expected response code to be 201. Got '%v'", rr.Code)
	}
	if expected.contentType != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be '%v'. Got '%v'",
			expected.contentType, rr.Header().Get("Content-type"))
	}
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expected.body != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expected.body, string(resBody))
	}
}

func TestNoContent(t *testing.T) {
	rr := httptest.NewRecorder()
	responder.NoContent(rr)

	if 204 != rr.Code {
		t.Errorf("Expected response code to be 204. Got '%v'", rr.Code)
	}
	resBody, _ := ioutil.ReadAll(rr.Body)
	if "" != string(resBody) {
		t.Errorf("Expected response body to be empty. Got '%v'", string(resBody))
	}
}

func responseErrorToString(message, error string, status int) string {
	return fmt.Sprintf("{\"message\":\"%v\",\"error\":\"%v\",\"status\":%v}\n", message, error, status)
}

func TestBadRequest(t *testing.T) {
	err := errors.New("test")
	rr := httptest.NewRecorder()
	responder.BadRequest(rr, err)

	if 400 != rr.Code {
		t.Errorf("Expected response code to be '400'. Got '%v'", rr.Code)
	}
	if "application/json" != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be 'application/json'. Got '%v'",
			rr.Header().Get("Content-type"))
	}
	expectedBody := responseErrorToString("test", "Bad Request", 400)
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expectedBody != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expectedBody, string(resBody))
	}
}

func TestNotFound(t *testing.T) {
	err := errors.New("test")
	rr := httptest.NewRecorder()
	responder.NotFound(rr, err)

	if 404 != rr.Code {
		t.Errorf("Expected response code to be '404'. Got '%v'", rr.Code)
	}
	if "application/json" != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be 'application/json'. Got '%v'",
			rr.Header().Get("Content-type"))
	}
	expectedBody := responseErrorToString("test", "Not Found", 404)
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expectedBody != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expectedBody, string(resBody))
	}
}

func TestMethodNotAllowed(t *testing.T) {
	err := errors.New("test")
	rr := httptest.NewRecorder()
	responder.MethodNotAllowed(rr, err)

	if 405 != rr.Code {
		t.Errorf("Expected response code to be '405'. Got '%v'", rr.Code)
	}
	if "application/json" != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be 'application/json'. Got '%v'",
			rr.Header().Get("Content-type"))
	}
	expectedBody := responseErrorToString("test", "Method Not Allowed", 405)
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expectedBody != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expectedBody, string(resBody))
	}
}

func TestInternalServerError(t *testing.T) {
	err := errors.New("test")
	rr := httptest.NewRecorder()
	responder.InternalServerError(rr, err)

	if 500 != rr.Code {
		t.Errorf("Expected response code to be '500'. Got '%v'", rr.Code)
	}
	if "application/json" != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be 'application/json'. Got '%v'",
			rr.Header().Get("Content-type"))
	}
	expectedBody := responseErrorToString("test", "Internal Server Error", 500)
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expectedBody != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expectedBody, string(resBody))
	}
}
