"""Verifies the indexer turned crawled pages into Elasticsearch documents, each carrying a
384-dim embedding produced by the embedding-service."""
import statistics

import requests

EMBEDDING_DIMS = 384


def test_index_exists_with_expected_mapping(es_base_url, elasticsearch_index):
    resp = requests.get(f"{es_base_url}/{elasticsearch_index}/_mapping", timeout=10)
    assert resp.status_code == 200, f"index {elasticsearch_index!r} does not exist"

    properties = resp.json()[elasticsearch_index]["mappings"]["properties"]

    assert properties["title"]["type"] == "text"
    assert properties["body"]["type"] == "text"
    assert properties["url"]["type"] == "keyword"

    embedding = properties["embedding"]
    assert embedding["type"] == "dense_vector"
    assert embedding["dims"] == EMBEDDING_DIMS
    assert embedding["similarity"] == "cosine"


def test_every_crawled_page_was_indexed(indexed_documents, crawled_pages):
    indexed_urls = {doc["_source"]["url"] for doc in indexed_documents}
    crawled_urls = {page["url"] for page in crawled_pages}

    assert indexed_urls == crawled_urls, (
        f"indexed URLs do not match crawled URLs\n"
        f"missing from index: {crawled_urls - indexed_urls}\n"
        f"unexpected in index: {indexed_urls - crawled_urls}"
    )


def test_documents_have_title_url_and_body(indexed_documents):
    for doc in indexed_documents:
        source = doc["_source"]
        assert source["title"], f"document {source.get('url')!r} has an empty title"
        assert source["url"].startswith("https://en.wikipedia.org/wiki/")
        assert source["body"], f"document {source['url']!r} has an empty body"
        assert "<" not in source["body"], f"document {source['url']!r} body still contains HTML"


def test_documents_have_384_dim_embeddings(indexed_documents):
    for doc in indexed_documents:
        source = doc["_source"]
        embedding = source.get("embedding")
        assert embedding is not None, f"document {source['url']!r} has no embedding"
        assert len(embedding) == EMBEDDING_DIMS, (
            f"document {source['url']!r} has a {len(embedding)}-dim embedding, want {EMBEDDING_DIMS}"
        )
        assert all(isinstance(value, float) for value in embedding)
        # A real sentence embedding has spread-out values; a stub/zero vector would not.
        assert statistics.pstdev(embedding) > 0, (
            f"document {source['url']!r} has a degenerate (near-constant) embedding vector"
        )
