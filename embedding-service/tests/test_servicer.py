"""Unit tests for EmbeddingServicer, using a fake model to avoid loading real weights."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from embedding_service.pb import embedding_pb2
from embedding_service.servicer import EmbeddingServicer


class _FakeModel:
    def __init__(self, vector: list[float]) -> None:
        self._vector = vector

    def encode(self, text: str) -> list[float]:
        return self._vector


def test_embed_returns_the_model_vector() -> None:
    servicer = EmbeddingServicer(_FakeModel([0.12, -0.45, 0.91]))
    context = MagicMock()

    response = servicer.Embed(embedding_pb2.EmbedRequest(text="distributed systems"), context)

    assert list(response.vector) == pytest.approx([0.12, -0.45, 0.91])
    context.abort.assert_not_called()


def test_embed_rejects_blank_text() -> None:
    servicer = EmbeddingServicer(_FakeModel([0.0]))
    context = MagicMock()
    context.abort.side_effect = RuntimeError("aborted")

    with pytest.raises(RuntimeError):
        servicer.Embed(embedding_pb2.EmbedRequest(text="   "), context)

    context.abort.assert_called_once()
