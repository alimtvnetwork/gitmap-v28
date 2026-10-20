# 05-embedded-web-uis: Zero-CDN Embedded Web Dashboards Specification

- **Spec ID:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md`
- **Subsystem:** Mailing Composer UI & Schedule/Notifier Settings UI
- **Status:** APPROVED
- **Related Specs:**
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`
- **Implementation Subtask:** `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/05-subtask-embedded-mailing-and-schedule-uis.md`

---

## 1. Executive Summary & UI Architecture

GitMap provides dual embedded, zero-CDN web dashboards to manage email drafting, mailbox monitoring, and event-driven schedule notification rules directly from a local browser.

### 1.1 Dual Embedded Dashboards
1. **Mailing Composer UI (`gitmap mail ui`):**
   - Served at `http://127.0.0.1:8745/mail/ui`.
   - Rich and plain-text composition, attachment selection, profile selection, template library, and delivery preview.
2. **Schedule & Notifier Settings UI (`gitmap schedule ui`):**
   - Served at `http://127.0.0.1:8746/schedule/ui`.
   - Task manager for scheduled mail checks, interval email dispatches, folder event watcher rule configuration, and notification audit logs.

### 1.2 Zero-CDN Invariant
- Both UIs are 100% self-contained single-page applications (SPAs).
- Zero external CSS or JavaScript dependencies (no CDNs, Google Fonts, or remote asset requests).
- All styling, fonts, icons (SVG), and client scripts are embedded as Go string constants or asset slices.
- Full offline operability in air-gapped environments.

---

## 2. Design System & Dark Theme Tokens

Both dashboards adhere to the GitMap dark theme palette:

| Token Name | Hex Value | Purpose |
|---|---|---|
| `--bg-base` | `#090d16` | Main window background canvas |
| `--bg-surface` | `#0f172a` | Card, container, and panel background |
| `--bg-surface-elevated` | `#1e293b` | Elevated panels, modal overlays, dropdown menus |
| `--border-subtle` | `#1e293b` | Primary card and input borders |
| `--border-accent` | `#38bdf8` | Focus rings, active tabs, highlights |
| `--text-primary` | `#f8fafc` | High-contrast headings and body text |
| `--text-muted` | `#94a3b8` | Subheadings, labels, placeholders, timestamps |
| `--accent-cyan` | `#38bdf8` | Primary CTA buttons, key metrics, active badges |
| `--accent-sky` | `#0284c7` | Secondary accents and link hovers |
| `--status-success` | `#10b981` | Successful deliveries and running watchers |
| `--status-warning` | `#f59e0b` | Debouncing queues and mailbox warnings |
| `--status-danger` | `#ef4444` | Delivery failures and error banners |

---

## 3. Modular File Architecture (< 300 Lines Each)

To adhere strictly to file size guidelines, UI code is modularized into discrete files under 300 lines:

### 3.1 Mailing Composer Module (`cli/cmdmail/`)
- `ui_server.go`: HTTP listener lifecycle, port negotiation, route registration, graceful shutdown.
- `ui_assets.go`: Embedded CSS stylesheet strings, SVG icons, and MIME-type dispatchers.
- `ui_html.go`: HTML skeleton, navigation header, composer grid, and template modal markup.
- `ui_js.go`: Vanilla ES2020 client script for form handling, AJAX calls, and masked preview.

### 3.2 Schedule Settings Module (`cli/cmdschedule/`)
- `ui_server.go`: HTTP listener lifecycle, port negotiation, route registration, graceful shutdown.
- `ui_assets.go`: Embedded stylesheet strings, badge renderers, and SVG icons.
- `ui_html.go`: HTML layout, watcher rules table, cron schedule drawer, and log viewer markup.
- `ui_js.go`: Client script for watcher rule CRUD, event toggles, debounce sliders, and polling.

---

## 4. REST API Endpoints

The embedded web servers host local JSON REST endpoints bound strictly to `127.0.0.1`:

### 4.1 Mailing REST API (`/api/mail/*`)
- `GET /api/mail/accounts`: Returns list of configured mailbox profiles with masked emails.
- `POST /api/mail/send`: Accepts draft payload (`to`, `cc`, `subject`, `body`, `isHTML`, `attachments`).
- `GET /api/mail/templates`: Returns saved templates for one-click insertion.
- `POST /api/mail/templates`: Creates or updates a reusable email template.
- `POST /api/mail/test-connection`: Verifies SMTP connectivity for selected mailbox profile.

### 4.2 Schedule & Watcher REST API (`/api/schedule/*`)
- `GET /api/schedule/jobs`: Lists all active scheduled email tasks and mailbox check intervals.
- `POST /api/schedule/jobs`: Creates or updates a scheduled email task.
- `DELETE /api/schedule/jobs?id={id}`: Cancels or removes a scheduled task.
- `GET /api/schedule/watchers`: Lists configured folder watcher rules with debounce parameters.
- `POST /api/schedule/watchers`: Creates or updates a folder watcher rule.
- `DELETE /api/schedule/watchers?id={id}`: Removes a folder watcher rule.
- `GET /api/schedule/logs`: Retrieves recent notification dispatches and execution status.

---

## 5. UI Views & Component Breakdown

### 5.1 Mailing Composer (`/mail/ui`)
1. **Header Bar:** Brand title, current sender profile badge, active port indicator, connection status.
2. **Account Switcher:** Dropdown selecting between configured SMTP profiles (e.g. `work [u***r@c***y.com]`).
3. **Recipient Input Area:**
   - Multi-chip input field for `To`, `Cc`, `Bcc`.
   - Real-time privacy toggle: preview with masked addresses or unmasked addresses.
4. **Editor Pane:**
   - Subject line input with template variable auto-complete.
   - Body mode switcher (`Plain Text` / `HTML Preview`).
   - Drag-and-drop attachment container with size indicator.
5. **Action Bar:**
   - `Send Email` (Cyan CTA button with loading spinner).
   - `Save Draft` (Subtle button).
   - `Discard` (Outline button).

### 5.2 Schedule & Notifier Settings (`/schedule/ui`)
1. **Scheduler Overview Cards:** Active watcher count, pending scheduled emails, 24h notification tally.
2. **Folder Watcher Rule Designer:**
   - Folder path text input with validation.
   - Event trigger checkboxes: `File Created` (`CREATE`), `File Modified` (`WRITE`), `File Deleted` (`REMOVE`).
   - Debounce window range slider (`5s` to `60s`, default `8s`).
   - Notifier recipient chip list (auto-masked by default).
   - Subject & Body notification template editors.
3. **Watcher Rules Table:** Displays active rules, monitored path, triggers, last triggered timestamp, status switch.
4. **Execution Audit Drawer:** Paginated table displaying trigger history, debounced event counts, and dispatch results.

---

## 6. Security, Localhost Binding & Graceful Termination

1. **Localhost Restriction:** Listeners bind strictly to `127.0.0.1` or `::1`. Public network binding (`0.0.0.0`) is prohibited.
2. **CORS & Origin Validation:** Rejects any cross-origin requests whose `Origin` or `Referer` headers do not match the local listener address.
3. **Browser Auto-Launch:** CLI provides `--no-open` flag to suppress automatic browser launching.
4. **Graceful Shutdown:** Intercepts `SIGINT` / `SIGTERM` and shuts down HTTP server within a 3-second timeout window.
