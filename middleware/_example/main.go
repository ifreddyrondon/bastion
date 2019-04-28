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

func defaultMiddle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		m, snoop := middleware.WrapResponseWriter(w)
		next.ServeHTTP(snoop, r)
		fmt.Println(m.Code)
		fmt.Println(m.Bytes)
	}
	return http.HandlerFunc(fn)
}

func copyWriterMiddle(next http.Handler) http.Handler {
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

func hijackWriterMiddle(next http.Handler) http.Handler {
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
	app.With(defaultMiddle).Get("/", h)
	app.With(copyWriterMiddle).Get("/copy", h)
	app.With(hijackWriterMiddle).Get("/hijack", h)
	app.Serve()
}
