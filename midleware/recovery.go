package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"net/http/httputil"

	"github.com/ifreddyrondon/gobastion/utils"
)

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				dump, _ := httputil.DumpRequest(r, false)
				panicInfo := fmt.Sprintf("[Recovery] panic recovered:: req: %v err: %+v\n", string(dump), err)
				fmt.Fprintf(os.Stderr, panicInfo)
				debug.PrintStack()
				utils.InternalServerError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
