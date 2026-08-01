# PostgreSQL

Not a service — just the schema PostgreSQL applies on first boot, plus a place to note anything
DB-specific for `crawler` (the only writer) and `indexer` (the only reader).

## Structure

* `init/001_pages.sql` — schema, mounted read-only into the official `postgres:16-alpine` image
  at `/docker-entrypoint-initdb.d`. Postgres runs every `.sql`/`.sh` file in that directory, in
  filename order, **only the first time** the data volume is created (i.e. never again once
  `postgres_data` exists — see below).

## Schema

```sql
CREATE TABLE IF NOT EXISTS pages (
    id BIGSERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    html TEXT NOT NULL,
    status INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

* `url` is `UNIQUE` — this is what makes the crawler's re-run behavior idempotent (`crawler`
  checks `Exists(url)` before fetching, so a page already saved is never re-fetched or
  re-inserted).
* `title`/`status` are nullable in principle, but `crawler` always sets them from the HTTP
  response; `indexer` doesn't use `pages.title` for the indexed document title (see
  [`indexer/README.md`](../indexer/README.md) — it re-extracts the title from the `<h1
  id="firstHeading">` heading in `html` instead).
* `html` is the full raw HTML of the fetched page, used later by `indexer` for title/body
  extraction. It is never cleaned or truncated here.

## Extending this

* **Adding a migration**: add a new `init/NNN_description.sql` file. It only runs on a fresh
  volume — to apply it to an existing environment, either `docker compose down -v` (drops all
  data, re-crawls from scratch) or run the SQL by hand against the running container.
* **Credentials/DB name** are hardcoded in `docker-compose.yml` (`wiki`/`wiki`/`wiki`), not
  sourced from `.env` — this is a local dev stack, not a deployment target. If that ever needs to
  change, update `docker-compose.yml`'s `postgres.environment` and every service's
  `DATABASE_URL`/connection env var together.
* Connect directly for debugging:
  ```sh
  docker compose exec postgres psql -U wiki -d wiki
  ```
