VENV := .venv
PYTHON := $(VENV)/bin/python

.PHONY: venv test lint generate

venv:
	python3 -m venv $(VENV)
	$(PYTHON) -m pip install --upgrade pip -q
	$(PYTHON) -m pip install -r requirements-dev.txt -q

test: venv
	$(PYTHON) -m pytest tests/ -v

lint: venv
	$(PYTHON) -m ruff check .
	$(PYTHON) -m ruff format --check .

# Regenerates the gRPC/protobuf stubs from the shared proto/embedding/embedding.proto.
# The generated grpc file imports its pb2 module with an absolute import; patch it to a
# relative import so it works as part of the embedding_service.pb package.
generate: venv
	$(PYTHON) -m grpc_tools.protoc -I ../proto/embedding \
		--python_out=embedding_service/pb \
		--grpc_python_out=embedding_service/pb \
		--pyi_out=embedding_service/pb \
		embedding.proto
	sed -i '' 's/^import embedding_pb2 as embedding__pb2$$/from . import embedding_pb2 as embedding__pb2/' embedding_service/pb/embedding_pb2_grpc.py
