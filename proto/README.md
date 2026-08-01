# proto

The single source of truth for the gRPC contract between `embedding-service` (the server) and
its two Go clients, `indexer` and `search-api`. Generated code is checked into each consumer's own
tree (so a plain `go build`/`docker build` never needs `protoc` installed); this directory holds
only the `.proto` source.

## Structure

* `embedding/embedding.proto` — one RPC: `Embed(text) -> vector[384]`.

```proto
service EmbeddingService {
  rpc Embed(EmbedRequest) returns (EmbedResponse);
}

message EmbedRequest { string text = 1; }
message EmbedResponse { repeated float vector = 1; }
```

`Embed` rejects blank/whitespace-only text with `INVALID_ARGUMENT` (see
[`embedding-service/README.md`](../embedding-service/README.md)).

## Where the generated code lives

| Consumer | Generated files | Regenerate with |
|---|---|---|
| `indexer` (Go client) | `indexer/internal/pb/embedding/*.pb.go` | see command below |
| `search-api` (Go client) | `search-api/internal/pb/embedding/*.pb.go` | see command below |
| `embedding-service` (Python server) | `embedding-service/embedding_service/pb/*_pb2*.py` | `make generate` in `embedding-service/` |

The `.proto`'s `option go_package` points at `indexer`'s import path; `search-api`'s copy is
generated with an explicit `-M`/`Mfile=path` override (protoc-gen-go otherwise refuses to emit
a file whose declared `go_package` doesn't match the output module):

```sh
# from the repository root
protoc -I proto \
  --go_out=indexer --go_opt=module=github.com/wikipedia-search-engine/indexer \
  --go-grpc_out=indexer --go-grpc_opt=module=github.com/wikipedia-search-engine/indexer \
  proto/embedding/embedding.proto

protoc -I proto \
  --go_out=search-api --go_opt=Membedding/embedding.proto=github.com/wikipedia-search-engine/search-api/internal/pb/embedding \
  --go-grpc_out=search-api --go-grpc_opt=Membedding/embedding.proto=github.com/wikipedia-search-engine/search-api/internal/pb/embedding \
  proto/embedding/embedding.proto
```

Requires `protoc`, `protoc-gen-go` and `protoc-gen-go-grpc` on `PATH` (`go install
google.golang.org/protobuf/cmd/protoc-gen-go@latest` and
`google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`).

## Extending this

Adding a field or RPC: edit the `.proto` here first, regenerate all three consumers with the
commands above (or `make generate` for the Python one), then wire the new field/method through
each consumer's own port interface — Go callers depend on a hand-written interface
(`indexer`'s `service/indexer.Embedder`, `search-api`'s `service/search.Embedder`) implemented by
a thin adapter in `internal/transport/embedding`, never on the generated client type directly; the
Python server implements the generated `EmbeddingServiceServicer` base class directly in
`embedding_service/servicer.py`. Keep the three generated trees in sync — there's no CI check
enforcing that today, only convention.
