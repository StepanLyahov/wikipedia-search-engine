"""Verifies GET /search: BM25 full-text search over indexed title/body fields."""
import re

import requests

STOPWORDS = {
    "the", "and", "for", "are", "but", "not", "you", "all", "can", "her", "was", "one",
    "our", "out", "day", "get", "has", "him", "his", "how", "man", "new", "now", "old",
    "see", "two", "way", "who", "boy", "did", "its", "let", "put", "say", "she", "too",
    "use", "with", "this", "that", "from", "have", "also", "been", "were", "which",
    "their", "wikipedia", "encyclopedia",
}


def _distinctive_body_word(doc: dict) -> str:
    """Pick a real word from doc's indexed body that isn't part of its own title, so a
    search for it can only match via the body field, not the title-boosted field."""
    title_words = set(re.findall(r"[a-zA-Z]+", doc["title"].lower()))
    body_words = re.findall(r"[a-zA-Z]+", doc["body"].lower())

    for word in body_words:
        if len(word) >= 5 and word not in STOPWORDS and word not in title_words:
            return word

    raise AssertionError(f"could not find a distinctive body word for {doc['url']!r}")


def test_search_requires_query(search_api_base_url):
    resp = requests.get(f"{search_api_base_url}/search", timeout=10)
    assert resp.status_code == 400
    assert "error" in resp.json()


def test_search_finds_seed_page_by_title(search_api_base_url, crawler_seed_url):
    resp = requests.get(f"{search_api_base_url}/search", params={"q": "Elasticsearch"}, timeout=10)
    assert resp.status_code == 200

    hits = resp.json()["hits"]
    assert hits, "expected at least one hit for a title-matching query"
    assert hits[0]["url"] == crawler_seed_url, (
        "title-boosted multi_match should rank the exact title match first, "
        f"got {hits[0]['url']!r} instead"
    )


def test_search_matches_body_field_not_just_title(search_api_base_url, indexed_documents, crawler_seed_url):
    seed_doc = next(doc["_source"] for doc in indexed_documents if doc["_source"]["url"] == crawler_seed_url)
    body_word = _distinctive_body_word(seed_doc)

    resp = requests.get(f"{search_api_base_url}/search", params={"q": body_word}, timeout=10)
    assert resp.status_code == 200

    hits = resp.json()["hits"]
    urls = [hit["url"] for hit in hits]
    assert crawler_seed_url in urls, (
        f"query {body_word!r} (drawn from the page's own body) did not surface it via /search; "
        "multi_match may not be searching the body field"
    )


def test_search_pagination_limits_and_offsets(search_api_base_url, indexed_documents):
    if len(indexed_documents) < 2:
        return  # pagination has nothing meaningful to prove with a single-document corpus

    all_hits = requests.get(
        f"{search_api_base_url}/search", params={"q": "wikipedia", "size": 10}, timeout=10
    ).json()["hits"]
    if len(all_hits) < 2:
        return  # not all indexed documents necessarily match this query

    page_one = requests.get(
        f"{search_api_base_url}/search", params={"q": "wikipedia", "from": 0, "size": 1}, timeout=10
    ).json()["hits"]
    page_two = requests.get(
        f"{search_api_base_url}/search", params={"q": "wikipedia", "from": 1, "size": 1}, timeout=10
    ).json()["hits"]

    assert len(page_one) == 1
    assert len(page_two) == 1
    assert page_one[0]["url"] != page_two[0]["url"], "from=0 and from=1 returned the same document"


def test_search_no_match_returns_empty_hits(search_api_base_url):
    resp = requests.get(
        f"{search_api_base_url}/search", params={"q": "zzz_no_such_term_xyz123"}, timeout=10
    )
    assert resp.status_code == 200
    assert resp.json()["hits"] == []


def test_search_results_sorted_by_score_descending(search_api_base_url):
    resp = requests.get(f"{search_api_base_url}/search", params={"q": "wikipedia"}, timeout=10)
    scores = [hit["score"] for hit in resp.json()["hits"]]
    assert scores == sorted(scores, reverse=True)
