"""gRPC servicer implementing EmbeddingService."""

from __future__ import annotations

import grpc

from .model import EmbeddingModel
from .pb import embedding_pb2, embedding_pb2_grpc


class EmbeddingServicer(embedding_pb2_grpc.EmbeddingServiceServicer):
    """Serves Embed requests by delegating text encoding to an EmbeddingModel."""

    def __init__(self, model: EmbeddingModel) -> None:
        self._model = model

    def Embed(self, request: embedding_pb2.EmbedRequest, context: grpc.ServicerContext) -> embedding_pb2.EmbedResponse:
        if not request.text.strip():
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "text is required")

        vector = self._model.encode(request.text)

        return embedding_pb2.EmbedResponse(vector=vector)
