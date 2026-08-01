package search

import "context"

// Embedder is the outbound gRPC dependency that turns a query into a dense vector.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}
