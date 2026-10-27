# Subtask Plan 04: Distributed SQLite Split-DB Multi-Master Replication HUD & Swarm Steerer (Improvements 881–900)

- **Parent Task:** `345-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `345-04`
- **Spec Path:** `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/04-improvements-861-to-900-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 89: Distributed SQLite Split-DB Multi-Master Replication HUD (881–890)
  - Cluster 90: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (891–900)

## Implementation Steps

1. **Distributed Multi-Master SQLite Replication Engine (881–885):**
   - Implement vector-clock CRDT state synchronization for distributed Split-DB instances in `pkg/db/cluster/`.
   - Stream SQLite WAL frames across cluster nodes via encrypted gRPC / TLS streams.
   - Implement two-phase distributed epoch commit logs and optimistic concurrency controls.

2. **Split-DB Replication HUD & Self-Healing Ledger Integrity (886–890):**
   - Build interactive terminal HUD displaying replication lag, WAL frame velocity, and node sync status in Lipgloss.
   - Implement Merkle tree hash comparisons for database page integrity audits.
   - Automate non-blocking hot vacuum compaction and rolling point-in-time recovery snapshots.

3. **Antigravity Multi-Agent Swarm Orchestration Engine (891–895):**
   - Implement autonomous subtask scheduling and task lease auctioning across agent worker nodes in `pkg/agent/swarm/`.
   - Enforce token budget expenditure guards and context window resource limits.
   - Detect agent task starvation, circular deadlock loops, and coordinate dynamic failover.

4. **Autonomous Deliverability Feedback Loop & Warmup Pace Steerer (896–900):**
   - Collect real-time bounce rates, spam complaints, and inbox placement telemetry into an autonomous feedback loop in `pkg/agent/steerer/`.
   - Dynamically adjust dispatch volume curves, inter-turn delays, and recipient pool allocations per domain.
   - Render swarm mission telemetry HUD with real-time health scores, active agent roles, and throttle velocity vectors.

## Verification Checklist

- [ ] Multi-master SQLite replication verifies zero data loss and deterministic conflict resolution across 3 simulated nodes.
- [ ] Replication HUD renders live WAL streaming statistics with ANSI tokens and zero external web dependencies.
- [ ] Swarm scheduler recovers leased subtasks within 10s of simulated agent worker failure.
- [ ] Warmup steerer throttles dispatch velocity down by 50% immediately upon receiving 3 consecutive spam placement signals.
