package bastion

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

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
				log.Printf("[recovery] req: %v err: %+v\n", string(dump), err)
				if err = json.NewRender(w).InternalServerError(err); err != nil {
					log.Printf("[recovery] err: %v", err)
				}
				return
			}
		}()
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
