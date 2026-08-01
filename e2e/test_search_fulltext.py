"""Verifies GET /search: BM25 full-text search over indexed title/body fields."""
import requests


def _distinctive_body_term(es_base_url: str, index: str, doc_id: str) -> str:
    """Ask Elasticsearch itself, via the term vectors API, which term in this document's
    indexed body field is least common across the whole corpus.

    This deliberately doesn't re-tokenize the raw text in Python: Lucene's standard analyzer
    keeps some substrings together in ways a naive word-split won't replicate (e.g. a citation
    URL like "n3.nabble.com" is indexed as the single token "nabble.com", not "nabble" and
    "com"), so a term picked by guessing tokenization can come back with zero hits even though
    it's genuinely in the page's body. Reading real per-document term stats from the terms
    Elasticsearch actually indexed sidesteps that entirely, and stays correct whether the
    corpus has 2 pages or 200: the term need not be globally rare, just rarer than every other
    term on this page.
    """
    resp = requests.get(
        f"{es_base_url}/{index}/_termvectors/{doc_id}",
        params={"fields": "body", "term_statistics": "true"},
        timeout=10,
    )
    resp.raise_for_status()

    terms = resp.json()["term_vectors"]["body"]["terms"]
    candidates = [term for term in terms if len(term) >= 6 and term.isalpha()]
    if not candidates:
        raise AssertionError(f"document {doc_id!r} has no eligible body terms")

    return min(candidates, key=lambda term: terms[term]["doc_freq"])


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


def test_search_matches_body_field_not_just_title(
    search_api_base_url, indexed_documents, crawler_seed_url, es_base_url, elasticsearch_index
):
    seed_doc_id = next(
        doc["_id"] for doc in indexed_documents if doc["_source"]["url"] == crawler_seed_url
    )
    body_term = _distinctive_body_term(es_base_url, elasticsearch_index, seed_doc_id)

    resp = requests.get(
        f"{search_api_base_url}/search", params={"q": body_term, "size": 50}, timeout=10
    )
    assert resp.status_code == 200

    hits = resp.json()["hits"]
    urls = [hit["url"] for hit in hits]
    assert crawler_seed_url in urls, (
        f"query {body_term!r} (this page's own least-common indexed body term) did not surface "
        "it via /search; multi_match may not be searching the body field"
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
