# Subtask 04: Multi-Node SQLite Sharded Cluster & Antigravity Fleet Choreographer

## Status: READY
- **Subtask Slug:** `04-sqlite-shards-and-antigravity-fleet`
- **Parent Plan:** `.ai-memory/plans/pending/342-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/04-improvements-561-to-600-security-audit-and-swarm-telemetry.md`
- **Target Range:** Improvements 581–600 (Clusters 59–60)
- **Target Subsystem:** `pkg/db/cluster/`, `pkg/agent/swarm/`
- **Exclusive Owned Files:**
  - `pkg/db/cluster/topology_mapper.go`
  - `pkg/db/cluster/wal_lag_sentinel.go`
  - `pkg/db/cluster/epoch_coordinator.go`
  - `pkg/db/cluster/vacuum_governor.go`
  - `pkg/db/cluster/query_router.go`
  - `pkg/db/cluster/shard_rebalancer.go`
  - `pkg/db/cluster/cluster_cockpit.go`
  - `pkg/agent/swarm/peer_discovery.go`
  - `pkg/agent/swarm/lease_arbitrator.go`
  - `pkg/agent/swarm/token_governor.go`
  - `pkg/agent/swarm/deadlock_detector.go`
  - `pkg/agent/swarm/output_verifier.go`
  - `pkg/agent/swarm/transcript_miner.go`
  - `pkg/agent/swarm/master_cockpit.go`
  - `cli/cmdcluster/shard_cmd.go`
  - `cli/cmdswarm/fleet_cmd.go`

---

## 1. Context & Objectives

This subtask delivers multi-node sharded SQLite database governance and autonomous Antigravity multi-agent fleet telemetry:
1. **Multi-Node SQLite Sharded Cluster Dashboard & Replication Sentinel (Cluster 59, Items 581–590):** Shard topology discovery, WAL frame streaming replication lag monitoring, distributed consensus epoch leases, cross-node auto-vacuum governance, sharded query routing, and cluster management cockpit.
2. **Autonomous Antigravity Agent Fleet Choreographer & Telemetry Bridge (Cluster 60, Items 591–600):** Agent peer discovery mesh, subtask lease arbitration, token budget & LLM cost governance, inter-agent message transit logging, deadlock detection, output quality verification, transcript AST indexing, and master executive cockpit.
3. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isMasterNode`, `isSynchronized`, `isLeaseActive`, `isWithinBudget`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Sharded SQLite Cluster Governance
- Implement shard topology discovery and key mapper in `pkg/db/cluster/topology_mapper.go`.
- Build streaming WAL frame replication lag monitor in `pkg/db/cluster/wal_lag_sentinel.go`.
- Implement epoch leader lease coordinator in `pkg/db/cluster/epoch_coordinator.go`.
- Construct non-blocking auto-vacuum governor and query router in `pkg/db/cluster/`.
- Wire shard cluster cockpit in `pkg/db/cluster/cluster_cockpit.go`.

### Step 2: Autonomous Antigravity Swarm Choreographer
- Implement agent peer discovery and heartbeat telemetry in `pkg/agent/swarm/peer_discovery.go`.
- Construct subtask lease arbitrator in `pkg/agent/swarm/lease_arbitrator.go` to prevent file collision.
- Build token budget consumption and LLM expense governor in `pkg/agent/swarm/token_governor.go`.
- Implement autonomous deadlock and idle swarm detector in `pkg/agent/swarm/deadlock_detector.go`.
- Build mechanical code quality audit engine in `pkg/agent/swarm/output_verifier.go`.

### Step 3: Transcript Miner & Master Executive Cockpit
- Build SQLite FTS5 conversation log indexer in `pkg/agent/swarm/transcript_miner.go`.
- Construct master command deck and fleet dashboard in `pkg/agent/swarm/master_cockpit.go`.
- Wire CLI commands into `cli/cmdcluster/shard_cmd.go` and `cli/cmdswarm/fleet_cmd.go`.
