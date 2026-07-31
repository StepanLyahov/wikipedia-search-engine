package indexer

import "context"

// Indexer is the public use-case contract for the indexing service.
type Indexer interface {
	Run(ctx context.Context) error
}

var _ Indexer = (*Service)(nil)
