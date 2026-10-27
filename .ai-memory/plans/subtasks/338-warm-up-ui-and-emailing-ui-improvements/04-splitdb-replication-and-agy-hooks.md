# Subtask 04: Split-DB Replication & AGY Telemetry Hooks

- **Parent Plan:** `.ai-memory/plans/pending/338-warm-up-ui-and-emailing-ui-improvements.md`
- **Spec Path:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/04-improvements-161-to-200-remediation-builder-splitdb.md`
- **Sequence:** 338
- **Subtask Index:** 04
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Objectives

Implement and verify Improvements 181 through 200 spanning Clusters 19 and 20:
1. **Cluster 19: Offline Split-DB Synchronization & Multi-Node Replication (Items 181–190):**
   - SQLite Write-Ahead Log (WAL) frame streamer and memory-mapped delta replicator.
   - Conflict-Free Replicated Data Types (CRDT) PN-Counters for distributed tally convergence.
   - P2P gossip mesh discovery and signed node heartbeat monitor.
   - Tamper-evident SHA-256 Merkle tree audit log replication in `mail_audit.db`.
   - Selective edge database partitioning, PII hashing, and offline buffer journal.
   - Distributed daily quota leases, automated zero-downtime vacuum HUD, schema migration synchronizer, and hot snapshot backup manager.

2. **Cluster 20: Antigravity IDE Live Telemetry, Sound Effects & Autonomous Agent Hooks (Items 191–200):**
   - Antigravity IDE status bar JSON-RPC bridge and metric popover menu.
   - Procedural Web Audio API sound synthesizer (`audio_engine.js`) generating event-driven harmonic feedback.
   - Real-time WebSocket telemetry event bus for live metrics streaming.
   - Autonomous agent execution webhooks and slash-command CLI bridge (`gitmap warmup /status`).
   - Cross-platform headless service daemon (`gitmap lemon warmup daemon`).
   - Seed conversation sentiment heatmap, panic recovery crash dump dispatcher, local virtual blackhole sink, and graduation certificate generator.

## Concrete Deliverables

- `pkg/db/replication/wal_streamer.go`: Memory-mapped WAL frame streamer and checksum validator.
- `pkg/db/crdt/pn_counter.go`: CRDT state-based PN-Counter implementation.
- `pkg/db/audit/merkle_tree.go`: SHA-256 Merkle tree builder and audit trail verifier.
- `pkg/ide/bridge/status_bar.go`: JSON-RPC status bar provider for Antigravity IDE.
- `web/static/js/audio_engine.js`: Pure browser-native Web Audio API synthesizer.
- `pkg/web/telemetry/ws_bus.go`: WebSocket hub and client connection supervisor.
- `pkg/agent/hooks/webhook_dispatcher.go`: Autonomous subagent event webhook engine.
- `pkg/daemon/service_manager.go`: Cross-platform OS service daemon wrapper.

## Quality Standards

- Leaf error sovereignty: Bubble `*apperror.AppError` directly; zero process exit calls in library packages.
- Strict positive boolean discipline (`isReplicating`, `isConverged`, `hasAudioEnabled`, `hasAgentHooks`).
- Functions strictly bounded between 8 and 15 lines.
- Zero-CDN requirement: completely self-contained scripts, tokens, and audio oscillators.
