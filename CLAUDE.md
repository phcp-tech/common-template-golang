# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Overview

Go 1.26 module template for creating new Go modules with best practices: a single User entity showing the standard adapter → service → infra/dao → domain/model layering, Gin, PostgreSQL via `dbsqlx` — raw SQL, no ORM. No `application.go`; `main.go` is the composition root. No migration step (no `migrate()` in `main.go`) — the schema is managed independently.

## Commands

```bash
go build
go test ./...
swag init   # regenerate Swagger docs
```

## Non-obvious things

- **No ORM.** `infra/dao/user_dao_impl.go` writes raw SQL against `*sqlx.DB`, using `db.Rebind()` and `dbsqlx.SortSql`/`PageSql` — copy this DAO's shape when adding a new entity to a project generated from this template.
- **No auth.** The template's only endpoint (`GET /usrapi/v1/users/list`) is public, to keep the reference example minimal — real projects generated from this template typically add `token.Authenticate()`/`auth.Authorize()` per the pattern used in the other microservices.
- **`config/mock_data.sql` doesn't exist** — unlike some of the other services in this workspace, this template has no seed-data file; `infra/dao/user_dao_impl_test.go` seeds its own in-memory SQLite rows directly.
