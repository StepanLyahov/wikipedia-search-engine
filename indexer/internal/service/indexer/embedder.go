package indexer

import "context"

// Embedder is the outbound gRPC dependency that turns text into a dense vector.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}
