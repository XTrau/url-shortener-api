# URL Shortener API

## Описание проекта

API для сокращения ссылок.

Проект является учебным и демонстрирует умения работы с технологиями на Golang. 

Данные в API хранятся в базе данных, также реализовано кеширование данных.

Пакеты Golang которые используются в проекте:

- net/http
- database/sql
- pgx (драйвер для PostgreSQL)
- golang-migrate (для миграций)
- go-redis/v9 (для Redis)

### Запуск в Docker compose

Windows:

```
docker-compose up -d
```

Linux:

```
docker compose up -d
```

### Пример задачи

Например есть длинная ссылка:

```
https://somesite.com/path/path2/path3/path4?param1=value1&param2=value2&param3=value3
```

Этот сервис сделает её короче:

```
https://service.com/Idhfayr73e
```

Если переходить на короткую ссылку, то сервис будет редиректить пользователя на основную ссылку.

### API Docs

GET `/` - фронтенд

POST `/short` - сокращает переданную ссылку, возвращает slug для короткой ссылки. В body необходимо передавать url который нужно сократить.

Пример body:
```
{ "url": "https://google.com" }
```

Пример ответа:

```
{ "slug":"IziXC5Az" }
```

GET `/s/{slug}` - делает редирект на ссылку к которой привязан переданный slug.