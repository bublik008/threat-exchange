# Threat Exchange Service

REST-сервис на Go для обмена информацией об угрозах, реализованный на основе предоставленной OpenAPI-спецификации.

## Стек технологий

- **Язык:** Go (1.27+)
- **База данных:** PostgreSQL 16
- **Драйвер БД:** `pgx` (pgxpool)
- **Кодогенерация:** `oapi-codegen` (генерация типов и серверного интерфейса по OpenAPI)
- **Инфраструктура:** Docker, Docker Compose

## Быстрый запуск (через Docker)

Весь проект (база данных + приложение) поднимается одной командой (обязательно наличие Docker и Docker Compose)

1. Запуск контейнеры в фоне:
   
```bash
   docker compose -f deployments/docker-compose.yml up --build -d
```

2. Применение миграции

```bash
docker exec -i threat-exchange-db psql -U threat -d threat_exchange <  migrations/0001_init_sql.up.sql
```

3. Сервис будет доступен по адресу: `http://localhost:8080`

## Тестирование

В проекте реализованы как юнит-, так и интеграционные тесты.

**Запуск юнит-тестов (быстрые, без БД):**

```bash
go test ./...
```

**Запуск интеграционных тестов (требуется запущенный Docker, так как используется `testcontainers`):**

```bash
go test -tags=integration ./internal/repository/postgres/ -v
```

## Примеры запросов (cURL)

**Создать угрозу (POST /info):**

```bash

curl -X POST http://localhost:8080/info \
  -H "Content-Type: application/json" \
  -d '{
    "md5": "d41d8cd98f00b204e9800998ecf8427e",
    "name": "Trojan.Win32.Generic",
    "severity": "high",
    "threat_type": "trojan",
    "description": "Generic trojan detected by signature",
    "tags": ["win32", "packed"]
  }'
```

**Получить информацию (GET /info/{md5}):**

```bash
curl http://localhost:8080/info/d41d8cd98f00b204e9800998ecf8427e
```

**Обновить информацию (PUT /info/{md5}):**

```bash

curl -X PUT http://localhost:8080/info/d41d8cd98f00b204e9800998ecf8427e \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Trojan.Win32.Updated",
    "severity": "critical"
  }'
```

## Архитектура проекта

Проект следует принципам чистой архитектуры. Зависимости идут сверху вниз, слои общаются через интерфейсы:

- `cmd/server` — точка входа, DI-контейнер (сборка зависимостей), graceful shutdown.
    
- `internal/api` — HTTP-транспорт. Реализация сгенерированного интерфейса, middleware (логирование), маппинг ошибок в HTTP-коды.
    
- `internal/service` — бизнес-логика. Валидация, нормализация MD5, работа с доменными ошибками.
    
- `internal/repository` — интерфейс хранилища. Реализация на `pgx` в подпакете `postgres`.
    
- `api/openapi.yaml` — исходная спецификация (Source of Truth).
    
- `migrations` — SQL-миграции для базы данных.
    

## Реализованные требования и дополнительные плюсы

### Дополнение модели угрозы

В исходную спецификацию были добавлены осмысленные поля:

- `severity` (enum: low, medium, high, critical) — уровень опасности. Обязательное поле при создании.
    
- `threat_type` — тип угрозы (текст).
    
- `description` — текстовое описание.
    
- `tags` — массив строк для классификации.
    
- `created_at` / `updated_at` — временные метки (заполняются сервером).
    

Все ограничения (maxLength, pattern, enum) описаны в спецификации. На уровне БД добавлены CHECK-ограничения для `severity`.

### Обработка ошибок и валидация

- Валидация формата MD5 (regex + приведение к нижнему регистру для избежания дублей).
    
- Проверка обязательных полей (name, severity).
    
- Корректный маппинг доменных ошибок в HTTP-статусы:
    
    - `400 Bad Request` — невалидный MD5 или JSON.
        
    - `404 Not Found` — запись не найдена.
        
    - `409 Conflict` — MD5 уже существует.
        
    - `500 Internal Server Error` — внутренние ошибки (с логированием реальной причины, без утечки наружу).
        

### Кодогенерация

Использован `oapi-codegen` с настройкой `std-http-server`. Это позволяет генерировать серверный интерфейс под стандартный `net/http` (без тяжелых фреймворков).

### Docker

Написан multi-stage `Dockerfile` (сборка на `golang:alpine`, запуск на чистом `alpine`). Настроен `docker-compose.yml` с healthcheck для базы данных и зависимостью запуска приложения от готовности БД.

### Дополнительная защита (TLS)

Сервис поддерживает работу по HTTPS. Если при запуске передать переменные окружения `TLS_CERT_FILE` и `TLS_KEY_FILE` с путями к сертификату и ключу, сервер автоматически переключится на защищенное соединение (используется `http.ListenAndServeTLS`).

## Примечания

- Для простоты миграции применяются вручную через `psql`. В production-окружении следует использовать инструменты миграций (например, `golang-migrate`).
    
- В коде настроено структурированное логирование (`slog`) и graceful shutdown.
    
- Комментарии к коду написаны на русском языке для облегчения понимания.