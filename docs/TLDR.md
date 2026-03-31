# Baby Tracker - TLDR

> Manual.md but just the headlines. Click through to [Manual.md](Manual.md) for the full story.

---

## System Overview
- Multi-platform baby activity tracker (Desktop + Web + PWA)
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
- CORS wildcard middleware
- Decode -> Validate -> Log -> Save -> Respond

## Modules
- **Feeds**: Bottle / Breast / Solid, quantity + duration
- **Sleep**: Nap / Night, start/end times, quality rating
- **Growth**: Weight (kg), Height (cm), Head Circ (cm)
- **Diapers**: Wet / Dirty / Mixed

## Web App
- CRA + react-router-dom
- `api.js` wraps fetch with `apiGet`/`apiPost`
- PWA: manifest + service worker + icons

## Config
- Root `.env` for Go (PORT, DATA_DIR)
- `web/.env` for React (REACT_APP_API_BASE)
- `make env` creates from examples

## Build
- `make dev` -- run everything
- `make test` -- Go tests
- `make build` -- compile all
- `make setup` -- first-time setup

## Tests
- Go: model + storage + API handler tests
- Web: test runner configured, no cases yet

## Roadmap

- v0.3.1: Security hardening -- auth, CORS, storage locking, file permissions (see [Security-Review.md](Security-Review.md))
- v0.4: Charts & analytics
- v0.5: Smart alerts & reminders
- v0.6: Export (CSV/PDF)
- v1.0: Multi-profile, dark mode
- v1.5: Cloud sync
- v2.0: Native mobile, voice integration
- v3.0: "Body Soul and Mind Tracker" -- generalized health

## Known Gaps
- No DELETE/PUT endpoints
- No web tests
- No CI/CD
- Placeholder PWA icons
- Duplicate layout code
