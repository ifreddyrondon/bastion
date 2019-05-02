# Middleware

Bastion's middlewares are just stdlib `net/http` middleware handlers. There is nothing special about them, which means 
the router and all the tooling is designed to be compatible and friendly with any middleware in the community.

## Recovery 

Gracefully absorb panics and prints the stack trace. It log the panic (and a backtrace) in `os.Stdout` by default, this can be change
with the `RecoveryLoggerOutput` functional option and returns a HTTP 500 (Internal Server Error) status if possible.

### Options 
- `RecoveryLoggerOutput(w io.Writer)` set the logger output writer. Default `os.Stdout`.

```go
package main

import (
	"os"
	
	"github.com/ifreddyrondon/bastion/middleware"
)

func main() {
	// default
	middleware.Recovery()

	// with options
	middleware.Recovery(
		middleware.RecoveryLoggerOutput(os.Stdout),
	)
}
```

## InternalError
InternalError intercept responses to verify their status and handle the error. It gets the response code and 
if it's >= 500 handles the error with a default error message without disclosure internal information. 
The real error keeps logged.

### Options 
- `InternalErrMsg(s string)` set default error message to be sent. Default "looks like something went wrong".
- `InternalErrLoggerOutput(w io.Writer)` set the logger output writer. Default `os.Stdout`.

```go
package main

import (
	"errors"
	
	"github.com/ifreddyrondon/bastion/middleware"
)

func main() {
	// default
	middleware.InternalError()

	// with options
	middleware.InternalError(
		middleware.InternalErrMsg(errors.New("well, this is awkward")),
	)
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