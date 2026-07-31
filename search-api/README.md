# Search API

Serves full-text search over indexed Wikipedia pages using Elasticsearch's built-in BM25
ranking. Configuration is supplied through environment variables; see the repository-level
`.env.example`.

## API

```text
GET /search?q=<query>&from=<offset>&size=<count>
```

`q` is required. `from` defaults to `0`, `size` defaults to `SEARCH_API_DEFAULT_SIZE` (10) and is
capped at `SEARCH_API_MAX_SIZE` (100).

Response:

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
* `internal/repository` — the search persistence port, its mock, and an Elasticsearch adapter
  (`multi_match` over `title^3`/`body^1` via the Search API);
* `internal/service/search` — the search use case;
* `internal/handler/http` — the `GET /search` HTTP handler, router and its mock port;
* `internal/app` — dependency composition and application lifecycle (HTTP server start/graceful
  shutdown);
* `cmd/search-api` — application startup only.

The service layer is independent of both the HTTP transport and Elasticsearch's HTTP API: both
are supplied through interfaces at startup. The project rules are in
`../GOLANG_SERVICE_PROMPT.md`.

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
```
