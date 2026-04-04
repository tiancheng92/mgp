# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Test
go test ./...
go test -run TestName ./...

# Lint (if golangci-lint is available)
golangci-lint run
```

## Architecture Overview

`mgp` is a Go library that wraps [gin](https://github.com/gin-gonic/gin) to add:
1. **Structured response rendering** with a unified `Result[T]` envelope
2. **Swagger/goswag code generation** from route metadata
3. **Error code system** with HTTP status mapping

### Core Types

- **`Engine`** (`engine.go`) — wraps `*gin.Engine`, collects routes/groups for swagger generation. `Handle()` returns a `Swagger` interface for chaining metadata.
- **`RouterGroup`** (`router_group.go`) — wraps `*gin.RouterGroup`, adds default swagger metadata (tags, accepts, produces, auth) that apply to all child routes.
- **`Context`** (`context.go`) — wraps `*gin.Context`. Handler helpers:
  - `HR(f)` — standard HTTP handler (func variations: `func()`, `func() error`, `func() (any, error)`)
  - `HD(f, filename)` — file download
  - `HP(f)` — reverse proxy
  - `HW(f)` — WebSocket
  - `HWP(f)` — WebSocket proxy
  - `BindBody/BindQuery/BindParams/BindHeader` — bind & validate, abort with 400 on failure (chainable)
  - `BindPaginateQuery` — parses `page`, `page_size`, `order`, `search` query params
- **`Route`** (`route.go`) — holds swagger metadata per route; implements `Swagger` interface
- **`Swagger`** (`swagger.go`) — interface for chaining swagger metadata onto routes

### Response Format

All responses use `Result[T]{Code, Msg, Data}`. Success: `code=200000`, `msg="Success"`. Errors use `errors.Coder` (HTTP status + business code + message).

Paginated responses wrap data in `ResultPaginateData{Items, Paginate}` — detected automatically when `data` implements `PaginateInterface`.

### Error System (`errors/` package)

Extends `github.com/pkg/errors` with an error code registry:
- `errors.Register(key, code, httpStatus, message)` — register a named error code
- `errors.WithCode(key, err/message)` — create a coded error
- `errors.ParseCoder(err)` — extract `Coder` from error for response rendering
- Default codes in `errors/default_error_code/`: `ErrClientParam` (400000/HTTP 400), `ErrServer` (500000/HTTP 500)

### Swagger Generation

`engine.GenerateSwagger()` / `GenerateSwagger(routes, groups, defaultResponses)` generates a `goswag.go` file (path configurable via `SetSwagFileName`). It inspects struct fields via reflection to produce swagger annotations.

Route swagger metadata is set via method chaining on the returned `Swagger` interface:
```go
engine.Handle("GET", "/users", handler).
    SwaggerSummary("List users").
    SwaggerTags("users").
    SwaggerQuery(UserQuery{}).
    SwaggerReturns(RT[[]User]())
```

### Return Type Helpers

- `RT[T](statusCode...)` — returns `*ReturnType` with `Result[T]` body
- `PRT[T](statusCode...)` — returns `*ReturnType` with paginated `Result[PaginateData[[]T]]` body
- `RTE(statusCode...)` — returns `*ReturnType` with no body (empty response)

### Configurable Defaults

`default_variable.go` exposes package-level setters:
- `SetSuccessMsg(msg)` / `SetSuccessCode(code)` — override success response values
- `SetSwagFileName(path)` — override generated swagger file path (default: `./goswag.go`)
