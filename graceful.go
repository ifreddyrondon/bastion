package bastion

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

// graceful shutdown
func graceful(server *http.Server, ctx context.Context) {
	<-ctx.Done()
	log.Printf("[app:shutdown]")
	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("[app:gracefully_err] %v", err))
	} else {
		log.Printf("[app:gracefully] stopped")
	}
}
