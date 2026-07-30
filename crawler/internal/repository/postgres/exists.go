package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

// Exists reports whether a page with the supplied URL is already stored.
func (r *PageRepository) Exists(ctx context.Context, url string) (bool, error) {
	query, args, err := r.builder.Select("1").From("pages").Where(sq.Eq{"url": url}).Limit(1).ToSql()
	if err != nil {
		return false, err
	}
	var found int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}
