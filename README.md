# Threat Exchange Service

REST-сервис на Go для обмена информацией об угрозах, реализованный на основе предоставленной OpenAPI-спецификации.

## 🛠 Стек технологий
- **Язык:** Go (1.27+)
- **База данных:** PostgreSQL 16
- **Драйвер БД:** `pgx` (pgxpool)
- **Кодогенерация:** `oapi-codegen` (генерация типов и серверного интерфейса по OpenAPI)
- **Инфраструктура:** Docker, Docker Compose

## Быстрый запуск (через Docker)

Весь проект (база данных + приложение) поднимается одной командой. Убедитесь, что у вас установлен Docker и Docker Compose v2.

1. Запустите контейнеры в фоне:
   
```bash
   docker compose -f deployments/docker-compose.yml up --build -d
```

2. Применение миграции
    
```bash
    docker exec -i threat-exchange-db psql -U threat -d threat_exchange < migrations/0001_init_sql.up.sql
```

3. Запуск сервера

```bash
    export DATABASE_URL="postgres://threat:threat@localhost:5432/threat_exchange?sslmode=disable"
    go run ./cmd/server
```