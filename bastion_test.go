package gobastion_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/config"
)

var server *http.ServeMux

func getServerForApp(app *gobastion.Bastion) *http.ServeMux {
	server = http.NewServeMux()
	server.Handle("/", gobastion.GetInternalRouter(app))
	return server
}

func executeRequest(server *http.ServeMux, req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	return res
}

func TestDefaultBastion(t *testing.T) {
	bastion := gobastion.New(nil)
	s := getServerForApp(bastion)
	req, _ := http.NewRequest("GET", "/ping", nil)
	res := executeRequest(s, req)

	if 200 != res.Code {
		t.Errorf("Expected response code to be 200'. Got '%v'", res.Code)
	}
	body, _ := ioutil.ReadAll(res.Body)
	if "pong" != string(body) {
		t.Errorf("Expected response body to be 'pong'. Got '%v'", string(body))
	}
}

func TestBastionHelloWorld(t *testing.T) {
	bastion := gobastion.New(nil)
	bastion.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		bastion.Send(w, res)
	})

	s := getServerForApp(bastion)
	req, _ := http.NewRequest("GET", "/hello", nil)
	res := executeRequest(s, req)
	expected := "{\"message\":\"world\"}\n"

	if 200 != res.Code {
		t.Errorf("Expected response code to be 200'. Got '%v'", res.Code)
	}
	body, _ := ioutil.ReadAll(res.Body)
	if expected != string(body) {
		t.Errorf("Expected response body to be %v. Got %v", expected, string(body))
	}
}

func TestBastionHelloWorldFromFile(t *testing.T) {
	cfg, _ := config.FromFile("./config/testdata/config_test.yaml")
	bastion := gobastion.New(cfg)
	bastion.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		bastion.Send(w, res)
	})

	s := getServerForApp(bastion)
	req, _ := http.NewRequest("GET", "/api/hello", nil)
	res := executeRequest(s, req)
	expected := "{\"message\":\"world\"}\n"

	if 200 != res.Code {
		t.Errorf("Expected response code to be 200'. Got '%v'", res.Code)
	}
	body, _ := ioutil.ReadAll(res.Body)
	if expected != string(body) {
		t.Errorf("Expected response body to be %v. Got %v", expected, string(body))
	}
}

func TestNewRouter(t *testing.T) {
	r := gobastion.NewRouter()
	if r == nil {
		t.Errorf("Expected bastion router not to be nil")
	}
}
