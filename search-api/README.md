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
