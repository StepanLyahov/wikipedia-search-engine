// Package postgres contains PostgreSQL repository adapters.
package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PageRepository is a PostgreSQL implementation of repository.PageRepository.
type PageRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}
