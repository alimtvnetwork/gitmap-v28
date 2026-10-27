# Subtask 02: Automated Rescue Daemons & Interactive Terminal Studio

## Status: READY
- **Subtask Slug:** `02-rescue-daemons-and-terminal-studio`
- **Parent Plan:** `.ai-memory/plans/pending/339-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Owner:** Spec Author Subagent 2 / TUI & Delivery Resilience Engineer
- **Target Subsystem:** `pkg/mail/rescue/`, `pkg/cli/tui/`, `pkg/mail/seed/`
- **Exclusive Owned Files:**
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/03-improvements-231-to-260-rescue-daemons-and-terminal-studio.md`
  - `pkg/mail/rescue/daemon_supervisor.go`
  - `pkg/cli/tui/studio_console.go`
  - `pkg/mail/seed/reply_matrix.go`
  - `cli/cmdmail/cmd_studio.go`

---

## 1. Context & Objectives

This subtask implements automated reputation protection daemons and rich command-line operator consoles:
1. **Automated Inbox Rescue Daemons (Cluster 24, Items 231–240):** Autonomous background supervisor monitoring deliverability dips, initiating synthetic seed interactions, executing cooling ramps, and evading spam traps.
2. **Reputation Healing Ramps (Cluster 25, Items 241–250):** Controlled step-by-step restoration protocols for damaged sender mailboxes, rebalancing traffic ratios toward high-engagement seed clusters.
3. **Interactive Terminal TUI Studio (Cluster 26, Items 251–260):** Full-screen Lipgloss-driven terminal dashboard (`gitmap warmup studio`) featuring real-time sparklines, interactive mailbox selection, delivery waterfall monitors, and hotkey actions.
4. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isRescueActive`, `hasReputationRestored`, `isStudioRunning`).
   - Typed error returns via `*apperror.AppError`.
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Autonomous Inbox Rescue Daemon
- Construct `pkg/mail/rescue/daemon_supervisor.go` with periodic health polling against sender mailboxes.
- Implement spam trap detection heuristics and immediate circuit isolation.

### Step 2: Reputation Healing Ramp Protocols
- Author gradual volume ramp algorithms with adaptive seed response weighting in `pkg/mail/seed/reply_matrix.go`.
- Synchronize recovery state flags directly into SQLite `warmup.db`.

### Step 3: Interactive Terminal Studio (`gitmap warmup studio`)
- Build terminal dashboard in `pkg/cli/tui/studio_console.go` utilizing Charmbracelet Lipgloss / Bubbletea.
- Integrate hotkeys: `[q]` quit, `[p]` pause warming, `[r]` force recovery ramp, `[d]` inspect diagnostics.

---

## 3. Verification & Acceptance Criteria
- [ ] Rescue daemon detects simulated 5% bounce anomaly within 15 seconds and trips isolation circuit.
- [ ] Terminal studio renders smoothly on standard 80x24 and 120x40 terminal windows with zero screen flicker.
- [ ] Seed reply matrix maintains realistic conversational threads across synthetic seed nodes.
