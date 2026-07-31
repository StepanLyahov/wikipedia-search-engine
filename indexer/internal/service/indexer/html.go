package indexer

import (
	"strings"

	"golang.org/x/net/html"
)

const firstHeadingID = "firstHeading"

// extractTitle returns the text content of the element with id="firstHeading",
// e.g. the <h1 id="firstHeading">Elasticsearch</h1> heading rendered by Wikipedia.
func extractTitle(source string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(source))
	depth := 0
	var text strings.Builder

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return strings.TrimSpace(text.String())
		}

		token := tokenizer.Token()

		switch tokenType {
		case html.StartTagToken, html.SelfClosingTagToken:
			if depth == 0 {
				if hasID(token, firstHeadingID) {
					depth = 1
				}

				continue
			}
			if tokenType == html.StartTagToken {
				depth++
			}
		case html.TextToken:
			if depth > 0 {
				text.WriteString(token.Data)
			}
		case html.EndTagToken:
			if depth == 0 {
				continue
			}

			depth--
			if depth == 0 {
				return strings.TrimSpace(text.String())
			}
		case html.CommentToken, html.DoctypeToken:
			continue
		default:
			continue
		}
	}
}

// extractBody removes tags, scripts and styles from source and joins the remaining
// text into a single whitespace-normalised line.
func extractBody(source string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(source))
	skipDepth := 0
	var text strings.Builder

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return strings.Join(strings.Fields(text.String()), " ")
		}

		token := tokenizer.Token()

		switch tokenType {
		case html.StartTagToken:
			if isSkippable(token.Data) {
				skipDepth++
			}
		case html.EndTagToken:
			if isSkippable(token.Data) && skipDepth > 0 {
				skipDepth--
			}
		case html.TextToken:
			if skipDepth == 0 {
				text.WriteString(token.Data)
				text.WriteString(" ")
			}
		case html.SelfClosingTagToken, html.CommentToken, html.DoctypeToken:
			continue
		default:
			continue
		}
	}
}

func isSkippable(tag string) bool {
	return tag == "script" || tag == "style"
}

func hasID(token html.Token, id string) bool {
	for _, attr := range token.Attr {
		if attr.Key == "id" && attr.Val == id {
			return true
		}
	}

	return false
}
