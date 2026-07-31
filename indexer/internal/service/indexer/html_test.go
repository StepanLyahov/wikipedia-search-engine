package indexer

import "testing"

func TestExtractTitle(t *testing.T) {
	tests := map[string]struct {
		html string
		want string
	}{
		"simple heading": {
			html: `<html><body><h1 id="firstHeading">Elasticsearch</h1></body></html>`,
			want: "Elasticsearch",
		},
		"nested span": {
			html: `<h1 id="firstHeading"><span class="mw-page-title-main">Elasticsearch</span></h1>`,
			want: "Elasticsearch",
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
