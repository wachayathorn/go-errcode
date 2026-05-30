# errcode

`errcode` is a minimal Go package for structured application errors in HTTP APIs.

It provides:

- Stable error codes for frontend/client handling
- HTTP status codes for response mapping
- Safe client-facing messages
- Internal root cause preservation with `Unwrap`
- Compatibility with `errors.Is` and `errors.As`

## Installation

```sh
go get github.com/wachayathorn/go-errcode
```

Import the package:

```go
import "github.com/wachayathorn/go-errcode"
```

## Core Types

```go
type Code string

type AppError struct {
	StatusCode int    `json:"status_code" example:"500"`
	ErrorCode  Code   `json:"error_code" example:"INTERNAL_SERVER_ERROR"`
	Message    string `json:"message" example:"Something went wrong"`
	Cause      error  `json:"-"`
}
```

`Cause` is excluded from JSON responses to avoid leaking internal implementation details to clients.

## Available Errors

```go
errcode.BadRequest
errcode.Unauthorized
errcode.Forbidden
errcode.NotFound
errcode.Duplicate
errcode.AlreadyExists
errcode.TooManyRequests
errcode.InvalidTokenError
errcode.InternalServerError
```

## Basic Usage

Return predefined errors directly:

```go
return errcode.NotFound
```

Customize the message:

```go
return errcode.NotFound.WithMessage("product not found")
```

Preserve an internal cause:

```go
return errcode.InternalServerError.WithCause(err)
```

Customize the message and preserve the cause:

```go
return errcode.InternalServerError.
	WithMessage("failed to create product").
	WithCause(err)
```

## Service Layer Example

```go
func (s *ProductService) GetProduct(ctx context.Context, id string) (Product, error) {
	if id == "" {
		return Product{}, errcode.BadRequest.WithMessage("id is required")
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}
```

## Repository Layer Example

```go
func (r *ProductRepository) GetByID(ctx context.Context, id string) (Product, error) {
	var product Product

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, price
		FROM products
		WHERE id = $1
	`, id).Scan(&product.ID, &product.Name, &product.Price)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, errcode.NotFound.WithMessage("product not found")
		}

		return Product{}, errcode.InternalServerError.WithCause(err)
	}

	return product, nil
}
```

## HTTP Handler Example

```go
func writeError(w http.ResponseWriter, err error) {
	var appErr *errcode.AppError
	if errors.As(err, &appErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Status())

		_ = json.NewEncoder(w).Encode(appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(errcode.InternalServerError)
}
```

Example response:

```json
{
  "status_code": 404,
  "error_code": "NOT_FOUND",
  "message": "product not found"
}
```

## Using `errors.As`

Use `errors.As` to extract `*AppError` from an error chain:

```go
var appErr *errcode.AppError
if errors.As(err, &appErr) {
	fmt.Println(appErr.Status())
	fmt.Println(appErr.ErrCode())
}
```

## Using `errors.Is`

Use `errors.Is` to check whether an error belongs to a specific app error code:

```go
err := errcode.NotFound.WithMessage("product not found")

if errors.Is(err, errcode.NotFound) {
	fmt.Println("handle not found")
}
```

This works because `AppError` implements custom matching by `ErrorCode`:

```go
func (e *AppError) Is(target error) bool {
	targetErr, ok := target.(*AppError)
	if !ok {
		return false
	}

	return e.ErrorCode == targetErr.ErrorCode
}
```

**Note**: The `Is()` method matches errors by `ErrorCode` only. Two `AppError` instances with the same `ErrorCode` will match via `errors.Is()`, even if they have different `StatusCode` or `Message`. This design allows error category matching regardless of customization.

Without this method, `errors.Is(err, errcode.NotFound)` would compare pointers. A new error created by `WithMessage` would not equal the global sentinel error.

## Root Cause Handling

Use `WithCause` when preserving an internal error:

```go
err := errcode.InternalServerError.WithCause(sql.ErrConnDone)
```

Then you can still check the root cause:

```go
if errors.Is(err, sql.ErrConnDone) {
	fmt.Println("database connection is done")
}
```

This works because `AppError` implements `Unwrap`:

```go
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}
```

## Best Practices

- Use `WithMessage` for safe client-facing messages.
- Use `WithCause` for internal errors that should be logged but not exposed.
- Do not expose raw database/cache/internal errors to clients.
- Use `errors.As` in HTTP handlers to map errors to responses.
- Use `errors.Is` when checking error categories.
- Do not mutate predefined global errors directly.

## Testing

Run the test suite:

```sh
go test ./...
```

## Avoid

Do not mutate global errors:

```go
errcode.NotFound.Message = "custom message"
```

Use this instead:

```go
return errcode.NotFound.WithMessage("custom message")
```

Do not expose internal causes in API responses:

```go
return errcode.InternalServerError.WithMessage(err.Error())
```

Use this instead:

```go
return errcode.InternalServerError.WithCause(err)
```

## Quick Pattern

```go
if invalidInput {
	return errcode.BadRequest.WithMessage("name is required")
}

if notFound {
	return errcode.NotFound.WithMessage("product not found")
}

if err != nil {
	return errcode.InternalServerError.WithCause(err)
}
```
