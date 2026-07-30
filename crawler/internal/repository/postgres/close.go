package postgres

// Close releases the repository connection pool.
func (r *PageRepository) Close() { r.pool.Close() }
