# Subtask 01: CLI Flags, Arguments, & Types Architecture

## Status: Completed
- **Parent Plan:** [293-gitmap-plans-command-and-fleet-sync.md](../../pending/293-gitmap-plans-command-and-fleet-sync.md)
- **Target Files:**
  - `cli/constants/constants_commands_extra.go`
  - `cli/cmdplans/plans_types.go`

---

## 1. Context & Objectives
Define typed constants for `CmdPlans` ("plans") and `CmdPlansAlias` ("plan"). Define `PlanStatusType` enums, `PlanItem`, `PlansSummary`, `PlansReport`, and `PlansOptions` structs with positive booleans (`isJSON`, `isYAML`, `isAll`, `hasMorePending`, `hasMoreCompleted`).

---

## 2. Implementation Checklist
- [ ] Add `CmdPlans = "plans"` and `CmdPlansAlias = "plan"` to `cli/constants/constants_commands_extra.go`.
- [ ] Create `cli/cmdplans/plans_types.go`.
- [ ] Ensure all functions $\le 15$ lines and files $\le 100$ lines.
