# Middleware

Bastion's middlewares are just stdlib `net/http` middleware handlers. There is nothing special about them, which means 
the router and all the tooling is designed to be compatible and friendly with any middleware in the community.

## Recovery 

Gracefully absorb panics and returns a HTTP 500 (Internal Server Error) status if possible.
The stacktrace can be handled through the callback function set by `RecoveryCallback` which receives the request 
and the panic arg as a params.

### Options 
- `RecoveryCallback(f func(req *http.Request, err error))` sets the callback function to handler the request when recovers from panics.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/bastion/middleware"
)

func main() {
	// default
	middleware.Recovery()

	// handler error 
	callback := func(req *http.Request, err error) {
		fmt.Printf("url: %v\n", req.URL.RequestURI())
		fmt.Printf("method: %v\n", req.Method)
		fmt.Printf("proto: %v\n", req.Proto)
		fmt.Println(err)
	}
	middleware.Recovery(middleware.RecoveryCallback(callback))
}
```

## InternalError
InternalError intercept responses to verify their status and handle the error. It gets the response code and 
if it's >= 500 handles the error with a default error message without disclosure internal information. 
The real error can be handled through the callback function `InternalErrCallback`.

### Options 
- `InternalErrMsg(s string)` set default error message to be sent. Default "looks like something went wrong".
- `InternalErrCallback(f func(int, io.Reader))` sets the callback function when internal error middleware catch a 500 error.

```go
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/ifreddyrondon/bastion/middleware"
)

func main() {
	// default
	middleware.InternalError()

	// with options
	middleware.InternalError(
		middleware.InternalErrMsg(errors.New("well, this is awkward")),
	)

	// handler error 
	handlerErr := func(code int, r io.Reader) {
		fmt.Printf("code: %v\n", code)
		var buf bytes.Buffer
		buf.ReadFrom(r)
		fmt.Printf(buf.String())
	}
	middleware.InternalError(middleware.InternalErrCallback(handlerErr))
}
```

## Logger
Logger is a middleware that logs the start and end of each request, along with some useful data about what was 
requested, what the response status was, and how long it took to return.

Alternatively, look at https://github.com/rs/zerolog#integration-with-nethttp.

### Options 
- `AttachLogger(log zerolog.Logger)` chain the logger with the middleware.
- `EnableLogReqIP()` show the request ip.
- `EnableLogUserAgent()` show the user agent of the request.
- `EnableLogReferer()` show referer of the request.
- `DisableLogMethod()` hide the request method.
- `DisableLogURL()` hide the request url.
- `DisableLogStatus()` hide the request status.
- `DisableLogSize()` hide the request size.
- `DisableLogDuration()` hide the request duration.
- `DisableLogRequestID()` hide the request id.

```go
package main

import (
	"github.com/ifreddyrondon/bastion/middleware"
)

func main() {
	// default
	middleware.Logger()

	// for full info in production
	middleware.Logger(
		middleware.EnableLogReqIP(),
		middleware.EnableLogUserAgent(),
		middleware.EnableLogReferer(),
	)
}
```

## Listing

Parses the url from a request and stores a [listing.Listing](https://github.com/ifreddyrondon/bastion/blob/master/middleware/listing/listing.go#L11) on the context, it can be accessed through middleware.GetListing.

Sample usage.. for the url: `/repositories/1?limit=10&offset=25`

```go
package main

import (
	"net/http"
	
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware"
)

func list(w http.ResponseWriter, r *http.Request) {
	listing, _ := middleware.GetListing(r.Context())
	// do something with listing
}

func main() {
	app := bastion.New()
	app.Use(middleware.Listing())
	app.Get("/repositories/{id}", list)
	app.Serve()
}
```

## WrapResponseWriter

What happens when it is necessary to know the http status code or the bytes written or even the response it self?
WrapResponseWriter provides an easy way to capture http related metrics from your application's http.Handlers or event 
hijack the response.

Sample usage.. The `defaultMiddleware` capture the metrics http status code and the bytes written, 
the `copyWriterMiddleware` captures the default metrics and creates a copy of the written content and 
the `hijackWriterMiddleware` does the same as the previous ones but don't flush the content. 

```go
package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware"
)

func h(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	w.Write([]byte("created"))
}

func defaultMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		m, snoop := middleware.WrapResponseWriter(w)
		next.ServeHTTP(snoop, r)
		fmt.Println(m.Code)
		fmt.Println(m.Bytes)
	}
	return http.HandlerFunc(fn)
}

func copyWriterMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var out bytes.Buffer
		m, snoop := middleware.WrapResponseWriter(w, middleware.WriteHook(middleware.CopyWriterHook(&out)))
		next.ServeHTTP(snoop, r)
		fmt.Println(m.Code)
		fmt.Println(m.Bytes)
		fmt.Println(out.String())
	}
	return http.HandlerFunc(fn)
}

func hijackWriterMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var out bytes.Buffer
		m, snoop := middleware.WrapResponseWriter(w, middleware.WriteHook(middleware.HijackWriteHook(&out)))
		next.ServeHTTP(snoop, r)
		fmt.Println(m.Code)
		fmt.Println(m.Bytes)
		fmt.Println(out.String())
	}
	return http.HandlerFunc(fn)
}

func main() {
	app := bastion.New()
	app.With(defaultMiddleware).Get("/", h)
	app.With(copyWriterMiddleware).Get("/copy", h)
	app.With(hijackWriterMiddleware).Get("/hijack", h)
	app.Serve()
}
```
 
 ## Auxiliary middlewares and more references
 
 For more references check [chi middlewares](https://github.com/go-chi/chi#middlewares)