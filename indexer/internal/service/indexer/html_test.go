package indexer

import "testing"

const wantElasticsearch = "Elasticsearch"

func TestExtractTitle(t *testing.T) {
	tests := map[string]struct {
		html string
		want string
	}{
		"simple heading": {
			html: `<html><body><h1 id="firstHeading">Elasticsearch</h1></body></html>`,
			want: wantElasticsearch,
		},
		"nested span": {
			html: `<h1 id="firstHeading"><span class="mw-page-title-main">Elasticsearch</span></h1>`,
			want: wantElasticsearch,
		},
		"missing heading": {
			html: `<html><body><p>No heading here</p></body></html>`,
			want: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := extractTitle(tt.html); got != tt.want {
				t.Fatalf("extractTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractBody(t *testing.T) {
	source := `
		<html>
		<head><style>.hidden { display: none; }</style></head>
		<body>
			<script>console.log("ignored")</script>
			<p>Elasticsearch is a search engine
			based on Apache Lucene.</p>
		</body>
		</html>
	`

	want := "Elasticsearch is a search engine based on Apache Lucene."
	if got := extractBody(source); got != want {
		t.Fatalf("extractBody() = %q, want %q", got, want)
	}
}

func TestExtractBodySkipsNoscriptContent(t *testing.T) {
	// golang.org/x/net/html tokenizes <noscript> as raw text, so a tracking pixel like this
	// (seen on real Wikipedia pages) would otherwise leak into the body as literal markup.
	source := `<p>Elasticsearch is a search engine.</p>
		<noscript><img src="https://en.wikipedia.org/wiki/Special:CentralAutoLogin/start" alt="" width="1" height="1"></noscript>`

	want := "Elasticsearch is a search engine."
	if got := extractBody(source); got != want {
		t.Fatalf("extractBody() = %q, want %q", got, want)
	}
}
