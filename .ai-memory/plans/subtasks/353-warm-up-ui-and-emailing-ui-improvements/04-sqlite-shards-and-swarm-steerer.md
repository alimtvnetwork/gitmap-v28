# Subtask Plan 04: Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1681–1700)

- **Parent Task:** `353-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `353-04`
- **Spec Path:** `02-spec/21-app/353-warm-up-ui-and-emailing-ui-improvements/04-improvements-1661-to-1700-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 169: Distributed SQLite Split-DB Multi-Master Replication HUD (1681–1690)
  - Cluster 170: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (1691–1700)

## Implementation Steps

1. **Distributed SQLite Split-DB Multi-Master Replication (1681–1685):**
   - Author vector-clock CRDT state synchronization across distributed worker replicas in `pkg/db/cluster/`.
   - Implement encrypted WAL frame streaming over mutual TLS 1.3 channels.
   - Author page-level Merkle tree integrity auditor to quickly identify corrupted or desynchronized 4096-byte database pages.
   - Establish distributed two-phase epoch commit coordinator and non-blocking hot VACUUM compaction engine.

2. **Replication Cockpit & Partition Management (1686–1690):**
   - Build interactive Bubble Tea replication cluster HUD in `pkg/cli/tui/cluster_hud.go`.
   - Implement deterministic Last-Write-Wins (LWW) conflict arbiter with Blake3 cryptographic tie-breaking.
   - Author rendezvous hash-partition shard rebalancer and point-in-time recovery (PITR) WAL log replayer.
   - Enable zero-downtime rolling schema migrations via shadow tables and column additions.

3. **Antigravity Multi-Agent Swarm Orchestrator (1691–1695):**
   - Author real-time multi-agent swarm cockpit HUD rendering active agents, assigned subtasks, and elapsed time in `pkg/cli/tui/swarm_hud.go`.
   - Design dynamic subtask auctioneer with expiring leases and automatic requeueing in `pkg/agent/swarm/auctioneer.go`.
   - Track token expenditure and context window quotas with automated compaction gates.
   - Detect circular deadlocks using Tarjan's SCC algorithm and implement sub-250ms embedded Raft consensus master elections.

4. **Autonomous Warmup Steerer & Unified Mission HUD (1696–1700):**
   - Close the deliverability feedback loop with autonomous PID throttle steering in `pkg/agent/steerer/`.
   - Generate multi-turn contextual synthetic conversations between AI personas to build authentic mailbox engagement.
   - Render multi-ISP reputation vector heatmaps (Google, Microsoft, Yahoo, Apple).
   - Author domain warmup graduation gate and full-screen unified mission telemetry dashboard.

## Verification Checklist

- [ ] Vector clocks maintain strict causal consistency across multi-node SQLite shards.
- [ ] Non-blocking hot VACUUM reclaims disk space without table lock timeouts.
- [ ] Antigravity subtask auctioneer reclaims expired leases within 60s.
- [ ] All Go structs and SQLite DDL use strictly positive `is`/`has` booleans.
- [ ] 100% relative repository paths maintained with zero absolute paths.
