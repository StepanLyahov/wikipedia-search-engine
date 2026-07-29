# Домашнее задание №2. Web Crawler

## Цель

Реализовать сервис, который умеет обходить Википедию и сохранять страницы в PostgreSQL.

## Что нужно сделать

Необходимо реализовать сервис: crawler

Сервис получает начальный URL: https://en.wikipedia.org/wiki/Elasticsearch

и должен:

1. скачать HTML страницы;
2. извлечь ссылки `/wiki/...`;
3. перейти по найденным ссылкам;
4. сохранить страницу в PostgreSQL;
5. избежать повторного скачивания.

## Что должно получиться после выполнения

После запуска crawler:

* PostgreSQL содержит до 100 страниц;
* каждая страница имеет:
  * url;
  * title;
  * html;
  * status;
  * created_at;
* отсутствуют дубликаты.

## Техническое описание

### Шаг 1

Создать таблицу:

```sql
CREATE TABLE pages (

    id BIGSERIAL PRIMARY KEY,

    url TEXT UNIQUE NOT NULL,

    title TEXT,

    html TEXT NOT NULL,

    status INTEGER,

    created_at TIMESTAMP

);
```

### Шаг 2

Скачать HTML через HTTP GET.

Например:

```text
GET

https://en.wikipedia.org/wiki/Elasticsearch
```

### Шаг 3

Извлечь ссылки:

```text
/wiki/Lucene

/wiki/Search_engine

/wiki/Information_retrieval

/wiki/Apache_Lucene
```

### Шаг 4

Реализовать BFS обход.

Поддержать:

```text
maxDepth = 2

maxPages = 100
```

### Шаг 5

Реализовать защиту от повторного обхода.

## API

```go
type Crawler interface {
    Crawl(seed string) error

}
```

## Критерии приемки

* страницы скачиваются;
* страницы сохраняются в PostgreSQL;
* отсутствуют дубликаты;
* crawler не зацикливается;
* соблюдаются ограничения:
  * maxDepth;
  * maxPages.