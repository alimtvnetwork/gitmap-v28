# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 339 Enterprise Spec)

- **Task Sequence:** 339
- **Task Slug:** `339-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 339 of the GitMap Warm-Up and Emailing UI Engine. Sequence 339 delivers a robust, zero-CDN Web Component architecture, time-travel deliverability debugging, resilient automated rescue daemons, terminal studios, multi-IP egress failover, design tokens, Split-DB WAL optimization, and autonomous Antigravity agent integration:
- Zero-CDN Shadow DOM custom web elements (`<warmup-metric-gauge>`, `<warmup-envelope-card>`, `<warmup-time-scrubber>`) with encapsulated glassmorphic token theming.
- Time-travel deliverability debugger allowing operators to scrub backwards and forwards through historical delivery events, inspect snapshot MTA states, and replay deliverability incidents.
- Adaptive warming gradient engine modeling logistic and sigmoid acceleration curves with dynamic ISP feedback throttling.
- Automated inbox placement rescue daemons, reputation repair protocols, and terminal Lipgloss visualizers.
- Dynamic multi-IP egress pools, local interface binding, per-IP health probing, and 3-state ISP circuit breakers (`Closed`, `Open`, `Half-Open`).
- Cohesive design system tokens, mathematical typography scales, WCAG 2.1 AAA high-contrast modes, and procedural Web Audio API synthesizers.
- Split-DB SQLite storage engine tuning (`warmup.db` and `mail_audit.db`) with WAL pragmas, automated checkpointing, memory caches, and SHA-256 Merkle tree audit logging.
- Antigravity autonomous agent event bus, pub/sub WebSockets, anomaly triage alerts, slash-command CLI listeners, and headless background daemons.

## Subtasks

- [ ] `01-web-components-and-time-travel.md` — Author specifications and implement Zero-CDN Shadow DOM Web Components, Design Tokens, Time-Travel Deliverability Debugger, and Adaptive Warming Gradient Engine (Improvements 201–230, Clusters 21–23).
- [ ] `02-rescue-daemons-and-terminal-studio.md` — Author specifications and implement Automated Inbox Rescue Daemons, Reputation Healing Ramps, and Interactive Terminal TUI Operator Consoles (Improvements 231–260, Clusters 24–26).
- [ ] `03-failover-and-design-tokens.md` — Author specifications and implement Multi-IP Egress Failover, ISP Circuit Breakers, Design System CSS Custom Properties, and WCAG High-Contrast Modes (Improvements 261–280, Clusters 27–28).
- [ ] `04-splitdb-tuning-and-agy-event-bus.md` — Author specifications and implement Dual Split-DB WAL Tuning, In-Memory Caches, Merkle Tree Auditing, Antigravity Agent Event Bus, and Telemetry Alerts (Improvements 281–300, Clusters 29–30).

## Architectural Clusters Breakdown

1. **Web Components & Time-Travel (Clusters 21–23 / Improvements 201–230):**
   - Zero-CDN Shadow DOM custom elements and glassmorphic card primitives in `pkg/web/components/custom_elements.go` and `web/static/js/components/`.
   - Deterministic temporal state scrubber and delivery incident replayer in `pkg/mail/timetravel/state_scrubber.go`.
   - Adaptive mathematical warming curve modeler and ISP feedback shaper in `pkg/mail/warmup/curve_modeler.go`.

2. **Rescue Daemons & Terminal Studio (Clusters 24–26 / Improvements 231–260):**
   - Automated deliverability rescue daemon and spam trap evasion coordinator in `pkg/mail/rescue/daemon_supervisor.go`.
   - Interactive Lipgloss terminal studio and operator control desk (`gitmap warmup studio`) in `pkg/cli/tui/studio_console.go`.
   - Autonomous seed reply generation matrix and persona response synchronizer in `pkg/mail/seed/reply_matrix.go`.

3. **Failover & Design Tokens (Clusters 27–28 / Improvements 261–280):**
   - Multi-IP outbound interface binder, health radar, and 3-state ISP circuit breakers in `pkg/mail/egress/interface_binder.go` and `pkg/mail/circuit/isp_breaker.go`.
   - Real-time automatic IP failover router and ASN-aware gateway dispatcher in `pkg/mail/egress/failover_router.go`.
   - Semantic CSS custom properties dictionary, fluid typography clamp curves, WCAG 2.1 AAA high-contrast styling, and SVG sprite sheets in `web/static/css/tokens.css`.

4. **Split-DB Tuning & AGY Event Bus (Clusters 29–30 / Improvements 281–300):**
   - Dual SQLite Split-DB engine (`warmup.db` and `mail_audit.db`), WAL auto-checkpointing, prepared statement cache, and Bloom filter suppression index in `pkg/db/splitdb/`.
   - Cryptographic SHA-256 Merkle tree transaction auditor and automated vacuum log compactor in `pkg/db/audit/merkle_tree.go`.
   - Real-time WebSocket event broker, signed webhook dispatcher, anomaly triage alerts, procedural Web Audio chimes, and headless daemon supervisor in `pkg/agent/bus/` and `pkg/agent/anomaly/`.

## Non-Negotiable Standards & Quality Gates

- **Error Return Sovereignty:** Every leaf function returning an error must bubble `*apperror.AppError` directly. Absolutely zero `os.Exit` or `cliexit.HandleError` invocations within internal library packages.
- **Positive Booleans Only:** Strict enforcement of `is*` and `has*` prefixes (`isFailoverTriggered`, `hasActiveAlerts`, `isCircuitClosed`, `hasValidTokens`). Negative booleans are completely banned.
- **Function Limit:** Functions <= 15 lines max, <= 8 lines preferred. Flatten all nested conditionals with early returns and guard clauses.
- **Strict Relative Pathing:** 100% repository-relative paths across all markdown files and Go packages.
- **Zero External Dependencies / CDNs:** Embedded CSS custom properties, self-contained SVG vectors, and native DOM custom elements.
