# Subtask 02: Content Entropy Synthesizers, ISP Circuit Breakers & Terminal Studio

## Status: READY
- **Subtask Slug:** `02-self-healing-circuits-and-tui`
- **Parent Plan:** `.ai-memory/plans/pending/340-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/03-improvements-331-to-360-entropy-circuits-and-tui.md`
- **Owner:** Spec Author Subagent 2 / Protocol & TUI Interface Engineer
- **Target Subsystem:** `pkg/mail/entropy/`, `pkg/mail/circuit/`, `pkg/cli/tui/`
- **Exclusive Owned Files:**
  - `pkg/mail/entropy/synthesizer.go`
  - `pkg/mail/entropy/spintax_engine.go`
  - `pkg/mail/circuit/breaker_fsm.go`
  - `pkg/mail/circuit/throttle_listener.go`
  - `pkg/cli/tui/operator_studio.go`
  - `pkg/cli/tui/circuit_gauge.go`
  - `cli/cmdwarmup/cmd_studio.go`

---

## 1. Context & Objectives

This subtask delivers automated content entropy, ISP fault tolerance, and interactive terminal studio controls:
1. **Autonomous Content Entropy Synthesizers (Cluster 34, Items 331–340):** Natural spintax variation, structural MIME diversification, and dynamic header entropy preventing spam-fingerprint collisions.
2. **Real-Time ISP Circuit Breakers & Throttle Listeners (Cluster 35, Items 341–350):** 3-state finite state machine responding to SMTP error codes (421, 451, 452) with automatic rate throttling.
3. **Interactive Terminal Operator Studio (Cluster 36, Items 351–360):** Comprehensive Lipgloss terminal UI with live charts, circuit status toggles, and manual override controls.
4. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isCircuitClosed`, `hasHighEntropy`, `isStudioRunning`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Content Entropy & Spintax Generation
- Implement dynamic spintax token evaluator in `pkg/mail/entropy/spintax_engine.go`.
- Construct structural entropy and header randomizer in `pkg/mail/entropy/synthesizer.go`.

### Step 2: ISP Circuit Breaker FSM
- Develop finite state machine (`Closed`, `Open`, `Half-Open`) in `pkg/mail/circuit/breaker_fsm.go`.
- Build real-time SMTP error code listener and backoff calculator in `pkg/mail/circuit/throttle_listener.go`.

### Step 3: Interactive Terminal Studio
- Construct full-screen Lipgloss terminal studio layout in `pkg/cli/tui/operator_studio.go`.
- Implement live ANSI circuit gauges and interactive keybindings in `pkg/cli/tui/circuit_gauge.go`.

### Step 4: CLI Integration
- Wire up interactive studio launch command in `cli/cmdwarmup/cmd_studio.go`.
