# Домашнее задание №6. Vector Search

## Цель

Реализовать поиск по смыслу с использованием embeddings.

## Что нужно сделать

Добавить новый endpoint:

```text
GET /semantic
```

Query params:

```text
q k
```

## Что должно происходить

Пользователь выполняет запрос:

```text
how search engines work
```

search-api:

1.Получает embedding запроса.

```text
grpc
↓
embedding-service
↓
[0.12,-0.42,...]
```

2.Выполняет kNN поиск в Elasticsearch.

3.Возвращает наиболее похожие документы.

## Что должно получиться

Запрос:

```text
GET /semantic?q=how search engines work
```

должен находить:

```text
Search Engine
Information Retrieval
Lucene
Elasticsearch
```

Даже если в документе нет точного совпадения слов.

## Техническое описание

### Шаг 1

Получить embedding запроса.

Например:

```text
how search engines work
↓
[0.12,-0.45,...]
```

### Шаг 2

Выполнить kNN запрос.

Пример:

```json
POST wiki_pages/_search

{
  "knn": {
    "field": "embedding",
    "query_vector": [0.12, -0.45, ...],
    "k": 10,
    "num_candidates": 100
  }
}
```

### Что означают параметры

```text
field
```

Поле с embeddings.

```text
query_vector
```

Embedding поискового запроса.

```text
k
```

Количество документов в ответе.

```text
num_candidates
```

Количество кандидатов, которые Elasticsearch рассматривает перед ранжированием.

Чем больше:
* тем выше качество;
* тем медленнее поиск.

### Шаг 3

Преобразовать ответ Elasticsearch.

Вернуть:

```json
{
"hits":[
    {
      "title":"Search Engine",
      "url":"https://en.wikipedia.org/wiki/Search_engine",
      "score":0.93
    }
]
}
```

## API

Request:

```text
GET /semantic?q=how search engines work&k=10
```

Response:

```json
{
"hits":[
    {
      "title":"Search Engine",
      "url": "https://en.wikipedia.org/wiki/Search_engine",
      "score":0.93
    }
]
}
```

## Критерии приемки
* search-api получает embedding через gRPC;
* выполняется kNN поиск;
* поиск находит документы по смыслу;
* документы возвращаются по убыванию score.