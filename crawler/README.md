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
