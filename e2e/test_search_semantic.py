"""Verifies GET /semantic: query embedding via the embedding-service gRPC call, followed by
a kNN search over the `embedding` field.

The crawler now discovers dozens of pages from the default seed (the "Elasticsearch" article),
including several very closely related ones (Kibana, Apache Lucene, OpenSearch, "Search engine
(computing)", ...). With that much real topical overlap, a generic query can legitimately put
one of those close relatives ahead of the Elasticsearch page itself -- that's correct semantic
behavior, not a bug. So instead of asserting the Elasticsearch page ranks exactly first, these
tests assert it's *findable* within a generous top-k for an on-topic query, and *not* ranked
first for a clearly unrelated one. Since this tests against live Wikipedia content, future
article edits could in principle shift a result in or out of that window; that's an accepted
trade-off of testing the real pipeline instead of canned fixtures.
"""
import requests

# Generous but not unbounded: enough headroom that a small score shuffle from a future Wikipedia
# edit won't flip the test, while still proving the match isn't just "everything is top-k".
ON_TOPIC_TOP_K = 20


def test_semantic_requires_query(search_api_base_url):
    resp = requests.get(f"{search_api_base_url}/semantic", timeout=10)
    assert resp.status_code == 400
    assert "error" in resp.json()


def test_semantic_finds_search_engine_topic_by_meaning(search_api_base_url, crawler_seed_url):
    """A query about search-engine technology, with no words quoted from the article, should
    still find the Elasticsearch page -- even if close relatives (Lucene, Kibana, "Search engine
    (computing)", ...) legitimately outrank it."""
    resp = requests.get(
        f"{search_api_base_url}/semantic",
        params={"q": "how search engines work", "k": ON_TOPIC_TOP_K},
        timeout=10,
    )
    assert resp.status_code == 200

    hits = resp.json()["hits"]
    urls = [hit["url"] for hit in hits]
    assert crawler_seed_url in urls, (
        f"expected the search-engine article within the top {ON_TOPIC_TOP_K} semantic results "
        f"for a search-engine query, got {urls!r}"
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


def test_semantic_k_limits_result_count(search_api_base_url, indexed_documents):
    resp = requests.get(
        f"{search_api_base_url}/semantic", params={"q": "how search engines work", "k": 1}, timeout=10
    )
    assert resp.status_code == 200
    assert len(resp.json()["hits"]) == 1

    # k above Elasticsearch's default response size (10) must still be honored: the knn clause's
    # k only ranks candidates internally, the top-level "size" is what actually bounds the
    # response -- easy to set one and forget the other.
    if len(indexed_documents) >= ON_TOPIC_TOP_K:
        resp = requests.get(
            f"{search_api_base_url}/semantic",
            params={"q": "how search engines work", "k": ON_TOPIC_TOP_K},
            timeout=10,
        )
        assert resp.status_code == 200
        assert len(resp.json()["hits"]) == ON_TOPIC_TOP_K


def test_semantic_results_sorted_by_score_descending(search_api_base_url):
    resp = requests.get(
        f"{search_api_base_url}/semantic",
        params={"q": "how search engines work", "k": ON_TOPIC_TOP_K},
        timeout=10,
    )
    scores = [hit["score"] for hit in resp.json()["hits"]]
    assert scores == sorted(scores, reverse=True)


def test_semantic_and_fulltext_agree_on_topic(search_api_base_url, crawler_seed_url):
    """Sanity-check that both search modes surface the same page for a query that is both a
    literal title match and a coherent semantic query: /search should rank it first (BM25's
    title boost is decisive for an exact title match), /semantic should at least find it."""
    fulltext_hits = requests.get(
        f"{search_api_base_url}/search", params={"q": "Elasticsearch"}, timeout=10
    ).json()["hits"]
    semantic_hits = requests.get(
        f"{search_api_base_url}/semantic",
        params={"q": "Elasticsearch", "k": ON_TOPIC_TOP_K},
        timeout=10,
    ).json()["hits"]

    assert fulltext_hits and fulltext_hits[0]["url"] == crawler_seed_url

    semantic_urls = [hit["url"] for hit in semantic_hits]
    assert crawler_seed_url in semantic_urls
