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

## Extending this

* **HTML extraction** (`internal/service/indexer/html.go`) uses `golang.org/x/net/html`'s
  low-level `Tokenizer`, not the full DOM `Parser`. That tokenizer hardcodes a handful of tags as
  *raw text* (`script`, `style`, `noscript`, `textarea`, `title`, ...) — their content, including
  any nested markup, comes through as a single literal `TextToken` rather than being parsed as
  child elements. `extractBody`'s `isSkippable` allowlist must cover every raw-text tag whose
  content shouldn't end up in the indexed body (a real Wikipedia `<noscript><img
  src="...campaign-tracking-pixel..."></noscript>` block was leaking into `body` until `noscript`
  was added — see `TestExtractBodySkipsNoscriptContent`). If indexed bodies ever contain stray
  `<...>` text again, this is the first place to look.
* **Embeddings**: `internal/service/indexer/embed.go` calls the `Embedder` port once per
  document (`title + " " + body`) and silently *skips* (logs + drops) any document whose call
  fails or whose vector length doesn't match `Config.EmbeddingDims` — it does not fail the whole
  run. If you need embedding failures to be fatal instead, change `embedDocuments`.
* **Changing the embedding model/dimension**: update `INDEXER_EMBEDDING_DIMS` (must match
  `embedding-service`'s model output size) — `ensure_index.go` renders the ES mapping's `dims`
  from this value, so the index mapping and the model stay in sync automatically as long as both
  are pointed at the same env var value.
* The ES document shape (`internal/domain/document.go`) is the contract with `search-api` — any
  field added here should be reflected in the mapping in `ensure_index.go` and, if it should be
  searchable/returned, in `search-api`'s `domain.Hit` and Elasticsearch adapters too.
