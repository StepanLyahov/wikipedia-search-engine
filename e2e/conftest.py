"""Fixtures that stand up the full docker-compose stack once per test session, then expose
Postgres/Elasticsearch/search-api clients to the test modules.

The pipeline under test is: crawler -> Postgres -> indexer (+ embedding-service) ->
Elasticsearch -> search-api. Rather than re-deriving each service's env-var defaults here,
fixtures read them back from `docker compose config`, so the tests always exercise whatever
configuration (.env overrides included) actually drives the stack.
"""
from __future__ import annotations

import os
import time
from urllib.parse import urlparse

import psycopg2
import psycopg2.extras
import pytest
import requests

import docker_utils as compose

SKIP_COMPOSE = os.environ.get("E2E_SKIP_COMPOSE", "").lower() in ("1", "true", "yes")
KEEP_STACK = os.environ.get("E2E_KEEP_STACK", "").lower() in ("1", "true", "yes")

CRAWLER_TIMEOUT = float(os.environ.get("E2E_CRAWLER_TIMEOUT", "180"))
INDEXER_TIMEOUT = float(os.environ.get("E2E_INDEXER_TIMEOUT", "180"))
SERVICE_HEALTHY_TIMEOUT = float(os.environ.get("E2E_SERVICE_HEALTHY_TIMEOUT", "180"))
HTTP_READY_TIMEOUT = float(os.environ.get("E2E_HTTP_READY_TIMEOUT", "60"))


def _wait_for_http(url: str, timeout: float) -> None:
    deadline = time.monotonic() + timeout
    last_error = None
    while time.monotonic() < deadline:
        try:
            requests.get(url, timeout=5)
            return
        except requests.RequestException as exc:
            last_error = exc
            time.sleep(1)
    raise TimeoutError(f"{url} was not reachable within {timeout}s: {last_error}")


@pytest.fixture(scope="session")
def compose_config():
    return compose.resolved_config()


@pytest.fixture(scope="session", autouse=True)
def compose_stack(compose_config):
    """Bring the whole stack up once, wait for the crawl -> index pipeline to finish, and
    tear it down at the end of the session (skippable via E2E_SKIP_COMPOSE/E2E_KEEP_STACK for
    faster local iteration on the tests themselves)."""
    if SKIP_COMPOSE:
        yield
        return

    compose.down(volumes=True)
    compose.up(build=True)

    try:
        compose.wait_for_healthy("postgres", SERVICE_HEALTHY_TIMEOUT)
        compose.wait_for_healthy("elasticsearch", SERVICE_HEALTHY_TIMEOUT)
        compose.wait_for_healthy("embedding-service", SERVICE_HEALTHY_TIMEOUT)
        compose.wait_for_one_shot_completion("crawler", CRAWLER_TIMEOUT)
        compose.wait_for_one_shot_completion("indexer", INDEXER_TIMEOUT)

        search_api_port = compose.published_port(compose_config, "search-api", 8080)
        _wait_for_http(f"http://localhost:{search_api_port}/search?q=healthcheck", HTTP_READY_TIMEOUT)

        yield
    finally:
        if not KEEP_STACK:
            compose.down(volumes=True)


@pytest.fixture(scope="session")
def crawler_seed_url(compose_config) -> str:
    env = compose.service_environment(compose_config, "crawler")
    return env["CRAWLER_SEED_URL"]


@pytest.fixture(scope="session")
def elasticsearch_index(compose_config) -> str:
    env = compose.service_environment(compose_config, "indexer")
    return env["ELASTICSEARCH_INDEX"]


@pytest.fixture(scope="session")
def es_base_url(compose_config) -> str:
    port = compose.published_port(compose_config, "elasticsearch", 9200)
    return f"http://localhost:{port}"


@pytest.fixture(scope="session")
def search_api_base_url(compose_config) -> str:
    port = compose.published_port(compose_config, "search-api", 8080)
    return f"http://localhost:{port}"


@pytest.fixture(scope="session")
def pg_conn(compose_config):
    database_url = compose.service_environment(compose_config, "crawler")["DATABASE_URL"]
    parsed = urlparse(database_url)
    postgres_port = compose.published_port(compose_config, "postgres", 5432)

    deadline = time.monotonic() + SERVICE_HEALTHY_TIMEOUT
    last_error = None
    while time.monotonic() < deadline:
        try:
            conn = psycopg2.connect(
                host="localhost",
                port=postgres_port,
                dbname=parsed.path.lstrip("/"),
                user=parsed.username,
                password=parsed.password,
            )
            break
        except psycopg2.OperationalError as exc:
            last_error = exc
            time.sleep(1)
    else:
        raise TimeoutError(f"could not connect to Postgres: {last_error}")

    yield conn
    conn.close()


@pytest.fixture(scope="session")
def crawled_pages(pg_conn) -> list[dict]:
    """All rows the crawler saved, fetched once per session and reused across test modules."""
    with pg_conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
        cur.execute("SELECT id, url, title, html, status, created_at FROM pages ORDER BY id")
        return [dict(row) for row in cur.fetchall()]


@pytest.fixture(scope="session")
def indexed_documents(compose_stack, es_base_url, elasticsearch_index) -> list[dict]:
    """All documents the indexer wrote to Elasticsearch, fetched once per session."""
    requests.post(f"{es_base_url}/{elasticsearch_index}/_refresh", timeout=10)

    resp = requests.post(
        f"{es_base_url}/{elasticsearch_index}/_search",
        json={"size": 1000, "query": {"match_all": {}}},
        timeout=10,
    )
    resp.raise_for_status()
    return resp.json()["hits"]["hits"]
