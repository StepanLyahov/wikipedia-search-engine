# Домашнее задание №3. Indexer

## Цель

Научиться извлекать текст из HTML и индексировать документы в Elasticsearch.

## Что нужно сделать

Реализовать сервис: indexer

Сервис должен:

1. читать документы из PostgreSQL;
2. очищать HTML;
3. извлекать title;
4. извлекать body;
5. создавать индекс;
6. индексировать документы в Elasticsearch.

## Что должно получиться после выполнения

После запуска:

* существует индекс: wiki_pages

Каждый документ имеет вид:

```json
{
"id":1,
"url":"https://en.wikipedia.org/wiki/Elasticsearch",
"title": "Elasticsearch",
"body": "Elasticsearch is a search engine ..."
}
```

## Техническое описание

### Шаг 1

Получить страницы из PostgreSQL:

```sql
SELECT *

FROM pages
```

### Шаг 2

Извлечь title.

Например:

```html
<h1 id="firstHeading">

Elasticsearch

</h1>
```

↓

```text
Elasticsearch
```

### Шаг 3

Извлечь body.

Например:

```html
<p>

Elasticsearch is a search engine based on Apache Lucene.

</p>
```

↓

```text
Elasticsearch is a search engine based on Apache Lucene.
```

Необходимо:

* удалить HTML теги;
* удалить script;
* удалить style;
* объединить текст в одну строку.

### Шаг 4

Создать индекс:

```json
PUT wiki_pages

{

"mappings": {

    "properties": {

      "id": {

        "type":"long"

      },

      "url": {

        "type":"keyword"

      },

      "title": {

        "type":"text"

      },

      "body": {

        "type":"text"

      }

    }

}

}
```

### Шаг 5

Выполнить Bulk Indexing.

Например:

```text
POST _bulk
```

Отправить сразу несколько документов.

## API

```go
type Indexer interface {

    Run(ctx context.Context) error

}
```

## Критерии приемки

* индекс создается автоматически;
* документы индексируются через Bulk API;
* title корректно извлекается;
* body корректно очищается;
* url сохраняется в Elasticsearch;
* повторный запуск индексатора не приводит к ошибкам.