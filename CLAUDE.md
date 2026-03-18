# Atur Dana Backend — Claude Code Guide

## Project Overview

Personal finance management REST API (Go) for tracking income/expense transactions and budgets. Uses JWT authentication and PostgreSQL.

## Tech Stack

| Concern | Library |
|---|---|
| Router | gorilla/mux |
| ORM | GORM + pgx (PostgreSQL) |
| Auth | golang-jwt/jwt v5 + bcrypt |
| Validation | go-playground/validator v10 |
| API Docs | swaggo/swag (Swagger 2.0) |
| Config | godotenv |

## Project Structure

```
cmd/main.go                  # Entry point — loads .env, init DB/validator, start :8080
internal/
  models/models.go           # GORM models: User, Transaction, Category, Budget
  handlers/                  # HTTP handlers: auth, transaction, budget
  routes/routes.go           # Route registration
  middleware/jwt.go          # JWT middleware for protected routes (/api/*)
  requests/                  # Request structs with validate tags
  responses/                 # Response structs + swagger wrapper types
  common/                    # Shared helpers: response, validator, constants, swagger types
  auth/jwt.go                # JWT token generation
  db/db.go                   # DB init + AutoMigrate
docs/                        # Swagger-generated (do not edit manually)
build/docker-compose.yml     # App + Postgres containers
Makefile                     # make docs / build / run
```

## Key Commands

```bash
make run          # Regenerate docs + run server
make build        # Regenerate docs + build to bin/aturdana
make docs         # Regenerate Swagger docs only

go build ./...    # Quick compile check
go vet ./...      # Static analysis
gofmt -w .        # Format all files
```

## Environment Setup

```bash
cp .env.example .env
# Set: DATABASE_URL, DB_NAME, DB_USER, DB_PASSWORD, JWT_SECRET_KEY
```

## Docker

```bash
docker network create pfm                          # One-time setup
docker-compose -f build/docker-compose.yml up -d  # Start app + postgres
```

## API Routes

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | No | Register user, returns JWT |
| POST | `/auth/login` | No | Login, returns JWT |
| GET | `/api/transactions` | Bearer | List transactions |
| POST | `/api/transactions` | Bearer | Create transaction |
| GET | `/api/budgets` | Bearer | List budgets |
| POST | `/api/budgets` | Bearer | Create budget |
| GET | `/swagger/` | No | Swagger UI |

## Code Conventions

### Handler Pattern
Every handler follows this exact flow:
1. `json.NewDecoder(r.Body).Decode(&request)`
2. `common.ValidateRequest(request)` → `common.JSONValidationError` on error
3. Business logic + DB calls via `db.DB`
4. `common.JSONResponse(w, status, data, message)` or `common.JSONError(w, status, message)`

### Swagger Annotations (required on every handler)
```go
// FuncName godoc
// @Summary      Short summary
// @Description  Longer description
// @Tags         TagName
// @Accept       json
// @Produce      json
// @Param        body body requests.XxxRequest true "Payload description"
// @Success      200 {object} responses.SwaggerXxxResponse
// @Failure      400 {object} common.SwaggerValidationErrorResponse
// @Failure      500 {object} common.SwaggerErrorResponse
// @Router       /path [method]
```

### Response Format
```json
// Success
{ "status": 200, "message": "...", "data": { ... } }

// Validation error
{ "status": "error", "errors": { "FieldName": "error message" } }
```

### Request Structs
```go
type CreateXxxRequest struct {
    Field string `json:"field" validate:"required,min=3"`
}
```

## Task Completion Checklist

After any change:
- [ ] `go build ./...` — ensure it compiles
- [ ] `gofmt -w .` — format code
- [ ] `go vet ./...` — static analysis
- [ ] `make docs` — if any handler/route/annotation changed
- [ ] New handlers must have Swagger annotations + route registered in `routes/routes.go`
