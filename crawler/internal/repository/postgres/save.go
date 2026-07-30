package postgres

import (
	"context"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
)

// Save persists a page and ignores duplicate URLs.
func (r *PageRepository) Save(ctx context.Context, page domain.Page) error {
	query, args, err := r.builder.Insert("pages").
		Columns("url", "title", "html", "status", "created_at").
		Values(page.URL, page.Title, page.HTML, page.Status, page.CreatedAt).
		Suffix("ON CONFLICT (url) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
