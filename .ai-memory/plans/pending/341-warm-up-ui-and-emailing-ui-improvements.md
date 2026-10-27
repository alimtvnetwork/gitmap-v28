# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 341 Enterprise Spec)

- **Task Sequence:** 341
- **Task Slug:** `341-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 341 of the GitMap Warm-Up and Emailing UI Engine. Sequence 341 advances email warming from foundational peer swarms into high-precision predictive deliverability, interactive AST canvas editing, enterprise anti-spam sentinel guards, Split-DB distributed replication, and live Antigravity multi-agent swarm steering:
- **Predictive Analytics & Bayesian Warming Curves:** Dynamic regression forecasting of ISP threshold cliffs, seed engagement decay, and non-linear volume ramps.
- **Interactive TUI Operator Studio & AST Canvas:** Full-featured Lipgloss terminal dashboard and DOM-like AST template visualizer with side-by-side semantic diffing.
- **Autonomous Anti-Spam Sentinel & Honeypot Avoidance:** Bayesian spam-trap probability scoring, pristine honeypot lexical scanners, DNSBL multi-provider queries, and automatic quarantine orchestration.
- **Polyglot Template AST Diffing & Bandit Scheduling:** Multi-engine template AST parsing, Shannon text entropy verification, i18n placeholder coherence, and Thompson-sampling multi-variant bandit schedules.
- **Split-DB Cross-Replica Sync & SQLite Maintenance:** Distributed WAL frame streaming with vector clocks, non-blocking incremental auto-vacuums, B-tree fragmentation auditing, and monthly reputation shard rotation.
- **Antigravity Multi-Agent Swarm Orchestrator & Live Steerer:** Real-time subagent supervisor, interactive operator task steering HUD, token budget allocation, bidirectional AGY event bus, and unified master cockpit.

## Subtasks

- [ ] `01-predictive-analytics-and-bayesian-curves.md` — Author specifications and implement Predictive Analytics, Bayesian Volume Curves, and Multi-ISP Threshold Cliff Detectors (Improvements 401–430, Clusters 41–43).
- [ ] `02-interactive-tui-and-canvas-inspector.md` — Author specifications and implement Interactive TUI Operator Studio, Live Visual Canvas Inspector, and Render Pipeline Diagnostics (Improvements 431–460, Clusters 44–46).
- [ ] `03-anti-spam-sentinel-and-splitdb-sync.md` — Author specifications and implement Autonomous Anti-Spam Trap Detection, Honeypot Avoidance Telemetry, Polyglot Campaign AST Diffing, and Split-DB Cross-Replica WAL Sync (Improvements 461–490, Clusters 47–49).
- [ ] `04-antigravity-swarm-orchestrator.md` — Author specifications and implement Antigravity Multi-Agent Swarm Supervisor, Live Task Steerer HUD, Token Cost Governance, and Master Cockpit (Improvements 491–500, Cluster 50).

## Architectural Clusters Breakdown

1. **Predictive Analytics & Bayesian Curves (Clusters 41–43 / Improvements 401–430):**
   - Bayesian threshold curve calculator and ISP reputation trendline modeler in `pkg/mail/predictive/bayesian_curve.go`.
   - Seed engagement decay modeling and ISP inbox placement forecasting in `pkg/mail/predictive/decay_modeler.go`.
   - Dynamic non-linear daily ramp schedule generator in `pkg/mail/predictive/ramp_scheduler.go`.

2. **Interactive TUI & Canvas Inspector (Clusters 44–46 / Improvements 431–460):**
   - Full-screen Lipgloss terminal studio and live queue monitor in `pkg/cli/tui/canvas_studio.go`.
   - Interactive DOM-like visual template inspector and layout debugger in `pkg/web/canvas/inspector.go`.
   - Microsecond rendering pipeline profiler and streaming buffer telemetry in `pkg/mail/telemetry/render_profiler.go`.

3. **Anti-Spam Sentinel & Split-DB Sync (Clusters 47–49 / Improvements 461–490):**
   - Bayesian spam-trap detection, lexical scanners, and DNSBL blocklist sentinel in `pkg/mail/trap/` and `pkg/mail/dnsbl/`.
   - Polyglot template AST parsing, structural diffing, and spintax entropy verification in `pkg/mail/template/`.
   - Distributed SQLite WAL frame replication, incremental auto-vacuum scheduler, and B-tree fragmentation auditing in `pkg/db/splitdb/`.

4. **Antigravity Swarm Orchestrator & Cockpit (Cluster 50 / Improvements 491–500):**
   - Real-time multi-agent supervisor, live task steerer, and token budget allocator in `pkg/agent/swarm/`.
   - Bidirectional AGY event bus bridge and conversational `/warmup` slash command handler in `pkg/agent/swarm/`.
   - Unified master cockpit and acoustic feedback engine in `pkg/agent/swarm/cockpit_presenter.go` and `pkg/agent/swarm/audio_engine.go`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Zero exits in leaf packages. All errors returned directly via typed `*apperror.AppError`.
- **Function Limit:** Functions <= 15 lines max (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Positive booleans with `is` or `has` prefixes exclusively (`isTrapSuspected`, `hasValidMxRecords`).
- **Zero CDN Dependencies:** 100% self-contained embedded web assets via Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all documentation and configurations.
