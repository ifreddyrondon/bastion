package gognar_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ifreddyrondon/gognar"
)

type testerReaderCloser struct {
	io.Reader
}

func (t testerReaderCloser) Close() error { return nil }

type address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func addressToBytes(address string, lat, lng float64) []byte {
	res := fmt.Sprintf(`{"address":"%v", "lat":%v, "lng":%v}`, address, lat, lng)
	return []byte(res)
}

func TestReadJSONWithDefinedStruct(t *testing.T) {
	container := new(address)
	expected := struct {
		address  string
		lat, lng float64
	}{"jorge matte gormaz", 1, 1}

	payload := addressToBytes(expected.address, expected.lat, expected.lng)
	input := testerReaderCloser{bytes.NewBuffer(payload)}
	err := gognar.ReadJSON(&input, &container)

	if err != nil {
		t.Fatalf("Expected err to be nil. Got '%v'", err)
	}
	if expected.address != container.Address {
		t.Errorf("Expected address to be '%v'. Got '%v'", expected.address, container.Address)
	}
	if expected.lat != container.Lat {
		t.Errorf("Expected lat to be '%v'. Got '%v'", expected.lat, container.Lat)
	}
	if expected.lng != container.Lng {
		t.Errorf("Expected lng to be '%v'. Got '%v'", expected.lng, container.Lng)
	}
}

func TestReadJSONWithMap(t *testing.T) {
	container := make(map[string]interface{})
	expected := struct {
		address  string
		lat, lng float64
	}{"jorge matte gormaz", 1, 1}

	payload := addressToBytes(expected.address, expected.lat, expected.lng)
	input := testerReaderCloser{bytes.NewBuffer(payload)}
	err := gognar.ReadJSON(&input, &container)

	if err != nil {
		t.Fatalf("Expected err to be nil. Got '%v'", err)
	}
	if expected.address != container["address"] {
		t.Errorf("Expected address to be '%v'. Got '%v'", expected.address, container["address"])
	}
	if expected.lat != container["lat"] {
		t.Errorf("Expected lat to be '%v'. Got '%v'", expected.lat, container["lat"])
	}
	if expected.lng != container["lng"] {
		t.Errorf("Expected lng to be '%v'. Got '%v'", expected.lng, container["lng"])
	}
}

func TestReadJSONError(t *testing.T) {
	container := make(map[string]interface{})
	input := testerReaderCloser{strings.NewReader("`")}
	expectedErr := "invalid character '`' looking for beginning of value"
	err := gognar.ReadJSON(&input, &container)

	if expectedErr != err.Error() {
		t.Fatalf("Expected err to be '%v'. Got '%v'", expectedErr, err.Error())
	}
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
	gognar.ResponseJson(rr, http.StatusOK, a)

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
	gognar.Send(rr, a)

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
	gognar.Created(rr, a)

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
	gognar.NoContent(rr)

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

func TestAbort(t *testing.T) {
	expected := struct {
		contentType string
		status      int
		error       string
		message     string
	}{
		"application/json",
		http.StatusBadRequest,
		"Not Found",
		"Test message",
	}

	rr := httptest.NewRecorder()
	gognar.Abort(rr, http.StatusBadRequest, expected.error, expected.message)

	if expected.status != rr.Code {
		t.Errorf("Expected response code to be '%v'. Got '%v'", expected.status, rr.Code)
	}
	if expected.contentType != rr.Header().Get("Content-type") {
		t.Errorf("Expected response Content-type to be '%v'. Got '%v'",
			expected.contentType, rr.Header().Get("Content-type"))
	}
	expectedBody := responseErrorToString(expected.message, expected.error, expected.status)
	resBody, _ := ioutil.ReadAll(rr.Body)
	if expectedBody != string(resBody) {
		t.Errorf("Expected response body to be '%v'. Got '%v'", expectedBody, string(resBody))
	}
}

func TestBadRequest(t *testing.T) {
	err := errors.New("test")
	rr := httptest.NewRecorder()
	gognar.BadRequest(rr, err)

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
	gognar.NotFound(rr, err)

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
	gognar.MethodNotAllowed(rr, err)

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
