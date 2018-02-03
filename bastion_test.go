package gobastion_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/config"
)

func TestDefaultBastion(t *testing.T) {
	bastion := gobastion.New(nil)
	e := gobastion.Tester(t, bastion)
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		Text().Equal("pong")
}

func TestBastionHelloWorld(t *testing.T) {
	bastion := gobastion.New(nil)
	bastion.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		bastion.Send(w, res)
	})

	expected := map[string]interface{}{"message": "world"}

	e := gobastion.Tester(t, bastion)
	e.GET("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
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

	expected := map[string]interface{}{"message": "world"}
	e := gobastion.Tester(t, bastion)
	e.GET("/api/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestNewRouter(t *testing.T) {
	r := gobastion.NewRouter()
	if r == nil {
		t.Errorf("Expected bastion router not to be nil")
	}
}
