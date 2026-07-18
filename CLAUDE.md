# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Overview

Go 1.26 module template for creating new Go modules with best practices: a single User entity showing the standard adapter → service → infra/dao → domain/model layering, Gin, SQLite via `dbsqlx` — raw SQL. `main.go` is the composition root, and the schema lives in `config/schema_sqlite.sql`, managed independently.

## Commands

```bash
go build
go test ./...
swag init   # regenerate Swagger docs
```

## Non-obvious things

- `infra/dao/user_dao_impl.go` writes raw SQL against `*sqlx.DB`, using `db.Rebind()` and `dbsqlx.SortSql`/`PageSql` — copy this DAO's shape when adding a new entity to a project generated from this template.
- The template's only endpoint (`GET /usrapi/v1/users/list`) is public, to keep the reference example minimal.
- `config/schema_sqlite.sql` holds both the `users` table DDL and 2 seed rows (Tom/Jerry) — it isn't applied automatically, so it must be run manually against `data/sqlite.db` before starting the app. `infra/dao/user_dao_impl_test.go` doesn't use this file; it seeds its own in-memory SQLite rows directly, against a matching schema.
