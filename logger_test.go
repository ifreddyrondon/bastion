package bastion_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func TestLogFromHandlerWithContext(t *testing.T) {
	t.Parallel()

	res := map[string]string{"response": "ok"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := bastion.LoggerFromCtx(r.Context())
		l.Info().Msg("handler")
		render.JSON.Send(w, res)
	})

	out := &bytes.Buffer{}
	app := bastion.New(bastion.DisablePrettyLogging(), bastion.LoggerOutput(out))
	app.Mount("/", handler)

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(200).JSON().
		Object().ContainsMap(res)

	assert.Contains(t, out.String(), `handler`)
}
