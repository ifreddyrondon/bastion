package gognar

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var r *http.ServeMux

func TestMain(m *testing.M) {
	r = http.NewServeMux()
	r.HandleFunc("/ping", pingHandler)
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestPingHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	res := executeRequest(req)
	if 200 != res.Code {
		t.Errorf("Expected response code to be 200'. Got '%v'", res.Code)
	}
	body, _ := ioutil.ReadAll(res.Body)
	if "pong" != string(body) {
		t.Errorf("Expected response body to be 'pong'. Got '%v'", string(body))
	}
}
