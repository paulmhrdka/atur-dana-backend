<h1 style="text-align: center;">Atur Dana BE</h1>
<br />

**Atur Dana** is a personal financial management app designed to help you manage your personal finances effectively. This repository serves as the backend for **Atur Dana**, written in [Golang](https://go.dev/).

## Features

- Track income and expenses with filtering, pagination, and date ranges
- Transaction summary with totals and category breakdown
- Category management
- Budgeting tools
- Secure user authentication (JWT)
- Structured logging with request ID tracing
- Health check endpoint with DB connectivity status
- Swagger API documentation

## Technologies Used

- **Programming Language:** Golang
- **Database:** PostgreSQL
- **API:** RESTful API
- **Authentication:** JWT
- **Frameworks and Libraries:** gorilla/mux, GORM, swaggo/swag

## Getting Started

### Prerequisites

- Go (version 1.20.1 or higher)
- PostgreSQL
- Git

## Project Structure
```bash
├── bin                 # Compiled binary
├── build               # Docker Compose configuration
├── cmd
│   └── main.go         # Entry point — loads .env, init DB/validator, graceful shutdown
├── docs                # Swagger-generated files + schema.sql + seed.sql
└── internal
    ├── auth            # JWT token generation
    ├── common          # Shared helpers: response, validator, constants, swagger types
    ├── db              # DB init + AutoMigrate
    ├── handlers        # HTTP handlers: auth, transaction, category, health
    ├── metrics         # Prometheus metrics
    ├── middleware      # JWT auth, request ID, structured logger
    ├── models          # GORM models: User, Transaction, Category, Budget
    ├── requests        # Request payload structs with validate tags
    ├── responses       # Response structs + swagger wrapper types
    └── routes          # API route definitions
```

## API Routes

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | No | Register user, returns JWT |
| POST | `/auth/login` | No | Login, returns JWT |
| GET | `/api/transactions` | Bearer | List transactions (filter by date, category, type; paginated) |
| POST | `/api/transactions` | Bearer | Create transaction |
| GET | `/api/transactions/summary` | Bearer | Income/expense totals + category breakdown |
| GET | `/api/categories` | Bearer | List categories (filter active_only) |
| GET | `/health` | No | Health check (DB connectivity) |
| GET | `/swagger/` | No | Swagger UI |

### Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/paulmhrdka/aturdana-backend.git
    cd aturdana-backend
    ```

2. **Set up environment variables:**

    Copy `.env.example` into `.env` and fill in your credentials:

    ```sh
    cp .env.example .env
    # Set: DATABASE_URL, DB_NAME, DB_USER, DB_PASSWORD, JWT_SECRET_KEY
    # Optional: LOG_FORMAT=json (default: text)
    ```

3. **Install dependencies:**

    ```sh
    go mod tidy
    ```

4. **Start the server:**

    ```sh
    make run   # regenerates Swagger docs then starts :8080
    # or
    go run cmd/main.go
    ```

5. **Load seed data (optional):**

    ```sh
    psql -U <user> -d <db> -f docs/seed.sql
    ```

### Docker

```sh
docker network create pfm
docker-compose -f build/docker-compose.yml up -d
```

## Contact

For any questions or feedback, please reach out to:

- **Email:** [mahardikapaul@gmail.com](mailto:mahardikapaul@gmail.com)
- **GitHub Issues:** [GitHub Issues Page](https://github.com/paulmhrdka/atur-dana-backend/issues)
