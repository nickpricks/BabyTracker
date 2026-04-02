# Makefile Reference

> Complete reference for all `make` targets in BabyTracker. Run `make` or `make help` to see a summary.

---

## How It Works

The Makefile loads `.env` at the root via `-include .env` + `export`, making all environment variables (PORT, DATA_DIR, API_KEY, etc.) available to every target automatically. The default goal is `help`, so running bare `make` prints the target list.

All targets are declared `.PHONY` — they always run regardless of file timestamps.

---

## Development

### `make dev`

Starts both the API server and the Vite web dev server concurrently.

- API runs on `:${PORT:-8080}` (default 8080)
- Web runs on `:3000`
- Uses `trap 'kill 0' EXIT` to kill both processes when you Ctrl+C
- Internally calls `make api` and `make web` as background jobs

```
make dev
```

### `make api`

Runs the Go API server via `go run ./cmd/api`. Uses PORT from `.env` (default 8080). The API serves REST endpoints at `/api/{resource}` and handles CORS for the web frontend.

```
make api
```

### `make desktop`

Runs the Fyne desktop GUI app via `go run ./cmd/desktop`. Reads/writes JSON data directly (no API needed). Requires a display — won't work in headless/SSH sessions.

**macOS note:** You'll see `ld: warning: ignoring duplicate libraries: '-lobjc'` — this is harmless.

```
make desktop
```

### `make web`

Starts the Vite dev server for the React PWA. Requires the API server to be running separately (or use `make dev` to run both).

- Runs on port 3000
- Hot module replacement (HMR) enabled
- Reads `web/.env` for `VITE_API_BASE` and `VITE_API_KEY`

```
make web
```

---

## Build

### `make build`

Builds all three targets: API binary, desktop binary, and web production bundle. Equivalent to running `build-api`, `build-desktop`, and `build-web` sequentially.

```
make build
```

### `make build-api`

Compiles the API server to `bin/api`.

```
make build-api
# Output: bin/api
```

### `make build-desktop`

Compiles the Fyne desktop app to `bin/desktop`. First build is slow (~30s) due to CGO/Fyne compilation; subsequent builds are cached.

```
make build-desktop
# Output: bin/desktop
```

### `make build-web`

Runs the Vite production build for the web app. Output goes to `web/build/` (configured via `build.outDir` in `vite.config.js`). Includes PWA service worker generation.

```
make build-web
# Output: web/build/
```

---

## Test

### `make test`

Runs all Go tests across `internal/` packages (models, storage, api, config). Does not include desktop tests (Fyne requires a display) or web tests.

```
make test
```

### `make test-cover`

Runs Go tests with coverage profiling, prints a per-function coverage report, then cleans up the `coverage.out` file.

```
make test-cover
# Example output:
# babytracker/internal/models/feed.go:25:    IsBottleFeed    100.0%
# babytracker/internal/storage/storage.go:45: SaveFeed       80.0%
# total:                                      (statements)   72.3%
```

### `make test-web`

Runs the web frontend tests via vitest. Tests cover API client, all 4 components, ErrorBoundary, and App routing (41 tests total).

```
make test-web
```

### `make test-all`

Runs everything in sequence, stopping on first failure:

1. Go tests with coverage report
2. Web tests (vitest)

Unlike running `make test` + `make test-web` separately, this chains them with `&&` so a Go test failure prevents web tests from running.

```
make test-all
```

---

## Lint / Tidy

### `make lint`

Runs `go vet ./...` across the entire Go codebase. Catches suspicious constructs, unreachable code, incorrect printf format strings, etc.

```
make lint
```

### `make lint-web`

Runs ESLint on the web frontend (`web/src/`). Checks for React best practices, unused variables, and JSX issues.

```
make lint-web
```

### `make tidy`

Runs `go mod tidy` to add missing and remove unused Go module dependencies from `go.mod` and `go.sum`.

```
make tidy
```

### `make update`

Updates all Go dependencies to their latest minor/patch versions, then tidies. This modifies `go.mod` and `go.sum`. Run `make test` afterward to verify nothing broke.

```
make update
# Equivalent to:
#   go get -u ./...
#   go mod tidy
```

---

## Setup

### `make setup`

Full first-time project setup. Runs three steps in order:

1. `make env` — creates `.env` files from templates
2. `make tidy` — tidies Go modules
3. `make install-web` — installs web dependencies

```
make setup
```

### `make env`

Creates `.env` files from `.env.example` templates if they don't already exist. Safe to run multiple times — won't overwrite existing files.

- Root `.env` — configures PORT, DATA_DIR, APP_TITLE, API_KEY, CORS_ORIGIN
- `web/.env` — configures VITE_API_BASE, VITE_API_KEY

```
make env
```

### `make install-web`

Installs web frontend dependencies via `bun install` in the `web/` directory.

```
make install-web
```

---

## Bench

### `make bench`

Generates 10,000 entries per module (40,000 total) for stress testing the JSON storage engine. Useful for testing UI performance with large datasets.

**Before generating:**
- Backs up existing `~/.babytracker/*.json` to `~/.babytracker/.backup/`

**Data generated:**
- `feeds.json` — 10k feed entries across ~10 months
- `sleep.json` — 10k sleep entries (naps + nights)
- `growth.json` — 10k growth measurements with realistic progression
- `diapers.json` — 10k diaper entries

```
make bench
```

### `make bench-restore`

Restores data from the backup created by `make bench`. Copies `~/.babytracker/.backup/*.json` back to `~/.babytracker/` and removes the backup directory.

Fails gracefully if no backup exists.

```
make bench-restore
```

---

## Clean

### `make clean`

Removes all build artifacts. Runs both `clean-bin` and `clean-web`.

```
make clean
```

### `make clean-bin`

Removes the `bin/` directory containing compiled Go binaries.

```
make clean-bin
```

### `make clean-web`

Removes the `web/build/` directory containing the Vite production build output.

```
make clean-web
```

---

## Help

### `make help`

Prints a formatted list of all targets with their `##` descriptions. This is the default target — running bare `make` triggers it.

Uses `MAKEFILE_LIST` (a built-in Make variable listing all parsed Makefiles) so it would also include targets from any files brought in via `include`.

```
make help
# or just:
make
```

---

## Quick Reference

| Target | What it does |
|---|---|
| `make` | Show help |
| `make setup` | First-time setup |
| `make dev` | Run API + web concurrently |
| `make test-all` | Run all tests |
| `make build` | Build everything |
| `make lint` | Vet Go code |
| `make bench` | Generate 40k test entries |
| `make bench-restore` | Restore original data |
| `make clean` | Remove build artifacts |
| `make update` | Update Go dependencies |
