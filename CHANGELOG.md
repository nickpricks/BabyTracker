# Changelog

All notable changes to BabyTracker will be documented in this file.

## [v0.3.2] — 2026-04-06

### Bug Fixes
- **FlexTime type** — custom `UnmarshalJSON` accepts both RFC3339 (`T10:30:00Z`) and timezone-less (`T10:30:00`) timestamps
- **Removed Z-hack** — web components no longer append `Z` to time strings (Feeds, Sleep, SusuPoty)
- **Desktop Recent Activity** — all 4 tabs now load and display real data (newest-first, all entries), refresh after each save
- **Auth test fix** — `api.test.js` properly mocks `VITE_API_KEY` for the no-auth test case

### New Features
- **DELETE/PUT endpoints** — all 4 resources support update and delete via API
- **Paginated API** — list endpoints return `{ items, total, limit, offset }` (newest-first), support `?limit=N&offset=M`
- **Infinite scroll (web)** — `useLoadMore` hook with IntersectionObserver, escalating pagination (5 → 50 after 3 loads)
- **Web API client** — added `apiPut`, `apiDelete`, per-resource `updateX`/`deleteX` functions
- **CORS** — allowed methods now include PUT and DELETE

### Docs & Roadmap
- **Roadmap restructured** — consolidated 3 stale roadmaps (README, Manual, TLDR) into single `docs/ROADMAP.md`
- **Version plan overhauled:**
  - v0.3 + v0.3.1 combined into "Full Platform + Security"
  - v1.0 (Storage) + v1.5 (Cloud Sync) merged → **v0.7 Storage & Sync**
  - v3.0 (Adult Mode) moved → **v0.8** (AFP activities as base)
  - v2.0 (Platform Expansion) → **v1.0** (true 1.0 milestone)
  - Added **v0.3.4** Desktop Polish (menus, about, shortcuts)
  - Added **v0.9** App Distribution (publishing, auto-update channel)
- **Full docs sweep** — all 8 docs files updated (versions, CORS, Tailwind, themes, vitest, Vite)
- **CLAUDE.md** — updated Go 1.26.0/Fyne v2.7.3, added Planning & Docs section
- **RESUME.md** added to `.gitignore`

## [v0.4] — 2026-04-02

- Tailwind CSS v4 redesign, Dashboard component, 3-theme system (Lullaby, Nursery_OS, Midnight Feed)
- Bench data generator (`cmd/bench/`)
- CORS rewrite (external handler wrapping mux)
- Makefile additions: `make bench`, `make bench-restore`, `make update`

## [v0.3.1] — 2026-03-30

- Security hardening: API auth (Bearer token), CORS origin lockdown, storage mutex, file permissions (0600), body size limits

## [v0.3] — 2026-03-28

- All 4 modules (Feed, Sleep, Growth, Diaper), standard Go layout, REST API, React SPA, PWA, tests

## [v0.2] — 2026-03-25

- Feed tracking module, FeedEntry model, JSON persistence, desktop form

## [v0.1] — 2026-03-20

- Initial Fyne window, basic UI scaffold
