# Subtask 04: Antigravity IDE Event Bus & Autonomous Deliverability Scorecards

## Status: READY
- **Subtask Slug:** `04-agy-event-bus-and-audit`
- **Parent Plan:** `.ai-memory/plans/pending/340-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/04-improvements-361-to-400-design-tokens-splitdb-agy.md`
- **Owner:** Spec Author Subagent 3 / Agent Platform & Audit Specialist
- **Target Subsystem:** `pkg/agent/bus/`, `pkg/agent/anomaly/`, `pkg/mail/audit/`, `pkg/db/audit/`
- **Exclusive Owned Files:**
  - `pkg/agent/bus/event_broker.go`
  - `pkg/agent/anomaly/triage_engine.go`
  - `pkg/agent/bus/webhook_dispatcher.go`
  - `pkg/mail/audit/scorecard_engine.go`
  - `pkg/mail/audit/compliance_matrix.go`
  - `pkg/db/audit/merkle_tree.go`
  - `pkg/mail/coordinator/self_healing_master.go`
  - `cli/cmdwarmup/slash_dispatcher.go`

---

## 1. Context & Objectives

This subtask delivers the Antigravity IDE event bus, autonomous deliverability audit scorecards, and cryptographic Merkle verification:
1. **Antigravity IDE Event Bus & Dynamic Workspace Hooks (Cluster 39, Items 381–390):** Bidirectional pub/sub WebSocket message broker, IDE activity throttling sensors, slash-command dispatcher (`/warmup`), statistical anomaly triage alerts, and signed webhooks.
2. **Autonomous Deliverability Scorecards & Historical Audit (Cluster 40, Items 391–400):** Multi-factor 0–100 deliverability scorecard, ISP compliance radar matrix, SHA-256 Merkle tree verification ledger, and autonomous self-healing coordinator.
3. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isWebSocketConnected`, `hasPassedAudit`, `isMerkleValid`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Antigravity IDE Event Broker & Slash Commands
- Construct bidirectional pub/sub WebSocket broker in `pkg/agent/bus/event_broker.go`.
- Implement statistical anomaly triage engine (Z-score calculation) in `pkg/agent/anomaly/triage_engine.go`.
- Build universal chat slash-command trigger parser in `cli/cmdwarmup/slash_dispatcher.go`.
- Implement signed HMAC-SHA256 webhook dispatcher in `pkg/agent/bus/webhook_dispatcher.go`.

### Step 2: Autonomous Deliverability Scorecards & Compliance Radar
- Construct 10-dimension algorithmic deliverability health scorecard engine in `pkg/mail/audit/scorecard_engine.go`.
- Implement ISP-specific compliance radar checking Google, Yahoo, and Microsoft guidelines in `pkg/mail/audit/compliance_matrix.go`.

### Step 3: Cryptographic Merkle Tree Audit Ledger & Self-Healing
- Build SHA-256 Merkle tree tamper-evident verification ledger in `pkg/db/audit/merkle_tree.go`.
- Coordinate closed-loop autonomous deliverability self-healing in `pkg/mail/coordinator/self_healing_master.go`.
