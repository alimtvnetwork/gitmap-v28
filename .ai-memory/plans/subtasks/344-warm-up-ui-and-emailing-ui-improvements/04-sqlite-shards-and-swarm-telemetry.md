# Subtask 04: SQLite Multi-Cluster Replicator & Antigravity Swarm Telemetry Bridge

## Status: READY
- **Subtask Slug:** `04-sqlite-shards-and-swarm-telemetry`
- **Parent Plan:** `.ai-memory/plans/pending/344-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/04-improvements-761-to-800-audit-governance-and-swarm-telemetry.md`
- **Target Range:** Improvements 781–800 (Clusters 79–80)
- **Target Subsystem:** `pkg/db/cluster/`, `pkg/db/merkle/`, `pkg/agent/swarm/`, `pkg/agent/telemetry/`
- **Exclusive Owned Files:**
  - `pkg/db/cluster/shard_topology.go`
  - `pkg/db/cluster/wal_streamer.go`
  - `pkg/db/cluster/two_phase_commit.go`
  - `pkg/db/cluster/vacuum_scheduler.go`
  - `pkg/db/merkle/page_tree.go`
  - `pkg/db/merkle/range_diff.go`
  - `pkg/agent/swarm/peer_mesh.go`
  - `pkg/agent/swarm/lease_arbiter.go`
  - `pkg/agent/swarm/token_guard.go`
  - `pkg/agent/swarm/deadlock_detector.go`
  - `pkg/agent/telemetry/bridge_hud.go`
  - `pkg/agent/telemetry/event_streamer.go`
  - `cli/cmdcluster/cluster_cmd.go`
  - `cli/cmdswarm/telemetry_cmd.go`

---

## 1. Context & Objectives

This subtask delivers enterprise distributed SQLite multi-cluster synchronization, Merkle-tree reconciliation, and the Antigravity multi-agent swarm telemetry bridge:
1. **SQLite Split-DB Multi-Cluster Replicator & Cross-Node Ledger Synchronization (Cluster 79, Items 781–790):** Consistent hash ring node routing, streaming WAL frame replication, page-level Merkle-tree diffing, two-phase distributed epoch commits, online non-blocking vacuum compaction, and split-db cluster management.
2. **Antigravity Multi-Agent Swarm Orchestrator & Live Warmup Telemetry Bridge (Cluster 80, Items 791–800):** Real-time multi-agent swarm mission cockpit, peer discovery gossip mesh, subtask leasing arbiters, token expenditure governors, automated circular deadlock detection, and live warmup telemetry bridge HUD.
3. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isLeaderNode`, `isShardHealthy`, `isSyncedWithMaster`, `isWithinBudget`, `isAgentAlive`, `isTelemetryStreaming`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: SQLite Split-DB Multi-Cluster Replicator
- Build consistent hash ring topology mapper in `pkg/db/cluster/shard_topology.go`.
- Implement streaming WAL frame replicator across cluster nodes in `pkg/db/cluster/wal_streamer.go`.
- Implement two-phase distributed epoch commit coordinator in `pkg/db/cluster/two_phase_commit.go`.
- Construct non-blocking online vacuum compaction scheduler in `pkg/db/cluster/vacuum_scheduler.go`.

### Step 2: Page-Level Merkle-Tree Conflict Resolution
- Build SQLite page-level cryptographic Merkle-tree generator in `pkg/db/merkle/page_tree.go`.
- Implement efficient range diffing and multi-master reconciliation in `pkg/db/merkle/range_diff.go`.

### Step 3: Antigravity Swarm Orchestration Engine
- Build peer discovery gossip mesh and heartbeat tracker in `pkg/agent/swarm/peer_mesh.go`.
- Implement subtask exclusivity lease arbiter in `pkg/agent/swarm/lease_arbiter.go`.
- Build token budget hard-cap guard in `pkg/agent/swarm/token_guard.go`.
- Implement circular dependency and agent deadlock detector in `pkg/agent/swarm/deadlock_detector.go`.

### Step 4: Live Warmup Telemetry Bridge HUD
- Implement terminal telemetry bridge HUD in `pkg/agent/telemetry/bridge_hud.go`.
- Build real-time warmup event streamer and buffer in `pkg/agent/telemetry/event_streamer.go`.

### Step 5: CLI Command Dispatch
- Wire multi-cluster management commands into `cli/cmdcluster/cluster_cmd.go`.
- Wire swarm telemetry bridge commands into `cli/cmdswarm/telemetry_cmd.go`.
