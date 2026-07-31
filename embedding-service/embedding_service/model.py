"""Wraps the sentence-transformers model used to turn text into embeddings."""

from __future__ import annotations

from sentence_transformers import SentenceTransformer


class EmbeddingModel:
    """Loads a SentenceTransformer model and encodes text into a dense vector."""

    def __init__(self, model_name: str) -> None:
        self._model = SentenceTransformer(model_name)

    def encode(self, text: str) -> list[float]:
        return self._model.encode(text).tolist()
