"""Verifies GET /semantic: query embedding via the embedding-service gRPC call, followed by
a kNN search over the `embedding` field.

These queries were picked empirically against the default seed (the "Elasticsearch" Wikipedia
article) plus whatever it links to (currently just "Main Page"): each pair below reliably ranks
the topically-closer page first, by a comfortable score margin, even though the query text does
not literally quote the page. If a future Wikipedia edit meaningfully changes either article's
wording, the margins here could shift — that is an accepted trade-off of testing against live
content rather than fixtures, per the project's genuinely-end-to-end test brief.
"""
import requests


def test_semantic_requires_query(search_api_base_url):
    resp = requests.get(f"{search_api_base_url}/semantic", timeout=10)
    assert resp.status_code == 400
    assert "error" in resp.json()


def test_semantic_finds_search_engine_topic_by_meaning(search_api_base_url, crawler_seed_url):
    """A query about search-engine technology, with no words quoted from the article, should
    still rank the Elasticsearch page first."""
    resp = requests.get(
        f"{search_api_base_url}/semantic", params={"q": "how search engines work"}, timeout=10
    )
    assert resp.status_code == 200

    hits = resp.json()["hits"]
    assert hits, "expected at least one semantic hit"
    assert hits[0]["url"] == crawler_seed_url, (
        f"expected the search-engine article to rank first for a search-engine query, got {hits[0]!r}"
    )


def test_semantic_distinguishes_unrelated_topic(search_api_base_url, indexed_documents, crawler_seed_url):
    """A query about a different topic (encyclopedias in general, not search technology
    specifically) should rank some other indexed page above the Elasticsearch article --
    proving the ranking tracks meaning rather than always returning the same page."""
    if len(indexed_documents) < 2:
        return  # nothing to contrast against with a single-document corpus

    resp = requests.get(
        f"{search_api_base_url}/semantic",
        params={"q": "history of encyclopedias and free knowledge projects"},
        timeout=10,
    )
    assert resp.status_code == 200

    hits = resp.json()["hits"]
    assert hits, "expected at least one semantic hit"
    assert hits[0]["url"] != crawler_seed_url, (
        "expected a page other than the search-engine article to rank first for an "
        "encyclopedia-history query"
    )


def test_semantic_k_limits_result_count(search_api_base_url):
    resp = requests.get(
        f"{search_api_base_url}/semantic", params={"q": "how search engines work", "k": 1}, timeout=10
    )
    assert resp.status_code == 200
    assert len(resp.json()["hits"]) == 1


def test_semantic_results_sorted_by_score_descending(search_api_base_url):
    resp = requests.get(
        f"{search_api_base_url}/semantic", params={"q": "how search engines work"}, timeout=10
    )
    scores = [hit["score"] for hit in resp.json()["hits"]]
    assert scores == sorted(scores, reverse=True)


def test_semantic_and_fulltext_agree_on_topic(search_api_base_url, crawler_seed_url):
    """Sanity-check that both search modes converge on the same page for a query that is both
    a literal title match and a coherent semantic query."""
    fulltext_hits = requests.get(
        f"{search_api_base_url}/search", params={"q": "Elasticsearch"}, timeout=10
    ).json()["hits"]
    semantic_hits = requests.get(
        f"{search_api_base_url}/semantic", params={"q": "Elasticsearch"}, timeout=10
    ).json()["hits"]

    assert fulltext_hits and fulltext_hits[0]["url"] == crawler_seed_url
    assert semantic_hits and semantic_hits[0]["url"] == crawler_seed_url
