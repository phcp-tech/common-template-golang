# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go 1.26 module template for creating new Go modules with best practices. This template provides a standard project structure with RESTful API, Swagger documentation, dependency injection, and common Go microservice patterns.

## Common Commands

```bash
# Update dependencies
go mod tidy

# Generate dependency injection code
wire ./pkg/injector

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
  - injector/  # Dependency injection configuration
    - wire.go           # Wire dependency injection definitions
    - injector.go       # Injector setup
    - wire_gen.go       # Generated wire code

docs/         # Documentation
  - docs.go    # Swagger documentation generation

config/       # Configuration files
  - app.toml   # Application configuration
```

## Architecture Highlights

### Standard Go Project Layout
- **Adapter Layer**: Handles HTTP requests/responses, route dispatching, parameter validation
- **Domain Layer**: Defines business domain models
- **Infrastructure Layer**: Data access implementations
- **Service Layer**: Business logic processing
- **Package Layer**: Utility functions, DTOs, and dependency injection

### Key Features
- **RESTful API**: Gin framework-based API with Swagger documentation
- **Dependency Injection**: Google Wire for managing dependencies
- **Swagger Documentation**: Automatic API documentation generation
- **Standard Patterns**: Follows best practices for Go microservices
- **Common Library Integration**: Uses common-library-golang for shared utilities

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
| `main.go` | Application entry point, initializes configuration, dependency injection, database migration, route registration |
| `adapter/user_restful_api.go` | User RESTful API implementation |
| `service/user_service.go` | User business logic |
| `domain/model/user.go` | User domain model |
| `pkg/injector/wire.go` | Dependency injection configuration |
| `docs/docs.go` | Swagger documentation generation |

## Development Notes

1. **Swagger Documentation**: Automatically generated API documentation available at `/swagger/index.html`
2. **Dependency Injection**: Uses Google Wire for clean dependency management
3. **Common Library**: Integrates with common-library-golang for shared utilities
4. **Standard Structure**: Follows industry best practices for Go microservices
5. **Extensible Design**: Easily extendable for new features and services
6. **Configuration Management**: Flexible configuration handling through `config/app.toml`

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/google/wire` - Dependency injection
- `github.com/phcp-tech/common-library-golang` - Common library (requires local replace)
- `gorm.io/gorm` - ORM
- `gorm.io/datatypes` - JSONB data type support
- `github.com/swaggo/files` - Swagger file serving
- `github.com/swaggo/gin-swagger` - Gin Swagger integration
- `github.com/swaggo/swag` - Swagger documentation generation
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/go-resty/resty/v2` - HTTP client
- `gopkg.in/yaml.v3` - YAML configuration parsing