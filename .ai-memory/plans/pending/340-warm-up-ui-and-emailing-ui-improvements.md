# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 340 Enterprise Spec)

- **Task Sequence:** 340
- **Task Slug:** `340-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 340 of the GitMap Warm-Up and Emailing UI Engine. Sequence 340 advances email warming into an autonomous multi-agent swarm ecosystem featuring cognitive sending rhythms, peer-to-peer fleet mesh networks, multi-theme accessible design systems, high-speed Split-DB SQLite sharding, and closed-loop self-healing deliverability:
- **Multi-Agent Deliverability Swarm:** 10 specialized micro-agents continuously collaborating to balance IP reputation, seed engagement, content entropy, and cadence optimization.
- **Cognitive Sending Rhythms & Queue Saturation:** Mathematical modeling of human circadian workflows and ISP reception queues with fatigue visualization.
- **P2P Fleet Mesh & Cross-Region Seeds:** Zero-trust peer gossip, cross-cloud worker dispatch, and low-latency seed inbox synchronization.
- **Design System Tokens & Accessible Motion:** 10-step perceptual color ramps, 5 production glassmorphic themes, WCAG 2.1 AAA high-contrast modes, and fluid typography.
- **Split-DB SQLite Sharding & Query Caches:** High-throughput WAL pragmas, monthly audit table sharding, LRU memory caches, Bloom filter suppressions, and SHA-256 Merkle tree verification.
- **Antigravity IDE Event Bus & Scorecards:** Bidirectional WebSocket broker, slash-command CLI triggers, audio deliverability synth, and multi-factor 0–100 health scorecards.

## Subtasks

- [ ] `01-swarm-rhythms-and-fleet-mesh.md` — Author specifications and implement Multi-Agent Collaborative Deliverability Swarm, Cognitive Sending Rhythms, and Multi-Node P2P Fleet Mesh (Improvements 301–330, Clusters 31–33).
- [ ] `02-self-healing-circuits-and-tui.md` — Author specifications and implement Autonomous Content Entropy Synthesizers, Real-Time ISP Circuit Breakers, and Interactive Terminal Operator Studio (Improvements 331–360, Clusters 34–36).
- [ ] `03-design-tokens-and-splitdb.md` — Author specifications and implement Design System Multi-Theme Palettes, Accessible Motion Tokens, and High-Throughput Split-DB SQLite Sharding (Improvements 361–380, Clusters 37–38).
- [ ] `04-agy-event-bus-and-audit.md` — Author specifications and implement Antigravity IDE Event Bus, Autonomous Deliverability Scorecards, and Tamper-Evident Merkle Ledgers (Improvements 381–400, Clusters 39–40).

## Architectural Clusters Breakdown

1. **Swarm Rhythms & Fleet Mesh (Clusters 31–33 / Improvements 301–330):**
   - Multi-agent collaboration protocols and consensus engine in `pkg/agent/swarm/consensus_engine.go`.
   - Circadian rhythm and Poisson send burst modeler in `pkg/mail/rhythm/circadian_modeler.go`.
   - P2P gossip peer router and seed inbox coordinator in `pkg/mail/mesh/peer_router.go`.

2. **Self-Healing Circuits & TUI (Clusters 34–36 / Improvements 331–360):**
   - Autonomous content entropy generator and spintax builder in `pkg/mail/entropy/synthesizer.go`.
   - Three-state ISP circuit breaker and throttle listener in `pkg/mail/circuit/breaker_fsm.go`.
   - Lipgloss interactive operator studio console in `pkg/cli/tui/operator_studio.go`.

3. **Design Tokens & Split-DB Sharding (Clusters 37–38 / Improvements 361–380):**
   - Mathematical color ramp and glassmorphic theme tokens in `pkg/web/tokens/color_ramp.go` and `web/static/css/tokens.css`.
   - Split-DB SQLite WAL pragmas, monthly sharding, and memory query cache in `pkg/db/splitdb/` and `pkg/db/cache/`.
   - In-memory Bloom filter suppression lookup in `pkg/mail/suppression/bloom_filter.go`.

4. **AGY Event Bus & Autonomous Audit (Clusters 39–40 / Improvements 381–400):**
   - Bidirectional WebSocket event broker and slash-command dispatcher in `pkg/agent/bus/` and `pkg/cli/cmdwarmup/`.
   - 10-dimension deliverability health scorecard and compliance radar in `pkg/mail/audit/scorecard_engine.go`.
   - SHA-256 Merkle tree verification ledger and self-healing coordinator in `pkg/db/audit/merkle_tree.go` and `pkg/mail/coordinator/`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Zero exits in leaf packages. All errors returned directly via typed `*apperror.AppError`.
- **Function Limit:** Functions <= 15 lines max (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Positive booleans with `is` or `has` prefixes exclusively (`isHealthy`, `hasPassedAudit`).
- **Zero CDN Dependencies:** 100% self-contained embedded web assets via Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all documentation and configurations.
