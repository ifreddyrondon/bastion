package bastion_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
)

func TestRecovery(t *testing.T) {
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
