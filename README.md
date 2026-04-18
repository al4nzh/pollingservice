# Polling Service Backend 


## Structure

- `cmd/api`: application entrypoint
- `internal/app`: dependency wiring and startup orchestration
- `internal/config`: environment-based runtime configuration
- `internal/server/http`: HTTP transport and route registration
- `internal/handler`: HTTP handlers (delivery layer)
- `internal/service`: business logic (use-case layer)
- `internal/domain`: core domain models



