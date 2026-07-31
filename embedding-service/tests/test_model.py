"""Integration test for the real embedding model: downloads and loads all-MiniLM-L6-v2."""

from __future__ import annotations

import pytest

from embedding_service.model import EmbeddingModel

EXPECTED_DIMS = 384


@pytest.fixture(scope="module")
def model() -> EmbeddingModel:
    return EmbeddingModel("all-MiniLM-L6-v2")


def test_encode_produces_a_384_dimensional_vector(model: EmbeddingModel) -> None:
    vector = model.encode("Elasticsearch is a search engine based on Apache Lucene.")

    assert len(vector) == EXPECTED_DIMS
    assert all(isinstance(value, float) for value in vector)


def test_encode_is_deterministic_for_the_same_text(model: EmbeddingModel) -> None:
    text = "Distributed systems coordinate independent computers as one system."

    assert model.encode(text) == model.encode(text)
