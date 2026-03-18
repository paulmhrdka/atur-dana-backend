# Task Completion Checklist

When completing a coding task in aturdana-backend:

1. **Run the app** to verify it compiles and starts:
   ```bash
   go run ./cmd/main.go
   ```
   Or just build check:
   ```bash
   go build ./...
   ```

2. **Regenerate Swagger docs** if any handler signatures, routes, or annotations changed:
   ```bash
   make docs
   ```

3. **Format code**:
   ```bash
   gofmt -w .
   ```

4. **Static analysis**:
   ```bash
   go vet ./...
   ```

5. **No test suite** exists in the project currently — no `go test` to run.

6. **If adding a new handler**: ensure Swagger annotations are added, route is registered in `routes/routes.go`, and request/response types are defined.
