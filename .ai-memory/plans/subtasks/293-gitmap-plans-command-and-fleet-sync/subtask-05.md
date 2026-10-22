# Subtask 05: Root Command Wiring & Dispatch Integration

## Status: Completed
- **Parent Plan:** [293-gitmap-plans-command-and-fleet-sync.md](../../pending/293-gitmap-plans-command-and-fleet-sync.md)
- **Target Files:**
  - `cli/cmdplans/plans.go`
  - `cli/cmd/rootutility.go`

---

## 1. Context & Objectives
Implement `RunPlans(args []string) error` in `cli/cmdplans/plans.go`. Parse options and limits, coordinate discovery and rendering. Wire `constants.CmdPlans` and `constants.CmdPlansAlias` in `cli/cmd/rootutility.go`.

---

## 2. Implementation Checklist
- [ ] Create `cli/cmdplans/plans.go`.
- [ ] Wire `CmdPlans` and `CmdPlansAlias` in `utilityToolEntries()`.
- [ ] Ensure functions $\le 15$ lines and early guard returns.
