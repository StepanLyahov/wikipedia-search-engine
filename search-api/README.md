# Search API

Serves full-text search (BM25) and semantic search (kNN over embeddings) across indexed
Wikipedia pages. Configuration is supplied through environment variables; see the
repository-level `.env.example`.

## API

### Keyword search

```text
GET /search?q=<query>&from=<offset>&size=<count>
```

`q` is required. `from` defaults to `0`, `size` defaults to `SEARCH_API_DEFAULT_SIZE` (10) and is
capped at `SEARCH_API_MAX_SIZE` (100). Runs a `multi_match` over `title^3`/`body^1`.

### Semantic search

```text
GET /semantic?q=<query>&k=<count>
```

`q` is required. `k` defaults to `SEARCH_API_DEFAULT_K` (10) and is capped at `SEARCH_API_MAX_K`
(100). Embeds the query via the `embedding-service` (gRPC) and runs a kNN search over the
`embedding` field, so it finds documents by meaning even without exact word matches. The number of
candidates Elasticsearch considers before ranking is controlled by `SEARCH_API_NUM_CANDIDATES`
(100) — higher improves quality at the cost of latency.

Both endpoints share the same response shape:

```json
{
  "hits": [
    {
      "title": "Distributed computing",
      "url": "https://en.wikipedia.org/wiki/Distributed_computing",
      "score": 13.42
    }
  ]
}
```

## Structure

* `config` — environment configuration;
* `internal/domain` — business entities;
* `internal/repository` — the search persistence port (`Search` for BM25, `VectorSearch` for kNN),
  its mock, and an Elasticsearch adapter;
* `internal/service/search` — the search and semantic-search use cases; depends on an `Embedder`
  port for query embeddings;
* `internal/transport/embedding` — the embedding-service gRPC client adapter;
* `internal/handler/http` — the `GET /search` and `GET /semantic` HTTP handlers, router and mock
  port;
* `internal/app` — dependency composition and application lifecycle (HTTP server start/graceful
  shutdown, embedding client teardown);
* `cmd/search-api` — application startup only.

The service layer is independent of the HTTP transport, Elasticsearch's HTTP API and the
embedding-service's gRPC API: all are supplied through interfaces at startup. The project rules
are in `../GOLANG_SERVICE_PROMPT.md`.

Run quality checks from this directory:

```sh
make test
make lint
```

`make lint` automatically installs the pinned linter version into `search-api/bin`; it does not
require adding anything to `PATH`.

Run the complete stack from the repository root:

```sh
docker compose up --build
```

Once the indexer has finished populating `wiki_pages`, query the API:

```sh
curl "http://localhost:8080/search?q=distributed%20systems"
curl "http://localhost:8080/semantic?q=how%20search%20engines%20work"
```

## Extending this

* **Adding an endpoint**: add a method to the `Handler` in its own file
  (`internal/handler/http/<name>.go`), register the route in `router.go`, and — if it needs a new
  capability from the service layer — extend the `Searcher` port (`searcher.go`) and its mock.
  Business logic belongs in `internal/service/search`, not in the handler; the handler only does
  request parsing/validation and response shaping.
* **`/search` and `/semantic` share one `Service`** (`internal/service/search`) and one
  Elasticsearch adapter (`internal/repository/elasticsearch`) rather than being split into two
  service packages — both are "the search use case" over the same index, just two query modes
  (`Search` = `multi_match`, `Semantic` = embed-then-`VectorSearch`). Keep new query modes on the
  same `Service`/`SearchIndex` port unless they need genuinely different dependencies.
* **`num_candidates`** (`SEARCH_API_NUM_CANDIDATES`) is a server-side quality/latency knob, not a
  request parameter — deliberately not exposed on `/semantic`, per the assignment's API shape. If
  you do want per-request control, it threads through `search.Config` → `Service.Semantic` →
  `repository.SearchIndex.VectorSearch`.
* **`VectorSearch`'s request body must set both `knn.k` and the top-level `size` to the same
  value** (`internal/repository/elasticsearch/vector_search.go`). Elasticsearch's `knn.k` only
  controls how many nearest neighbors get ranked internally; the top-level `size` is what
  actually bounds how many come back, and it defaults to 10 if omitted — so `k=20` would
  silently return only 10 hits without `Size: k` alongside it. If you add another field to the
  request body, keep this pairing in mind; `e2e/test_search_semantic.py`'s
  `test_semantic_k_limits_result_count` guards against regressing it.
* The embedding-service gRPC client (`internal/transport/embedding`) implements the
  `search.Embedder` port the same way `indexer`'s does — see
  [`../proto/README.md`](../proto/README.md) for the shared contract and regeneration commands if
  the `.proto` changes.
