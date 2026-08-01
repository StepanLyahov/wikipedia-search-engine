// Package embedding contains the embedding-service gRPC client adapter.
package embedding

import (
	"time"

	pb "github.com/wikipedia-search-engine/search-api/internal/pb/embedding"
	"google.golang.org/grpc"
)

// Client adapts the embedding-service gRPC API to the search Embedder port.
type Client struct {
	conn    *grpc.ClientConn
	client  pb.EmbeddingServiceClient
	timeout time.Duration
}
