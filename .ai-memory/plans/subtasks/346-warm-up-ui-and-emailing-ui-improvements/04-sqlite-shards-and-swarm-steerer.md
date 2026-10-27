# Subtask Plan 04: Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 981–1000)

- **Parent Task:** `346-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `346-04`
- **Spec Path:** `02-spec/21-app/346-warm-up-ui-and-emailing-ui-improvements/04-improvements-961-to-1000-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 99: Distributed SQLite Split-DB Multi-Master Replication HUD (981–990)
  - Cluster 100: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (991–1000)

## Implementation Steps

1. **Distributed Multi-Master SQLite Shard Replication Engine (981–985):**
   - Implement vector-clock CRDT state synchronization across distributed Split-DB instances in `pkg/db/cluster/`.
   - Stream encrypted WAL frames over mutual TLS streams with non-blocking backpressure.
   - Implement distributed epoch commit coordinators and optimistic concurrency validation.

2. **Split-DB Replication HUD & Merkle Page Verification (986–990):**
   - Author terminal replication HUD tracking WAL frame velocity, replica lag, and sync status in Lipgloss.
   - Implement page-level Merkle tree hash verification for database integrity auditing.
   - Provide automated rolling hot vacuum compaction and point-in-time recovery log snapshots.

3. **Antigravity Multi-Agent Swarm Orchestrator (991–995):**
   - Architect autonomous subtask scheduling and task lease auctioning across agent worker nodes in `pkg/agent/swarm/`.
   - Enforce per-agent token budget expenditure guards and context window management.
   - Implement circular deadlock detection and autonomous worker respawn / failover.

4. **Deliverability Feedback Loop & Autonomous Warmup Steerer (996–1000):**
   - Ingest live bounce rates, spam placements, and seed reply sentiment into feedback governor in `pkg/agent/steerer/`.
   - Dynamically adjust dispatch volume curves, inter-turn delays, and tenant recipient pools per domain.
   - Render swarm mission telemetry HUD with real-time health scores, active agent roles, and throttle velocity vectors.

## Verification Checklist

- [ ] Distributed SQLite multi-master replication resolves write conflicts deterministically across 3 cluster nodes.
- [ ] Replication HUD displays live WAL streaming telemetry with Lipgloss zero-CDN terminal formatting.
- [ ] Swarm scheduler detects unresponsive agents and re-auctions leased tasks within 10 seconds.
- [ ] Autonomous steerer throttles dispatch velocity down by 50% immediately upon receiving 3 consecutive spam signals.
