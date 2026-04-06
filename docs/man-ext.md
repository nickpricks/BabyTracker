# Baby Tracker - Third-Party Dependencies

> A guide to every external dependency: what it does, why we use it, and what to watch out for.

---

## Go Dependencies (from `go.mod`)

### Fyne (`fyne.io/fyne/v2` v2.7.3) -- Desktop GUI Framework

**What it is**: Fyne is a cross-platform GUI toolkit for Go, built on OpenGL. It provides widgets, layouts, data binding, theming, and application lifecycle management. Think of it as Go's answer to Qt or Electron -- but native, compiled, and without a browser engine.

**Why we use it**: Fyne lets us build a native desktop app entirely in Go, sharing the same `internal/` packages as the API server. No CGO bridge to a C++ toolkit, no JavaScript runtime -- just Go.

**What we use from Fyne**:
- `fyne.App` / `fyne.Window` -- Application and window lifecycle (see `internal/desktop/app.go`)
- `container.NewAppTabs` / `container.NewTabItem` -- Tabbed interface layout
- `container.NewVBox` / `container.NewHBox` -- Vertical and horizontal box layouts
- `widget.NewForm` / `widget.FormItem` -- Form construction with labeled inputs
- `widget.NewEntry` / `widget.NewMultiLineEntry` -- Text input fields
- `widget.NewSelect` -- Dropdown selection widgets
- `widget.NewButton` -- Action buttons
- `widget.NewCard` -- Grouped content panels with titles
- `widget.NewLabel` / `widget.NewSeparator` -- Display elements
- `data/binding` -- Reactive data binding (`binding.NewString()`, `binding.NewFloat()`, `binding.FloatToString()`)
- `theme.AccountIcon()` -- Default app icon

**The Fyne Pain**:

Fyne requires **system-level graphics dependencies** that are NOT installed by `go get`. This is the #1 stumbling block for new developers:

**macOS**:
- Requires Xcode Command Line Tools (`xcode-select --install`)
- The linker emits `ld: warning: ignoring duplicate libraries: '-lobjc'` on every build -- this is **harmless** and comes from Fyne's Objective-C bridge. You cannot suppress it. Just ignore it.
- First build is slow (~30-60s) because it compiles the OpenGL bindings

**Linux**:
- Requires: `gcc`, `libgl1-mesa-dev`, `xorg-dev` (Ubuntu/Debian) or equivalent
- On headless servers (CI/CD), you need a virtual display (`xvfb-run`) for tests that touch the GUI
- Wayland support is improving but X11 is more reliable as of v2.6

**Windows**:
- Requires a C compiler: MSYS2 with MinGW-w64 or TDM-GCC
- The `gcc` must be in PATH
- Cross-compilation from macOS/Linux to Windows is possible but painful

**Common Fyne gotchas**:
- `widget.NewEntryWithData(binding)` creates a two-way bound entry, but the binding must be the correct type (`binding.String` for text, `binding.FloatToString(binding.Float)` for numeric entries)
- `app.New()` must be called exactly once; calling it again creates a second app instance
- `window.ShowAndRun()` blocks the goroutine -- it's the event loop. Put it last in `main()`.
- Fyne's coordinate system uses `float32` (`fyne.NewSize(800, 600)`), not pixels -- it's DPI-independent
- `SetCloseIntercept` overrides the default close behavior; you MUST call `window.Close()` inside it or the window won't close

**Version note**: v2.7.3 is current as of this project (master branch). Fyne v2 has breaking changes from v1 (import path changed from `fyne.io/fyne` to `fyne.io/fyne/v2`).

---

### gorilla/mux (`github.com/gorilla/mux` v1.8.1) -- HTTP Router

**What it is**: A powerful HTTP request router and dispatcher for Go. Extends the standard `net/http` with path variables, regex constraints, middleware support, and method-based routing.

**Why we use it**: The standard library's `http.ServeMux` (prior to Go 1.22) doesn't support path parameters like `/api/feeds/{id}` or method-based routing (GET vs POST on the same path). gorilla/mux fills this gap cleanly.

**What we use from mux**:
- `mux.NewRouter()` -- Creates the router instance
- `r.HandleFunc(path, handler).Methods(methods...)` -- Registers handlers with HTTP method constraints
- `mux.Vars(r)` -- Extracts path variables (e.g., `{id}`) from the request
- `r.Use(middleware)` -- Registers middleware functions
- `mux.CORSMethodMiddleware(r)` -- Built-in CORS method handling

**gorilla/mux gotchas**:
- **OPTIONS and middleware**: `mux.CORSMethodMiddleware` adds `Access-Control-Allow-Methods` headers, but your custom CORS middleware must explicitly handle OPTIONS requests. If you register a route with `.Methods("GET", "POST")`, OPTIONS requests won't match that route unless you also add `"OPTIONS"` to the methods list OR handle it in middleware before routing. Our `router.go` handles this by intercepting OPTIONS in the CORS middleware and returning 200 immediately.
- **Path variables are strings**: `mux.Vars(r)["id"]` returns a string. You must `strconv.Atoi()` it yourself.
- **Route order matters**: More specific routes should be registered before general ones. Our regex constraint `{id:[0-9]+}` prevents ambiguity.
- **gorilla/mux maintenance status**: The Gorilla project was archived in late 2022, then community-maintained. It works fine for our use case, but Go 1.22+ added native path parameters to `http.ServeMux`, making mux less necessary for new projects.

---

### Indirect Go Dependencies

These are pulled in transitively by Fyne. You don't import them directly, but they're in `go.sum`:

| Dependency | Purpose |
|-----------|---------|
| `fyne.io/systray` | System tray integration for desktop apps |
| `BurntSushi/toml` | TOML config parsing (Fyne's internal config) |
| `davecgh/go-spew` | Debug pretty-printer (used by testify) |
| `fredbi/uri` | URI parsing library |
| `fsnotify/fsnotify` | File system notifications (Fyne's file watchers) |
| `fyne-io/gl-js` | WebGL bindings for Fyne's web target |
| `fyne-io/glfw-js` | GLFW bindings for Fyne's web target |
| `fyne-io/image` | Image format support beyond stdlib |
| `fyne-io/oksvg` | SVG rendering for Fyne icons |
| `go-gl/gl` | OpenGL bindings -- the core rendering backend |
| `go-gl/glfw` | GLFW window/input library -- creates the actual OS window |
| `go-text/render` | Text rendering engine |
| `go-text/typesetting` | Font layout and shaping |
| `godbus/dbus` | D-Bus IPC (Linux desktop integration) |
| `hack-pad/go-indexeddb` | IndexedDB bindings (Fyne web storage) |
| `hack-pad/safejs` | Safe JavaScript interop (Fyne web target) |
| `jeandeaual/go-locale` | System locale detection |
| `jsummers/gobmp` | BMP image format support |
| `nfnt/resize` | Image resizing (icon scaling) |
| `nicksnyder/go-i18n` | Internationalization support |
| `pmezard/go-difflib` | Diff library (used by testify) |
| `rymdport/portal` | XDG Desktop Portal integration (Linux file dialogs) |
| `srwiley/oksvg` | Another SVG library (Fyne's fallback renderer) |
| `srwiley/rasterx` | SVG rasterization |
| `stretchr/testify` | Testing assertions and mocks |
| `yuin/goldmark` | Markdown parser (Fyne's rich text widget) |
| `golang.org/x/image` | Extended image format support |
| `golang.org/x/net` | Extended networking (HTTP/2, etc.) |
| `golang.org/x/sys` | Low-level OS primitives |
| `golang.org/x/text` | Unicode text processing |
| `gopkg.in/yaml.v3` | YAML parsing |

---

## JavaScript Dependencies (from `web/package.json`)

### React (`react` ^18.2.0) -- UI Library

**What it is**: The dominant JavaScript library for building component-based user interfaces. Uses a virtual DOM for efficient updates and JSX for declarative UI templates.

**What we use**: Functional components with hooks (`useState`, `useEffect`). No class components, no Redux, no complex state management -- just local component state and prop drilling.

**React patterns in this project**:
- Each page component (Feeds, Sleep, Growth, SusuPoty) manages its own form state via `useState`
- API calls happen in `useEffect` (on mount) and form submission handlers
- Controlled inputs: every `<input>` and `<select>` has `value` + `onChange`
- Feedback and error display via conditional rendering

### React DOM (`react-dom` ^18.2.0) -- DOM Renderer

**What it is**: The package that connects React to the browser DOM. Provides `createRoot()` for React 18's concurrent rendering.

**Our usage**: `createRoot(document.getElementById("root")).render(<App />)` in `index.jsx`.

### React Router DOM (`react-router-dom` ^6.22.3) -- Client-Side Routing

**What it is**: The standard routing library for React SPAs. Provides declarative route definitions, navigation components, and URL parameter handling.

**What we use**:
- `BrowserRouter` -- HTML5 history-based routing (wraps the entire app)
- `Routes` / `Route` -- Declarative route definitions in `Routes.jsx`
- `Navigate` -- Redirect from `/` to `/feeds`
- `Link` -- Client-side navigation (no full page reload) in `Main.jsx`

**Why `Link` instead of `<a href>`**: Using `<a>` causes a full page reload, killing React's SPA state and causing a flash. `<Link>` does client-side navigation, preserving state and enabling instant transitions. This was a bug fix in v0.3.

### Vite (`vite` ^8.0.3) -- Build Toolchain

**What it is**: A fast, modern build tool that uses native ES modules for development and Rollup for production builds. Replaces Create React App (CRA) and webpack.

**Why we migrated**: CRA is unmaintained and uses webpack, which has slow cold starts. Vite provides near-instant dev server startup via native ESM and faster HMR.

**What it provides**:
- Native ESM dev server with instant HMR (no bundling during development)
- Rollup-based production builds (minification, tree shaking, chunk hashing)
- `@vitejs/plugin-react` for JSX/React Fast Refresh
- `@tailwindcss/vite` plugin for Tailwind CSS v4 integration
- `vite-plugin-pwa` for service worker / PWA manifest generation
- Environment variables via `import.meta.env` (prefixed with `VITE_`)

**Vite gotchas**:
- Environment variables must be prefixed with `VITE_` (not `REACT_APP_`) to be exposed to the browser
- Vite reads `web/.env` natively; the root `.env` is only for Go targets
- JSX must be in `.jsx` files (Vite does not transform `.js` files containing JSX)
- `import.meta.env.MODE` gives `"development"` or `"production"` (replaces `process.env.NODE_ENV`)

### Tailwind CSS (`tailwindcss` ^4.2.2) -- Utility-First CSS Framework

**What it is**: A utility-first CSS framework that provides low-level utility classes for building custom designs directly in markup.

**Our usage**: Tailwind CSS v4, integrated via `@tailwindcss/vite` plugin. Used for all web UI styling -- layout, spacing, typography, colors, responsive design, and dark mode. Theme-aware via CSS custom properties (e.g., `text-fg-heading`, `bg-surface`, `border-line`).

**v4-specific notes**: Tailwind v4 uses a CSS-first configuration model (no `tailwind.config.js`). Configuration lives in the CSS files themselves using `@theme` directives.

### Vitest (`vitest` ^4.1.2) -- Test Runner

**What it is**: A Vite-native test framework, API-compatible with Jest. Uses the same Vite config and transform pipeline, so tests run with the same module resolution as the app.

**Our usage**: 43 tests across `api.js`, all 4 components, `ErrorBoundary`, and `App` routing. Uses `@testing-library/react` for component testing and `jsdom` as the DOM environment. Run with `make test-web` or `bun run test`.

### ESLint (`eslint` ^8.57.0) + Plugin (`eslint-plugin-react` ^7.33.2) -- Linter

**What it is**: JavaScript/JSX static analysis tool that catches bugs, enforces style, and prevents common mistakes.

**Our config**: Minimal. Run with `make lint-web`.

---

## Build Tools

### bun (used instead of npm)

**What it is**: A fast JavaScript runtime, package manager, and bundler. Drop-in replacement for npm/yarn with significantly faster install times.

**Our usage**: `bun install` (via `make install-web`) and `bun run dev` / `bun run test` (via Makefile targets). The `bun.lock` file is gitignored.

**Note**: Vite handles all bundling -- bun serves purely as the package manager and script runner.

### Make (GNU Make)

**What it is**: The classic build automation tool. Our `Makefile` is the single entry point for all development, build, test, and deployment operations.

**Key Makefile features**:
- `-include .env` + `export` loads environment variables for all targets
- `.PHONY` declarations prevent file-based caching of target execution
- `trap 'kill 0' EXIT` in `make dev` ensures both API and web servers stop together
- `## comments` after targets are extracted by `make help` for self-documentation
