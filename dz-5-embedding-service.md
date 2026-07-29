# Домашнее задание №5. Embedding Service

## Цель

Научиться генерировать embeddings для документов и сохранять их в Elasticsearch.

## Что нужно сделать

Реализовать сервис: embedding-service

Сервис должен предоставлять gRPC API: Embed(text) -> vector

## Что должно получиться после выполнения

Документ в Elasticsearch:

```json
{
"id":1,
"url": "https://en.wikipedia.org/wiki/Elasticsearch",
"title": "Elasticsearch",
"body": "Elasticsearch is a search engine ...",
"embedding":[0.12, -0.45, 0.91, ...]
}
```

## Техническое описание

### Шаг 1

Реализовать gRPC сервис.

```proto
syntax = "proto3";
package embedding;
service EmbeddingService {
  rpc Embed(
      EmbedRequest
  )

  returns (
      EmbedResponse
  );
}
```

Request:

```proto
message EmbedRequest {
    string text = 1;
}
```

Response:

```proto
message EmbedResponse {
    repeated float vector = 1;
}
```

### Шаг 2

Выбрать модель.
Рекомендуется:

```text
all-MiniLM-L6-v2
```

Причины:
* бесплатная;
* небольшая;
* быстро работает;
* размер embedding:

```text
384
```

### Шаг 3

Indexер должен обращаться к embedding-service.

Например:

```text
Indexer
↓
grpc
↓
Embedding Service
↓
[0.12,-0.43,...]
↓
Indexer
↓
Elasticsearch
```

### Шаг 4

Добавить новое поле в mapping.

```json
PUT wiki_pages

{
"mappings": {
    "properties": {
      "embedding": {
        "type":"dense_vector",
        "dims":384,
        "index":true,
        "similarity":"cosine"
      }
    }
}
}
```

## Почему embeddings генерируются отдельно

В рамках курса запрещается использовать встроенную генерацию embeddings в Elasticsearch.

Причины:
1.Не все версии Elasticsearch поддерживают эту возможность.
2.Функциональность зависит от лицензии.
3.В production-системах ML модель часто выносится в отдельный сервис.

## Критерии приемки
* реализован gRPC сервис;
* embeddings имеют размер 384;
* indexer получает embeddings через gRPC;
* embeddings сохраняются в Elasticsearch.