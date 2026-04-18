# Polling Service Backend Architecture

This repository now includes a clean backend baseline in Go.

## Structure

- `cmd/api`: application entrypoint
- `internal/app`: dependency wiring and startup orchestration
- `internal/config`: environment-based runtime configuration
- `internal/server/http`: HTTP transport and route registration
- `internal/handler`: HTTP handlers (delivery layer)
- `internal/service`: business logic (use-case layer)
- `internal/domain`: core domain models

## Run

```bash
go run ./cmd/api
```

Default server port: `8080`

## Health Endpoint

```bash
curl http://localhost:8080/health
```

Example response:

```json
{
  "status": "ok",
  "service": "pollingservice",
  "env": "development",
  "timestamp": "2026-04-09T00:00:00Z"
}
```

## Environment Variables

- `APP_ENV` (default: `development`)
- `HTTP_PORT` (default: `8080`)

## Suggested Next Layers

1. Add database adapter in `internal/repository`.
2. Add domain-specific use cases in `internal/service`.
3. Add middleware for logging, request IDs, auth.
4. Add graceful shutdown and context-aware startup.
5. Add tests for handlers and services.
