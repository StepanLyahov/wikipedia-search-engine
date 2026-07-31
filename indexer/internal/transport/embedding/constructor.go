package embedding

import (
	"time"

	pb "github.com/wikipedia-search-engine/indexer/internal/pb/embedding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient dials the embedding-service and returns a ready-to-use gRPC client adapter.
func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, client: pb.NewEmbeddingServiceClient(conn), timeout: timeout}, nil
}
