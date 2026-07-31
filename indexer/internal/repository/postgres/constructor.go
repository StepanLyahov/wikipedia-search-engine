package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPageRepository creates a PostgreSQL page repository.
func NewPageRepository(ctx context.Context, databaseURL string) (*PageRepository, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &PageRepository{pool: pool, builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}, nil
}
