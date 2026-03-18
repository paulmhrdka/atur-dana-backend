# Code Style & Conventions

## Naming
- **Packages**: lowercase single word (handlers, routes, middleware, common, models, requests, responses, auth, db)
- **Exported symbols**: PascalCase (Register, Login, JSONResponse)
- **Unexported**: camelCase
- **Constants**: PascalCase (Income, Expense)

## Handler Pattern
Every handler follows this structure:
1. Decode JSON body into request struct
2. `common.ValidateRequest(request)` → `common.JSONValidationError` on error
3. Business logic + DB calls
4. `common.JSONResponse(w, statusCode, data, message)` or `common.JSONError(w, statusCode, message)`

## Swagger Annotations
All handlers have godoc Swagger annotations:
```go
// FuncName godoc
// @Summary      ...
// @Description  ...
// @Tags         TagName
// @Accept       json
// @Produce      json
// @Param        body body requests.XxxRequest true "..."
// @Success      200 {object} responses.SwaggerXxxResponse
// @Failure      400 {object} common.SwaggerValidationErrorResponse
// @Failure      500 {object} common.SwaggerErrorResponse
// @Router       /path [method]
```

## Request Structs
Use struct tags: `json:"field" validate:"required,min=3"`

## Response Format
```json
{ "status": 200, "message": "...", "data": {...} }
// Error:
{ "status": "error", "errors": { "Field": "message" } }
```

## GORM Models
Embed `gorm.Model` for id/created_at/updated_at/deleted_at. Use `gorm:"..."` tags.

## No test files in project currently.
