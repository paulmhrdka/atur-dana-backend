# Suggested Commands

## Run & Build
```bash
make run           # Generate docs + run server (go run ./cmd/main.go)
make build         # Generate docs + build binary to bin/aturdana
make docs          # Regenerate Swagger docs only
go run ./cmd/main.go   # Run without regenerating docs
```

## Docker
```bash
docker-compose -f build/docker-compose.yml up      # Start app + postgres
docker-compose -f build/docker-compose.yml down     # Stop
docker network create pfm                           # Required once (network is external)
```

## Go Tools
```bash
go mod tidy        # Clean up dependencies
go vet ./...       # Static analysis
gofmt -w .         # Format all files
gofmpt -l .        # List files that need formatting
```

## Swagger
```bash
# Install swag if not present:
go install github.com/swaggo/swag/cmd/swag@latest
$(go env GOPATH)/bin/swag init --generalInfo cmd/main.go --dir . --output docs --parseDependency --parseInternal
```

## Environment Setup
```bash
cp .env.example .env
# Fill in: DATABASE_URL, DB_NAME, DB_USER, DB_PASSWORD, JWT_SECRET_KEY
```

## No test commands found in project (no *_test.go files)
