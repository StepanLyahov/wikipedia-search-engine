package mocks

import "context"

func (m *PageRepository) Exists(ctx context.Context, url string) (bool, error) {
	return m.ExistsFunc(ctx, url)
}
