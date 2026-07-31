"""Environment-based configuration for the embedding-service."""

from __future__ import annotations

import os
from dataclasses import dataclass


@dataclass(frozen=True)
class Config:
    grpc_port: int
    model_name: str
    max_workers: int

    @classmethod
    def from_env(cls) -> Config:
        return cls(
            grpc_port=int(os.environ.get("GRPC_PORT", "50051")),
            model_name=os.environ.get("EMBEDDING_MODEL_NAME", "all-MiniLM-L6-v2"),
            max_workers=int(os.environ.get("GRPC_MAX_WORKERS", "10")),
        )
