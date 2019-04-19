package bastion_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/bastion"
)

func TestDefaultBastion(t *testing.T) {
	t.Parallel()

	app := bastion.New()
	e := bastion.Tester(t, app)
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		Text().Equal("pong")
}

func TestBastionHelloWorld(t *testing.T) {
	t.Parallel()

	app := bastion.New()
	app.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		render.JSON.Send(w, map[string]string{"message": "hello bastion"})
	})

	expected := map[string]interface{}{"message": "hello bastion"}

	e := bastion.Tester(t, app)
	e.GET("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestNotFound(t *testing.T) {
	t.Parallel()
	expected := map[string]interface{}{
		"error":   "Not Found",
		"message": "resource /abc not found",
		"status":  404,
	}
	app := bastion.New()
	e := bastion.Tester(t, app)
	e.GET("/abc").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(expected)
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()
	app := bastion.New()
	app.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		render.JSON.Send(w, map[string]string{"message": "hello bastion"})
	})
	expected := map[string]interface{}{
		"error":   "Method Not Allowed",
		"message": "method POST not allowed for resource /hello",
		"status":  405,
	}
	e := bastion.Tester(t, app)
	e.POST("/hello").
		Expect().
		Status(http.StatusMethodNotAllowed).
		JSON().Object().Equal(expected)
}

func TestMountMiddlewareAfterSetup(t *testing.T) {
	t.Parallel()
	m := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	app := bastion.New()
	app.Use(m)
	app.Get("/", h)
	server := httptest.NewServer(app)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("ok")
}

func TestNewRouter(t *testing.T) {
	t.Parallel()

	r := bastion.NewRouter()
	assert.NotNil(t, r)
}
