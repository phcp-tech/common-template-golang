# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go 1.26 module template for creating new Go modules with best practices. This template provides a standard project structure with RESTful API, Swagger documentation, and common Go microservice patterns.

## Common Commands

```bash
# Update dependencies
go mod tidy

# Generate Swagger documentation
swag init

# Build the service
go build

# Run the service
./template
```

## Directory Structure

```
adapter/      # RESTful API routing layer
  - swagger_restful_api.go  # Swagger API
  - user_restful_api.go    # User API
  
domain/       # Domain layer
  model/      # Domain models
    - user.go           # User model

infra/        # Infrastructure layer (data access implementation)
  dao/        # DAO implementations
    - user_dao_impl.go   # User DAO implementation
    - user_dao_interface.go # User DAO interface

service/      # Business logic layer
  - user_service.go     # User business logic

pkg/          # Utility packages and DTOs
  - dto/       # Data transfer objects
    - user_dto.go       # User DTO

docs/         # Documentation
  - docs.go    # Swagger documentation generation

config/       # Configuration files
  - app.toml   # Application configuration
```

## Layer Architecture & Dependency Rules

```
┌──────────────────────────────────────────────────────────────┐
│  main.go + application.go   (composition root)               │
│  Wires all layers; registers pkg/metrics snapshot functions  │
└────────────────────────┬─────────────────────────────────────┘
                         │ depends on all layers
┌────────────────────────▼─────────────────────────────────────┐
│  adapter/               REST API (Gin) + CLI (Cobra)          │
│  Thin: validates input, calls service, returns HTTP response  │
└────────────────────────┬─────────────────────────────────────┘
                         │ depends on service only
┌────────────────────────▼─────────────────────────────────────┐
│  service/               Business workflow 
└──────────┬──────────────────────────┬────────────────────────┘
           │ depends on               │ depends on
┌──────────▼──────────┐  ┌───────────▼────────────────────────┐
│  domain/model/      │◄─┤  infra/                            │
│  Pure filter logic  │  │  dao/ (MySQL, ClickHouse, in-mem)  │
│  No infra imports   │  │  redis/, kafka/                    │
└──────────┬──────────┘  └───────────┬────────────────────────┘
           │                         │
           └────────────┬────────────┘
                        │ all layers depend on
┌───────────────────────▼──────────────────────────────────────┐
│  pkg/                 Shared utilities — zero domain imports │
│  metrics, global, dto, util                                  │
└──────────────────────────────────────────────────────────────┘
```

### Dependency Rules

| Layer | May import | Must NOT import |
|---|---|---|
| `adapter` | `service`, `pkg` | `domain`, `infra` directly |
| `service` | `domain`, `infra`, `pkg` | `adapter` |
| `domain/model` | `pkg` only | `infra`, `service`, `adapter` |
| `infra` | `domain` (interfaces), `pkg` | `service`, `adapter` |
| `pkg` | other `pkg/` sub-packages, external libs | any project layer above |

**Key rule**: `pkg/` is a pure utility layer. It cannot import `infra` or `domain`. When `pkg/metrics` needs DAO snapshot data, `application.go` registers function variables (`SavedTicksFn`, `SymbolsFn`, etc.) at startup — this avoids a `pkg → infra` upward dependency.

## Architecture Highlights

### Standard Go Project Layout
- **Adapter Layer**: Handles HTTP requests/responses, route dispatching, parameter validation
- **Domain Layer**: Defines business domain models
- **Infrastructure Layer**: Data access implementations
- **Service Layer**: Business logic processing
- **Package Layer**: Utility functions and DTOs

### Key Features
- **RESTful API**: Gin framework-based API with Swagger documentation
- **Swagger Documentation**: Automatic API documentation generation
- **Standard Patterns**: Follows best practices for Go microservices
- **Common Library Integration**: Uses common-library-golang-internal for shared utilities

### Data Models
Core entities:
- `User`: User model with basic user information

### API Routes
```
/api/v1/user/*    # User operations
```

## Key Files

| File | Purpose |
|------|---------|
| `main.go` | Application entry point: initializes config, creates Application, starts, waits for exit |
| `application.go` | Composition root: wires infrastructures, DAOs, and services |
| `adapter/user_restful_api.go` | User RESTful API implementation |
| `service/user_service.go` | User business logic |
| `domain/model/user.go` | User domain model |
| `docs/docs.go` | Swagger documentation generation |

## Development Notes

1. **Swagger Documentation**: Automatically generated API documentation available at `/swagger/index.html`
2. **Common Library**: Integrates with common-library-golang-internal for shared utilities
3. **Standard Structure**: Follows industry best practices for Go microservices
4. **Extensible Design**: Easily extendable for new features and services
5. **Configuration Management**: Flexible configuration handling through `config/app.toml`

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/phcp-tech/common-library-golang-internal` - Common library (requires local replace)
- `gorm.io/gorm` - ORM
- `gorm.io/datatypes` - JSONB data type support
- `github.com/swaggo/files` - Swagger file serving
- `github.com/swaggo/gin-swagger` - Gin Swagger integration
- `github.com/swaggo/swag` - Swagger documentation generation
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/go-resty/resty/v2` - HTTP client
- `gopkg.in/yaml.v3` - YAML configuration parsing