FROM python:3.11-slim

WORKDIR /app

ENV PYTHONUNBUFFERED=1 \
    HF_HOME=/app/.cache/huggingface \
    EMBEDDING_MODEL_NAME=all-MiniLM-L6-v2

COPY requirements.txt ./
# Install the CPU-only torch build first: this container never uses a GPU, and the default
# PyPI wheel otherwise pulls several GB of unused CUDA libraries.
RUN pip install --no-cache-dir torch --index-url https://download.pytorch.org/whl/cpu \
    && pip install --no-cache-dir -r requirements.txt

# Bake the model weights into the image so the container never needs network access at runtime.
RUN python -c "import os; from sentence_transformers import SentenceTransformer; SentenceTransformer(os.environ['EMBEDDING_MODEL_NAME'])"

COPY embedding_service ./embedding_service

RUN addgroup --system app && adduser --system --ingroup app app \
    && chown -R app:app /app
USER app

EXPOSE 50051
CMD ["python", "-m", "embedding_service.server"]
