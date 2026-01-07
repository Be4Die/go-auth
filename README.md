# Go Auth Service

Лёгкий, производственный готовый сервис аутентификации на Go (Gin, PostgreSQL, JWT).

## Возможности
- Регистрация и вход пользователя (bcrypt, JWT access/refresh)
- Чистая архитектура: domain → usecase → transport/infrastructure
- Логирование через `slog`, конфигурация из env
- Миграции через `docker-entrypoint-initdb.d`
- Готовые Dockerfile и docker-compose с healthchecks
- Тесты: unit (password, jwt, usecases), интеграционные (Postgres repo), HTTP-хэндлеры
- CI (GitHub Actions): сборка и прогон тестов, покрытие

## Быстрый старт
```sh
docker-compose up -d
```
Сервис поднимется на `http://localhost:8080`.

Эндпоинты:
- `POST /api/v1/auth/register` `{email, password}` → `201`
- `POST /api/v1/auth/login` `{email, password}` → `200` с токенами
- `GET /health` → `200`

## Конфигурация
См. `.env.example`. Ключевые переменные:
- `HTTP_PORT`, `DATABASE_URL`
- `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`

## Разработка и тесты
```sh
go test ./... -race -cover
```
Интеграционные тесты используют `DATABASE_URL`. В CI Postgres поднимается сервисом.

## Архитектура
- `internal/domain` — сущности и порты
- `internal/app/usecase` — бизнес-кейс регистрации/логина
- `internal/infrastructure/postgres` — репозиторий пользователей
- `internal/security` — `password` (bcrypt) и `jwt`
- `internal/transport/http` — Gin хэндлеры

## Продакшн
- Non-root пользователь в образе, healthcheck
- Надёжные миграции, отказ от DDL в коде
- Логи JSON в `production`, `GIN_MODE=release`

## Планы
- Хранение и ротация refresh-токенов
- Rate limiting, метрики, трейсинг

