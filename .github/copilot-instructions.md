# AI Agent Instructions for Fuego De Quasar Project

## Project Overview
This is a Go-based API service implementing a communication interception and location system based on triangulation principles. The project follows a clean architecture pattern with clear separation of concerns.

## Core Components

### 1. Calculation Engine (`internal/platform/calculos/`)
- `trilateracion.go`: Implements trilateration algorithms for locating signal sources
- `recuperarmensaje.go`: Handles message reconstruction from multiple intercepted signals

### 2. API Layer
- Uses Gin framework for HTTP routing
- Endpoints are defined in `handlers/rutas.go`
- Server configuration in `cmd/main.go`

### 3. Repository Layer (`internal/platform/repository/`)
- Interface-based design for data access abstraction
- Mock-friendly architecture for testing

## Key Patterns and Conventions

### Error Handling
Example from `recuperarmensaje.go`:
```go
if !found {
    return "", errors.New("no se pudo reconstruir ningún mensaje")
}
```

### Type Conversions
The project handles both float32 and float64 precision. See `GetLocation()` in `trilateracion.go` for conversion patterns.

### Configuration
Environment variables are used for configuration (e.g., PORT in `main.go`).

## Development Workflow

### Running the Service
```bash
go run cmd/main.go
```

### Project Structure
```
├── cmd/
│   └── main.go           # Entry point
├── handlers/
│   └── rutas.go          # HTTP routes
└── internal/
    └── platform/
        ├── calculos/     # Core business logic
        └── repository/   # Data access layer
```

## Integration Points
- HTTP API (Gin framework)
- Environment configuration
- Repository interface for data persistence

## Testing
The codebase is designed for testability:
- Repository interface enables mocking
- Pure functions in calculation engine
- Separate business logic from HTTP handlers

## Common Operations
- Adding new API endpoints: Extend `handlers/rutas.go`
- Implementing new calculations: Add to `internal/platform/calculos/`
- Data persistence changes: Implement through repository interface