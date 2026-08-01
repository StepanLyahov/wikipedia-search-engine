package crawler

import "testing"

func TestWikiLinks(t *testing.T) {
	source := `<html><body>
		<a href="/wiki/Lucene">relative link</a>
		<a href="https://en.wikipedia.org/wiki/Apache_Lucene">absolute link</a>
		<a href="https://en.wikipedia.org/wiki/Software_categories#Categorization_approaches">absolute link with anchor</a>
		<a href="https://en.wikipedia.org/wiki/Special:Random">absolute namespaced link</a>
		<a href="/wiki/Category:Software">relative namespaced link</a>
		<a href="/w/index.php?title=Elasticsearch&action=edit">query string link</a>
		<a href="https://en.m.wikipedia.org/wiki/Lucene">different host</a>
		<a href="https://github.com/elastic/elasticsearch">external site</a>
	</body></html>`

	want := []string{"/wiki/Lucene", "/wiki/Apache_Lucene"}

	got := wikiLinks(source)
	if len(got) != len(want) {
		t.Fatalf("wikiLinks() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("wikiLinks() = %v, want %v", got, want)
		}
	}
}

const wantElasticsearchPath = "/wiki/Elasticsearch"

func TestWikiArticlePath(t *testing.T) {
	tests := map[string]struct {
		href     string
		wantPath string
		wantOK   bool
	}{
		"relative article": {
			href: wantElasticsearchPath, wantPath: wantElasticsearchPath, wantOK: true,
		},
		"absolute article": {
			href: "https://en.wikipedia.org" + wantElasticsearchPath, wantPath: wantElasticsearchPath, wantOK: true,
		},
		"absolute article with anchor": {
			href: "https://en.wikipedia.org/wiki/Elasticsearch#History", wantOK: false,
		},
		"relative namespaced page": {
			href: "/wiki/Special:Random", wantOK: false,
		},
		"absolute namespaced page": {
			href: "https://en.wikipedia.org/wiki/Category:Databases", wantOK: false,
		},
		"query string": {
			href: "/w/index.php?title=Elasticsearch&action=history", wantOK: false,
		},
		"different wikipedia host": {
			href: "https://en.m.wikipedia.org/wiki/Elasticsearch", wantOK: false,
		},
		"external site": {
			href: "https://github.com/elastic/elasticsearch", wantOK: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			path, ok := wikiArticlePath(tt.href)
			if ok != tt.wantOK {
				t.Fatalf("wikiArticlePath(%q) ok = %v, want %v", tt.href, ok, tt.wantOK)
			}
			if ok && path != tt.wantPath {
				t.Fatalf("wikiArticlePath(%q) = %q, want %q", tt.href, path, tt.wantPath)
			}
		})
	}
}
