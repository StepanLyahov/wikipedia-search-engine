package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/wikipedia-search-engine/indexer/internal/domain"
)

// statusOK marks pages that were downloaded successfully and carry indexable content.
const statusOK = 200

// FetchAll returns every successfully crawled page.
func (r *PageRepository) FetchAll(ctx context.Context) ([]domain.Page, error) {
	query, args, err := r.builder.
		Select("id", "url", "html").
		From("pages").
		Where(sq.Eq{"status": statusOK}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pages := make([]domain.Page, 0)
	for rows.Next() {
		var page domain.Page
		if err := rows.Scan(&page.ID, &page.URL, &page.HTML); err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, rows.Err()
}
