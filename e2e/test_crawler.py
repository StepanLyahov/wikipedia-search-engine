"""Verifies the crawler fetched pages and saved them correctly into Postgres."""
from datetime import datetime, timedelta, timezone


def test_crawler_saved_at_least_one_page(crawled_pages):
    assert len(crawled_pages) >= 1, "crawler did not save any pages into Postgres"


def test_seed_page_was_saved_correctly(crawled_pages, crawler_seed_url):
    by_url = {page["url"]: page for page in crawled_pages}
    assert crawler_seed_url in by_url, f"seed URL {crawler_seed_url!r} was never saved"

    seed_page = by_url[crawler_seed_url]
    assert seed_page["status"] == 200
    assert seed_page["title"], "seed page has no title"
    assert seed_page["html"] and "<" in seed_page["html"], "seed page has no HTML content"


def test_pages_have_valid_wikipedia_urls(crawled_pages):
    for page in crawled_pages:
        assert page["url"].startswith("https://en.wikipedia.org/wiki/"), (
            f"unexpected URL saved by crawler: {page['url']!r}"
        )


def test_pages_have_no_duplicate_urls(crawled_pages):
    urls = [page["url"] for page in crawled_pages]
    assert len(urls) == len(set(urls)), "crawler saved duplicate URLs"


def test_pages_were_saved_recently(crawled_pages):
    now = datetime.now(timezone.utc)
    for page in crawled_pages:
        created_at = page["created_at"]
        if created_at.tzinfo is None:
            created_at = created_at.replace(tzinfo=timezone.utc)
        assert now - created_at < timedelta(hours=1), (
            f"page {page['url']!r} has a stale created_at ({created_at}); "
            "was this test run against an old, already-populated database?"
        )
