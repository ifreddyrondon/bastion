package middleware_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		panicArg interface{}
	}{
		{
			"recovery with err panic call",
			errors.New("testing recovery"),
		},
		{
			"recovery with string panic call",
			"testing recovery",
		},
		{
			"recovery with empty panic call",
			500,
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
				"message": "looks like something went wrong!",
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
