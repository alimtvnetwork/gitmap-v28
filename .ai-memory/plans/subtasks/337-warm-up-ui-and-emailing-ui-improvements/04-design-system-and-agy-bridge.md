# Subtask 04: Design System Tokens, SQLite Split-DB Telemetry & Antigravity Bridge

## Status: PLANNED
- **Subtask Slug:** `04-design-system-and-agy-bridge`
- **Parent Plan:** `.ai-memory/plans/pending/337-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/04-improvements-61-to-100-security-tokens-and-agy-integration.md`
- **Owner:** Platform & Integrations Subagent
- **Target Subsystem:** `pkg/web/tokens/`, `pkg/store/`, `pkg/ide/`, `pkg/bridge/`, `pkg/daemon/`
- **Improvements Covered:** 71 to 100 (Clusters 8, 9, and 10)
- **Exclusive Owned Files:**
  - `pkg/web/tokens/theme_registry.go`
  - `pkg/web/tokens/sparkline_generator.go`
  - `pkg/store/telemetry_driver.go`
  - `pkg/store/wal_monitor.go`
  - `pkg/store/audit_chain.go`
  - `pkg/ide/agy_bridge.go`
  - `pkg/bridge/task_notifier.go`
  - `pkg/bridge/slash_handler.go`
  - `pkg/daemon/lemon_service.go`
  - `pkg/web/handlers/executive_dashboard.go`

---

## 1. Executive Summary & Context

Subtask 04 delivers the underlying design tokens, high-performance database infrastructure, and Google Antigravity (AGY) developer ecosystem integration:
1. **Design System Tokens & Themes (Improvements 71–80):** TokyoNight and Arctic light CSS variables, WCAG AA/AAA automated contrast compliance, JetBrains Mono font scaling, density toggles, and inline SVG sparklines.
2. **SQLite Split-DB Telemetry & Audit Logs (Improvements 81–90):** Microsecond query profiler, WAL checkpoint health dials, Casbin RBAC audit logging, and tamper-evident SHA-256 hash chains.
3. **Antigravity IDE Integration & Automation Daemon (Improvements 91–100):** AGY status bar telemetry bridge, task completion email triggers, slash commands (`/mail status`), subagent swarm alerts, and cross-platform headless daemon (`gitmap lemon warmup`).

---

## 2. File Modification & Structural Ownership

```text
pkg/web/
├── tokens/
│   ├── theme_registry.go       [NEW] - Zero-runtime CSS custom properties registry
│   └── sparkline_generator.go  [NEW] - High-speed SVG inline sparkline generator
└── handlers/
    └── executive_dashboard.go  [NEW] - Panoramic executive KPI command center endpoint
pkg/store/
├── telemetry_driver.go         [NEW] - Database wrapper recording microsecond query latencies
├── wal_monitor.go              [NEW] - SQLite WAL file checkpoint and fragmentation dial
└── audit_chain.go              [NEW] - Tamper-evident SHA-256 chained audit ledger
pkg/ide/
└── agy_bridge.go               [NEW] - Local IPC socket provider for Antigravity IDE
pkg/bridge/
├── task_notifier.go            [NEW] - Long-running background task completion emailer
└── slash_handler.go            [NEW] - Antigravity slash-command dispatcher
pkg/daemon/
└── lemon_service.go            [NEW] - Cross-platform headless automation service
```

---

## 3. Technical Requirements & Architecture

### 3.1 Design System Tokens & WCAG AA Contrast
- TokyoNight dark theme: `--tn-bg-dark: #1a1b26`, `--tn-bg-card: #24283b`, `--tn-cyan: #7dcfff`.
- Arctic light theme: `--al-bg-canvas: #f8f9fc`, `--al-bg-card: #ffffff`, `--al-text-dark: #2e3440`.
- All text/background combinations programmatically verified for minimum 4.5:1 contrast.
- Zero-CDN architecture: embeds subset WOFF2 binaries and inline SVG vectors.

### 3.2 Split-DB Telemetry & Tamper-Evident Chaining
- Tracks query latency in memory ring buffer; flags queries > 5ms with `EXPLAIN QUERY PLAN`.
- Monitors SQLite WAL checkpoint states (`PASSIVE`, `FULL`, `RESTART`, `TRUNCATE`).
- Outbound audit log entries cryptographically chained:
  `EntryHash = SHA256(PreviousEntryHash || Timestamp || MailboxID || Recipient || MessageID)`.

### 3.3 Antigravity Bridge & Headless Lemon Daemon
- IPC bridge serving status JSON to IDE status bar over local pipe.
- Slash command router handling `/mail status` and `/mail ramp` requests directly from agent chat.
- Headless OS service daemon (`gitmap lemon warmup`) handling scheduled deliveries with automatic fault recovery.

---

## 4. Verification & Testing Checklist

- [ ] Unit tests in `pkg/web/tokens/theme_registry_test.go` verifying WCAG AA contrast math.
- [ ] Split-DB tests validating SHA-256 audit chaining and tamper detection.
- [ ] Database telemetry tests verifying statement cache hit/miss tracking.
- [ ] AGY IPC bridge tests verifying JSON payload serialization and socket reconnection.
- [ ] Coding guidelines audit: <= 15 line functions, positive booleans (`hasPassedWcagAa`, `isTamperEvident`), error return sovereignty.
