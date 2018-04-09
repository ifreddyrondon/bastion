package bastion_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogWithBastionLogger(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}
	app := bastion.New(bastion.Options{LoggerWriter: out, NoPrettyLogging: true})

	app.Logger.Info().Msg("main")
	assert.Contains(t, out.String(), `"main"`)
}

func TestLogFromHandlerWithContext(t *testing.T) {
	t.Parallel()

	rensponse := map[string]string{"response": "ok"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := bastion.LoggerFromCtx(r.Context())
		l.Info().Msg("handler")
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

	assert.Contains(t, out.String(), `handler`)
}
