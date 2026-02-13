# 👶 Baby Tracker

A minimalist yet mighty multi-platform app to track your baby's day-to-day growth 📈, nourishment 🍼, sleep patterns 😴, and the noble chronicles of 🚽 Susu-Poty 🧻. Built in Go with a [Fyne](https://fyne.io) desktop GUI, a React web app, and a REST API — all sharing the same backend logic.

Whether you're a caregiver, parent, or curious builder, this project aims to balance usability with technical learning—while keeping things fun and purpose-driven 🎯.

---

## 🛠️ Prerequisites

- [Go](https://go.dev/dl/) (v1.22 or later)
- [Node.js](https://nodejs.org/) (v18 or later) — for the web frontend
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

### React Web App (`web/.env`)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | React dev server port |
| `REACT_APP_API_BASE` | `http://localhost:8080/api` | API URL the web app connects to |

The Makefile automatically loads `.env` from the project root, so `make api` and `make desktop` pick up your values. The React dev server reads `web/.env` natively (CRA built-in).

Config is also available in Go code via `internal/config`:

```go
cfg, _ := config.Load()
fmt.Println(cfg.APIPort)  // "8080"
fmt.Println(cfg.DataDir)  // "/Users/you/.babytracker"
```

---

## 📋 Makefile Commands

Run `make` or `make help` to see all available commands:

### Development (run)

| Command | Description |
|---------|-------------|
| `make dev` | Run API server + web dev server concurrently |
| `make api` | Run only the Go API server |
| `make desktop` | Run the Fyne desktop app |
| `make web` | Run the React dev server (requires API running) |

### Build

| Command | Description |
|---------|-------------|
| `make build` | Build everything (Go binaries + web production build) |
| `make build-api` | Build the API server binary to `bin/api` |
| `make build-desktop` | Build the desktop app binary to `bin/desktop` |
| `make build-web` | Build the React app for production to `web/build/` |

### Test

| Command | Description |
|---------|-------------|
| `make test` | Run all Go tests |
| `make test-v` | Run all Go tests in verbose mode |
| `make test-cover` | Run Go tests with coverage report |
| `make test-web` | Run React tests |

### Lint & Tidy

| Command | Description |
|---------|-------------|
| `make lint` | Vet Go code (`go vet`) |
| `make lint-web` | Lint React code (eslint) |
| `make tidy` | Run `go mod tidy` |

### Setup & Clean

| Command | Description |
|---------|-------------|
| `make setup` | Full project setup (env files + Go tidy + web deps) |
| `make env` | Create `.env` files from examples (won't overwrite) |
| `make install-web` | Install web dependencies (`npm install`) |
| `make clean` | Remove all build artifacts (`bin/` + `web/build/`) |
| `make clean-bin` | Remove Go binaries only |
| `make clean-web` | Remove web build output only |

---

## 🏗️ Project Structure

```
BabyTracker/
├── cmd/
│   ├── desktop/main.go        # Fyne desktop entry point
│   └── api/main.go            # HTTP API server entry point
├── internal/
│   ├── models/                # Shared data models
│   │   ├── feed.go            #   FeedEntry + type constants
│   │   ├── sleep.go           #   SleepEntry + type/quality constants
│   │   ├── growth.go          #   GrowthEntry
│   │   └── diaper.go          #   DiaperEntry + type constants
│   ├── storage/               # JSON file persistence (~/.babytracker/)
│   │   └── storage.go         #   Generic load/save for all entities
│   ├── api/                   # HTTP handlers & router
│   │   ├── router.go          #   Gorilla mux routes + CORS
│   │   ├── handlers.go        #   Feed endpoints
│   │   ├── sleep_handlers.go  #   Sleep endpoints
│   │   ├── growth_handlers.go #   Growth endpoints
│   │   └── diaper_handlers.go #   Diaper endpoints
│   └── desktop/               # Fyne UI
│       ├── app.go             #   App lifecycle & window setup
│       ├── layout.go          #   Tabbed layout
│       └── tabs/              #   One file per feature tab
├── web/                       # React SPA
│   ├── public/                #   HTML, manifest, icons
│   └── src/
│       ├── components/        #   Feeds, Sleep, Growth, SusuPoty
│       ├── api.js             #   API client functions
│       ├── config.js          #   API base URL config
│       └── ...                #   App, Routes, Main layout
├── Makefile                   # Build, run, test commands
├── Notes.md                   # Detailed architecture notes
└── go.mod
```

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

---

## 🌱 Core Modules

All four modules are fully implemented across desktop, web, and API:

| Module | What it tracks |
|--------|---------------|
| **Feeds** | 🍼 Type (bottle/breast/solid), date, time, quantity, notes, duration |
| **Sleep** | 😴 Type (nap/night), start/end time, duration, quality, notes |
| **Growth** | 📏 Weight (kg), height (cm), head circumference (cm), notes |
| **Susu-Poty** | 🧷 Type (wet/dirty/mixed), date, time, notes |

### API Endpoints

Each module exposes three endpoints:

| Method | Pattern | Example |
|--------|---------|---------|
| `GET` | `/api/{module}` | `GET /api/feeds` — list all |
| `POST` | `/api/{module}` | `POST /api/sleep` — log new entry |
| `GET` | `/api/{module}/{id}` | `GET /api/growth/3` — get by ID |

Modules: `feeds`, `sleep`, `growth`, `diapers`

### Data Storage

JSON files in `~/.babytracker/`:
- `feeds.json`, `sleep.json`, `growth.json`, `diapers.json`

---

## 📱 Mobile (PWA)

The web app is a Progressive Web App — installable on Android 🤖 and iOS 🍎:

1. Open the web app on your phone's browser
2. Tap "Add to Home Screen" (or browser menu > Install)
3. The app runs in standalone mode, like a native app

For app store distribution, [Capacitor](https://capacitorjs.com/) can wrap the build.

---

## 🔮 Roadmap

| Version | Status | Features |
|---------|--------|----------|
| `v0.1` | ✅ Done | Initial Fyne window with basic UI |
| `v0.2` | ✅ Done | Modular architecture, Feed Tracker with persistence |
| `v0.3` | ✅ Done | All 4 modules complete, standard Go layout, API for all features, React connected to API, PWA, tests |
| `v0.4` | 🔜 Planned | History views with charts, pattern analytics |
| `v0.5` | 🔜 Planned | Reminders, notifications |
| `v1.0` | ⏳ Future | Multi-profile support, exportable reports, dark mode |
| `v2.0+` | 🚀 Vision | Adult Mode: rebranded as Body Soul and Mind Tracker |

---

## 📚 Resources

- [Go Tour](https://go.dev/tour/) — Interactive introduction to Go
- [Fyne Documentation](https://docs.fyne.io/) — GUI framework guide
- [React Documentation](https://react.dev/) — Web framework
- [PWA Guide](https://web.dev/progressive-web-apps/) — Progressive Web Apps
- [Gorilla Mux](https://github.com/gorilla/mux) — HTTP router

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

---
*👣 "In the grand chronicles of parenthood, every logged feed and tracked nap becomes a story of love, care, and growth."* ✨
