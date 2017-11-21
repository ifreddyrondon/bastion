package bastion_test

import (
	"net/http/httptest"
	"testing"

	"io/ioutil"
	"net/http"

	"os"

	"github.com/ifreddyrondon/gobastion"
)

var server *http.ServeMux

func TestMain(m *testing.M) {
	server = http.NewServeMux()
	server.Handle("/", bastion.NewBastion().Router)
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	return res
}

func TestBastionWithDefaultConfig(t *testing.T) {
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
