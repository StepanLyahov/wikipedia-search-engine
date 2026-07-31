"""Entrypoint for the embedding-service gRPC server."""

from __future__ import annotations

import logging
from concurrent import futures

import grpc

from .config import Config
from .model import EmbeddingModel
from .pb import embedding_pb2_grpc
from .servicer import EmbeddingServicer

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s %(message)s")
logger = logging.getLogger("embedding-service")


def serve() -> None:
    cfg = Config.from_env()

    logger.info("embedding-service starting model=%s port=%d", cfg.model_name, cfg.grpc_port)
    model = EmbeddingModel(cfg.model_name)
    logger.info("model loaded model=%s", cfg.model_name)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=cfg.max_workers))
    embedding_pb2_grpc.add_EmbeddingServiceServicer_to_server(EmbeddingServicer(model), server)
    server.add_insecure_port(f"[::]:{cfg.grpc_port}")

    server.start()
    logger.info("embedding-service started port=%d", cfg.grpc_port)

    server.wait_for_termination()


if __name__ == "__main__":
    serve()
