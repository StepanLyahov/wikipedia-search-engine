// Package embedding contains the embedding-service gRPC client adapter.
package embedding

import (
	"time"

	pb "github.com/wikipedia-search-engine/indexer/internal/pb/embedding"
	"google.golang.org/grpc"
)

// Client adapts the embedding-service gRPC API to the indexer Embedder port.
type Client struct {
	conn    *grpc.ClientConn
	client  pb.EmbeddingServiceClient
	timeout time.Duration
}
