package bastion

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/pkg/errors"
)

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch t := r.(type) {
				case error:
					err = errors.WithStack(t)
				case string:
					err = errors.WithStack(errors.New(t))
				default:
					err = errors.New(fmt.Sprint(t))
				}
				dump, _ := httputil.DumpRequest(req, true)
				fmt.Fprintf(os.Stderr, fmt.Sprintf("[recovery] req: %v err: %+v\n", string(dump), err))
				json.NewRenderer(w).InternalServerError(err)
				return
			}
		}()
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
