# BabyTracker - Claude Code Context

Multi-platform baby activity tracker: Go API + Fyne desktop + React PWA.

## Commands

```bash
make setup          # First-time setup (env files + deps)
make dev            # API (:8080) + React (:3000) concurrently
make api            # Go API server only
make desktop        # Fyne desktop app only
make web            # React dev server only (needs API running)
make test           # Go tests: go test ./internal/...
make test-cover     # Go tests with coverage
make test-web       # React tests (no test cases yet)
make lint           # go vet ./...
make lint-web       # eslint web/src
make build          # Build all (bin/api, bin/desktop, web/build)
```

## Architecture

```
cmd/desktop/main.go       -> Fyne desktop entry point
cmd/api/main.go           -> HTTP API server entry point
internal/models/          -> Shared data models (feed, sleep, growth, diaper)
internal/storage/         -> Generic JSON persistence (~/.babytracker/)
internal/api/             -> REST handlers + gorilla/mux router
internal/config/          -> Env-based config (PORT, DATA_DIR, APP_TITLE)
internal/desktop/         -> Fyne UI (app.go, layout.go, tabs/)
web/src/                  -> React SPA (CRA, react-router-dom)
```

Data flow: Desktop writes JSON directly; Web/PWA -> API server -> JSON files.

## Key Files

- `internal/storage/storage.go` -- Generic JSON engine using `loadJSON[T any]`/`saveJSON[T any]`
- `internal/api/router.go` -- Route registration + CORS middleware
- `internal/api/handlers.go` -- Feed handlers (pattern template for all modules)
- `web/src/api.js` -- `apiGet`/`apiPost` wrappers around fetch
- `web/src/config.js` -- API base URL from REACT_APP_API_BASE

## Code Conventions

- Go module: `babytracker` (Go 1.24, Fyne v2, gorilla/mux)
- Web: React 18 + CRA, bun for package management
- Storage: JSON files in `~/.babytracker/` (feeds.json, sleep.json, growth.json, diapers.json)
- ID generation: `max(existing IDs) + 1` -- not `len + 1`
- API routes: `/api/{resource}` (GET list, POST create), `/api/{resource}/{id}` (GET by ID)
- CORS: wildcard `*` via middleware in router.go
- Date format: `YYYY-MM-DD` (string), Time: `time.Time` (Go) / ISO string (JSON)
- Handlers follow pattern: decode -> validate required fields -> log -> save -> respond

## Gotchas

- `web/src/config.js` line 1 can have a spurious path prefix causing eslint parse errors -- it's real content, not a display artifact
- CRA requires `"eslintConfig": { "extends": ["react-app"] }` in web/package.json for JSX parsing
- macOS Fyne builds emit `ld: warning: ignoring duplicate libraries: '-lobjc'` -- harmless, ignore it
- gorilla/mux OPTIONS requests bypass middleware unless the route explicitly allows OPTIONS method (handled in router.go CORS middleware)
- Desktop app uses `storage` package directly (no HTTP); web app goes through the API
- `layout.go` duplicates tab creation from `app.go` -- `CreateMainLayout()` exists but `App.CreateMainContent()` is actually used

## Environment

- `.env` at root: `PORT` (API, default 8080), `DATA_DIR`, `APP_TITLE`
- `web/.env`: React env vars (`REACT_APP_API_BASE`)
- Run `make env` to create from `.env.example` templates
- Makefile loads root `.env` automatically via `-include .env` + `export`

## Testing

- Go tests: `go test ./internal/...` (models, storage, API handlers)
- Storage tests use temp dirs via `t.TempDir()`
- API handler tests use `httptest.NewRecorder()`
- Web tests: placeholder only (`make test-web`), no test cases written yet
