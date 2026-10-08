TodoList API (Go)

A RESTful task-management API written in Go, with JWT authentication, role-based authorization, two-factor authentication (TOTP), and an OpenAPI-documented interface. Built as a layered service with dependency injection to keep every component independently testable.

Highlights
- Layered architecture: handlers → services → repositories → models, with clear boundaries between HTTP, business logic, and persistence.
- Dependency injection with uber-go/dig, so components are wired in one place and easy to swap or mock in tests.
- Authentication & authorization: JWT-based auth, custom Gin middleware, and role-based route protection (user/admin).
- Account security: email confirmation, configurable password rules, password reset by token, TOTP two-factor authentication with QR-code setup, and account lockout controls.
- API-first documentation: OpenAPI spec served alongside an interactive Swagger UI; oapi-codegen is configured for generating types and Gin server code from the spec.
- Structured logging via the standard library's log/slog (JSON output).
- Zero-setup persistence: GORM with a pure-Go SQLite driver (no CGO required).


## Project Structure

```
.
├── main.go            # Composition root: DI container, middleware, route registration
├── handlers/          # HTTP layer: request parsing, validation, responses
├── services/          # Business logic and middleware (authN/authZ), DI container setup
├── repositories/      # Data access and utilities (token manager, validators, access rules)
├── models/            # Domain models and roles
├── Templates/         # Templates (e.g. email content)
├── docs/              # OpenAPI specification (openapi.yml)
└── codegen.yaml       # oapi-codegen configuration
```

## Roadmap (In progress)

- [ ] Unit tests for services and repositories using injected fakes
- [ ] Handler-level tests with `httptest`
- [ ] Dockerfile and CI pipeline (GitHub Actions: build, vet, test)
- [ ] Pagination for task listing
- [ ] Environment-based configuration (port, JWT secret, database path)
