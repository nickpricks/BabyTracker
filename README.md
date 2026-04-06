# 👶 Baby Tracker

A minimalist yet mighty multi-platform app to track your baby's day-to-day growth 📈, nourishment 🍼, sleep patterns 😴, and the noble chronicles of 🚽 Susu-Poty 🧻. Built in Go with a [Fyne](https://fyne.io) desktop GUI, a React web app with Tailwind CSS theming, and a REST API — all sharing the same backend logic.

Whether you're a caregiver, parent, or curious builder, this project aims to balance usability with technical learning—while keeping things fun and purpose-driven 🎯.

---

## 🛠️ Prerequisites

- [Go](https://go.dev/dl/) (v1.24 or later)
- [Bun](https://bun.sh/) — for the web frontend (not npm/yarn)
- [Make](https://www.gnu.org/software/make/) — included on macOS/Linux; on Windows use `choco install make` or run commands manually
- [Fyne dependencies](https://docs.fyne.io/started/) — C compiler and system graphics libs (see Fyne docs for your OS)

---

## 🚀 Quick Start

```bash
# Clone and enter the project
git clone <your-repo-url>
cd BabyTracker

# One-command setup: creates .env files, tidies Go modules, installs web deps
make setup

# Run everything (API + web dev server)
make dev
```

Then open http://localhost:3000 in your browser.

For the desktop app instead:
```bash
make desktop
```

---

## ⚙️ Configuration

Configuration is done via environment variables. Use `.env` files for convenience:

```bash
# Create .env files from examples (safe — won't overwrite existing)
make env
```

### Go Backend (`.env` in project root)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | API server port |
| `DATA_DIR` | `~/.babytracker` | Absolute path for JSON data storage |
| `APP_TITLE` | `Baby Tracker` | Desktop window title |
| `API_KEY` | *(empty)* | Bearer token for API auth (empty = no auth) |
| `CORS_ORIGIN` | `http://localhost:3000` | Allowed CORS origin |

### React Web App (`web/.env`)

| Variable | Default | Description |
|----------|---------|-------------|
| `VITE_API_BASE` | `http://localhost:8080/api` | API URL the web app connects to |
| `VITE_API_KEY` | *(empty)* | API key sent as Bearer token |

The Makefile automatically loads `.env` from the project root, so `make api` and `make desktop` pick up your values. Vite reads `web/.env` natively.

Config is also available in Go code via `internal/config`:

```go
cfg, _ := config.Load()
fmt.Println(cfg.APIPort)  // "8080"
fmt.Println(cfg.DataDir)  // "/Users/you/.babytracker"
```

---

## 📋 Makefile Commands

Run `make` or `make help` to see all available commands. Full reference: [docs/make.md](docs/make.md)

### Development (run)

| Command | Description |
|---------|-------------|
| `make dev` | Run API server + web dev server concurrently |
| `make api` | Run only the Go API server |
| `make desktop` | Run the Fyne desktop app |
| `make web` | Run the Vite dev server (requires API running) |

### Build

| Command | Description |
|---------|-------------|
| `make build` | Build everything (Go binaries + web production build) |
| `make build-api` | Build the API server binary to `bin/api` |
| `make build-desktop` | Build the desktop app binary to `bin/desktop` |
| `make build-web` | Build the web app for production to `web/build/` |

### Test

| Command | Description |
|---------|-------------|
| `make test` | Run all Go tests |
| `make test-cover` | Run Go tests with coverage report |
| `make test-web` | Run web tests (vitest) |
| `make test-all` | Run all tests (Go + coverage + web), stops on first failure |

### Lint, Tidy & Update

| Command | Description |
|---------|-------------|
| `make lint` | Vet Go code (`go vet`) |
| `make lint-web` | Lint web code (`bun run lint`) |
| `make tidy` | Run `go mod tidy` |
| `make update` | Update all Go deps to latest minor/patch |

### Setup, Clean & Bench

| Command | Description |
|---------|-------------|
| `make setup` | Full project setup (env files + Go tidy + web deps) |
| `make env` | Create `.env` files from examples (won't overwrite) |
| `make install-web` | Install web dependencies (`bun install`) |
| `make clean` | Remove all build artifacts (`bin/` + `web/build/`) |
| `make bench` | Generate 10k entries per module for stress testing (backs up data first) |
| `make bench-restore` | Restore data from bench backup |

---

## 🏗️ Project Structure

```
BabyTracker/
├── cmd/
│   ├── api/main.go            # HTTP API server entry point
│   ├── desktop/main.go        # Fyne desktop entry point
│   └── bench/main.go          # Bench data generator (10k entries)
├── internal/
│   ├── models/                # Shared data models (feed, sleep, growth, diaper)
│   ├── storage/               # JSON file persistence (~/.babytracker/)
│   ├── config/                # Environment-based configuration
│   ├── api/                   # HTTP handlers & gorilla/mux router
│   └── desktop/               # Fyne UI (app, layout, tabs)
├── web/                       # React SPA + PWA (Vite + Tailwind v4)
│   └── src/
│       ├── components/        #   Dashboard, Feeds, Sleep, Growth, SusuPoty
│       ├── themes/            #   Lullaby, Nursery_OS, Midnight Feed
│       ├── themes.js          #   Theme definitions + useTheme hook
│       └── index.css          #   Tailwind @theme token mapping
├── docs/                      # Project documentation
├── Makefile                   # Build, run, test commands
└── go.mod
```

---

## 📖 Documentation

Detailed technical documentation lives in the [`docs/`](docs/) folder:

| Document | What's inside |
|----------|---------------|
| [**TLDR.md**](docs/TLDR.md) | The entire project in headline form — scan in 60 seconds |
| [**Manual.md**](docs/Manual.md) | Full technical manual — architecture deep dive, module reference, roadmap |
| [**man.md**](docs/man.md) | Function & file reference — every file, every function described |
| [**man-ext.md**](docs/man-ext.md) | Third-party dependency guide — Fyne pain, gorilla/mux gotchas |
| [**make.md**](docs/make.md) | Makefile reference — every target with usage and examples |
| [**CLAUDE.md**](docs/CLAUDE.md) | Claude Code context — commands, conventions, known gotchas |
| [**Security-Review.md**](docs/Security-Review.md) | Security audit — 23 findings with severity ratings and fix status |

---

## 🧱 Architecture

```
+-------------------+       HTTP/JSON        +-------------------+
|   React Web App   | <------------------->  |    Go API Server  |
|   (web/)          |                        |   (cmd/api)       |
+-------------------+                        +-------------------+
                                                      |
+-------------------+                        +-------------------+
|   Fyne Desktop    | ---direct Go calls---> |  internal/models  |
|   (cmd/desktop)   |                        |  internal/storage |
+-------------------+                        +-------------------+
                                                      |
+-------------------+                                 v
|   Mobile (PWA)    | --- same as web -->    ~/.babytracker/*.json
+-------------------+
```

- **Desktop**: Calls `internal/storage` directly for native performance
- **Web + Mobile**: React SPA calls the Go API server over HTTP
- **Shared**: All platforms use the same models and storage format

For the full architecture deep dive, see [docs/Manual.md](docs/Manual.md).

---

## 🎨 Themes

The web app ships with three themes, built on the same CSS variable architecture as Floor-Tracker for future unification:

| Theme | Modes | Aesthetic |
|-------|-------|-----------|
| **Lullaby** (default) | Light + Dark | Warm cream nursery / soft blue night mode |
| **Nursery_OS** | Dark only | Cyberpunk baby monitor with neon glows |
| **Midnight Feed** | Dark only | Ultra-dim amber for 3am use |

Each tracking module has its own color identity (sage for feeds, lavender for sleep, coral for growth, blue for diapers). Theme selection persists in localStorage.

---

## 🌱 Core Modules

Four tracking modules, fully implemented across desktop, web, and API:

| Module | What it tracks |
|--------|---------------|
| **Feeds** | 🍼 Type (bottle/breast/solid), date, time, quantity, notes, duration |
| **Sleep** | 😴 Type (nap/night), start/end time, duration, quality, notes |
| **Growth** | 📏 Weight (kg), height (cm), head circumference (cm), notes |
| **Susu-Poty** | 🧷 Type (wet/dirty/mixed), date, time, notes |

Each module exposes REST endpoints: `GET /api/{module}`, `POST /api/{module}`, `GET /api/{module}/{id}`

Data stored as JSON files in `~/.babytracker/` — see [docs/man.md](docs/man.md) for the full API and function reference.

---

## 📱 Mobile (PWA)

The web app is a Progressive Web App — installable on Android 🤖 and iOS 🍎:

1. Open the web app on your phone's browser
2. Tap "Add to Home Screen" (or browser menu > Install)
3. The app runs in standalone mode, like a native app

For app store distribution, [Capacitor](https://capacitorjs.com/) can wrap the build.

---

## 🔮 Roadmap

See **[docs/ROADMAP.md](docs/ROADMAP.md)** for the full roadmap, tech debt tracker, and architectural evolution path.

**Current:** v0.4 (Tailwind redesign, dashboard, themes) | **Next:** v0.3.2 (bug fixes) → v0.4.1 (multi-child) → v0.5 (life journal) → v0.8 (adult mode)

---

## 💡 Why This Project?

Because caregiving deserves clean tools. Because learning Go should be hands-on. And because nothing says "I'm a full-stack developer" like tracking both baby food and baby poop in the same GUI.

- **Practical Tool**: Genuine utility for parents and caregivers 👪
- **Learning Vehicle**: Hands-on Go development with real-world complexity 🐹
- **Architecture Demo**: Clean code principles and modular design ✨
- **Community Building**: Open source collaboration with purpose 🤝

---

## 🙌 Powered By

Using Fyne & Go 🐹
Crafted with 👶, ☕️, and code 💻 by **Nick**
With a little 🤏🧸🐭 little help from
- [Microsoft Copilot](https://copilot.microsoft.com)
- [Claude](https://claude.ai/)
- [Gemini](https://gemini.google.com/)
- [Grok](https://grok.com/)
- [Dia](https://www.diabrowser.com)

No [ChatGPT](https://chatgpt.com)'s were harmed building this. 🫡

---
*👣 "In the grand chronicles of parenthood, every logged feed and tracked nap becomes a story of love, care, and growth."* ✨
