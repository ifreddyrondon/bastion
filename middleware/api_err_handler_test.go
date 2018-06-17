package middleware_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/stretchr/testify/assert"
)

func TestAPIErrHandlerCatch500DefaultMsg(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	expectedRes := map[string]interface{}{
		"message": "looks like something went wrong!",
		"error":   "Internal Server Error",
		"status":  500,
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(500).JSON().
		Object().ContainsMap(expectedRes)

	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"component":"api_error_handler`)
	assert.Contains(t, out.String(), `"response":"this should be logged`)
}

func TestAPIErrHandlerCatch500CustomMsg(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{
		LoggerWriter:     out,
		NoPrettyLogging:  true,
		API500ErrMessage: "test",
	})
	app.APIRouter.Mount("/", handler)

	expectedRes := map[string]interface{}{
		"message": "test",
		"error":   "Internal Server Error",
		"status":  500,
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(500).JSON().
		Object().ContainsMap(expectedRes)

	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"component":"api_error_handler`)
	assert.Contains(t, out.String(), `"response":"this should be logged`)
}

func TestAPIErrHandlerNot500(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this should be flushed"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200).Body().Equal("this should be flushed")

	assert.NotContains(t, out.String(), `"component":"api_error_handler`)
}

func TestAPIErrHandlerFailRenderWhen500(t *testing.T) {
	bastion.DefaultRender = func(http.ResponseWriter) render.Engine {
		return &mockRenderEngine{}
	}

	teardown := func() {
		bastion.DefaultRender = json.NewRender
	}
	defer teardown()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200)
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"component":"api_error_handler"`)
	assert.Contains(t, out.String(), `"error":"error render"`)
}
