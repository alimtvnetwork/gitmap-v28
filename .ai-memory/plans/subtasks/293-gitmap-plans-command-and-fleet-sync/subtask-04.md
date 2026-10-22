# Subtask 04: Tokyo Night Lipgloss Table & Metric Pills Renderer

## Status: Pending
- **Parent Plan:** [293-gitmap-plans-command-and-fleet-sync.md](../../pending/293-gitmap-plans-command-and-fleet-sync.md)
- **Target Files:**
  - `cli/cmdplans/plans_render.go`

---

## 1. Context & Objectives
Construct Tokyo Night tables using `termout.PrintTable` (`BorderRounded`, `RowStripe: true`), `TokyoOKPill`, `TokyoWarnPill`, and `TokyoMetricPill`. Render completed and pending tables with clear pending overflow guidance.

---

## 2. Implementation Checklist
- [ ] Create `cli/cmdplans/plans_render.go`.
- [ ] Implement `RenderTTY(report PlansReport)`.
- [ ] Implement table builders for completed and pending lists.
- [ ] Implement overflow guidance renderer.
