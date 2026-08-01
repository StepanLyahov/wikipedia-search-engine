package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
)

// embeddingField is the dense_vector field populated by the indexer's embedding-service calls.
const embeddingField = "embedding"

type vectorSearchRequestBody struct {
	// Size must match k: the knn clause's k only controls how many nearest neighbors
	// Elasticsearch ranks internally, not how many of them come back in the response.
	Size int      `json:"size"`
	KNN  knnQuery `json:"knn"`
}

type knnQuery struct {
	Field         string    `json:"field"`
	QueryVector   []float32 `json:"query_vector"`
	K             int       `json:"k"`
	NumCandidates int       `json:"num_candidates"`
}

// VectorSearch runs a kNN query against the embedding field and returns hits ranked by
// Elasticsearch's cosine similarity score, which is the query's default sort order.
func (c *SearchIndex) VectorSearch(ctx context.Context, vector []float32, k, numCandidates int) ([]domain.Hit, error) {
	reqBody := vectorSearchRequestBody{
		Size: k,
		KNN:  knnQuery{Field: embeddingField, QueryVector: vector, K: k, NumCandidates: numCandidates},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal vector search request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+c.index+"/_search", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build vector search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute vector search request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read vector search response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vector search: status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed searchResponseBody
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode vector search response: %w", err)
	}

	hits := make([]domain.Hit, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		hits = append(hits, domain.Hit{Title: hit.Source.Title, URL: hit.Source.URL, Score: hit.Score})
	}

	return hits, nil
}
