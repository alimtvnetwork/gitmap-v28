# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 342 Enterprise Spec)

- **Task Sequence:** 342
- **Task Slug:** `342-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 342 of the GitMap Warm-Up and Emailing UI Engine (Improvements 501–600). Sequence 342 expands warm-up infrastructure to connection pool velocity governors, interactive Bubble Tea engagement inspectors, real-time MX DNS sentinels, forensic bounce classifiers, multi-node SQLite sharded replication clusters, and autonomous Antigravity multi-agent fleet telemetry bridges:
- **Connection Pools & Velocity Regulation (Improvements 501–530, Clusters 51–53):** SMTP socket multiplexing, persistent TLS connection pooling, adaptive leaky bucket rate limiters, ISP throttle feedback controllers, and connection starvation sentinels.
- **Engagement Feedback & Bubbletea Inspector (Improvements 531–560, Clusters 54–56):** Interactive full-screen Bubbletea terminal inspector, seed response feedback tracking, sentiment extraction on incoming warmup replies, and render latency profiling HUDs.
- **Real-Time MX DNS Health & Forensic Bounce Analysis (Improvements 561–580, Clusters 57–58):** Deep MX host latency monitoring, reverse DNS PTR alignment, 50+ DNSBL blacklist queries, SPF/DKIM/DMARC/STS verification, RFC 3463 DSN parser, Bayesian semantic bounce categorization, and automatic circuit-breaking suppression list managers.
- **Multi-Node SQLite Sharded Clusters & AGY Fleet Choreographer (Improvements 581–600, Clusters 59–60):** Sharded SQLite database governance across physical cluster nodes, streaming WAL frame replication lag sentinels, epoch lease arbiters, incremental auto-vacuum governors, AGY peer discovery mesh, token cost governance, subtask lease arbiters, transcript AST miners, and master executive cockpit.

## Subtasks

- [ ] `01-connection-pools-and-velocity-regulator.md` — Author specifications and implement SMTP Connection Pools, TLS Multiplexers, and Velocity Rate Governors (Improvements 501–530, Clusters 51–53).
- [ ] `02-engagement-feedback-and-bubbletea-inspector.md` — Author specifications and implement Interactive Bubbletea Engagement Studio, Seed Interaction Telemetry, and Visual Inspector (Improvements 531–560, Clusters 54–56).
- [ ] `03-mx-sentinel-and-bounce-forensics.md` — Author specifications and implement Real-Time MX DNS Health Sentinel, Reverse IP Verification, Forensic DSN Parser, and Bayesian Bounce Categorizer (Improvements 561–580, Clusters 57–58).
- [ ] `04-sqlite-shards-and-antigravity-fleet.md` — Author specifications and implement Multi-Node SQLite Sharded Cluster Dashboard, WAL Replication Sentinel, Antigravity Agent Fleet Choreographer, Token Governance, and Master Swarm Cockpit (Improvements 581–600, Clusters 59–60).

## Architectural Clusters Breakdown

1. **Connection Pools & Velocity Regulation (Clusters 51–53 / Improvements 501–530):**
   - High-performance SMTP connection pooling and persistent keep-alive multiplexers in `pkg/mail/pool/`.
   - Adaptive leaky bucket rate limiter and throttle regulator in `pkg/mail/velocity/`.
   - Socket latency profiling and TLS handshake monitors in `pkg/mail/telemetry/`.

2. **Engagement Feedback & Bubbletea Inspector (Clusters 54–56 / Improvements 531–560):**
   - Full-screen interactive terminal inspector and queue dashboard in `pkg/cli/tui/bubbletea_inspector.go`.
   - Real-time seed inbox response tracker and sentiment analysis in `pkg/mail/engagement/`.
   - Microsecond template render profiling and visual canvas debuggers in `pkg/web/canvas/`.

3. **MX Sentinel & Bounce Forensics (Clusters 57–58 / Improvements 561–580):**
   - Real-time MX DNS latency sentinel, PTR alignment, DNSBL queries, and SPF/DKIM/DMARC/STS/BIMI verifiers in `pkg/mail/mx/`.
   - RFC 3463 / 3464 DSN bounce parser, Bayesian semantic categorizer, ISP error translator, and automated suppression list manager in `pkg/mail/bounce/`.
   - Feedback Loop (FBL) ARF analyzer and live forensic bounce diagnostic HUD in `pkg/mail/bounce/`.

4. **SQLite Shards & Antigravity Fleet (Clusters 59–60 / Improvements 581–600):**
   - Multi-node SQLite shard discovery, WAL streaming lag sentinels, epoch lease coordinators, and zero-downtime backup orchestrators in `pkg/db/cluster/`.
   - Antigravity agent peer discovery, subtask lease arbitration, token cost governance, inter-agent event bus, and deadlock remediation in `pkg/agent/swarm/`.
   - Conversation transcript full-text miner and master swarm orchestrator cockpit in `pkg/agent/swarm/`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Zero exits in leaf packages. All errors returned directly via typed `*apperror.AppError`.
- **Function Limit:** Functions <= 15 lines max (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Positive booleans with `is` or `has` prefixes exclusively (`isReachable`, `isSynchronized`, `hasValidMxRecords`).
- **Zero CDN Dependencies:** 100% self-contained embedded web assets via Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all documentation and configurations.
