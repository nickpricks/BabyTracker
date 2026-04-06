# BabyTracker Security & Vulnerability Review

**Date:** 2026-03-27 (updated from 2026-03-10 original)
**Scope:** Full codebase review -- Go API server, Fyne desktop app, React web frontend, JSON file storage, configuration
**Reviewed by:** 7-agent parallel review (code quality, silent failure analysis, comment quality, security drill, mentor review, implementation verification, comprehensive PR review)

---

## Executive Summary

| Severity | Count |
|---|---|
| Critical | 3 |
| High | 8 |
| Medium | 11 |
| Low | 7 |
| Informational | 5 |
| **Total** | **34** |

The 7-agent review (2026-03-27) confirmed all 23 original findings and uncovered **11 additional issues**, bringing the total to 34. As of v0.4, **12 findings have been fully fixed** and 1 partially fixed: API key authentication (FINDING-01), configurable CORS with localhost wildcard (FINDING-03), request body size limit (FINDING-08), storage mutex (FINDING-12), directory/file permissions (FINDING-13), CRA-to-Vite migration (FINDING-17), no silent data destruction (FINDING-24), atomic file writes (FINDING-25), thread-safe lazy init (FINDING-26), fetch error display (FINDING-28), Error Boundary (FINDING-29), hermetic API tests (FINDING-35), and PII file permissions partially addressed (FINDING-09). Remaining critical/high items: desktop errors invisible to users (FINDING-32), no TLS (FINDING-04), PII encryption at rest (FINDING-09), and no authorization (FINDING-02).

---

## 1. Authentication & Authorization

### FINDING-01: No Authentication on API Endpoints
- **Status:** [x] Fixed (2026-03-27) -- Bearer token auth middleware added; key from `API_KEY` env var; skipped when empty for dev convenience
- **Severity:** Critical
- **Agents flagged:** 4/7
- **Files:** `internal/api/router.go`, `internal/config/config.go`, `cmd/api/main.go`, `web/src/api.js`
- **Description:** The entire API server has zero authentication. Every endpoint (GET and POST for `/api/feeds`, `/api/sleep`, `/api/growth`, `/api/diapers`) is accessible to any client that can reach the server. There is no token, API key, session cookie, or any other auth mechanism.
- **Recommended Fix:** Implement authentication middleware. At minimum, add API key or bearer token authentication checked via middleware. For a single-user app, a static shared secret in `.env` is infinitely better than nothing.

### FINDING-02: No Authorization / Access Control
- **Severity:** High
- **Agents flagged:** 4/7
- **Files:** `internal/api/router.go` (lines 10-48)
- **Description:** Even if authentication were added, there is no concept of users or access control. All data is shared globally. Any authenticated user would have full read/write access to all records.
- **Recommended Fix:** If multi-user support is planned, implement user-scoped data isolation. If single-user, at minimum add authentication to prevent unauthorized access.

---

## 2. CORS & Network Security

### FINDING-03: Wildcard CORS Origin (`Access-Control-Allow-Origin: *`)
- **Status:** [x] Fixed (2026-03-27, updated v0.4) -- configurable `CORS_ORIGIN` env var (default `http://localhost:3000`); v0.4 rewrote CORS as an external `corsHandler` wrapping the mux router (so OPTIONS preflight is intercepted before mux's method matching), with localhost wildcard matching (any `http://localhost:*` port accepted when configured origin is localhost)
- **Severity:** High
- **Agents flagged:** 6/7
- **Files:** `internal/api/router.go`, `internal/config/config.go`
- **Description:** The CORS middleware sets `Access-Control-Allow-Origin: *`, meaning any website on the internet can make API requests to this server from a user's browser. Combined with FINDING-01, any malicious website could silently read or modify baby tracking data.
- **Recommended Fix:** Replace `*` with a specific allowed origin from an environment variable (`CORS_ORIGIN`). Default to `http://localhost:3000` for development.

### FINDING-04: No TLS / HTTP Only
- **Severity:** High
- **Files:** `cmd/api/main.go` (line 25)
- **Description:** The server uses `http.ListenAndServe` (plaintext HTTP). All data including baby health information is transmitted unencrypted. This is especially concerning over WiFi.
- **Recommended Fix:** Use `http.ListenAndServeTLS` with proper certificates, or deploy behind a reverse proxy (Caddy, nginx) that terminates TLS.

### FINDING-05: Missing Security Headers
- **Severity:** Medium
- **Files:** `internal/api/router.go` (lines 14-25)
- **Description:** The API responses lack standard security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Strict-Transport-Security`, `Content-Security-Policy`.
- **Recommended Fix:** Add a middleware that sets security headers on all responses.

### FINDING-06: Server Binds to All Interfaces
- **Severity:** Medium
- **Files:** `cmd/api/main.go` (line 25)
- **Description:** `http.ListenAndServe(":"+cfg.APIPort, r)` binds to `0.0.0.0` (all network interfaces). The unauthenticated API is accessible to every device on the LAN.
- **Recommended Fix:** Default to `127.0.0.1:`+port. Add a `BIND_ADDR` config option for intentional exposure.

---

## 3. Input Validation & Injection

### FINDING-07: Insufficient Server-Side Input Validation
- **Severity:** Medium
- **Agents flagged:** 5/7
- **Files:** `internal/api/handlers.go` (lines 35-51), `sleep_handlers.go`, `growth_handlers.go`, `diaper_handlers.go`
- **Description:** Input validation is minimal -- only checks for empty required string fields. No validation of: date format (could be `"banana"`), type values against allowed constants, numeric ranges (weight = -999), or string length limits.
- **Recommended Fix:** Validate date against `YYYY-MM-DD`. Validate `type` against defined constants (`FeedTypeBottle`, etc.). Add reasonable bounds for numeric fields. Cap string lengths.

### FINDING-08: No Request Body Size Limit
- **Status:** [x] Fixed (2026-03-27) -- `http.MaxBytesReader` middleware (1MB) applied to all requests in router.go
- **Severity:** High (upgraded from Medium -- 4/7 agents flagged)
- **Files:** `internal/api/router.go`
- **Description:** `json.NewDecoder(r.Body).Decode(&entry)` reads the full request body without size limits. A malicious client can send a multi-gigabyte payload and exhaust server memory.
- **Recommended Fix:** Wrap `r.Body` with `http.MaxBytesReader(w, r.Body, 1<<20)` (1MB) before decoding. Can be applied as middleware.

### FINDING-30: Desktop App Bypasses API Validation
- **Severity:** Medium
- **Agents flagged:** 2/7
- **Files:** `internal/desktop/tabs/feeds.go`, `sleep.go`, `growth.go`, `susupoty.go`
- **Description:** The desktop app writes directly to storage, bypassing the API's validation checks. Users can save entries with empty type, empty date, or other invalid data that the API handlers would reject. This creates inconsistent data that the web UI will display.
- **Recommended Fix:** Add validation in desktop button handlers before saving, or move validation into the storage layer so both paths enforce it.

---

## 4. Data Storage & Integrity

### FINDING-12: Race Condition in Storage (Read-Modify-Write Without Locking)
- **Status:** [x] Fixed (2026-03-27) -- `sync.Mutex` added to `StorageManager`; all `Save*` functions lock before read-modify-write. `getStorage()` uses `sync.Once` (FINDING-26 bonus).
- **Severity:** High
- **Agents flagged:** 7/7 (highest consensus finding)
- **Files:** `internal/storage/storage.go`
- **Description:** Every `Save*` function follows a read-modify-write pattern with zero concurrency protection. Two concurrent POST requests will both read the same file, both append, and the second write silently overwrites the first -- losing one entry. The API server handles requests concurrently via gorilla/mux.
- **Recommended Fix:** Add a `sync.Mutex` per resource (or one global mutex) to `StorageManager`. Lock during the entire read-modify-write cycle. For cross-process safety, consider `flock`.

### FINDING-24: Storage Silently Destroys All Data on Corrupt JSON Load *(NEW)*
- **Status:** [x] Fixed (2026-03-27) -- `Save*` functions now return error on corrupt/unreadable JSON instead of silently replacing with empty slice; only file-not-exists is treated as empty
- **Severity:** Critical
- **Agents flagged:** 2/7
- **Files:** `internal/storage/storage.go`
- **Description:** Every `Save*` function catches load errors by falling back to an empty slice: `if err != nil { feeds = []FeedEntry{} }`. If the JSON file is corrupt (partial write from FINDING-25, disk error, truncation), the save function silently replaces the entire file with a single new entry -- **permanently destroying all existing data**. There is no log, no error surfaced, no recovery path.
- **Impact:** A parent logs a feeding and silently loses their entire feed history. They will not know until they notice the history is gone.
- **Recommended Fix:** Propagate load errors to the caller. Only treat `os.IsNotExist` as "empty file" for first-run. If the file exists but is unreadable, return an error -- do not overwrite.

### FINDING-25: Non-Atomic File Writes *(NEW)*
- **Status:** [x] Fixed (2026-03-27) -- `saveJSON` now writes to `.tmp` file then does atomic `os.Rename`; file permissions tightened to `0600`
- **Severity:** High
- **Agents flagged:** 3/7
- **Files:** `internal/storage/storage.go`
- **Description:** `os.WriteFile` directly truncates and overwrites the target file. If the process crashes, loses power, or the disk fills mid-write, the file will contain partial/invalid JSON. On next load, this triggers FINDING-24 (silent total data loss).
- **Recommended Fix:** Write to a temp file first, then atomically rename:
  ```go
  tmpPath := filePath + ".tmp"
  os.WriteFile(tmpPath, data, 0600)
  os.Rename(tmpPath, filePath) // atomic on same filesystem
  ```

### FINDING-26: `getStorage()` Lazy Initialization Not Thread-Safe *(NEW)*
- **Status:** [x] Fixed (2026-03-27) -- `getStorage()` now uses `sync.Once` for thread-safe lazy initialization
- **Severity:** Medium
- **Agents flagged:** 4/7
- **Files:** `internal/storage/storage.go`
- **Description:** `getStorage()` reads and writes the package-level `globalStorage` variable without synchronization. Two goroutines could both see `nil` and both call `NewStorageManager()`, causing double directory creation or inconsistent state.
- **Recommended Fix:** Use `sync.Once` for lazy initialization.

### FINDING-13: Data Directory Created With Overly Permissive Permissions
- **Status:** [x] Fixed (2026-03-27) -- directory permissions changed to `0700`, file permissions to `0600`
- **Severity:** Low
- **Files:** `internal/storage/storage.go`
- **Description:** The data directory is created with `0755` permissions and data files with `0644`. On multi-user systems, other users can read the baby's health data.
- **Recommended Fix:** Use `0700` for the directory and `0600` for data files.

### FINDING-14: No Path Validation on DATA_DIR Environment Variable
- **Severity:** Low
- **Files:** `internal/config/config.go` (lines 39-40)
- **Description:** The `DATA_DIR` environment variable is used directly without validation. A misconfigured value could point to sensitive system directories.
- **Recommended Fix:** Validate that `DATA_DIR` is an absolute path without `..` traversal sequences.

---

## 5. Data Exposure & PII

### FINDING-09: PII Stored in Plaintext JSON Without Encryption
- **Status:** [ ] Partially fixed (2026-03-27) -- file permissions now `0600`, directory `0700` (via FINDING-13/FINDING-25 fixes); encryption at rest not yet implemented
- **Severity:** High
- **Agents flagged:** 3/7
- **Files:** `internal/storage/storage.go`
- **Description:** All baby health data is stored as plaintext JSON in `~/.babytracker/`. This data constitutes sensitive health information about a minor. Files are created with `0644` permissions (world-readable).
- **Recommended Fix:** Set file permissions to `0600` (owner-only). Consider encrypting data at rest. Document that the data directory contains sensitive health information.

### FINDING-10: Internal Errors Leaked to Clients
- **Severity:** Medium
- **Agents flagged:** 5/7
- **Files:** `internal/api/handlers.go` (line 28), `sleep_handlers.go` (line 18), `growth_handlers.go` (line 18), `diaper_handlers.go` (line 18)
- **Description:** On internal errors, the raw `err.Error()` is returned in the JSON response. This leaks internal file paths (`/Users/nick/.babytracker/feeds.json: permission denied`), file system structure, and Go error details.
- **Recommended Fix:** Log the real error server-side. Return a generic `"internal server error"` to the client.

### FINDING-11: Sensitive Data Logged to stdout
- **Severity:** Low
- **Files:** `internal/api/handlers.go` (line 45), `sleep_handlers.go` (line 34), `growth_handlers.go` (line 34), `diaper_handlers.go` (line 34)
- **Description:** Every POST handler logs the full entry with `log.Printf("Log Feed: %+v\n", feed)`, dumping all fields including free-text notes to stdout.
- **Recommended Fix:** Log only non-sensitive fields (ID, type, date) or log at debug level only.

---

## 6. API Reliability

### FINDING-27: `jsonResponse` Silently Discards Encoding Errors *(NEW)*
- **Severity:** Medium
- **Agents flagged:** 3/7
- **Files:** `internal/api/handlers.go` (line 20)
- **Description:** `_ = json.NewEncoder(w).Encode(payload)` explicitly suppresses the encoding error. If encoding fails, the HTTP response will have a 200 status but an empty or partial body. The client receives truncated JSON with no indication of failure.
- **Recommended Fix:** Log the error: `if err := json.NewEncoder(w).Encode(payload); err != nil { log.Printf("ERROR: failed to encode response: %v", err) }`

### FINDING-31: No Graceful Shutdown *(NEW)*
- **Severity:** Medium
- **Agents flagged:** 2/7
- **Files:** `cmd/api/main.go` (line 25)
- **Description:** `http.ListenAndServe` blocks until error. SIGINT/SIGTERM kills the process immediately. In-flight requests may be dropped and file writes (FINDING-25) could be interrupted mid-save, corrupting JSON data.
- **Recommended Fix:** Use `http.Server.Shutdown(ctx)` with signal handling to drain in-flight requests.

---

## 7. Denial of Service

### FINDING-15: No Rate Limiting
- **Severity:** Medium
- **Files:** `internal/api/router.go` (all routes)
- **Description:** There is no rate limiting on any endpoint. An attacker could flood the API with POST requests, rapidly filling disk space with JSON entries.
- **Recommended Fix:** Add rate limiting middleware (e.g., 60 requests/minute per IP).

### FINDING-16: No Pagination on List Endpoints
- **Severity:** Medium
- **Files:** `internal/api/handlers.go` (line 26), all list handlers
- **Description:** All list endpoints return every record ever stored. After a year of use, this could be thousands of entries in a single response.
- **Recommended Fix:** Implement pagination (`?page=1&limit=50`) or limit to most recent N entries.

---

## 8. Frontend Security

### FINDING-28: React Fetch Errors Silently Swallowed *(NEW)*
- **Status:** [x] Fixed (2026-03-27) -- all 4 components now set `fetchError` state and display error message when API fetch fails
- **Severity:** High
- **Agents flagged:** 2/7
- **Files:** `web/src/components/Feeds.jsx`, `Sleep.jsx`, `Growth.jsx`, `SusuPoty.jsx`
- **Description:** All four components have identical `catch { // API may not be running - that's okay }` blocks that swallow every fetch error -- network failures, server 500s, CORS rejections, JSON parse errors, even JavaScript runtime errors in the response processing. The user sees "No entries logged yet" with no indication the API failed.
- **Impact:** A parent opens the app, sees an empty list, and re-logs entries that already exist -- creating duplicates. Or worse, assumes the data was lost.
- **Recommended Fix:** Set an error state and display it: "Could not load recent entries. Is the API server running?"

### FINDING-29: No React Error Boundary *(NEW)*
- **Status:** [x] Fixed (2026-03-27) -- `ErrorBoundary` component wraps `<AppRoutes />` in `App.jsx`; shows fallback UI with refresh button
- **Severity:** Medium
- **Agents flagged:** 1/7
- **Files:** `web/src/App.jsx`, `web/src/components/ErrorBoundary.jsx`
- **Description:** No React Error Boundary exists anywhere in the component tree. If any component throws during rendering (e.g., `entries.slice` on a null value from a malformed API response), the entire app white-screens with no recovery path.
- **Recommended Fix:** Add an Error Boundary wrapping `<AppRoutes />` with a meaningful fallback UI.

### FINDING-19: XSS Defense Relies on React's Default Escaping
- **Severity:** Low
- **Files:** `web/src/components/Feeds.jsx`, `Growth.jsx`, `Sleep.jsx`, `SusuPoty.jsx`
- **Description:** React components render API data directly in JSX. React's default behavior escapes strings, so this is **not currently exploitable**. However, the lack of server-side sanitization (FINDING-07) means any future use of unsafe HTML rendering would make XSS trivially exploitable.
- **Recommended Fix:** Continue using React's default JSX escaping. Consider server-side type validation as defense in depth.

### FINDING-20: Error Messages from API Displayed to User
- **Severity:** Informational
- **Files:** `web/src/components/Feeds.jsx`, all components
- **Description:** API error messages are displayed directly to users. Combined with FINDING-10 (internal errors leaked), this could show internal server details.
- **Recommended Fix:** Show user-friendly error messages; log detailed errors to the browser console.

---

## 9. Desktop Application Security

### FINDING-32: Desktop Errors Invisible to Users *(NEW)*
- **Severity:** High
- **Agents flagged:** 6/7
- **Files:** `internal/desktop/tabs/feeds.go` (line 94), `sleep.go` (line 101), `growth.go` (line 69), `susupoty.go` (line 73)
- **Description:** All four desktop tabs use `fmt.Printf("Error saving ...: %v\n", err)` when storage fails. In a desktop GUI application, there is no visible terminal. The user clicks "Log Feed," the save fails silently, and the entry is permanently lost with zero feedback.
- **Impact:** A parent taps "Log Feed" after a 3am feeding, the save fails, the entry is gone forever.
- **Recommended Fix:** Use Fyne's `dialog.ShowError(err, window)`. Requires passing the window reference to tab constructors.

### FINDING-33: Desktop `binding.Get()` Errors Universally Suppressed *(NEW)*
- **Severity:** Low
- **Agents flagged:** 1/7
- **Files:** All desktop tab files (every `binding.Get()` call)
- **Description:** Every call to `binding.Get()` discards the error with `_`. If a binding is corrupted, the entry is saved with zero-values (empty strings, 0.0 floats) without the user knowing.
- **Recommended Fix:** Check at least critical bindings and show an error if they fail.

### FINDING-34: Sleep Duration Negative for Overnight Sleep *(NEW)*
- **Severity:** Low
- **Agents flagged:** 3/7
- **Files:** `internal/desktop/tabs/sleep.go` (line 83), `web/src/components/Sleep.jsx`
- **Description:** `endTime.Sub(startTime).Minutes()` produces a negative value when sleep crosses midnight (e.g., 22:00 to 06:00). Both desktop and web are affected. No validation is performed.
- **Recommended Fix:** Check for negative duration and add 24 hours, or validate and show an error.

---

## 10. Testing & Code Integrity

### FINDING-35: API Handler Tests Hit Production Data Directory *(NEW)*
- **Status:** [x] Fixed (2026-03-27) -- `testRouter(t)` helper now calls `storage.Init(t.TempDir())` for hermetic tests
- **Severity:** High
- **Agents flagged:** 7/7 (highest consensus finding, tied with FINDING-12)
- **Files:** `internal/api/handlers_test.go`
- **Description:** `setupTestRouter` creates a `StorageManager` and immediately discards it (`_ = sm`). Tests use the real global storage (`~/.babytracker/`), meaning they read/write to the user's actual data. Tests are non-hermetic, flaky, and can corrupt real data. The storage tests (`storage_test.go`) correctly use `t.TempDir()` but the API tests do not.
- **Recommended Fix:** Override `globalStorage` using `storage.Init(t.TempDir())` in test setup.

---

## 11. Dependencies & Supply Chain

### FINDING-17: react-scripts 5.0.1 Contains Known CVEs
- **Status:** [x] Fixed (2026-03-27) -- migrated from CRA (react-scripts) to Vite 8 + vitest. PWA via vite-plugin-pwa.
- **Severity:** Medium (development dependency)
- **Files:** `web/package.json`, `web/vite.config.js`
- **Description:** `react-scripts@5.0.1` is outdated and pulls in transitive dependencies with known vulnerabilities (`nth-check@1.0.2` -- CVE-2021-3803 ReDoS, `svgo@1.3.2` -- ReDoS). CRA is no longer actively maintained by Meta.
- **Recommended Fix:** Migrate from CRA to Vite or another modern build tool.

### FINDING-18: No Lock File for npm Dependencies
- **Severity:** Low
- **Files:** `web/package.json`, `.gitignore` (lines 31-34)
- **Description:** The `.gitignore` excludes all lock files (`bun.lock`, `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`). Builds are not reproducible -- different installs may pull different (potentially compromised) dependency versions.
- **Recommended Fix:** Commit `bun.lock` (the project uses bun). Remove it from `.gitignore`.

---

## 12. Service Worker & PWA Security

### FINDING-21: Service Worker May Cache Sensitive Data
- **Severity:** Informational
- **Files:** `web/src/serviceWorkerRegistration.js`, `web/src/index.jsx` (line 15)
- **Description:** CRA's default service worker caches assets. Cached data persists in Cache Storage, meaning health data could remain on a shared device even after clearing cookies.
- **Recommended Fix:** Ensure the service worker only caches static assets, not API responses. Document shared-device concerns.

---

## 13. Configuration & Secrets

### FINDING-22: .env Files Properly Excluded from Git
- **Severity:** Informational (positive finding)
- **Files:** `.gitignore` (lines 20-21)
- **Description:** `.gitignore` correctly excludes `.env` and `web/.env` while keeping `.env.example` tracked. No secrets were ever committed to git history. This is correctly configured.

### FINDING-23: PORT Environment Variable Not Validated
- **Severity:** Informational
- **Files:** `internal/config/config.go` (line 34)
- **Description:** `PORT` is used as a raw string without validation. Non-numeric values would cause a clear runtime error from `ListenAndServe`.
- **Recommended Fix:** Validate that `PORT` is numeric and in the valid range (1-65535) during config loading.

---

## Summary — Open Items First

### Open — High Priority

| Finding | Severity | Area | Effort | Impact |
|---------|----------|------|--------|--------|
| Show desktop errors in Fyne dialog (FINDING-32) | High | Desktop | Medium | Users see save failures |
| Enable TLS or document localhost-only (FINDING-04) | High | API | Medium | Encrypts data in transit |
| Encrypt PII at rest (FINDING-09) | High | Storage | Medium | Remaining: perms fixed, encryption not yet |

### Open — Medium Priority

| Finding | Area | Effort |
|---------|------|--------|
| Validate all input fields server-side (FINDING-07) | API | Medium |
| Implement access control (FINDING-02) | API | Medium |
| Stop leaking internal errors to clients (FINDING-10) | API | Small |
| Log `jsonResponse` encoding errors (FINDING-27) | API | Trivial |
| Add security headers middleware (FINDING-05) | API | Small |
| Bind to localhost by default (FINDING-06) | API | Trivial |
| Add graceful shutdown (FINDING-31) | API | Small |
| Add desktop-side validation (FINDING-30) | Desktop | Small |
| Add rate limiting (FINDING-15) | API | Medium |

### Open — Low / Informational

| Finding | Area |
|---------|------|
| Add pagination to list endpoints (FINDING-16) | API |
| Commit `bun.lock` (FINDING-18) | Web |
| Reduce logged PII (FINDING-11) | API |
| Validate DATA_DIR (FINDING-14) | Config |
| Sanitize displayed error messages on frontend (FINDING-20) | Web |
| Handle overnight sleep duration (FINDING-34) | Desktop + Web |
| Check desktop `binding.Get()` errors (FINDING-33) | Desktop |
| XSS defense-in-depth via server-side type validation (FINDING-19) | Web |
| Service worker cache scope (FINDING-21) | Web |
| Validate PORT range (FINDING-23) | Config |

---

### Fixed — 13 findings resolved (2026-03-27)

| Finding | What was done |
|---------|---------------|
| FINDING-01 | Bearer token auth middleware; `API_KEY` env var |
| FINDING-03 | Configurable `CORS_ORIGIN`, default `http://localhost:3000`; v0.4: external `corsHandler` wrapping mux, localhost wildcard matching |
| FINDING-08 | `http.MaxBytesReader` middleware, 1MB limit |
| FINDING-12 | `sync.Mutex` on `StorageManager` for all `Save*` functions |
| FINDING-13 | Directory permissions `0700`, file permissions `0600` |
| FINDING-17 | Migrated from CRA (react-scripts) to Vite 8 + vitest + vite-plugin-pwa |
| FINDING-24 | `Save*` propagates load errors instead of silently replacing with empty slice |
| FINDING-25 | Atomic writes via temp file + `os.Rename`; file perms `0600` |
| FINDING-26 | `sync.Once` for thread-safe `getStorage()` lazy init |
| FINDING-28 | All 4 React components show fetch error messages instead of silent empty list |
| FINDING-29 | `ErrorBoundary` wraps `<AppRoutes />` with fallback UI |
| FINDING-35 | API handler tests use `t.TempDir()` for hermetic isolation |
| FINDING-09 | Partial: file perms `0600`, dir `0700` (encryption at rest not yet implemented) |

---

## Change Log

| Date | Action |
|---|---|
| 2026-03-10 | Initial security review (23 findings) |
| 2026-03-27 | 7-agent parallel review: added 11 new findings (FINDING-24 through FINDING-35), upgraded FINDING-08 from Medium to High, added agent consensus counts, reorganized sections |
| 2026-03-27 | Fixed FINDING-01 (API key auth), FINDING-24 (no silent data destruction), FINDING-25 (atomic writes + 0600 perms), FINDING-13 (dir perms 0700), FINDING-35 (hermetic tests). Partial fix for FINDING-09 (perms only, no encryption). |
| 2026-03-27 | Fixed FINDING-12 (mutex), FINDING-03 (CORS origin), FINDING-08 (body limit), FINDING-26 (sync.Once), FINDING-28 (fetch error display), FINDING-29 (Error Boundary), FINDING-17 (CRA->Vite migration). Added 41 web tests (vitest). |
| 2026-04-06 | v0.4 docs sweep: updated executive summary to reflect all 12+1 fixes (was stale at 5+1); updated FINDING-03 fix description with v0.4 CORS rewrite details (external corsHandler, localhost wildcard); corrected .js -> .jsx file references in FINDING-19, FINDING-20, FINDING-34 (Vite migration changed extensions). |
