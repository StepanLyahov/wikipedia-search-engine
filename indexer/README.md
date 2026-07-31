# Indexer

Reads crawled pages from PostgreSQL, extracts a clean title and body from their HTML, and
bulk-indexes them into Elasticsearch. Configuration is supplied through environment variables;
see the repository-level `.env.example`.

## Structure

* `config` — environment configuration;
* `internal/domain` — business entities;
* `internal/repository` — persistence interfaces, mocks, a PostgreSQL page reader (Squirrel) and
  an Elasticsearch document index adapter (Bulk API over HTTP);
* `internal/service/indexer` — indexing use case, HTML title/body extraction and its interfaces;
* `internal/app` — dependency composition and application lifecycle;
* `cmd/indexer` — application startup only.

The service layer is independent of `pgx` and Elasticsearch's HTTP API: both are supplied
through interfaces at startup. The project rules are in `../GOLANG_SERVICE_PROMPT.md`.

Indexing is idempotent: documents are indexed by their PostgreSQL id, and index creation
tolerates an already-existing `wiki_pages` index, so re-running the indexer never fails.

Run quality checks from this directory:

```sh
make test
make lint
```

`make lint` automatically installs the pinned linter version into `indexer/bin`; it does not
require adding anything to `PATH`.

Run the complete stack from the repository root:

```sh
docker compose up --build
```

The indexer is a one-shot container: it runs once after the crawler finishes and exits.
Inspect the index with:

```sh
curl http://localhost:9200/wiki_pages/_search?size=1&pretty
```
