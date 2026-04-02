# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Multi-platform baby tracker: Go API + Fyne desktop + React PWA.
Full docs: docs/Manual.md, docs/man.md

## Commands

```bash
make setup          # First-time setup (env files + deps)
make dev            # API (:8080) + Web (:3000) concurrently
make api            # Go API server only
make desktop        # Fyne desktop app only
make web            # Vite dev server only (needs API running)
make test           # Go tests: go test ./internal/...
make test-web       # Web tests (vitest)
make test-cover     # Go tests with coverage report
make test-all       # All tests (Go + coverage + web)
make lint           # go vet ./...
make lint-web       # bun run lint (web/src)
make build          # Build all (bin/api, bin/desktop, web/build)
```

Single test: `go test -run TestFeedEntry ./internal/models/`

## Architecture

```
cmd/api/main.go           -> HTTP API entry point
cmd/desktop/main.go       -> Fyne desktop entry point
internal/models/          -> Domain models (feed, sleep, growth, diaper)
internal/storage/         -> Generic JSON engine (Go generics: loadJSON[T], saveJSON[T])
internal/api/             -> REST handlers + gorilla/mux router
internal/config/          -> Env config (PORT, DATA_DIR, APP_TITLE, API_KEY, CORS_ORIGIN)
internal/desktop/         -> Fyne UI (app.go, layout.go, tabs/)
web/src/                  -> React SPA (Vite + react-router-dom)
```

Desktop writes JSON directly; Web/PWA -> API -> JSON files in ~/.babytracker/

## Key Conventions

- Go module: `babytracker` (Go 1.24.4, Fyne v2.6.1, gorilla/mux v1.8.1)
- Storage: ~/.babytracker/*.json, IDs = max(existing) + 1, mutex-protected, atomic writes
- API: /api/{resource} (GET/POST), /api/{resource}/{id} (GET)
- Auth: Bearer token via API_KEY env var (empty = no auth)
- CORS: configurable origin via CORS_ORIGIN (default: http://localhost:3000)
- Body limit: 1MB max via http.MaxBytesReader middleware
- Handlers: decode -> validate -> log -> save -> respond
- Web: React 18, Vite, bun (not npm), vitest for testing

## Gotchas

- macOS Fyne: `ld: warning: ignoring duplicate libraries: '-lobjc'` is harmless
- gorilla/mux OPTIONS bypasses middleware unless route allows OPTIONS (handled in router.go)
- layout.go duplicates tab creation from app.go — CreateMainLayout() exists but App.CreateMainContent() is used
- Component files use .jsx extension (Vite requires this for JSX)

## Testing

- Go tests cover: models (100%), config (91%), api (46%), storage (67%)
- No tests for `internal/desktop/` (Fyne UI — requires display)
- Storage tests use `t.TempDir()` for isolation
- API handler tests use `httptest.NewRecorder()` + `t.TempDir()` for hermetic tests
- Web tests: 43 tests across api.js, all 4 components, ErrorBoundary, App routing (vitest)

## Environment

- Root `.env`: PORT (8080), DATA_DIR, APP_TITLE, API_KEY, CORS_ORIGIN
- `web/.env`: VITE_API_BASE, VITE_API_KEY
- `make env` creates from .env.example templates
