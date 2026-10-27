# Subtask 04: Antigravity Multi-Agent Swarm Orchestrator & Live Task Steerer

## Status: READY
- **Subtask Slug:** `04-antigravity-swarm-orchestrator`
- **Parent Plan:** `.ai-memory/plans/pending/341-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/04-improvements-461-to-500-security-audit-and-agy-bridge.md`
- **Target Range:** Improvements 491–500 (Cluster 50)
- **Target Subsystem:** `pkg/agent/swarm/`, `cli/cmdwarmup/`
- **Exclusive Owned Files:**
  - `pkg/agent/swarm/supervisor.go`
  - `pkg/agent/swarm/steerer.go`
  - `pkg/agent/swarm/budget_controller.go`
  - `pkg/agent/swarm/bus_bridge.go`
  - `pkg/agent/swarm/convergence_monitor.go`
  - `pkg/agent/swarm/task_decomposer.go`
  - `pkg/agent/swarm/slash_handler.go`
  - `pkg/agent/swarm/memory_store.go`
  - `pkg/agent/swarm/audio_engine.go`
  - `pkg/agent/swarm/cockpit_presenter.go`
  - `cli/cmdwarmup/swarm_cmd.go`

---

## 1. Context & Objectives

This subtask delivers the Antigravity multi-agent swarm orchestrator, live dynamic task steerer, and executive master cockpit:
1. **Real-Time Swarm Supervisor & Task Steering (Items 491–492):** Supervise specialized AGY subagents with dynamic human-in-the-loop task steering, pausing, and directive injection.
2. **Token Cost Governance & AGY Bus Bridge (Items 493–494):** Enforce strict LLM token budgets with automated tier downgrading and maintain low-latency WebSocket inter-agent communication.
3. **Autonomous Convergence, Micro-Batches & Slash Commands (Items 495–497):** Swarm deadlock detection, self-healing restarts, 50-recipient micro-batch decomposition, and `/warmup` conversational slash command handlers.
4. **Persistent Learning, Audio Telemetry & Master Cockpit (Items 498–500):** Cross-session heuristic memory store, pure Web Audio API synthesizer feedback, and unified master terminal/web command cockpit.
5. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isAgentActive`, `isInterruptionAck`, `isSwarmOperational`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Real-Time Multi-Agent Supervisor & Live Steerer
- Construct supervisor state machine managing agent lifecycles in `pkg/agent/swarm/supervisor.go`.
- Implement interactive operator steering HUD and directive dispatcher in `pkg/agent/swarm/steerer.go`.

### Step 2: Token Budget Governance & Event Bus WebSocket Bridge
- Implement LLM token usage tracking and tier downgrading in `pkg/agent/swarm/budget_controller.go`.
- Build bidirectional WebSocket bus bridge in `pkg/agent/swarm/bus_bridge.go`.

### Step 3: Swarm Convergence, Task Decomposer & Slash Dispatcher
- Implement deadlock detection and restart manager in `pkg/agent/swarm/convergence_monitor.go`.
- Build micro-batch task divider in `pkg/agent/swarm/task_decomposer.go`.
- Implement conversational `/warmup` slash command handler in `pkg/agent/swarm/slash_handler.go`.

### Step 4: Persistent Memory Ledger, Audio Synthesizer & Master Cockpit
- Construct cross-session heuristic memory store in `pkg/agent/swarm/memory_store.go`.
- Build Web Audio / terminal bell acoustic feedback engine in `pkg/agent/swarm/audio_engine.go`.
- Deliver the grand 4-quadrant unified master cockpit in `pkg/agent/swarm/cockpit_presenter.go` and `cli/cmdwarmup/swarm_cmd.go`.
