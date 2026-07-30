package crawler

import (
	"strings"

	"golang.org/x/net/html"
)

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
			if attr.Key == "href" && strings.HasPrefix(attr.Val, "/wiki/") && !strings.ContainsAny(attr.Val, "#:?") {
				links = append(links, attr.Val)
				break
			}
		}
	}
}
