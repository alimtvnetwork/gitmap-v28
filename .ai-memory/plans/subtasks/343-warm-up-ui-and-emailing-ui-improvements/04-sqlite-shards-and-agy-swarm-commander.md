# Subtask 04: Distributed SQLite Sharded Replicator & Autonomous Antigravity Swarm Commander

## Status: READY
- **Subtask Slug:** `04-sqlite-shards-and-agy-swarm-commander`
- **Parent Plan:** `.ai-memory/plans/pending/343-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/04-improvements-661-to-700-sharding-audit-and-agy-fleet.md`
- **Target Range:** Improvements 681–700 (Clusters 69–70)
- **Target Subsystem:** `pkg/db/sharding/`, `pkg/db/merkle/`, `pkg/agent/swarm/`
- **Exclusive Owned Files:**
  - `pkg/db/sharding/shard_topology.go`
  - `pkg/db/sharding/wal_streamer.go`
  - `pkg/db/sharding/epoch_lease.go`
  - `pkg/db/sharding/vacuum_compact.go`
  - `pkg/db/merkle/tree_builder.go`
  - `pkg/db/merkle/conflict_resolver.go`
  - `pkg/agent/swarm/commander_hud.go`
  - `pkg/agent/swarm/peer_mesh.go`
  - `pkg/agent/swarm/task_arbiter.go`
  - `pkg/agent/swarm/budget_guard.go`
  - `pkg/agent/swarm/deadlock_watcher.go`
  - `pkg/agent/swarm/transcript_indexer.go`
  - `cli/cmdcluster/shard_manage_cmd.go`
  - `cli/cmdswarm/commander_cmd.go`

---

## 1. Context & Objectives

This subtask delivers carrier-grade distributed SQLite sharding with cryptographic Merkle-tree conflict resolution and the next-generation Antigravity Swarm Mission Commander:
1. **Distributed SQLite Sharded Replicator & Merkle-Tree Conflict Resolver (Cluster 69, Items 681–690):** Physical shard routing, streaming WAL frame replication, Merkle-tree range diffing, two-phase distributed epoch commits, and non-blocking vacuum compaction.
2. **Autonomous Antigravity Swarm Orchestration Engine & Real-Time Mission Commander (Cluster 70, Items 691–700):** Real-time multi-agent swarm cockpit, peer discovery heartbeat mesh, subtask leasing arbitration, token/LLM expenditure guards, automated circular deadlock detection, and SQLite FTS5 transcript indexing.
3. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isLeaderNode`, `isShardHealthy`, `isSyncedWithMaster`, `isWithinBudget`, `isAgentAlive`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Distributed SQLite Sharding & Replication
- Implement consistent hashing shard topology mapper in `pkg/db/sharding/shard_topology.go`.
- Build streaming WAL frame replicator across cluster nodes in `pkg/db/sharding/wal_streamer.go`.
- Implement distributed epoch leader lease arbiter in `pkg/db/sharding/epoch_lease.go`.
- Construct non-blocking online vacuum compaction scheduler in `pkg/db/sharding/vacuum_compact.go`.

### Step 2: Merkle-Tree Conflict Resolution
- Build SQLite page-level cryptographic Merkle-tree generator in `pkg/db/merkle/tree_builder.go`.
- Implement multi-master divergence detector and last-write-wins (LWW) resolver in `pkg/db/merkle/conflict_resolver.go`.

### Step 3: Antigravity Swarm Mission Commander
- Implement interactive terminal mission commander HUD in `pkg/agent/swarm/commander_hud.go`.
- Build peer discovery gossip mesh and node heartbeat monitor in `pkg/agent/swarm/peer_mesh.go`.
- Implement subtask exclusivity lease arbiter in `pkg/agent/swarm/task_arbiter.go`.
- Build token budget hard-cap guard in `pkg/agent/swarm/budget_guard.go`.
- Implement circular dependency and agent deadlock detector in `pkg/agent/swarm/deadlock_watcher.go`.
- Build conversation transcript FTS5 search indexer in `pkg/agent/swarm/transcript_indexer.go`.

### Step 4: CLI Command Dispatch
- Wire cluster sharding commands into `cli/cmdcluster/shard_manage_cmd.go`.
- Wire swarm commander commands into `cli/cmdswarm/commander_cmd.go`.
