// Package domain contains crawler business entities.
package domain

import "time"

// Page is the domain representation of a downloaded Wikipedia article.
type Page struct {
	URL       string
	Title     string
	HTML      string
	Status    int
	CreatedAt time.Time
}
