# Subtask 01: Multi-Agent Deliverability Swarm, Cognitive Rhythms & P2P Fleet Mesh

## Status: READY
- **Subtask Slug:** `01-swarm-rhythms-and-fleet-mesh`
- **Parent Plan:** `.ai-memory/plans/pending/340-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/02-improvements-301-to-330-swarm-rhythms-and-fleet-mesh.md`
- **Owner:** Spec Author Subagent 1 / Swarm & Distributed Systems Engineer
- **Target Subsystem:** `pkg/agent/swarm/`, `pkg/mail/rhythm/`, `pkg/mail/mesh/`
- **Exclusive Owned Files:**
  - `pkg/agent/swarm/consensus_engine.go`
  - `pkg/agent/swarm/agent_roles.go`
  - `pkg/mail/rhythm/circadian_modeler.go`
  - `pkg/mail/rhythm/queue_visualizer.go`
  - `pkg/mail/mesh/peer_router.go`
  - `pkg/mail/mesh/seed_sync.go`
  - `cli/cmdwarmup/cmd_swarm.go`

---

## 1. Context & Objectives

This subtask delivers the autonomous multi-agent swarm and peer-to-peer mesh layers for Sequence 340:
1. **Multi-Agent Deliverability Swarm (Cluster 31, Items 301–310):** Establish 10 specialized deliverability micro-agents exchanging consensus telemetry to auto-balance outbound ramps and inbox placement.
2. **Cognitive Sending Rhythms & Queue Saturation (Cluster 32, Items 311–320):** Model human circadian patterns, Poisson burst distributions, and real-time ISP reception queue fatigue.
3. **Multi-Node P2P Fleet Mesh (Cluster 33, Items 321–330):** Decentralized peer-to-peer fleet coordination, zero-trust peer gossip, and cross-region seed inbox synchronization.
4. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isSwarmActive`, `hasValidConsensus`, `isPeerConnected`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Multi-Agent Swarm Framework & Roles
- Define the 10 micro-agent roles and state machine in `pkg/agent/swarm/agent_roles.go`.
- Implement consensus voting and state synchronization in `pkg/agent/swarm/consensus_engine.go`.

### Step 2: Cognitive Circadian Rhythms & Queue Saturation
- Construct mathematical circadian curve modeler in `pkg/mail/rhythm/circadian_modeler.go`.
- Build real-time ISP queue friction visualizer in `pkg/mail/rhythm/queue_visualizer.go`.

### Step 3: P2P Mesh Routing & Seed Sync
- Develop zero-trust peer gossip router in `pkg/mail/mesh/peer_router.go`.
- Implement cross-cloud seed inbox state synchronizer in `pkg/mail/mesh/seed_sync.go`.

### Step 4: CLI Swarm Management
- Expose swarm status and peer inspection commands in `cli/cmdwarmup/cmd_swarm.go`.
