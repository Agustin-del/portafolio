# AGENTS.md - Portafolio Development Guide

## Project Overview

This is a personal portfolio website built with Go using the [Echo](https://echo.labstack.com/) web framework and [templ](https://templ.guide/) for type-safe HTML templates. The project page is under development, with a planned file-system-like view to navigate local project repos and display file contents read-only.

## Build, Run, and Test Commands

### Working Directory
All commands should be run from the `src/` directory.

### Development
```bash
# Run with hot reload (uses air)
cd src && air

# Or manually:
cd src && go build -o tmp/server main.go && ./tmp/main
```

### Production Build
```bash
cd src && go build -o bin/server main.go
```

### Linting
```bash
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
- **Files**: snake_case (e.g., `main.go`, `db.go`)
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
- Use `errors.Is()` for error checking (see `main.go:151`)
- Log errors with context using `c.Logger().Error("message", "error", err)`
- Handle errors early with early returns

### Echo Framework Patterns
- Define custom HTTP error handler as shown in `main.go:73-88`
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
- Define struct types for data models (see `pages.Proyecto` in `proyectos.templ` for project data)
- Use meaningful field names in Spanish (e.g., `De`, `Asunto`, `Mensaje`)
- Avoid `any`/interface{} when concrete types suffice

### Security
- Never commit secrets (environment variables, API keys)
- Use environment variables for sensitive data (`os.Getenv()`)
- Implement Content-Security-Policy headers (see `main.go:15`)
- Validate all form inputs before processing

### Database
- Database code lives in `db/` package
- Currently uses SQLite (`github.com/mattn/go-sqlite3`)
- Initialize connections in `db.Init()`
- Schema defined in `db/db.go` and `pseudo-schema.sql`

### Static Assets
- Serve static files with `e.Static()` (see `main.go:102-104`)
- Organized in `static/imagenes`, `static/estilos`, `static/scripts`

## Project Structure

```
src/
├── main.go          # Main application entry point
├── main_test.go     # Application tests
├── pseudo-schema.sql # SQL schema reference
├── .air.toml        # Hot reload configuration
├── go.mod
├── go.sum
├── cron/            # Cron job logic (e.g., email sending)
├── data/            # SQLite database storage
├── db/              # Database initialization and queries
│   └── db.go
├── proyectos/       # Local project repositories (cloned recursively)
│   └── <project-name>/
│       ├── descripcion.md
│       ├── descripcion_detallada.md
│       └── src/      # Source code to browse
├── ui/
│   ├── layout.templ   # Base layout template
│   ├── layout_templ.go (generated)
│   └── pages/
│       ├── inicio.templ
│       ├── proyectos.templ  # Project page (skeleton, in development)
│       ├── contacto.templ
│       └── error.templ
├── static/
│   ├── estilos/
│   ├── scripts/
│   └── imagenes/
├── tmp/             # Temporary build artifacts (air output)
└── vendor/         # Vendored dependencies
```

## Environment Variables
Required for email functionality:
- `USUARIO_SMTP` - SMTP username
- `PASS_SMTP` - SMTP password

## Project Page Integration (In Development)

### Project Storage
- Projects are stored as local git repositories, cloned recursively manually when updated.
- Store project repos in `src/proyectos/` directory.
- Each project must have the following structure:
  - `descripcion.md`: Brief description (e.g., "Built with Go using Echo and templ")
  - `descripcion_detallada.md`: Detailed personal description of system design
  - `src/`: Source code folder to be browsable in the file tree view
- Project files are NOT served as static assets (static/ is for global assets like CSS/images)

### Planned View
- Left sidebar: File tree navigation for the selected project (like a file explorer).
- Right panel: Read-only display of the selected file's content, styled like a text editor.

### Caching Strategy
- All project data is loaded into memory at server startup by scanning the `src/proyectos/` directory.
- Data structure: `map[string]map[string][]byte` where:
  - Key (level 1): Project name (e.g., "mi-proyecto")
  - Key (level 2): File path relative to project root (e.g., "descripcion.md", "src/main.go")
  - Value: File content as `[]byte` for direct streaming to HTTP response
- Once loaded, the data resides in memory with no TTL or cache expiration.
- Since projects are updated manually (infrequent changes), the server must be restarted to reflect updates in project repos.
- No mutex needed: data is loaded once at startup before the server starts accepting requests, so there are no concurrent writes.

### Memory Considerations
- To avoid excessive memory usage:
  - Exclude unnecessary directories/files from loading: `.git`, `node_modules`, `vendor`, build artifacts, binary files.
  - Store only text file contents under a reasonable size limit (e.g., 1MB).
  - Monitor memory usage if adding large projects.

## Common Tasks

### Adding a New Page
1. Create new `.templ` file in `ui/pages/`
2. Run `templ generate` in `src/`
3. Add route in `main.go`
4. Add static assets if needed

### Adding a New Project
1. Clone the project repo recursively into the designated projects directory (e.g., `src/proyectos/<project-name>`).
2. Restart the server to reload the in-memory cache, or trigger a manual cache reload.
3. Verify the project appears in the file tree navigation on the `/proyectos` page.

## Important Notes

- After editing any `.templ` file, run:
  ```bash
  templ generate
  ```
- Generated files (`*_templ.go`) must not be edited manually.
- All commands must be run from the `src/` directory.
- Never commit vendored dependencies, local database files, or temporary build artifacts to version control.
