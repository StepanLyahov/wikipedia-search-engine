# Embedding Service

Turns text into a 384-dimensional dense vector using the `all-MiniLM-L6-v2`
[sentence-transformers](https://www.sbert.net/) model, served over gRPC. The indexer calls this
service so embedding generation stays out of Elasticsearch, per the course rules (not all ES
versions/licenses support built-in embeddings, and production systems typically isolate the ML
model in its own service).

## API

Defined in [`../proto/embedding/embedding.proto`](../proto/embedding/embedding.proto):

```proto
service EmbeddingService {
  rpc Embed(EmbedRequest) returns (EmbedResponse);
}

message EmbedRequest { string text = 1; }
message EmbedResponse { repeated float vector = 1; }
```

`Embed` rejects blank text with `INVALID_ARGUMENT`.

## Structure

* `embedding_service/pb` — generated protobuf/gRPC stubs (checked in; regenerate with `make
  generate`);
* `embedding_service/model.py` — loads the SentenceTransformer model and encodes text;
* `embedding_service/servicer.py` — the gRPC servicer, depends only on an object with an
  `encode(text) -> list[float]` method;
* `embedding_service/server.py` — process entrypoint: builds and starts the gRPC server;
* `embedding_service/config.py` — environment configuration;
* `tests/test_servicer.py` — fast unit tests using a fake model;
* `tests/test_model.py` — loads the real model and asserts the vector is 384-dimensional.

## Configuration

| Env var                | Default             |
|-------------------------|----------------------|
| `GRPC_PORT`             | `50051`              |
| `EMBEDDING_MODEL_NAME`  | `all-MiniLM-L6-v2`   |
| `GRPC_MAX_WORKERS`      | `10`                 |

## Development

```sh
make test
make lint
```

Both targets create/reuse a local `.venv` and install `requirements-dev.txt` into it.

## Docker

The image bakes the model weights in at build time (`HF_HOME=/app/.cache/huggingface`), so the
container needs no network access at runtime:

```sh
docker compose up --build
```

## Extending this

* **Changing the model**: set `EMBEDDING_MODEL_NAME` — but the *dimension* is baked into three
  other places that must move together: `INDEXER_EMBEDDING_DIMS` (env, defaults to `384`,
  drives the Elasticsearch `dense_vector` mapping in `indexer`), and any test asserting a fixed
  vector length (`tests/test_model.py` here, `e2e/test_indexer.py` in the end-to-end suite). A
  model with a different output size needs a fresh `wiki_pages` index (`docker compose down -v`)
  since Elasticsearch's `dense_vector` dimension is fixed at index-creation time.
* **The gRPC contract** lives in [`../proto/embedding/embedding.proto`](../proto/embedding/embedding.proto),
  shared with the Go clients in `indexer` and `search-api` — see
  [`../proto/README.md`](../proto/README.md) before adding an RPC or field, and regenerate all
  three consumers together.
* **`servicer.py` depends on a narrow protocol** (anything with `encode(text) -> list[float]`),
  not on `SentenceTransformer` directly — that's what makes `tests/test_servicer.py` fast (a fake
  model, no PyTorch import). Keep that seam if you change how embeddings are produced.
