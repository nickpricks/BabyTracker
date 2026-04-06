# Baby Tracker - Function & File Reference

> Line-by-line reference for every file and function in the codebase.

---

## Entry Points

### `cmd/api/main.go`

The HTTP API server entry point.

- **`main()`** -- Loads config, initializes storage with the configured data directory, sets up the gorilla/mux router via `api.SetupRouter()`, and starts the HTTP server on the configured port. Logs the listen address and data directory on startup. Fatal-exits on config or storage init failure.

### `cmd/desktop/main.go`

The Fyne desktop application entry point.

- **`main()`** -- Loads config, initializes storage, creates the desktop `App` via `desktop.NewApp()`, logs the data directory, and calls `app.Run()` which blocks until the window is closed. Fatal-exits on config, storage, or app init failure.

---

## Internal Packages

### `internal/config/config.go`

Centralized environment-based configuration.

- **`Config` struct** -- Holds `APIPort` (string), `DataDir` (string), `AppTitle` (string), `APIKey` (string), `CORSOrigin` (string). All fields populated from environment variables with fallback defaults.

- **Constants**: `DefaultAPIPort = "8080"`, `DefaultDataDir = ".babytracker"`, `DefaultAppTitle = "Baby Tracker"`, `DefaultCORSOrigin = "http://localhost:3000"`

- **`Load() (*Config, error)`** -- Reads `PORT`, `DATA_DIR`, `APP_TITLE`, `API_KEY`, and `CORS_ORIGIN` from environment variables. For `DATA_DIR`, falls back to `$HOME/.babytracker` if not set. Returns error only if `os.UserHomeDir()` fails.

- **`envOr(key, fallback string) string`** -- Helper that returns the environment variable value or the fallback. Unexported.

---

### `internal/models/feed.go`

Feed entry data model.

- **`FeedEntry` struct** -- Fields: `ID` (int), `Date` (string, YYYY-MM-DD), `Time` (time.Time), `Type` (string), `Quantity` (float64, ml/oz), `Notes` (string), `Duration` (int, minutes). All fields have `json` tags for serialization.

- **Feed type constants**: `FeedTypeBottle = "Bottle"`, `FeedTypeBreastLeft = "Breast (Left)"`, `FeedTypeBreastRight = "Breast (Right)"`, `FeedTypeBreastBoth = "Breast (Both)"`, `FeedTypeSolid = "Solid Food"`

- **`(f *FeedEntry) IsBottleFeed() bool`** -- Returns true if `Type == FeedTypeBottle`.

- **`(f *FeedEntry) IsBreastFeed() bool`** -- Returns true if Type is any of the three breast variants (Left, Right, Both).

- **`(f *FeedEntry) HasQuantity() bool`** -- Returns true if the feed type is Bottle or Solid Food (types that have a measurable quantity).

### `internal/models/sleep.go`

Sleep entry data model.

- **`SleepEntry` struct** -- Fields: `ID` (int), `Date` (string), `StartTime` (time.Time), `EndTime` (time.Time), `Duration` (int, minutes), `Type` (string), `Quality` (string), `Notes` (string).

- **Sleep type constants**: `SleepTypeNap = "Nap"`, `SleepTypeNight = "Night"`

- **Quality constants**: `SleepQualityGood = "Good"`, `SleepQualityFair = "Fair"`, `SleepQualityPoor = "Poor"`

- **`(s *SleepEntry) IsNap() bool`** -- Returns true if `Type == SleepTypeNap`.

- **`(s *SleepEntry) IsNightSleep() bool`** -- Returns true if `Type == SleepTypeNight`.

### `internal/models/growth.go`

Growth measurement data model.

- **`GrowthEntry` struct** -- Fields: `ID` (int), `Date` (string), `Weight` (float64, kg), `Height` (float64, cm), `HeadCircumference` (float64, cm, json tag: `head_circ`), `Notes` (string).

- **`(g *GrowthEntry) HasWeight() bool`** -- Returns true if `Weight > 0`.

- **`(g *GrowthEntry) HasHeight() bool`** -- Returns true if `Height > 0`.

- **`(g *GrowthEntry) HasHeadCircumference() bool`** -- Returns true if `HeadCircumference > 0`.

### `internal/models/diaper.go`

Diaper change data model.

- **`DiaperEntry` struct** -- Fields: `ID` (int), `Date` (string), `Time` (time.Time), `Type` (string), `Notes` (string).

- **Diaper type constants**: `DiaperTypeWet = "Wet"`, `DiaperTypeDirty = "Dirty"`, `DiaperTypeMixed = "Mixed"`

- **`(d *DiaperEntry) IsWet() bool`** -- Returns true if Type is Wet or Mixed (Mixed counts as both).

- **`(d *DiaperEntry) IsDirty() bool`** -- Returns true if Type is Dirty or Mixed.

---

### `internal/storage/storage.go`

Generic JSON file persistence engine.

- **`StorageManager` struct** -- Holds `dataDir` (string) for the storage directory path.

- **`NewStorageManager() (*StorageManager, error)`** -- Creates a StorageManager using the default directory `$HOME/.babytracker`. Calls `NewStorageManagerWithDir`.

- **`NewStorageManagerWithDir(dataDir string) (*StorageManager, error)`** -- Creates a StorageManager with the given directory, ensuring it exists via `os.MkdirAll(dataDir, 0755)`.

- **`Init(dataDir string) error`** -- Initializes the global singleton `StorageManager` with an explicit directory. Should be called from `main()` before any other storage operations.

- **`getStorage() (*StorageManager, error)`** -- Returns the global singleton, lazy-initializing with defaults if `Init()` was never called. Unexported.

- **`loadJSON[T any](sm *StorageManager, filename string) ([]T, error)`** -- Generic function that reads a JSON file from the data directory and unmarshals it into a slice of `T`. Returns an empty slice (not nil) if the file doesn't exist. Unexported.

- **`saveJSON[T any](sm *StorageManager, filename string, items []T) error`** -- Generic function that marshals a slice of `T` to indented JSON and writes it to the data directory. Uses `os.WriteFile` with `0644` permissions. Unexported.

- **`nextID(ids []int) int`** -- Scans a slice of existing IDs, finds the maximum, and returns `max + 1`. Returns 1 for an empty slice. Unexported.

- **`SaveFeed(feed *models.FeedEntry) error`** -- Loads existing feeds, generates a new ID, appends the entry, and saves. Assigns the generated ID to the feed's `ID` field.

- **`LoadFeeds() ([]models.FeedEntry, error)`** -- Returns all feed entries from `feeds.json`.

- **`SaveSleep(entry *models.SleepEntry) error`** -- Same pattern as SaveFeed, writes to `sleep.json`.

- **`LoadSleep() ([]models.SleepEntry, error)`** -- Returns all sleep entries from `sleep.json`.

- **`SaveGrowth(entry *models.GrowthEntry) error`** -- Same pattern, writes to `growth.json`.

- **`LoadGrowth() ([]models.GrowthEntry, error)`** -- Returns all growth entries from `growth.json`.

- **`SaveDiaper(entry *models.DiaperEntry) error`** -- Same pattern, writes to `diapers.json`.

- **`LoadDiapers() ([]models.DiaperEntry, error)`** -- Returns all diaper entries from `diapers.json`.

- **`GetDataDirectory() (string, error)`** -- Returns the storage directory path from the global singleton. Useful for logging/diagnostics.

---

### `internal/api/router.go`

HTTP route registration, middleware, and CORS handling.

- **`SetupRouter(cfg *config.Config) http.Handler`** -- Creates a gorilla/mux router, attaches a 1MB request body size limit middleware, optionally adds Bearer-token auth middleware (skipped if `cfg.APIKey` is empty; OPTIONS requests bypass auth). Registers all 12 endpoints (3 per module: list, create, get-by-id). Returns the router wrapped in `corsHandler()`, which means the return type is `http.Handler` (not `*mux.Router`), because CORS must intercept OPTIONS preflight before mux rejects it with 405.

- **`corsHandler(corsOrigin string, next http.Handler) http.Handler`** -- Wraps a handler with CORS headers. Sets `Access-Control-Allow-Origin` to the configured origin (from `CORS_ORIGIN` env var, default `http://localhost:3000`), `Access-Control-Allow-Headers: Content-Type, Authorization`, and `Access-Control-Allow-Methods: GET,POST,OPTIONS`. If the configured origin is localhost, accepts any `http://localhost:*` origin for dev convenience. Returns 200 immediately for OPTIONS preflight requests. Unexported.

### `internal/api/handlers.go`

Feed endpoint handlers and shared utilities.

- **`jsonResponse(w http.ResponseWriter, status int, payload interface{})`** -- Sets `Content-Type: application/json`, writes the status code, and JSON-encodes the payload. Used by all handlers across all handler files. Unexported.

- **`handleListFeeds(w, r)`** -- GET `/api/feeds`. Calls `storage.LoadFeeds()`, returns the full list as JSON array.

- **`handleLogFeed(w, r)`** -- POST `/api/feeds`. Decodes JSON body into `FeedEntry`, validates that `type` and `date` are non-empty, logs the entry, calls `storage.SaveFeed()`, returns 201 with the saved entry (including generated ID).

- **`handleGetFeed(w, r)`** -- GET `/api/feeds/{id}`. Extracts `id` from URL path via `mux.Vars()`, loads all feeds, linear-scans for matching ID, returns 404 if not found.

### `internal/api/sleep_handlers.go`

Sleep endpoint handlers.

- **`handleListSleep(w, r)`** -- GET `/api/sleep`. Returns all sleep entries.

- **`handleLogSleep(w, r)`** -- POST `/api/sleep`. Validates `date` and `type` required. Saves and returns 201.

- **`handleGetSleep(w, r)`** -- GET `/api/sleep/{id}`. Linear scan by ID, 404 if not found.

### `internal/api/growth_handlers.go`

Growth endpoint handlers.

- **`handleListGrowth(w, r)`** -- GET `/api/growth`. Returns all growth entries.

- **`handleLogGrowth(w, r)`** -- POST `/api/growth`. Validates only `date` required (all measurements optional). Saves and returns 201.

- **`handleGetGrowth(w, r)`** -- GET `/api/growth/{id}`. Linear scan by ID, 404 if not found.

### `internal/api/diaper_handlers.go`

Diaper endpoint handlers.

- **`handleListDiapers(w, r)`** -- GET `/api/diapers`. Returns all diaper entries.

- **`handleLogDiaper(w, r)`** -- POST `/api/diapers`. Validates `date` and `type` required. Saves and returns 201.

- **`handleGetDiaper(w, r)`** -- GET `/api/diapers/{id}`. Linear scan by ID, 404 if not found.

---

### `internal/desktop/app.go`

Fyne application lifecycle management.

- **`App` struct** -- Holds `fyneApp` (fyne.App) and `window` (fyne.Window).

- **`NewApp() *App`** -- Creates a new Fyne application, sets the icon to `theme.AccountIcon()`, creates the main window titled "Baby Tracker" at 800x600, centers it on screen.

- **`(a *App) CreateMainContent() fyne.CanvasObject`** -- Builds four tab items (Feeds, Sleep, Growth, Susu-Poty) by calling each `tabs.CreateXxxTab()`, assembles them into `container.NewAppTabs` with top tab placement.

- **`(a *App) SetupWindow()`** -- Sets the main content, marks the window as master (closing it exits the app), installs a close interceptor.

- **`(a *App) Run()`** -- Calls `SetupWindow()` then `window.ShowAndRun()`. Blocks until the window is closed.

- **`(a *App) GetWindow() fyne.Window`** -- Returns the main window reference.

- **`(a *App) GetApp() fyne.App`** -- Returns the Fyne app instance.

### `internal/desktop/layout.go`

Alternative layout constructor (currently unused by `App.Run()`).

- **`CreateMainLayout() *container.AppTabs`** -- Creates the same four-tab layout as `App.CreateMainContent()`. This function exists as a standalone alternative but is not called in the current code path.

### `internal/desktop/tabs/feeds.go`

Fyne feed tracking form.

- **Constants**: `dateFormat = time.DateOnly`, `timeFormat = time.TimeOnly` (shared across all tabs in this package)

- **`CreateFeedsTab() *fyne.Container`** -- Builds a complete feed logging form with: data bindings for date/time/quantity/notes, a `Select` widget for feed type (5 options), form entries, a "Log Feed" button that parses time, constructs a `FeedEntry`, and calls `storage.SaveFeed()`. Includes "Quick Bottle" and "Quick Breast" buttons that pre-fill the form. Has a placeholder section for recent feeds display.

### `internal/desktop/tabs/sleep.go`

Fyne sleep tracking form.

- **`CreateSleepTab() *fyne.Container`** -- Sleep logging form with: type select (Nap/Night), quality select (Good/Fair/Poor), date/start-time/end-time entries, notes. The log button computes duration from end-start times. Includes "Quick Nap" and "Quick Night" buttons. Recent sleep placeholder.

### `internal/desktop/tabs/growth.go`

Fyne growth tracking form.

- **`CreateGrowthTab() *fyne.Container`** -- Growth measurement form with: date entry, weight/height/head circumference float entries, notes. No quick action buttons (growth entries are less frequent). Recent measurements placeholder.

### `internal/desktop/tabs/susupoty.go`

Fyne diaper tracking form.

- **`CreateSusuPotyTab() *fyne.Container`** -- Diaper logging form with: type select (Wet/Dirty/Mixed), date/time entries, notes. Includes "Quick Wet" and "Quick Dirty" buttons. Card title: "The Susu-Poty Chronicles". Recent changes placeholder.

---

## Web Application

### `web/src/index.jsx`

React entry point. Renders `<App />` inside `<StrictMode>`, mounts to `#root` via `createRoot`.

### `web/src/App.jsx`

Root component. Initializes the theme system via `useTheme()`, wraps the app in `<BrowserRouter>` -> `<MainLayout theme={theme}>` -> `<ErrorBoundary>` -> `<AppRoutes>`.

### `web/src/Main.jsx`

- **`ThemeSwitcher({ theme })`** -- Theme and color mode selector. Renders a `<select>` for theme choice (Lullaby, Nursery_OS, Midnight Feed) and light/dark/system toggle buttons (hidden for dark-only themes). Unexported.

- **`MainLayout({ children, theme })`** -- Layout component with a sticky header ("Baby Tracker" + ThemeSwitcher), a `<main>` content area (max-width 3xl, centered), and a fixed bottom navigation bar. Navigation uses `NavLink` with active-state highlighting for Home, Feeds, Sleep, Growth, and Diapers. Styled with Tailwind CSS utility classes and theme-aware CSS custom properties.

### `web/src/Routes.jsx`

- **`AppRoutes()`** -- Defines routes: `/` renders `<Dashboard />`, four feature routes (`/feeds`, `/sleep`, `/growth`, `/susupoty`), and a `*` catch-all 404.

### `web/src/themes.js`

Theme definitions and state management.

- **`THEMES`** -- Object defining three themes: `lullaby` (Nursery family, light+dark), `nursery-os` (Cyberpunk family, dark-only), `midnight-feed` (Nursery family, dark-only). Each theme has `id`, `name`, `family`, `darkOnly`, `cssClass`, and `previewColors`.

- **`getStoredTheme()`** / **`getStoredColorMode()`** -- Read persisted theme/mode from localStorage (defaults: "lullaby", "system").

- **`applyTheme(themeId, colorMode)`** -- Applies the theme CSS class and dark mode to `document.documentElement`, persists choices to localStorage.

- **`useTheme()`** -- React hook that manages theme and color mode state, applies on mount and change, and listens for system color-scheme changes. Returns `{ themeId, setThemeId, colorMode, setColorMode }`.

### `web/src/config.js`

- **`API_BASE`** -- Exported constant. Uses `VITE_API_BASE` env var (via `import.meta.env`) or defaults to `http://localhost:8080/api`.

### `web/src/api.js`

HTTP client layer for the React app.

- **`apiGet(path)`** -- Fetches `${API_BASE}${path}` with GET, throws on non-OK response, returns parsed JSON. Unexported.

- **`apiPost(path, body)`** -- Fetches with POST, `Content-Type: application/json`, stringified body. Attempts to parse error response on failure. Unexported.

- **`getFeeds()`** / **`logFeed(feed)`** -- Feed API calls.
- **`getSleep()`** / **`logSleep(entry)`** -- Sleep API calls.
- **`getGrowth()`** / **`logGrowth(entry)`** -- Growth API calls.
- **`getDiapers()`** / **`logDiaper(entry)`** -- Diaper API calls.

### `web/src/components/Dashboard.jsx`

- **`SummaryCard({ icon, title, color, linkTo, children })`** -- Reusable card with icon, title, and a "+" link to the module's logging page. Unexported.

- **`MiniList({ items, empty })`** -- Renders a list of recent entries or an empty-state message. Unexported.

- **`Dashboard()`** -- Home page component. Fetches the last 5 entries from each module (3 for growth) on mount via `Promise.all`. Displays a greeting, a quick-entry strip with shortcut buttons for each module, and a 2-column summary grid of `SummaryCard` components showing recent feeds, sleep, growth measurements, and diaper changes. Shows a loading spinner while fetching.

### `web/src/components/ErrorBoundary.jsx`

- **`ErrorBoundary`** -- React error boundary (class component). Catches rendering errors in child components and displays a fallback UI instead of crashing the whole app.

### `web/src/components/index.jsx`

Barrel export file. Re-exports `Feeds`, `Growth`, `Sleep`, and `SusuPoty` as named exports.

### `web/src/components/Feeds.jsx`

- **`Feeds()`** -- Complete feed logging component. Manages form state via `useState` (feedType, date, time, quantity, notes, feedback, error, recentFeeds, loading). Fetches recent feeds on mount via `useEffect`. Submit handler calls `logFeed()`, shows feedback for 3 seconds, and resets the form. "Quick Bottle" and "Quick Breast" buttons pre-fill type and timestamp. Displays last 10 feeds in reverse chronological order.

### `web/src/components/Sleep.jsx`

- **`Sleep()`** -- Sleep logging component. Same pattern as Feeds with fields for type (Nap/Night), quality (Good/Fair/Poor), date, start time, end time, notes. Computes duration client-side. Quick Nap / Quick Night buttons.

### `web/src/components/Growth.jsx`

- **`Growth()`** -- Growth measurement component. Fields for date, weight, height, head circumference, notes. No quick action buttons.

### `web/src/components/SusuPoty.jsx`

- **`SusuPoty()`** -- Diaper logging component. Type select (Wet/Dirty/Mixed), date, time, notes. Quick Wet / Quick Dirty buttons.

### PWA Support

PWA functionality (service worker, manifest) is handled by `vite-plugin-pwa` in the Vite config. The old CRA-based `serviceWorkerRegistration.js` has been removed.

---

## Configuration & Build Files

### `Makefile`

Build system with phony targets for dev, build, test, lint, setup, bench, and clean. Loads `.env` via `-include` and exports all variables. Default target is `help`, which greps for `##` comments. Includes `bench` (generates 10k entries per module, backs up existing data) and `bench-restore` (restores from bench backup) targets for performance testing.

### `go.mod`

Go module `babytracker`, Go 1.26.0. Direct dependencies: `fyne.io/fyne/v2` (v2.7.3), `github.com/gorilla/mux` (v1.8.1).

### `.env.example`

Template for root environment file. Documents PORT, DATA_DIR, APP_TITLE, API_KEY, and CORS_ORIGIN.

### `web/.env.example`

Template for web environment file. Documents VITE_API_BASE and VITE_API_KEY.

### `web/package.json`

Vite-based React 18 project. Scripts: dev, build, preview, test (vitest run), test:watch (vitest), lint. Dependencies: react, react-dom, react-router-dom. Dev dependencies: @tailwindcss/vite, @vitejs/plugin-react, tailwindcss v4, vite v8, vitest, vite-plugin-pwa, @testing-library/{jest-dom,react,user-event}, eslint, eslint-plugin-react, jsdom.

### `web/public/manifest.json`

PWA manifest with `standalone` display mode, app name, and icon references (192px, 512px).

### `.gitignore`

Ignores: tmp, IDE files, node_modules, .env files, Go binaries (bin/), web build output, lock files.
