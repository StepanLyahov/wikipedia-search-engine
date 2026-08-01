# Crawler

Downloads Wikipedia pages breadth-first and persists them in PostgreSQL. Configuration is supplied through environment variables; see the repository-level `.env.example`.

## Structure

* `config` — environment configuration;
* `internal/domain` — business entities;
* `internal/repository` — persistence interfaces, mocks and a PostgreSQL adapter using Squirrel;
* `internal/service/crawler` — crawling use case and its interfaces;
* `internal/transport/wikipedia` — HTTP adapter for Wikipedia;
* `internal/app` — dependency composition and application lifecycle;
* `cmd/crawler` — application startup only.

The service layer is independent of `pgx` and `net/http`: both are supplied through interfaces at startup. The project rules are in `../GOLANG_SERVICE_PROMPT.md`.

Run quality checks from this directory:

```sh
make test
make lint
```

`make lint` automatically installs the pinned linter version into `crawler/bin`; it does not require adding anything to `PATH`.

Run the complete stack from the repository root:

```sh
docker compose up --build
```

The crawler is a one-shot container: it exits after it reaches the queue limit or `CRAWLER_MAX_PAGES`. Inspect stored pages with:

```sh
docker compose exec postgres psql -U wiki -d wiki -c 'SELECT id, title, status, url FROM pages;'
```

## Extending this

* **Link discovery** (`internal/service/crawler/html.go`, `wikiLinks`/`wikiArticlePath`) accepts
  both page-relative anchors (`href="/wiki/..."`) and the absolute
  `href="https://en.wikipedia.org/wiki/..."` links most in-body Wikipedia content actually uses
  today — an earlier version only matched the relative form, so it silently crawled just one
  extra page (`/wiki/Main_Page`, linked from the page chrome, not the article body) regardless of
  `CRAWLER_MAX_DEPTH`/`CRAWLER_MAX_PAGES`. `wikiArticlePath` is the single place that decides
  whether an `href` is a followable article link; it excludes `Special:`, `Category:`, `Help:`,
  `Portal:`, `File:`, etc. namespaced pages, `#`/`?` fragments/queries, and any other host. If you
  add another accepted link form (e.g. a different Wikipedia language edition), extend that
  function and its table-driven test in `html_test.go` rather than `wikiLinks` itself.
* **Idempotency** relies entirely on the `pages.url` `UNIQUE` constraint plus the
  `Exists`-before-`Fetch` check in `Crawl` (`internal/service/crawler/crawl.go`) — a page that's
  already stored is skipped without inspecting its links again, so re-running the crawler never
  expands the frontier further than the first run did. If you change link discovery, keep this
  property in mind.
* New per-page fields (e.g. capturing HTTP headers, redirect chains) belong on `domain.Page`
  (`internal/domain/page.go`), threaded through `transport/wikipedia.Client.Fetch` and persisted
  via `repository/postgres` — add a migration in [`../postgres/README.md`](../postgres/README.md)
  for any new column.
