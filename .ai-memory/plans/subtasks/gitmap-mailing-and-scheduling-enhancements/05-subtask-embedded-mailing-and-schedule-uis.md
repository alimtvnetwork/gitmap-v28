# Subtask 05: Embedded Mailing Composer UI & Schedule Settings UI

- **Subtask ID:** `subtask-05-embedded-mailing-and-schedule-uis`
- **Parent Task:** `gitmap-mailing-and-scheduling-enhancements`
- **Assigned Worker:** Worker 03
- **Status:** PENDING

---

## 1. Owned Files (Exclusive)

The following files are strictly owned by Worker 03:
- `cli/cmdmail/ui_server.go`
- `cli/cmdmail/ui_assets.go`
- `cli/cmdmail/ui_html.go`
- `cli/cmdmail/ui_js.go`
- `cli/cmdschedule/ui_server.go`
- `cli/cmdschedule/ui_assets.go`
- `cli/cmdschedule/ui_html.go`
- `cli/cmdschedule/ui_js.go`

---

## 2. Read-Only Context Files

- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
- `cli/constants/constants.go`
- `cli/termout/termout.go`

---

## 3. Subtask Objectives & Step-by-Step Execution Plan

### Step 1: Modular Zero-CDN Server Skeletons
1. Implement `cli/cmdmail/ui_server.go`:
   - Bind HTTP server strictly to `127.0.0.1:8745` (with dynamic increment if port is occupied).
   - Register route handlers for `/mail/ui` and `/api/mail/*`.
   - Provide graceful shutdown with `context.WithTimeout(context.Background(), 3*time.Second)`.
   - Implement `--no-open` flag check and local browser launch helper.
2. Implement `cli/cmdschedule/ui_server.go`:
   - Bind HTTP server strictly to `127.0.0.1:8746` (with dynamic increment if port is occupied).
   - Register route handlers for `/schedule/ui` and `/api/schedule/*`.
   - Provide graceful shutdown and browser launch integration.

### Step 2: Dark Theme Assets & Templates (`ui_assets.go`, `ui_html.go`)
1. Implement `cli/cmdmail/ui_assets.go` and `cli/cmdschedule/ui_assets.go`:
   - Store dark theme CSS string constants using palette tokens:
     - Canvas: `#090d16`, Cards: `#0f172a`, Borders: `#1e293b`.
     - Accents: Cyan `#38bdf8`, Sky `#0284c7`.
     - Status: Green `#10b981`, Amber `#f59e0b`, Red `#ef4444`.
   - Store inline SVG icons for actions (mail, send, paperclip, clock, folder, bell, check).
   - Guarantee zero remote URL references (no Google Fonts, FontAwesome, or CDN scripts).
2. Implement `cli/cmdmail/ui_html.go` and `cli/cmdschedule/ui_html.go`:
   - Render responsive single-page application structure.
   - Mailing Composer layout: account switcher, recipient pills, subject, body mode tabs, attachment container.
   - Schedule Settings layout: watcher metrics, folder watcher form, rules table, execution audit log drawer.

### Step 3: Client-Side Interaction & REST Integration (`ui_js.go`)
1. Implement `cli/cmdmail/ui_js.go`:
   - Fetch `/api/mail/accounts` to populate sender profiles (with masked email preview).
   - Form submission handler for `/api/mail/send` with draft preview and field validation.
   - Template inserter modal logic.
2. Implement `cli/cmdschedule/ui_js.go`:
   - Fetch active rules from `/api/schedule/watchers` and populate interactive rule cards.
   - Folder rule creation handler: path validation, event checkboxes, debounce slider (5s - 60s), recipient chip editor.
   - Delete/toggle rule state with live optimistic UI update.
   - Fetch `/api/schedule/logs` for real-time audit event updates.

### Step 4: Strict File Size Governance (< 300 Lines Each)
- Verify every individual Go file is strictly under 300 lines.
- Extract any helper structs or formatters if file length approaches 250 lines.

---

## 4. Verification & Quality Gates

1. **Zero-CDN Check:**
   - Scan all generated asset strings for `http://` or `https://` (excluding localhost callbacks).
2. **Localhost Guard:**
   - Verify server refuses binding on non-loopback interfaces.
3. **Graceful Exit:**
   - Send `Interrupt` signal to running UI process; assert exit code 0 within 3 seconds.
4. **UI Response Status:**
   - `GET http://127.0.0.1:8745/mail/ui` returns HTTP 200 with HTML content-type.
   - `GET http://127.0.0.1:8746/schedule/ui` returns HTTP 200 with HTML content-type.
