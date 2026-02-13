# Baby Tracker - Project Notes
**Last Updated:** 2026-02-13

---

## Project Purpose

Baby Tracker is a multi-platform application for tracking a baby's daily activities:
- **Feeds** (bottle, breast, solid)
- **Sleep** (naps, night sleep, quality)
- **Growth** (weight, height, head circumference)
- **Diaper changes** (the "Susu-Poty Chronicles")

Available as a desktop app (Fyne/Go), web app (React), and mobile PWA.

---

## Architecture

```
BabyTracker/
  cmd/
    desktop/main.go           # Fyne desktop entry point
    api/main.go               # HTTP API server entry point
  internal/
    models/                   # Shared data models (all 4 features)
      feed.go, sleep.go, growth.go, diaper.go
      *_test.go               # Model unit tests
    storage/                  # JSON file persistence (~/.babytracker/)
      storage.go              # Generic save/load for all entity types
      storage_test.go         # Storage round-trip tests
    api/                      # HTTP handlers + router (package api)
      router.go               # Gorilla mux routes + CORS middleware
      handlers.go             # Feed handlers
      sleep_handlers.go       # Sleep handlers
      growth_handlers.go      # Growth handlers
      diaper_handlers.go      # Diaper handlers
      handlers_test.go        # API handler tests
    desktop/                  # Fyne UI (package desktop)
      app.go                  # App lifecycle, window setup
      layout.go               # Main tabbed layout
      tabs/
        feeds.go              # Feed logging form + quick actions
        sleep.go              # Sleep logging form + quick actions
        growth.go             # Growth measurement form
        susupoty.go           # Diaper logging form + quick actions
  web/                        # React SPA (calls API server)
    public/
      index.html              # HTML entry + PWA meta tags
      manifest.json           # PWA manifest (installable on mobile)
      icon-192.png, icon-512.png
    src/
      index.jsx               # React entry + service worker registration
      App.jsx                 # Router wrapper
      Main.jsx                # Navigation layout (React Router Links)
      Routes.jsx              # Route definitions
      api.js                  # API client (feeds, sleep, growth, diapers)
      config.js               # API_BASE configuration
      serviceWorkerRegistration.js  # PWA offline support
      components/
        Feeds.js              # Feed form, connected to API, recent entries
        Sleep.js              # Sleep form, connected to API, recent entries
        Growth.js             # Growth form, connected to API, recent entries
        SusuPoty.js           # Diaper form, connected to API, recent entries
        index.js              # Barrel exports
```

### Data Flow

```
Fyne Desktop --> storage (direct JSON read/write)
React Web    --> Go API Server --> storage (JSON read/write)
Mobile PWA   --> Go API Server --> storage (same as web)
```

All platforms share `internal/models` and `internal/storage`.

---

## Running the App

```bash
# Desktop (Fyne GUI)
go run ./cmd/desktop

# API Server (for web + mobile)
go run ./cmd/api
# Runs on http://localhost:8080

# Web Frontend (React dev server)
cd web && npm start
# Runs on http://localhost:3000, proxies API to :8080

# Build web for production
cd web && npx react-scripts build

# Run all tests
go test ./internal/...
```

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/feeds | List all feeds |
| POST | /api/feeds | Log a feed |
| GET | /api/feeds/{id} | Get feed by ID |
| GET | /api/sleep | List all sleep entries |
| POST | /api/sleep | Log sleep |
| GET | /api/sleep/{id} | Get sleep entry by ID |
| GET | /api/growth | List all growth entries |
| POST | /api/growth | Log growth |
| GET | /api/growth/{id} | Get growth entry by ID |
| GET | /api/diapers | List all diaper entries |
| POST | /api/diapers | Log diaper change |
| GET | /api/diapers/{id} | Get diaper entry by ID |

---

## Data Models

- **FeedEntry**: ID, Date, Time, Type (Bottle/Breast/Solid), Quantity, Notes, Duration
- **SleepEntry**: ID, Date, StartTime, EndTime, Duration, Type (Nap/Night), Quality (Good/Fair/Poor), Notes
- **GrowthEntry**: ID, Date, Weight (kg), Height (cm), HeadCircumference (cm), Notes
- **DiaperEntry**: ID, Date, Time, Type (Wet/Dirty/Mixed), Notes

Storage: JSON files in `~/.babytracker/` (feeds.json, sleep.json, growth.json, diapers.json)

---

## Mobile / PWA

The web app is configured as a Progressive Web App:
- `manifest.json` with standalone display mode
- Service worker registration for offline caching
- App icons for home screen installation
- Responsive layout (max-width 500px forms)

Users can "Add to Home Screen" on Android/iOS for an app-like experience.

For app store distribution, Capacitor can wrap the React build:
```bash
cd web
npx cap init
npx cap add android
npx cap add ios
```

---

## Development History

### v0.1 (Initial)
- Fyne desktop app with feeds tab only
- FeedEntry model + JSON storage

### v0.2.1 (Web + API scaffold)
- React web app scaffolded with feeds form
- Go API server with feed endpoints
- Web form not yet connected to API

### v0.3 (Current - Full Implementation)
- **Restructured** to standard Go layout (`cmd/`, `internal/`)
- **All 4 features complete**: Feeds, Sleep, Growth, Diapers
  - Models with type constants and helper methods
  - Generic JSON storage with proper ID generation
  - REST API endpoints for all features
  - Desktop Fyne forms with quick actions for all features
  - React web forms connected to API with recent entries display
- **PWA support**: manifest, service worker, icons
- **Navigation fixed**: React Router Links instead of `<a href>`
- **Tests**: Model tests, storage round-trip tests, API handler tests

### Future
- History views and charts
- Reminders and notifications
- Multi-profile support
- Export features (CSV, PDF)
- CI/CD integration
- Replace placeholder PWA icons with proper branding

---

## References

- [Fyne Documentation](https://docs.fyne.io/)
- [Go Modules Layout](https://go.dev/doc/modules/layout)
- [React Documentation](https://react.dev/)
- [PWA Guide](https://web.dev/progressive-web-apps/)
