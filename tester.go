package bastion

import (
	"net/http"
	"testing"

	"gopkg.in/gavv/httpexpect.v1"
)

// Tester is an end-to-end testing helper for bastion handlers.
// It receives a reporter testing.T and http.Handler as params.
func Tester(t *testing.T, bastion *Bastion) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(bastion.r),
			Jar:       httpexpect.NewJar(),
		},

		// use fatal failures
		Reporter: httpexpect.NewAssertReporter(t),

		// use verbose logging
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
