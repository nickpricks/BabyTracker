# Baby Tracker - TLDR

> Manual.md but just the headlines. Click through to [Manual.md](Manual.md) for the full story.

---

## System Overview
- Multi-platform baby activity tracker (Desktop + Web + PWA)
- Stack: Go 1.24, Fyne v2.6.1, React 18, Vite, Tailwind CSS v4, bun
- Shared Go core, four tracking domains

## Architecture
- Fyne desktop -> direct storage I/O
- React web -> HTTP API -> storage I/O
- All data in `~/.babytracker/*.json`

## Shared Core (`internal/`)
- `models/` -- domain entities (FeedEntry, SleepEntry, GrowthEntry, DiaperEntry)
- `storage/` -- generic JSON engine with Go generics
- `config/` -- env-based config (PORT, DATA_DIR, APP_TITLE)
- `api/` -- REST handlers + gorilla/mux CORS middleware
- `desktop/` -- Fyne tabbed UI with data-binding forms

## Storage Engine
- Generic `loadJSON[T]` / `saveJSON[T]`
- Append-only, full-file rewrite
- ID = `max(existing) + 1`
- Global singleton, lazy init

## API Layer
- GET/POST per resource, GET by ID
- CORS configurable-origin middleware
- Decode -> Validate -> Log -> Save -> Respond

## Modules
- **Feeds**: Bottle / Breast / Solid, quantity + duration
- **Sleep**: Nap / Night, start/end times, quality rating
- **Growth**: Weight (kg), Height (cm), Head Circ (cm)
- **Diapers**: Wet / Dirty / Mixed

## Web App
- Vite + React 18 + Tailwind CSS v4, managed with bun
- `api.js` wraps fetch with `apiGet`/`apiPost`
- Dashboard with 3 switchable themes (v0.4)
- PWA: manifest + service worker + icons

## Config
- Root `.env` for Go (PORT, DATA_DIR, API_KEY, CORS_ORIGIN)
- `web/.env` for Vite (VITE_API_BASE, VITE_API_KEY)
- `make env` creates from examples

## Build
- `make dev` -- run everything
- `make test` -- Go tests
- `make build` -- compile all
- `make setup` -- first-time setup

## Tests
- Go: model + storage + API handler tests
- Web: 43 vitest tests (api, components, ErrorBoundary, routing)

## Roadmap

See **[ROADMAP.md](ROADMAP.md)** for the full roadmap and known gaps.
