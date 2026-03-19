# Atur Dana - Project Overview

## Purpose
Personal finance management REST API for tracking income/expense transactions and budgets. Built as a backend service with JWT authentication.

## Tech Stack
- **Language**: Go 1.20
- **Module**: `atur-dana`
- **Router**: gorilla/mux v1.8.1
- **ORM**: GORM v1.25.10 with PostgreSQL (pgx v5.5.5)
- **Auth**: golang-jwt/jwt v5.2.1 + bcrypt
- **Validation**: go-playground/validator v10.20.0
- **Docs**: swaggo/swag (Swagger 2.0)
- **Config**: godotenv v1.5.1

## Architecture
```
cmd/main.go                  # Entry point - loads env, init DB, graceful shutdown :8080
internal/
  models/models.go           # GORM models: User, Transaction, Category, Budget
  handlers/                  # HTTP handlers: auth.go, transaction.go, category.go, health.go
  routes/routes.go           # Route registration (mux subrouters)
  middleware/                # jwt.go, logger.go (slog + request ID), request_id.go
  metrics/metrics.go         # Prometheus metrics
  requests/                  # Request structs with validate tags
  responses/                 # Response structs + swagger types
  common/                    # Shared: response.go, validator.go, constant.go, swagger_types.go
  auth/jwt.go                # JWT generation helper
  db/db.go                   # DB init, AutoMigrate
docs/                        # Swagger-generated files + schema.sql + seed.sql
build/docker-compose.yml     # Docker Compose (app + postgres)
Makefile                     # make docs / build / run
```

## Data Models
- **User**: id, username (unique), password_hash, email (unique)
- **Transaction**: id, user_id, type (income/expense), amount, description, category_id, date
- **Category**: id, name (unique), is_active
- **Budget**: id, user_id, category_id, amount, start_date, end_date

## API Routes
- `POST /auth/register` — register + return JWT
- `POST /auth/login` — login + return JWT
- `GET /api/transactions` — list with filter (date range, category, type) + pagination
- `POST /api/transactions` — create transaction
- `GET /api/transactions/summary` — income/expense totals + category breakdown
- `GET /api/categories` — list categories (filter: active_only)
- `GET/POST /api/budgets` — protected
- `GET /health` — health check with DB connectivity
- `GET /swagger/` — Swagger UI

## Observability
- Structured logging via `log/slog` (text or JSON via LOG_FORMAT env var)
- Request ID middleware (X-Request-ID header)
- Prometheus metrics skeleton (`internal/metrics/metrics.go`)
- Graceful shutdown on SIGINT/SIGTERM with 15s timeout
- HTTP server timeouts: ReadTimeout=15s, WriteTimeout=15s
