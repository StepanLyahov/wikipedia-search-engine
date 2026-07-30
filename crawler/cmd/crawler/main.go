// Command crawler starts the Wikipedia crawler application.
package main

import (
	"context"
	"log"
	"os"

	"github.com/wikipedia-search-engine/crawler/internal/app"
)

func main() {
	application, err := app.New(context.Background())
	if err != nil {
		log.Fatalf("create application: %v", err)
	}
	runErr := application.Run(context.Background())
	application.Close()
	if runErr != nil {
		log.Printf("crawl failed: %v", runErr)
		os.Exit(1)
	}
}
