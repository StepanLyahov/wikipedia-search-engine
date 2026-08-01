# End-to-end tests

Exercises the whole pipeline against the real `docker-compose.yml` stack — no mocks:

```
crawler -> Postgres -> indexer (+ embedding-service, gRPC) -> Elasticsearch -> search-api
```

The suite brings the stack up itself, waits for the crawl and index one-shot jobs to finish,
then asserts against the real Postgres rows, the real Elasticsearch documents/mapping, and the
real `/search` and `/semantic` HTTP endpoints of `search-api`.

## What's covered

* `test_crawler.py` — the crawler saved the seed page (and only real `en.wikipedia.org/wiki/*`
  URLs), with a title, HTML body and recent timestamp.
* `test_indexer.py` — the `wiki_pages` index exists with the expected mapping (including the
  384-dim `dense_vector` embedding field with cosine similarity), every crawled page was
  indexed, and every document carries a real (non-degenerate) 384-dim embedding.
* `test_search_fulltext.py` — `GET /search`: title-boosted ranking, matching on the body field
  (not just the title), pagination, empty results for no match, descending score order.
* `test_search_semantic.py` — `GET /semantic`: query embedding + kNN search finds the
  topically-relevant page for queries that don't literally quote it, `k` limits result count,
  descending score order.

## Running

```sh
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
pytest
```

This runs `docker compose down -v` (clean slate) then `docker compose up -d --build` from the
repository root, waits for the pipeline to finish, runs the tests, then tears the stack down
again. The first run builds every image (the `embedding-service` image bakes in the
`all-MiniLM-L6-v2` model weights, which takes a few minutes); subsequent runs reuse Docker's
layer cache and are much faster.

### Useful environment variables

| Variable | Default | Purpose |
|---|---|---|
| `E2E_SKIP_COMPOSE` | unset | Skip `docker compose up/down` entirely and test against a stack you started yourself (fast iteration on the tests). |
| `E2E_KEEP_STACK` | unset | Leave the stack running after the test session (skip the final `docker compose down -v`). |
| `E2E_CRAWLER_TIMEOUT` | `180` | Seconds to wait for the crawler container to exit. |
| `E2E_INDEXER_TIMEOUT` | `180` | Seconds to wait for the indexer container to exit. |
| `E2E_SERVICE_HEALTHY_TIMEOUT` | `180` | Seconds to wait for postgres/elasticsearch/embedding-service to report healthy. |
| `E2E_HTTP_READY_TIMEOUT` | `60` | Seconds to wait for `search-api` to accept HTTP requests. |

Any `docker compose` variable from the repository-root `.env`/`.env.example` (e.g.
`CRAWLER_SEED_URL`, `SEARCH_API_PORT`) also applies here, since the fixtures read the *resolved*
compose config (`docker compose config`) rather than hardcoding defaults.

## Notes on test design

* Fixtures read ports, the DB DSN, the crawl seed URL and the ES index name back out of
  `docker compose config`, so the tests always match whatever actually drives the stack instead
  of duplicating `.env` defaults.
* With the default seed, the crawler now discovers around a hundred real pages (Kibana, Apache
  Lucene, Full-text search, "Search engine (computing)", ...), not just the seed — several of
  them genuinely close in topic to the "Elasticsearch" article. The semantic-search tests are
  written for that reality: instead of asserting the Elasticsearch page ranks exactly first for
  an on-topic query (fragile once several legitimately-related pages exist to compete for rank
  1), they assert it's *findable* within a generous top-k, and that a clearly unrelated query
  ranks something else first. Since this tests against live Wikipedia content, a future article
  edit could in principle shift a result in or out of that window; that's an accepted trade-off
  of testing the real pipeline instead of canned fixtures.
* `test_search_fulltext.py`'s body-field test asks Elasticsearch's term vectors API which
  indexed term is rarest across the corpus, rather than re-tokenizing the raw text in Python —
  Lucene's standard analyzer doesn't split everything the way a naive word regex would (e.g. a
  citation URL like `n3.nabble.com` is indexed as one token, `nabble.com`, not `nabble` + `com`),
  so guessing tokenization can pick a "distinctive" word that was never actually indexed as such.
  Asking Elasticsearch directly sidesteps that class of mismatch entirely.
* `test_semantic_k_limits_result_count` deliberately checks a `k` above 10: Elasticsearch's kNN
  `k` only controls how many candidates get ranked internally, not how many come back — the
  top-level `size` does that, and it's easy to set one and forget the other (this test exists
  because that was a real bug here, fixed in `search-api/internal/repository/elasticsearch/vector_search.go`).
