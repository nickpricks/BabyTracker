# 🗺️ Roadmap

## ✅ Completed

| Version | Milestone | Key Deliverables |
|---------|-----------|-----------------|
| v0.1 | 🏗️ Foundation | Fyne window, basic UI scaffold |
| v0.2 | 🍼 Feed Tracking | FeedEntry model, JSON persistence, desktop form |
| v0.3 | 🚀 Full Platform + Security | All 4 modules, standard Go layout, REST API, React SPA, PWA, API auth, CORS lockdown, storage mutex, file permissions ([details](Security-Review.md)) |
| v0.4 | 🎨 Tailwind Redesign | Dashboard, 3-theme system, bench tooling, CORS rewrite |

## 🔜 Up Next

| Version | Milestone | Key Deliverables |
|---------|-----------|-----------------|
| **v0.3.2** | 🐛 **Bug Fixes & Debt** | Desktop recent activity (wire up `Load*()` in all 4 tabs), API timezone parsing (custom `FlexTime` unmarshaler, remove `Z` workaround), DELETE/PUT endpoints |
| **v0.3.3** | 🧪 **Test Coverage** | API get-by-ID handlers, storage init paths, desktop `fyne.test` headless tests, web edge cases |
| **v0.3.4** | 🖥️ **Desktop Polish** | App menus (Help, About, version info), keyboard shortcuts (Ctrl+R refresh) |
| **v0.4.0** | 👶👶 **Multi-Child** | `children.json` profiles, directory-per-child storage (`~/.babytracker/{child}/`), API routes `/api/{child}/{resource}`, child switcher in web + desktop |
| v0.4.1 | 📝 Docs Sweep | All docs dir document - thoroly examined and re-written if need be - according to wahtever purpose they surve |
| **v0.4.2** | 📦 **Export/Import** | Versioned JSON envelope with metadata, child-aware bundles, ID remapping on import, web UI for download/upload |
| **v0.4.3** | 🧒 **Toddler Modules** | Meals/Nutrition (replaces feeds), Potty Training (replaces diapers), stage-based module visibility (infant/toddler) |
| **v0.4.4** | 🏆 **Milestones & Firsts** | First words, first steps, custom milestones -- timestamped, media-attachable, feeds directly into Life Journal |
| **v0.4.5** | 📝 **Docs Sweep** | Update all docs to reflect multi-child, export/import, toddler modules, milestones |
| **v0.5** | 📖 **Life Journal** | Configurable daily/weekly summary view, milestone cards, shareable screenshot cards, history views with charts |

## 🔮 Future

| Version | Milestone | Key Deliverables |
|---------|-----------|-----------------|
| v0.6 | 🔔 Smart Alerts | Configurable reminders, feeding interval notifications, sleep schedule suggestions |
| v0.7 | 💾 Storage & Sync | SQLite migration, indexing, query engine + optional cloud backend (Firebase), real-time sync, conflict resolution, backup/restore |
| v0.8 | 🧘 Adult Mode | Rebranded as "Body Soul and Mind Tracker" -- generalized health tracking using AFP activities as base, plugin architecture |
| v0.9 | 📦 App Distribution | App publishing (macOS .app bundle, Homebrew tap), auto-update channel, Check for Updates |
| v1.0 | 📱 Platform Expansion | Native mobile via Capacitor, widgets, voice integration |

## 🏛️ Architectural Evolution Path

- **Current (v0.4)**: Flat JSON files per child, single-user, local-first
- **Near-term (v0.5-v0.6)**: Multi-child directories, export/import bundles, toddler modules
- **Mid-term (v0.7)**: SQLite + Firebase sync layer
- **Long-term (v0.8+)**: Generalized tracker, cloud-first, plugin SDK

## ⚠️ Known Limitations & Technical Debt

1. 🚫 No DELETE or PUT endpoints -- entries cannot be edited or removed via API
2. 🔁 `layout.go` duplicates tab creation logic already in `app.go`
3. 📋 Web components have no shared form abstraction (each is ~200 lines of similar code)
4. 🗄️ No database migration path -- switching from JSON requires manual data conversion
5. 📴 Service worker caches the shell but not API responses (no offline data access)
6. 🖼️ Placeholder PWA icons -- need proper branding assets
7. ⚙️ No CI/CD pipeline configured
8. 🚦 Go API has no request rate limiting
9. 👻 Desktop "Recent Activity" sections are static placeholders
