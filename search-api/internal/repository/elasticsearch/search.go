package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wikipedia-search-engine/search-api/internal/domain"
	"github.com/wikipedia-search-engine/search-api/internal/repository"
)

// searchFields weights title over body so title matches rank higher, per the multi_match query.
var searchFields = []string{"title^3", "body^1"}

type searchRequestBody struct {
	From  int         `json:"from"`
	Size  int         `json:"size"`
	Query searchQuery `json:"query"`
}

type searchQuery struct {
	MultiMatch multiMatchQuery `json:"multi_match"`
}

type multiMatchQuery struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}

type searchResponseBody struct {
	Hits struct {
		Hits []struct {
			Score  float64 `json:"_score"`
			Source struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Search runs a title/body multi_match query against the index and returns hits ranked by
// Elasticsearch's BM25 score, which is the query's default sort order.
func (c *SearchIndex) Search(ctx context.Context, query string, from, size int) ([]domain.Hit, error) {
	reqBody := searchRequestBody{
		From: from,
		Size: size,
		Query: searchQuery{
			MultiMatch: multiMatchQuery{Query: query, Fields: searchFields},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal search request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+c.index+"/_search", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute search request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read search response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search: status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed searchResponseBody
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	hits := make([]domain.Hit, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		hits = append(hits, domain.Hit{Title: hit.Source.Title, URL: hit.Source.URL, Score: hit.Score})
	}

	return hits, nil
}

var _ repository.SearchIndex = (*SearchIndex)(nil)
