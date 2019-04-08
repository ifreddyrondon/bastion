package bastion_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func TestLoggerForDevelopment(t *testing.T) {
	t.Parallel()

	res := map[string]string{"response": "ok"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().Send(w, res)
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.DisablePrettyLogging(), bastion.LoggerOutput(out))
	app.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200).JSON().
		Object().ContainsMap(res)
	assert.Contains(t, out.String(), `"module":"bastion"`)
	assert.Contains(t, out.String(), `"status":200`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"url":"/"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"level":"info`)
	assert.NotContains(t, out.String(), `"user_agent"`)
}

func TestLoggerRequestLevelErrorForStatusGreaterThan500(t *testing.T) {
	// request with http >= 500 should be tagged as error
	t.Parallel()

	response400 := map[string]interface{}{"message": "test", "error": "Bad Request", "status": 400}
	handler400 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().BadRequest(w, errors.New("test"))
	})

	response500 := map[string]interface{}{"message": "looks like something went wrong", "error": "Internal Server Error", "status": 500}
	handler500 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().InternalServerError(w, errors.New("test"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.DisablePrettyLogging(), bastion.LoggerOutput(out))
	app.Mount("/400", handler400)
	app.Mount("/500", handler500)

	e := bastion.Tester(t, app)
	e.GET("/400").Expect().Status(400).JSON().
		Object().ContainsMap(response400)
	assert.Contains(t, out.String(), `"module":"bastion"`)
	assert.Contains(t, out.String(), `"status":400`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"url":"/400"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"level":"info`)
	assert.NotContains(t, out.String(), `"user_agent"`)

	e.GET("/500").Expect().Status(500).JSON().
		Object().ContainsMap(response500)
	assert.Contains(t, out.String(), `"module":"bastion"`)
	assert.Contains(t, out.String(), `"status":500`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"url":"/500"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"level":"error`)
	assert.NotContains(t, out.String(), `"user_agent"`)
}

func TestLoggerRequestForProductionAppendMoreInfo(t *testing.T) {
	// production should append extra info to the log, like ip, user_agent and referer
	t.Parallel()

	response500 := map[string]interface{}{"message": "looks like something went wrong", "error": "Internal Server Error", "status": 500}
	handler500 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().InternalServerError(w, errors.New("test"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.DisablePrettyLogging(), bastion.LoggerOutput(out), bastion.Env("production"))
	app.Mount("/500", handler500)

	e := bastion.Tester(t, app)
	e.GET("/500").WithHeader("User-Agent", "Mozilla").Expect().
		Status(500).JSON().
		Object().ContainsMap(response500)
	assert.Contains(t, out.String(), `"module":"bastion"`)
	assert.Contains(t, out.String(), `"status":500`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"url":"/500"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"level":"error`)
	// extra info
	assert.Contains(t, out.String(), `"user_agent":"Mozilla"`)
}

func TestLoggerRequestErrorLvl(t *testing.T) {
	// Error lvl should only print for >= 500 http status
	t.Parallel()

	rensponse200 := map[string]string{"response": "ok"}
	handler200 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().Send(w, rensponse200)
	})

	response400 := map[string]interface{}{"message": "test", "error": "Bad Request", "status": 400}
	handler400 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().BadRequest(w, errors.New("test"))
	})

	response500 := map[string]interface{}{"message": "looks like something went wrong", "error": "Internal Server Error", "status": 500}
	handler500 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().InternalServerError(w, errors.New("test"))
	})

	out := &bytes.Buffer{}
	app := bastion.New(
		bastion.DisablePrettyLogging(),
		bastion.LoggerOutput(out),
		bastion.LoggerLevel(bastion.ErrorLevel),
	)
	app.Mount("/200", handler200)
	app.Mount("/400", handler400)
	app.Mount("/500", handler500)

	e := bastion.Tester(t, app)
	e.GET("/200").Expect().Status(200).JSON().
		Object().ContainsMap(rensponse200)
	assert.NotContains(t, out.String(), `"status":200`)

	e.GET("/400").Expect().Status(400).JSON().
		Object().ContainsMap(response400)
	assert.NotContains(t, out.String(), `"status":400`)

	e.GET("/500").WithHeader("User-Agent", "Mozilla").Expect().
		Status(500).JSON().
		Object().ContainsMap(response500)
	assert.Contains(t, out.String(), `"module":"bastion"`)
	assert.Contains(t, out.String(), `"status":500`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"url":"/500"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"level":"error`)
}

func TestLoggerRequestPrettyLogging(t *testing.T) {
	// When pretty logging is active it does not put any into the buffer
	t.Parallel()

	rensponse := map[string]string{"response": "ok"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NewJSON().Send(w, rensponse)
	})

	out := &bytes.Buffer{}
	app := bastion.New()
	app.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200)
	assert.Empty(t, out.Len())
}
