package bastion_test

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/render/json"

	"github.com/ifreddyrondon/bastion"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name                    string
		panicArg                interface{}
		expectedMessageResponse string
	}{
		{
			"recovery with err panic call",
			errors.New("testing recovery"),
			"testing recovery",
		},
		{
			"recovery with string panic call",
			"testing recovery",
			"testing recovery",
		},
		{
			"recovery with empty panic call",
			500,
			"500",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(tc.panicArg)
			})

			app := bastion.New(bastion.Options{})
			app.APIRouter.Mount("/", handler)

			expectedRes := map[string]interface{}{
				"message": tc.expectedMessageResponse,
				"error":   "Internal Server Error",
				"status":  500,
			}

			e := bastion.Tester(t, app)
			e.GET("/").Expect().Status(500).JSON().
				Object().ContainsMap(expectedRes)
		})
	}
}

func TestRecoveryLogRequestGET(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(500).JSON()
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"component":"recovery"`)
	assert.Contains(t, out.String(), `"error":"test"`)
	assert.Contains(t, out.String(), `"req":{"url":"/","method":"GET","proto":"HTTP/1.1","host":"","headers":{},"body":""}`)
}

func TestRecoveryLogRequestWithHeaders(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("User-Agent", "Mozilla").Expect().Status(500).JSON()
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"component":"recovery"`)
	assert.Contains(t, out.String(), `"error":"test"`)
	assert.Contains(t, out.String(), `"req":{"url":"/","method":"GET","proto":"HTTP/1.1","host":"","headers":{"user-agent":"Mozilla"}`)
}

func TestRecoveryLogRequestPOST(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	payload := map[string]string{"hello": "world"}

	e := bastion.Tester(t, app)
	e.POST("/").WithJSON(payload).
		Expect().Status(500).JSON()
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"component":"recovery"`)
	assert.Contains(t, out.String(), `"error":"test"`)
	assert.Contains(t, out.String(), `"req":{"url":"/","method":"POST","proto":"HTTP/1.1","host":"","headers":{"content-type":"application/json; charset=utf-8"},"body":"{\"hello\":\"world\"}"}`)
}

func TestLoggerRequestForDevelopment(t *testing.T) {
	// default
	t.Parallel()

	rensponse := map[string]string{"response": "ok"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewRender(w).Send(rensponse); err != nil {
			require.NotNil(t, err)
		}
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200).JSON().
		Object().ContainsMap(rensponse)
	assert.Contains(t, out.String(), `"app":"bastion"`)
	assert.Contains(t, out.String(), `"status":200`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"url":"/"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"level":"info`)
	assert.NotContains(t, out.String(), `"user_agent"`)
}

func TestLoggerRequesLevelErrorForStatusGreaterThan500(t *testing.T) {
	// request with http >= 500 should be tagged as error
	t.Parallel()

	response400 := map[string]interface{}{"message": "test", "error": "Bad Request", "status": 400}
	handler400 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewRender(w).BadRequest(errors.New("test")); err != nil {
			require.NotNil(t, err)
		}
	})

	response500 := map[string]interface{}{"message": "test", "error": "Internal Server Error", "status": 500}
	handler500 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewRender(w).InternalServerError(errors.New("test")); err != nil {
			require.NotNil(t, err)
		}
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})
	app.APIRouter.Mount("/400", handler400)
	app.APIRouter.Mount("/500", handler500)

	e := bastion.Tester(t, app)
	e.GET("/400").Expect().Status(400).JSON().
		Object().ContainsMap(response400)
	assert.Contains(t, out.String(), `"app":"bastion"`)
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
	assert.Contains(t, out.String(), `"app":"bastion"`)
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

	response500 := map[string]interface{}{"message": "test", "error": "Internal Server Error", "status": 500}
	handler500 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewRender(w).InternalServerError(errors.New("test")); err != nil {
			require.NotNil(t, err)
		}
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, Env: "production", NoPrettyLogging: true})
	app.APIRouter.Mount("/500", handler500)

	e := bastion.Tester(t, app)
	e.GET("/500").WithHeader("User-Agent", "Mozilla").Expect().
		Status(500).JSON().
		Object().ContainsMap(response500)
	assert.Contains(t, out.String(), `"app":"bastion"`)
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
		if err := json.NewRender(w).Send(rensponse200); err != nil {
			require.NotNil(t, err)
		}
	})

	response400 := map[string]interface{}{"message": "test", "error": "Bad Request", "status": 400}
	handler400 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewRender(w).BadRequest(errors.New("test")); err != nil {
			require.NotNil(t, err)
		}
	})

	response500 := map[string]interface{}{"message": "test", "error": "Internal Server Error", "status": 500}
	handler500 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewRender(w).InternalServerError(errors.New("test")); err != nil {
			require.NotNil(t, err)
		}
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, LoggerLevel: bastion.ErrorLevel, NoPrettyLogging: true})
	app.APIRouter.Mount("/200", handler200)
	app.APIRouter.Mount("/400", handler400)
	app.APIRouter.Mount("/500", handler500)

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
	assert.Contains(t, out.String(), `"app":"bastion"`)
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
		if err := json.NewRender(w).Send(rensponse); err != nil {
			require.NotNil(t, err)
		}
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200)
	assert.Empty(t, out.Len())
}
