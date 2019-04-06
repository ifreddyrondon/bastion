package bastion_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func TestLogWithBastionLogger(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}
	app := bastion.New(bastion.NoPrettyLogging(), bastion.LoggerOutput(out))

	app.Logger.Info().Msg("main")
	assert.Contains(t, out.String(), `"main"`)
}

func TestLogFromHandlerWithContext(t *testing.T) {
	t.Parallel()

	res := map[string]string{"response": "ok"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := bastion.LoggerFromCtx(r.Context())
		l.Info().Msg("handler")
		render.NewJSON().Send(w, res)
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.NoPrettyLogging(), bastion.LoggerOutput(out))
	app.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200).JSON().
		Object().ContainsMap(res)

	assert.Contains(t, out.String(), `handler`)
}
