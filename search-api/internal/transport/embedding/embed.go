package embedding

import (
	"context"

	pb "github.com/wikipedia-search-engine/search-api/internal/pb/embedding"
	"github.com/wikipedia-search-engine/search-api/internal/service/search"
)

// Embed calls the embedding-service to turn text into a dense vector.
func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Embed(ctx, &pb.EmbedRequest{Text: text})
	if err != nil {
		return nil, err
	}

	return resp.GetVector(), nil
}

var _ search.Embedder = (*Client)(nil)
