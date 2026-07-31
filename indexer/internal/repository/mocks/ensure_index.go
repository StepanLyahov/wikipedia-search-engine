package mocks

import "context"

func (m *DocumentIndex) EnsureIndex(ctx context.Context) error { return m.EnsureIndexFunc(ctx) }
