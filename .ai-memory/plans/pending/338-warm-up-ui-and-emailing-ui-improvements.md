# Plan: Warm-Up UI and Emailing UI Improvements (V8 Enterprise Spec)

- **Task Sequence:** 338
- **Task Slug:** `338-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 338 of the GitMap Warm-Up and Emailing UI Engine. Sequence 338 scales the platform into a multi-node, resilient deliverability command deck featuring:
- Interactive network topology canvas with direct-manipulation Bezier delivery routing and AI spam score heuristics.
- High-concurrency mailbox mesh scheduling with adaptive round-robin rotation, real-time quota gauges, and terminal Lipgloss TUI monitoring.
- Autonomous RFC 3463 bounce diagnosis, dynamic ISP token-bucket rate shaping, automated mailbox quarantine circuit breakers, and cooling ramp recovery.
- Drag-and-drop email template studio with fluid typography scaling, client viewport emulation, and inline MJML compilation.
- Offline SQLite Split-DB synchronization (`warmup.db` and `mail_audit.db`), CRDT counter convergence, and WAL delta replication.
- Antigravity IDE workspace status bar telemetry, procedural Web Audio API audio chimes, and autonomous subagent execution webhooks.

## Subtasks

- [ ] `01-interactive-node-canvas-and-ai-copilot.md` — Author specifications and implement Interactive Node Canvas, Dynamic Delivery Topology, AI Deliverability Co-Pilot, and Multi-Persona Seed Simulation Physics (Improvements 101–130, Clusters 11–13).
- [ ] `02-mailbox-mesh-and-terminal-tui.md` — Author specifications and implement Multi-Mailbox Mesh Orchestration, Visual Rotation Schedules, and Terminal TUI Operator Consoles (Improvements 131–160, Clusters 14–16).
- [ ] `03-bounce-remediation-and-template-builder.md` — Author specifications and implement Autonomous Bounce Remediation, Smart Throttling, and Visual Email Template Studio with Live Component Builder (Improvements 161–180, Clusters 17–18).
- [ ] `04-splitdb-replication-and-agy-hooks.md` — Author specifications and implement Offline Split-DB Synchronization, Multi-Node Replication, Antigravity IDE Live Telemetry, Procedural Audio, and Autonomous Agent Hooks (Improvements 181–200, Clusters 19–20).

## Architectural Clusters Breakdown

1. **Topology Canvas & AI Co-Pilot (Clusters 11–13 / Improvements 101–130):**
   - Force-directed physics canvas for visual outbound SMTP routing in `pkg/web/canvas/topology_mesh.go`.
   - Local LLM prompt synthesizer and Shannon token entropy calculator in `pkg/mail/ai/entropy_scorer.go`.
   - Multi-persona seed simulation engine and realistic reply thread generator in `pkg/mail/seed/persona_simulator.go`.

2. **Mailbox Mesh & Terminal TUI (Clusters 14–16 / Improvements 131–160):**
   - Multi-mailbox pool balancer and dynamic quota reservation coordinator in `pkg/mail/pool/mesh_coordinator.go`.
   - Lipgloss-driven terminal dashboard and interactive CLI monitor (`gitmap warmup tui`) in `pkg/cli/tui/warmup_dashboard.go`.
   - DNS record health checks, automated PTR/rDNS verifiers, and TLS cipher suit testing in `pkg/mail/dns/rdns_checker.go`.

3. **Bounce Remediation & Template Studio (Clusters 17–18 / Improvements 161–180):**
   - RFC 3463 status parser, per-ISP token bucket shaper, and quarantine circuit breaker in `pkg/mail/bounce/status_parser.go` and `pkg/mail/circuit/quarantine_breaker.go`.
   - Zero-runtime CSS grid email canvas, responsive viewport emulator, and dark mode preview shader in `pkg/web/builder/canvas_engine.go` and `pkg/web/builder/dark_mode_shader.go`.
   - Dynamic Mustache variable validator, bulletproof VML button generator, and inline image optimizer in `pkg/mail/template/token_parser.go`.

4. **Split-DB Replication & AGY Hooks (Clusters 19–20 / Improvements 181–200):**
   - SQLite WAL frame streaming, CRDT PN-Counter state mergers, and Merkle tree audit logging in `pkg/db/replication/wal_streamer.go` and `pkg/db/audit/merkle_tree.go`.
   - Antigravity IDE status bar JSON-RPC bridge and procedural Web Audio API sound engine in `pkg/ide/bridge/status_bar.go` and `web/static/js/audio_engine.js`.
   - Autonomous agent execution webhooks, slash-command router, and headless daemon manager in `pkg/agent/hooks/webhook_dispatcher.go` and `pkg/daemon/service_manager.go`.

## Non-Negotiable Standards & Quality Gates

- **Error Return Sovereignty:** Every leaf function returning an error must bubble `*apperror.AppError` directly. Absolutely zero `os.Exit` or `cliexit.HandleError` invocations within internal library packages.
- **Positive Booleans Only:** Strict enforcement of `is*` and `has*` prefixes (`isThrottled`, `hasValidCanvas`, `isConverged`). Negative booleans are completely banned.
- **Function Limit:** Functions <= 15 lines max, <= 8 lines preferred. Flatten all nested conditionals with early returns and guard clauses.
- **Strict Relative Pathing:** 100% repository-relative paths across all markdown files and Go packages.
- **Zero External Dependencies / CDNs:** Embedded CSS custom properties, self-contained SVG vectors, and native DOM event listeners.
