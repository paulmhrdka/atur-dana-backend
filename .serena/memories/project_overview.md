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
cmd/main.go                  # Entry point - loads env, init DB, start server :8080
internal/
  models/models.go           # GORM models: User, Transaction, Category, Budget
  handlers/                  # HTTP handlers: auth.go, transaction.go, budget.go
  routes/routes.go           # Route registration (mux subrouters)
  middleware/jwt.go          # JWT auth middleware for protected routes
  requests/                  # Request structs with validate tags
  responses/                 # Response structs + swagger types
  common/                    # Shared: response.go, validator.go, constant.go, swagger_types.go
  auth/jwt.go                # JWT generation helper
  db/db.go                   # DB init, AutoMigrate
docs/                        # Swagger-generated files (swagger.json, swagger.yaml, docs.go)
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
- `GET/POST /api/transactions` — protected
- `GET/POST /api/budgets` — protected
- `GET /swagger/` — Swagger UI
