package gobastion

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
)

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(res Responder) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					dump, _ := httputil.DumpRequest(r, false)
					panicInfo := fmt.Sprintf("[Recovery] panic recovered:: req: %v err: %+v\n", string(dump), err)
					fmt.Fprintf(os.Stderr, panicInfo)
					debug.PrintStack()
					res.InternalServerError(w, fmt.Errorf("%s", err))
					return
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
