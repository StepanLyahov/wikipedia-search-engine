package crawler

import (
	"strings"

	"golang.org/x/net/html"
)

const wikipediaOrigin = "https://en.wikipedia.org"

// wikiLinks returns the same-site article paths (e.g. "/wiki/Lucene") linked from source,
// accepting both page-relative hrefs and the absolute https://en.wikipedia.org/wiki/... hrefs
// most in-body Wikipedia links use today.
func wikiLinks(source string) []string {
	tokenizer := html.NewTokenizer(strings.NewReader(source))
	links := make([]string, 0)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		if tokenType != html.StartTagToken && tokenType != html.SelfClosingTagToken {
			continue
		}
		token := tokenizer.Token()
		if token.Data != "a" {
			continue
		}
		for _, attr := range token.Attr {
			if attr.Key != "href" {
				continue
			}
			if path, ok := wikiArticlePath(attr.Val); ok {
				links = append(links, path)
			}

			break
		}
	}
}

// wikiArticlePath extracts the "/wiki/..." path from an anchor href, if it points to a plain
// en.wikipedia.org article. It accepts both a page-relative href ("/wiki/Lucene") and an
// absolute one ("https://en.wikipedia.org/wiki/Lucene"), normalizing either to the relative
// form. Namespaced pages (Special:, Category:, File:, Help:, ...), URLs with an anchor or query
// string, and links to any other host are excluded.
func wikiArticlePath(href string) (string, bool) {
	path := strings.TrimPrefix(href, wikipediaOrigin)

	if !strings.HasPrefix(path, "/wiki/") || strings.ContainsAny(path, "#:?") {
		return "", false
	}

	return path, true
}
