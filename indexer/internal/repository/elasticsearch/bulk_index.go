package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/wikipedia-search-engine/indexer/internal/domain"
	"github.com/wikipedia-search-engine/indexer/internal/repository"
)

type bulkAction struct {
	Index struct {
		Index string `json:"_index"`
		ID    string `json:"_id"`
	} `json:"index"`
}

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			Status int `json:"status"`
			Error  *struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

// BulkIndex indexes the supplied documents in a single Elasticsearch Bulk API request.
// Documents are indexed by their Postgres id, so re-indexing the same document overwrites
// it in place instead of failing, keeping repeated indexer runs error-free.
func (c *DocumentIndex) BulkIndex(ctx context.Context, documents []domain.Document) error {
	if len(documents) == 0 {
		return nil
	}

	body, err := bulkRequestBody(c.index, documents)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+c.index+"/_bulk", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bulk index: status %d: %s", resp.StatusCode, string(respBody))
	}

	var bulkResp bulkResponse
	if err := json.Unmarshal(respBody, &bulkResp); err != nil {
		return fmt.Errorf("decode bulk response: %w", err)
	}

	if bulkResp.Errors {
		return firstBulkError(bulkResp)
	}

	return nil
}

func bulkRequestBody(index string, documents []domain.Document) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	for _, document := range documents {
		var action bulkAction
		action.Index.Index = index
		action.Index.ID = strconv.FormatInt(document.ID, 10)

		actionLine, err := json.Marshal(action)
		if err != nil {
			return nil, err
		}
		buf.Write(actionLine)
		buf.WriteByte('\n')

		docLine, err := json.Marshal(document)
		if err != nil {
			return nil, err
		}
		buf.Write(docLine)
		buf.WriteByte('\n')
	}

	return &buf, nil
}

func firstBulkError(resp bulkResponse) error {
	for _, item := range resp.Items {
		if item.Index.Error != nil {
			return fmt.Errorf("bulk item failed: %s: %s", item.Index.Error.Type, item.Index.Error.Reason)
		}
	}

	return fmt.Errorf("bulk index reported errors")
}

var _ repository.DocumentIndex = (*DocumentIndex)(nil)
