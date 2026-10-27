# Plan: Warm-Up UI and Emailing UI Improvements

- **Task Sequence:** 337
- **Task Slug:** `337-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This plan orchestrates the complete design, architecture, and modular implementation of 100 concrete improvements for the GitMap Warm-Up and Emailing UI Engine. The system delivers an enterprise-grade outbound deliverability control plane inspired by modern autonomous manager dashboards (such as Google Antigravity Manager). It features progressive warmup ramp curves, automated seed list interactions, ISP placement health dials, multi-mailbox composer studios, terminal CLI counterparts (`gitmap warmup` and `gitmap mail`), DNS authentication diagnostics (SPF/DKIM/DMARC/BIMI), TokyoNight design tokens, microsecond SQLite Split-DB telemetry, and native Antigravity IDE automation hooks.

## Subtasks

- [ ] `01-warmup-telemetry-and-ramp-curves.md` — Author specifications and implement Warm-Up Ramp Curves, Schedule Timelines, Deliverability Health Gauges, and Seed List Simulation Physics (Improvements 01–30, Clusters 1–3).
- [ ] `02-composer-studio-and-multimailbox-hud.md` — Author specifications and implement Multi-Mailbox Management, Outbound Composer Studio, and Inbox/Bounce Telemetry Monitors (Improvements 31–50, Clusters 4–5).
- [ ] `03-terminal-cli-and-security-diagnostics.md` — Author specifications and implement Terminal UI CLI Tools, SPF/DKIM/DMARC Authentication Diagnostics, and Spam Score HUD (Improvements 51–70, Clusters 6–7).
- [ ] `04-design-system-and-agy-bridge.md` — Author specifications and implement TokyoNight/Arctic Design System Tokens, SQLite Split-DB Telemetry, and Antigravity IDE Automation Bridge (Improvements 71–100, Clusters 8–10).

## Planned Deliverables & Architectural Clusters

1. **Warmup Telemetry & Ramp Engine (Clusters 1–3 / Improvements 01–30):**
   - Mathematical ramp curve calculators (Linear, Exponential, S-Curve, Fibonacci, Wavefront) in `pkg/mail/warmup/ramp_curves.go`.
   - Tri-color SVG radial dial deliverability gauges and multi-ISP placement breakdowns (Google Workspace, M365, Yahoo, iCloud) in `pkg/web/widgets/deliverability_gauge.go`.
   - Automated seed interaction physics, read velocity simulation, and reply thread weaving in `pkg/mail/seed/reply_simulator.go`.

2. **Composer Studio & Multi-Mailbox Control (Clusters 4–5 / Improvements 31–50):**
   - Multi-mailbox health card deck, rotation round-robin schedulers, and quota gauges in `pkg/mail/mailbox/pool_manager.go`.
   - Zero-dependency markdown/rich-text email composer studio with variable interpolation and real-time preview in `pkg/web/composer/studio.go`.
   - Bounce categorization engine (Hard, Soft, Policy, Challenge) and real-time telemetry streaming in `pkg/mail/telemetry/bounce_analyzer.go`.

3. **Terminal CLI & Security Diagnostics (Clusters 6–7 / Improvements 51–70):**
   - High-speed Lipgloss terminal commands: `gitmap warmup status`, `gitmap warmup start`, `gitmap mail send`, and `gitmap mail doctor`.
   - Comprehensive DNS authentication diagnostics: recursive SPF depth matrix, multi-selector DKIM alignment, DMARC RUA parser, and BIMI/VMC verification.
   - Offline SpamAssassin / Rspamd heuristic score predictor and blackhole trap detection scanner.

4. **Design System & Antigravity Bridge (Clusters 8–10 / Improvements 71–100):**
   - TokyoNight dark mode and Arctic light mode CSS tokens, WCAG AA contrast compliance, and JetBrains Mono typography hierarchy.
   - SQLite Split-DB microsecond query telemetry, WAL checkpoint health dials, Casbin RBAC audit logs, and tamper-evident SHA-256 hash chains.
   - Antigravity IDE workspace status bar indicators, task completion email notification triggers, slash-command bridges, and cross-platform headless automation daemon (`gitmap lemon warmup`).

## Quality Gates & Verification Standards

- **Error Management:** All leaf functions return typed errors (`*apperror.AppError`) directly. Zero `os.Exit` or `cliexit.HandleError` calls in leaf packages.
- **Boolean Discipline:** 100% positive booleans (`is`/`has` prefixes exclusively). No negative boolean flags.
- **Function Limit:** Functions strictly <= 15 lines max (<= 8 lines preferred). Flattened conditionals with zero nested if statements.
- **Relative Path Hygiene:** Strict repository-relative paths across all markdown specifications, plans, and Go imports.
- **Zero-CDN Invariant:** Completely self-contained web assets; all SVG icons, WOFF2 fonts, and CSS tokens served locally.
