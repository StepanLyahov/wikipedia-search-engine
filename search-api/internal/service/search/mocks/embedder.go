package mocks

import "context"

// Embedder is a test mock for search.Embedder.
type Embedder struct {
	EmbedFunc func(context.Context, string) ([]float32, error)
}

func (m *Embedder) Embed(ctx context.Context, text string) ([]float32, error) {
	return m.EmbedFunc(ctx, text)
}
