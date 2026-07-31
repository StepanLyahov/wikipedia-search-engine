// Command search-api starts the Wikipedia BM25 search HTTP API.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wikipedia-search-engine/search-api/internal/app"
)

func main() {
	application, err := app.New(context.Background())
	if err != nil {
		log.Fatalf("create application: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	runErr := application.Run(ctx)
	application.Close()
	stop()
	if runErr != nil {
		log.Printf("search-api failed: %v", runErr)
		os.Exit(1)
	}
}
