package wikipedia

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wikipedia-search-engine/crawler/internal/domain"
	service "github.com/wikipedia-search-engine/crawler/internal/service/crawler"
	"golang.org/x/net/html"
)

const maxHTMLSize = 10 << 20

// Fetch downloads a Wikipedia page and converts it to a domain entity.
func (c *Client) Fetch(ctx context.Context, rawURL string) (domain.Page, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return domain.Page{}, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	response, err := c.client.Do(req)
	if err != nil {
		return domain.Page{}, err
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxHTMLSize))
	if err != nil {
		return domain.Page{}, err
	}
	return domain.Page{URL: rawURL, Title: pageTitle(string(body)), HTML: string(body), Status: response.StatusCode, CreatedAt: time.Now().UTC()}, nil
}

func pageTitle(source string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(source))
	inTitle := false
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return ""
		case html.StartTagToken:
			inTitle = tokenizer.Token().Data == "title"
		case html.TextToken:
			if inTitle {
				return strings.TrimSpace(tokenizer.Token().Data)
			}
		case html.EndTagToken:
			if tokenizer.Token().Data == "title" {
				inTitle = false
			}
		case html.SelfClosingTagToken, html.CommentToken, html.DoctypeToken:
			continue
		default:
			return ""
		}
	}
}

var _ service.Fetcher = (*Client)(nil)
