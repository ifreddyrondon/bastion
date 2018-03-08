package bastion_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/config"
	"github.com/ifreddyrondon/bastion/render"
)

func TestDefaultBastion(t *testing.T) {
	app := bastion.New(nil)
	e := bastion.Tester(t, app)
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		Text().Equal("pong")
}

func TestBastionHelloWorld(t *testing.T) {
	app := bastion.New(nil)
	app.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		render.JSONRender(w).Send(res)
	})

	expected := map[string]interface{}{"message": "world"}

	e := bastion.Tester(t, app)
	e.GET("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestBastionHelloWorldFromFile(t *testing.T) {
	cfg, _ := config.FromFile("./config/testdata/config_test.yaml")
	app := bastion.New(cfg)
	app.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		render.JSONRender(w).Send(res)
	})

	expected := map[string]interface{}{"message": "world"}
	e := bastion.Tester(t, app)
	e.GET("/api/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestNewRouter(t *testing.T) {
	r := bastion.NewRouter()
	if r == nil {
		t.Errorf("Expected bastion router not to be nil")
	}
}
