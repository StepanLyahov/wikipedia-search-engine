# Домашнее задание №4. BM25 Search API

## Цель

Реализовать полнотекстовый поиск по документам, используя встроенный алгоритм ранжирования Elasticsearch — BM25.

После выполнения задания пользователь должен иметь возможность выполнять поисковые запросы и получать наиболее релевантные документы.

## Что нужно сделать

Необходимо реализовать сервис: search-api

Сервис должен предоставлять REST API: GET /search

Query parameters:

```text
q      - поисковый запрос
from   - смещение
size   - количество документов
```

## Что должно получиться после выполнения

Запрос: GET /search?q=distributed systems

должен возвращать:

```json
{
  "hits": [
    {
      "title": "Distributed computing",
      "url": "https://en.wikipedia.org/wiki/Distributed_computing",
      "score": 13.42
    }
  ]
}
```

## Техническое описание

### Шаг 1

Реализовать endpoint: GET /search

### Шаг 2

Выполнить запрос в Elasticsearch.

Использовать: multi_match

Пример:

```json
POST wiki_pages/_search

{
  "query": {
    "multi_match": {
      "query":
      "distributed systems",
      "fields": [
        "title^3",
        "body^1"
      ]
    }
  }
}
```

### Что означает

```text
title^3
```

Поле title в 3 раза важнее.

Например:

Документ:

```text
Title:

Distributed Systems

Body:

Lorem ipsum...
```

должен иметь больший score, чем:

```text
Title:
    Random article
Body:
    distributed systems...
```

### Шаг 3

Поддержать пагинацию.
Использовать:

```json
{
  "from":0,
  "size":10
}
```

## API

Request:

```text
GET /search?q=elasticsearch&from=0&size=10
```

Response:

```json
{
"hits":[
    {
      "title":"Elasticsearch",
      "url":"https://en.wikipedia.org/wiki/Elasticsearch",
      "score":12.3
    }
]
}
```

## Критерии приемки

* поиск работает;
* поиск выполняется одновременно по:
  * title;
  * body;
* title имеет больший вес;
* результаты отсортированы по score;
* поддерживается пагинация.
