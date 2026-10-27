# Subtask Plan 04: Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1181–1200)

- **Parent Task:** `348-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `348-04`
- **Spec Path:** `02-spec/21-app/348-warm-up-ui-and-emailing-ui-improvements/04-improvements-1161-to-1200-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 119: Distributed SQLite Split-DB Multi-Master Replication HUD (1181–1190)
  - Cluster 120: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (1191–1200)

## Implementation Steps

1. **Distributed Multi-Master SQLite Shard Replication Engine (1181–1185):**
   - Implement vector-clock CRDT state synchronization across distributed Split-DB instances in `pkg/db/cluster/`.
   - Stream encrypted WAL frames over mutual TLS streams with non-blocking backpressure and ChaCha20-Poly1305 encryption.
   - Implement distributed epoch commit coordinators and optimistic concurrency validation.
   - Build page-level Merkle tree integrity auditors and non-blocking hot VACUUM compaction engines.

2. **Split-DB Replication HUD & Merkle Page Verification (1186–1190):**
   - Author terminal replication HUD tracking WAL frame velocity, replica lag, and sync status in Lipgloss.
   - Implement deterministic Last-Write-Wins (LWW) conflict arbitration with tie-breaker ordering.
   - Implement dynamic hash-partition shard rebalancers, point-in-time recovery (PITR) log replay, and zero-downtime rolling schema migrations.
   - Persist replication metadata into `repodb/pipeline.db` using positive schema flags (`isClusterHealthy`, `hasQuorum`, `isDegraded`).

3. **Antigravity Multi-Agent Swarm Orchestrator (1191–1195):**
   - Architect autonomous subtask scheduling and task lease auctioning across agent worker nodes in `pkg/agent/swarm/`.
   - Enforce per-agent token budget expenditure guards and context window management.
   - Implement Tarjan-based circular deadlock detection and autonomous worker respawn / failover.
   - Implement embedded Raft consensus module with sub-250ms leader elections.

4. **Deliverability Feedback Loop & Autonomous Warmup Steerer (1196–1200):**
   - Ingest live bounce rates, spam placements, and seed reply sentiment into feedback governor in `pkg/agent/steerer/`.
   - Dynamically adjust dispatch volume curves, inter-turn delays, and tenant recipient pools per domain.
   - Generate multi-turn synthetic conversation threads between seed personas.
   - Author TrueColor terminal reputation scorecard visualizer, autonomous graduation certifier, and unified executive mission telemetry HUD.

## Verification Checklist

- [ ] Distributed SQLite multi-master replication resolves write conflicts deterministically across 3 cluster nodes.
- [ ] Replication HUD displays live WAL streaming telemetry with Lipgloss zero-CDN terminal formatting.
- [ ] Swarm scheduler detects unresponsive agents and re-auctions leased tasks within 10 seconds.
- [ ] Autonomous steerer throttles dispatch velocity down by 50% immediately upon receiving 3 consecutive spam signals.
- [ ] Full-screen mission HUD consolidates swarm, database, and deliverability telemetry without UI latency.
