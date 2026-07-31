package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// indexMappingTemplate is the wiki_pages mapping required by the indexing use case.
// %d is the embedding vector's dimensionality, fixed by the embedding-service's model.
const indexMappingTemplate = `{
	"mappings": {
		"properties": {
			"id": {"type": "long"},
			"url": {"type": "keyword"},
			"title": {"type": "text"},
			"body": {"type": "text"},
			"embedding": {"type": "dense_vector", "dims": %d, "index": true, "similarity": "cosine"}
		}
	}
}`

// resourceAlreadyExists is the Elasticsearch error type returned when the index already exists,
// which makes index creation idempotent across indexer re-runs.
const resourceAlreadyExists = "resource_already_exists_exception"

type indexErrorResponse struct {
	Error struct {
		Type string `json:"type"`
	} `json:"error"`
}

// EnsureIndex creates the document index if it does not already exist.
func (c *DocumentIndex) EnsureIndex(ctx context.Context) error {
	mapping := fmt.Sprintf(indexMappingTemplate, c.embeddingDims)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/"+c.index, bytes.NewBufferString(mapping))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return nil
	}

	var errResp indexErrorResponse
	if json.Unmarshal(body, &errResp) == nil && errResp.Error.Type == resourceAlreadyExists {
		return nil
	}

	return fmt.Errorf("create index: status %d: %s", resp.StatusCode, string(body))
}
