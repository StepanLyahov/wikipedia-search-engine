package app

// Close releases application resources.
func (a *App) Close() {
	a.logger.Info("search-api service stopping")
	a.logger.Sync()
}
