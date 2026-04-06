# Baby Tracker - Technical Manual

> The definitive engineering reference for Baby Tracker: a multi-platform, event-driven activity tracking system built on a shared Go core with heterogeneous frontend delivery.

---

## 1. System Overview

Baby Tracker is a **polyglot, multi-runtime application** designed to track infant care activities across four domains: feeding, sleep, growth measurements, and diaper changes. The system employs a **shared-core architecture** where a single Go codebase powers both a native desktop GUI (via Fyne) and a RESTful HTTP API consumed by a React single-page application and Progressive Web App.

### Design Philosophy

- **Single Source of Truth**: All platforms read/write the same JSON data store
- **Zero-Infrastructure Deployment**: No database server, no Docker, no cloud dependencies -- just flat files
- **Progressive Enhancement**: Desktop for power users, web for convenience, PWA for mobile-first caregivers
- **Domain-Driven Modeling**: Each tracking domain (feeds, sleep, growth, diapers) is a self-contained module with its own model, storage operations, API handlers, desktop tab, and web component

---

## 2. Architecture Deep Dive

### 2.1 High-Level System Architecture

```
                                    +---------------------------+
                                    |     Data Layer            |
                                    |  ~/.babytracker/*.json    |
                                    +---------------------------+
                                         ^            ^
                                         |            |
                              Direct I/O |            | Direct I/O
                                         |            |
+-------------------+           +--------+--------+   |
|   Fyne Desktop    |---------->|  internal/      |   |
|   Application     |  Go calls |  storage/       |   |
|   (cmd/desktop)   |           |  (Generic JSON) |   |
+-------------------+           +-----------------+   |
                                                      |
+-------------------+           +-----------------+   |
|   React Web App   |  HTTP     |  Go API Server  |---+
|   (web/src)       |---------->|  (cmd/api)      |
+-------------------+  JSON     |  gorilla/mux    |
        |                       +-----------------+
        |
+-------------------+
|   Mobile PWA      |
|   (same as web)   |
+-------------------+
```

### 2.2 The Shared Core (`internal/`)

The `internal/` directory is the **nucleus** of the entire system. It contains zero platform-specific code -- every package here is consumed by both the desktop and API entry points.

| Package | Purpose | Consumed By |
|---------|---------|-------------|
| `internal/models` | Domain entity definitions, type constants, helper methods | Everyone |
| `internal/storage` | Generic JSON serialization engine with ID generation | Desktop (direct), API (via handlers) |
| `internal/config` | Environment-based configuration with sensible defaults | Both entry points |
| `internal/api` | HTTP route registration, CORS middleware, request/response handlers | API server only |
| `internal/desktop` | Fyne window management, tabbed layout, form construction | Desktop only |
| `cmd/bench` | HTTP load-testing / benchmarking tool for the API | Standalone CLI |

### 2.3 Data Flow Patterns

**Desktop Path** (lowest latency):
```
User Input -> Fyne Form -> models.XxxEntry{} -> storage.SaveXxx() -> JSON file
```

**Web/API Path** (network-mediated):
```
User Input -> React Form -> fetch(POST) -> API Handler -> json.Decode -> storage.SaveXxx() -> JSON file
```

**Read Path** (both platforms):
```
storage.LoadXxx() -> json.Unmarshal -> []models.XxxEntry -> Display
```

### 2.4 The Storage Engine

The storage layer is built on **Go generics** (`loadJSON[T any]`, `saveJSON[T any]`), providing type-safe serialization for any model type with zero boilerplate duplication. Key design decisions:

- **Append-only writes**: Each save loads the full dataset, appends, and rewrites. This is intentional for simplicity at the current scale (hundreds to low thousands of records).
- **Max-ID generation**: `nextID()` scans all existing IDs to find the maximum, then increments. This prevents ID collisions after deletions (unlike `len(items) + 1`).
- **Global singleton with lazy init**: `getStorage()` initializes on first use, or can be explicitly initialized via `Init(dataDir)` from `main()`.
- **Atomic directory creation**: `os.MkdirAll` ensures the data directory exists before any read/write.

### 2.5 The API Layer

The REST API is a **resource-oriented** HTTP service built on gorilla/mux:

| Endpoint Pattern | Methods | Handler Pattern |
|------------------|---------|----------------|
| `/api/{resource}` | GET, POST | List all / Create new |
| `/api/{resource}/{id}` | GET | Retrieve by ID |

**CORS Strategy**: An external CORS handler wraps the gorilla/mux router, setting `Access-Control-Allow-Origin` to the configured `CORS_ORIGIN` (default: `http://localhost:3000`) and handling OPTIONS preflight requests. This allows the React dev server on `:3000` to talk to the API on `:8080` without proxy configuration.

**Handler Architecture**: Every handler follows the same disciplined pattern:
1. Decode request body (POST) or extract path params (GET by ID)
2. Validate required fields (return 400 on failure)
3. Log the operation
4. Delegate to `storage` package
5. Return JSON response with appropriate status code

### 2.6 The Desktop Application

The Fyne desktop app is structured as a **tabbed interface** with one tab per tracking domain:

- `App` struct wraps `fyne.App` + `fyne.Window`
- `CreateMainContent()` builds four `TabItem`s, each delegating to `tabs.CreateXxxTab()`
- Each tab is a self-contained form using Fyne's **data binding** system (`binding.NewString()`, `binding.NewFloat()`) for reactive two-way form state
- **Quick action buttons** pre-fill common entries (e.g., "Quick Bottle", "Quick Nap") to reduce data entry friction for sleep-deprived parents

### 2.7 The Web Application

A **Vite** React project styled with **Tailwind CSS v4**, with client-side routing via `react-router-dom`:

- `App.jsx` wraps the router
- `Main.jsx` provides the navigation layout (header + nav links)
- `Routes.jsx` maps URL paths to components
- `Dashboard.jsx` provides an at-a-glance overview of recent activity across all modules
- `api.js` is the HTTP client layer (generic `apiGet`/`apiPost` wrappers)
- `config.js` centralizes the API base URL (configurable via `VITE_API_BASE`)
- Each component (Feeds, Sleep, Growth, SusuPoty) is a complete form with state management, submission, validation feedback, and a "recent entries" display
- **Theme system**: 3 themes -- Lullaby (default), Nursery_OS, and Midnight Feed -- switchable at runtime

**PWA Features**:
- `manifest.json` for standalone display mode and home screen installation
- Service worker registration for offline caching
- App icons at 192px and 512px

---

## 3. Module Reference

### 3.1 Feeds Module

Tracks bottle feedings, breastfeeding sessions, and solid food intake.

**Model fields**: ID, Date, Time, Type, Quantity, Notes, Duration

**Feed types**: Bottle, Breast (Left), Breast (Right), Breast (Both), Solid Food

**Helper methods**:
- `IsBottleFeed()` / `IsBreastFeed()` -- classify by type for analytics
- `HasQuantity()` -- determines if the feed type has a measurable quantity

**Validation (API)**: Requires `type` and `date`

### 3.2 Sleep Module

Tracks nap and night sleep sessions with quality assessment.

**Model fields**: ID, Date, StartTime, EndTime, Duration, Type, Quality, Notes

**Sleep types**: Nap, Night
**Quality levels**: Good, Fair, Poor

**Helper methods**: `IsNap()`, `IsNightSleep()`

**Duration calculation**: Computed client-side as `EndTime - StartTime` in minutes

**Validation (API)**: Requires `date` and `type`

### 3.3 Growth Module

Records periodic growth measurements.

**Model fields**: ID, Date, Weight (kg), Height (cm), HeadCircumference (cm), Notes

**Helper methods**: `HasWeight()`, `HasHeight()`, `HasHeadCircumference()` -- check if measurements were recorded (> 0)

**Validation (API)**: Requires `date` only (all measurements are optional per entry)

### 3.4 Diaper Module (Susu-Poty)

Tracks diaper changes with type classification.

**Model fields**: ID, Date, Time, Type, Notes

**Diaper types**: Wet, Dirty, Mixed

**Helper methods**: `IsWet()`, `IsDirty()` -- note that Mixed returns true for both

**Validation (API)**: Requires `date` and `type`

---

## 4. Configuration System

The configuration layer (`internal/config`) reads from environment variables with cascading defaults:

| Variable | Default | Scope | Description |
|----------|---------|-------|-------------|
| `PORT` | `8080` | API server | HTTP listen port |
| `DATA_DIR` | `~/.babytracker` | Both | Absolute path for JSON data files |
| `APP_TITLE` | `Baby Tracker` | Desktop | Window title |
| `VITE_API_BASE` | `http://localhost:8080/api` | Web | API endpoint URL |
| `API_KEY` | *(empty)* | API server | Bearer token for auth (empty = no auth) |
| `CORS_ORIGIN` | `http://localhost:3000` | API server | Allowed CORS origin |

**Loading chain**: Makefile `-include .env` + `export` makes root `.env` available to all Go targets. Vite reads `web/.env` natively.

---

## 5. Build System

The Makefile provides a **comprehensive build pipeline**:

### Development
| Target | Action |
|--------|--------|
| `make dev` | Concurrent API + web dev servers |
| `make api` | API server via `go run` |
| `make desktop` | Fyne desktop via `go run` |
| `make web` | React dev server via `bun start` |

### Build
| Target | Output |
|--------|--------|
| `make build-api` | `bin/api` |
| `make build-desktop` | `bin/desktop` |
| `make build-web` | `web/build/` |

### Quality
| Target | Tool |
|--------|------|
| `make test` | `go test ./internal/...` |
| `make test-web` | vitest (43 tests) |
| `make test-cover` | Go coverage report |
| `make test-all` | All tests (Go + coverage + web) |
| `make lint` | `go vet ./...` |
| `make lint-web` | eslint |

### Setup
| Target | Action |
|--------|--------|
| `make setup` | Full first-time setup |
| `make env` | Create `.env` files from examples |
| `make install-web` | `bun install` for web deps |

---

## 6. Testing Strategy

### Go Tests
- **Model tests** (`*_test.go` in `internal/models/`): Validate struct behavior, type constants, helper methods
- **Storage tests** (`storage_test.go`): Round-trip serialization tests using `t.TempDir()` for isolation
- **API handler tests** (`handlers_test.go`): HTTP-level tests using `httptest.NewRecorder()` and `httptest.NewRequest()`

### Web Tests
- Vitest runner configured (`make test-web`)
- 43 tests across api.js, all 4 components, ErrorBoundary, App routing (vitest)

---

## 7. Roadmap

See **[ROADMAP.md](ROADMAP.md)** for the full roadmap, tech debt tracker, and architectural evolution path.

---

## 8. Security Considerations

- **Bearer token auth** via `API_KEY` env var (empty = no auth, for local dev)
- CORS: configurable origin via `CORS_ORIGIN` env var (default: `http://localhost:3000`)
- Request body limit: 1MB max via `http.MaxBytesReader` middleware
- JSON data files stored with `0600` permissions (owner-only)
- `.env` files are gitignored to prevent credential leakage

---

## 9. Performance Characteristics

- **Storage**: O(n) reads (full file scan), O(n) writes (rewrite entire file). Perfectly adequate for baby tracking volumes (typically < 50 entries/day across all modules).
- **API**: Single-threaded handler execution with gorilla/mux routing. Handles concurrent requests via Go's goroutine model.
- **Desktop**: Native rendering via Fyne/OpenGL. Startup time dominated by GUI initialization (~1-2s on macOS).
- **Web**: Vite development build with HMR. Production builds benefit from code splitting and service worker caching.

---

## 10. Known Limitations & Technical Debt

See [ROADMAP.md](ROADMAP.md) for known limitations and technical debt.
