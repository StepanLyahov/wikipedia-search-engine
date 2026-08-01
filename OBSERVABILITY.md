# Observability

How to look under the hood of the running stack (`docker compose up`): browse the crawled pages
in Postgres with a GUI client, inspect Elasticsearch's indices/documents/mapping, and tail every
service's logs live.

## 1. PostgreSQL — via DBeaver

Postgres is published on the host, so any desktop client works, not just DBeaver — same
parameters either way.

**Connection parameters** (defaults from `.env.example`; adjust the port if you overrode
`POSTGRES_PORT`):

| Field | Value |
|---|---|
| Host | `localhost` |
| Port | `5432` |
| Database | `wiki` |
| Username | `wiki` |
| Password | `wiki` |
| JDBC URL | `jdbc:postgresql://localhost:5432/wiki` |

**In DBeaver**: `Database` → `New Database Connection` → `PostgreSQL` → fill in the fields above
→ `Test Connection...` → `Finish`.

**Example queries**, once connected (the `pages` table is the crawler's only output — see
[`postgres/README.md`](postgres/README.md) for the schema):

```sql
-- Everything the crawler saved
SELECT id, url, title, status, created_at FROM pages ORDER BY id;

-- How many pages were crawled
SELECT count(*) FROM pages;

-- Pages that failed to fetch (non-2xx HTTP status)
SELECT id, url, status FROM pages WHERE status < 200 OR status >= 300;

-- Most recently crawled
SELECT url, title, created_at FROM pages ORDER BY created_at DESC LIMIT 10;

-- Find a specific page
SELECT * FROM pages WHERE url = 'https://en.wikipedia.org/wiki/Elasticsearch';
```

## 2. Elasticsearch

Elasticsearch itself doesn't ship a web UI (that's what Kibana is for; this stack skips it to
stay light). Two ways to look inside the `wiki_pages` index:

### a) curl / REST API — no extra setup

```sh
# List all indices, with health/size/doc count
curl "http://localhost:9200/_cat/indices?v"

# Cluster health
curl "http://localhost:9200/_cluster/health?pretty"

# Mapping: confirms title/body are text, url is keyword, embedding is a
# dense_vector with dims=384 and similarity=cosine
curl "http://localhost:9200/wiki_pages/_mapping?pretty"

# Document count
curl "http://localhost:9200/wiki_pages/_count"

# First 5 documents as stored (includes the embedding vector — long output)
curl "http://localhost:9200/wiki_pages/_search?size=5&pretty"

# One document by id
curl "http://localhost:9200/wiki_pages/_doc/1?pretty"

# The same BM25 query GET /search runs under the hood (title^3 / body^1)
curl -X POST "http://localhost:9200/wiki_pages/_search" \
  -H 'Content-Type: application/json' \
  -d '{"query": {"multi_match": {"query": "elasticsearch", "fields": ["title^3", "body"]}}}'

# A kNN query like GET /semantic runs under the hood — needs a real 384-float query_vector;
# easiest is to copy one straight out of an indexed document:
curl -s "http://localhost:9200/wiki_pages/_search?size=1" \
  | python3 -c "import json,sys; print(json.load(sys.stdin)['hits']['hits'][0]['_source']['embedding'])"
# then paste that array into:
curl -X POST "http://localhost:9200/wiki_pages/_search" \
  -H 'Content-Type: application/json' \
  -d '{"knn": {"field": "embedding", "query_vector": [ ...paste here... ], "k": 5, "num_candidates": 50}}'

# Delete the index (e.g. to force the indexer to rebuild it from scratch on its next run)
curl -X DELETE "http://localhost:9200/wiki_pages"
```

Add `?pretty` to any of these for indented JSON, or pipe through `python3 -m json.tool` / `jq`.

### b) A browser UI — optional, on demand

If you'd rather click around than curl, run [Elasticvue](https://elasticvue.com/): a static
single-page app with no backend of its own — it talks straight to Elasticsearch's REST API from
your browser. It's a one-off container, not part of `docker-compose.yml`, so it can't interfere
with the rest of the stack:

```sh
docker run -d --name elasticvue --rm -p 8085:8080 cars10/elasticvue
```

Open <http://localhost:8085>, add a cluster with URL `http://localhost:9200` (the same host port
Elasticsearch already publishes), and you get a full index/mapping/document/query browser.

`docker-compose.yml` already enables CORS on Elasticsearch (`http.cors.enabled`,
`http.cors.allow-origin: "*"`) specifically so a browser tool like this can call the API directly
— if Elasticvue instead shows a "You have to setup CORS" prompt, it means you're pointing it at an
Elasticsearch that predates that setting; run `docker compose up -d --force-recreate
elasticsearch` (or `docker compose down -v && docker compose up --build`) to pick it up.

```sh
docker stop elasticvue   # when you're done with it
```

## 3. Live logs per service

`docker compose logs -f <service>` streams a container's stdout/stderr in real time (`-f` =
follow, like `tail -f`). Run these from the repository root while the stack is up:

```sh
docker compose logs -f search-api          # long-running: keeps streaming as requests come in
docker compose logs -f embedding-service
docker compose logs -f elasticsearch
docker compose logs -f postgres

docker compose logs -f crawler             # one-shot: exits once the crawl finishes, so this
docker compose logs -f indexer             # just replays what already happened, then hangs open
```

Useful variations:

```sh
# Every service at once, interleaved (prefixed with the service name)
docker compose logs -f

# Last 200 lines, then keep following
docker compose logs -f --tail=200 search-api

# Several services in one stream
docker compose logs -f indexer search-api

# Only what's happened in the last 10 minutes, then keep following
docker compose logs -f --since=10m search-api
```

`Ctrl+C` stops *following* — it does not stop the container.

All three Go services (`crawler`, `indexer`, `search-api`) log structured JSON via `zap`. Drop the
`service-1 |` prefix and pipe through `jq` to pretty-print each entry, if you have `jq` installed:

```sh
docker compose logs -f --no-log-prefix search-api | jq .
```
