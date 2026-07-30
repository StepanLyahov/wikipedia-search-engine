package app

// Close releases application resources.
func (a *App) Close() {
	a.logger.Info("crawler service stopping")

	a.pages.Close()
	a.logger.Sync()
}
