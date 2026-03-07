# AGENTS.md - Portafolio Development Guide

## Project Overview

This is a personal portfolio website built with Go using the [Echo](https://echo.labstack.com/) web framework and [templ](https://templ.guide/) for type-safe HTML templates.

## Build, Run, and Test Commands

### Working Directory
All commands should be run from the `src/` directory.

### Development
```bash
# Run with hot reload (uses air)
cd src && air

# Or manually:
cd src && go build -o tmp/main server.go && ./tmp/main
```

### Production Build
```bash
cd src && go build -o bin/server server.go
```

### Testing
```bash
# Run all tests
cd src && go test ./...

# Run a single test (use -run flag with test name pattern)
cd src && go test -run TestName ./...

# Run tests with verbose output
cd src && go test -v ./...

# Run tests with coverage
cd src && go test -cover ./...
```

### Linting
```bash
# Run golangci-lint (if installed)
cd src && golangci-lint run

# Run go vet
cd src && go vet ./...

# Format code
cd src && go fmt ./...
```

### Dependencies
```bash
# Install dependencies
cd src && go mod download

# Add a new dependency
cd src && go get <package>

# Tidy go.mod
cd src && go mod tidy
```

## Code Style Guidelines

### General
- Use 2 spaces for indentation (not tabs)
- Keep lines under 100 characters when practical
- Add TODO comments for incomplete code (e.g., `//TODO: implement X`)

### Naming Conventions
- **Files**: snake_case (e.g., `server.go`, `db.go`)
- **Functions**: PascalCase for exported, camelCase for unexported (e.g., `render()`, `enviarEmail()`)
- **Variables**: camelCase (e.g., `de`, `pass`, `auth`)
- **Constants**: PascalCase for exported, camelCase for unexported (e.g., `csp`)
- **Packages**: lowercase, short names (e.g., `db`, `ui`)

### Imports
- Group imports in the following order (blank line between groups):
  1. Standard library
  2. Third-party packages
  3. Internal/packages
- Use import aliases only when necessary
- Example:
  ```go
  import (
      "errors"
      "fmt"
      "net/http"
      "net/smtp"
      "os"

      "github.com/a-h/templ"
      "github.com/labstack/echo/v4"
      "github.com/labstack/echo/v4/middleware"

      "portafolio/ui/pages"
  )
  ```

### Error Handling
- Return errors explicitly; avoid using `_` for error values in production code
- Use `errors.Is()` for error checking (see `server.go:151`)
- Log errors with context using `c.Logger().Error("message", "error", err)`
- Handle errors early with early returns

### Echo Framework Patterns
- Define custom HTTP error handler as shown in `server.go:73-88`
- Use middleware for cross-cutting concerns (logging, recovery, CSP)
- Return `nil` from handlers when response is already written via `render()`
- Check for `Hx-Request` header for HTMX support

### Templ Templates
- Template files use `.templ` extension
- Generated Go files have `_templ.go` suffix (do not edit)
- Run `templ generate` after modifying `.templ` files
- Keep templates in `ui/pages/` for pages, `ui/` for components
- Use consistent indentation within templ files

### Types
- Define struct types for data models (see `pages.Mail` usage)
- Use meaningful field names in Spanish (e.g., `De`, `Asunto`, `Mensaje`)
- Avoid `any`/interface{} when concrete types suffice

### Security
- Never commit secrets (environment variables, API keys)
- Use environment variables for sensitive data (`os.Getenv()`)
- Implement Content-Security-Policy headers (see `server.go:15`)
- Validate all form inputs before processing

### Database
- Database code lives in `db/` package
- Currently uses SQLite (`github.com/mattn/go-sqlite3`)
- Initialize connections in `db.Init()`

### Static Assets
- Serve static files with `e.Static()` (see `server.go:102-104`)
- Organized in `static/imagenes`, `static/estilos`, `static/scripts`

## Project Structure

```
src/
├── server.go          # Main application entry point
├── db/
│   └── db.go          # Database initialization
├── ui/
│   ├── layout.templ   # Base layout template
│   ├── layout_templ.go (generated)
│   └── pages/
│       ├── inicio.templ
│       ├── proyectos.templ
│       ├── contacto.templ
│       └── error.templ
├── static/
│   ├── estilos/
│   ├── scripts/
│   └── imagenes/
├── go.mod
└── .air.toml          # Hot reload configuration
```

## Environment Variables
Required for email functionality:
- `USUARIO_SMTP` - SMTP username
- `PASS_SMTP` - SMTP password

## Common Tasks

### Adding a New Page
1. Create new `.templ` file in `ui/pages/`
2. Run `templ generate` in `src/`
3. Add route in `server.go`
4. Add static assets if needed

### Running a Specific Test
```bash
cd src && go test -v -run TestFunctionName ./package/path
```
