package bastion_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
)

func TestRecovery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("testing recovery")
		panic(err)
	})

	app := bastion.New(nil)
	app.APIRouter.Mount("/", handler)

	expectedRes := map[string]interface{}{
		"message": "testing recovery",
		"error":   "Internal Server Error",
		"status":  500,
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().Status(500).JSON().
		Object().ContainsMap(expectedRes)
}
