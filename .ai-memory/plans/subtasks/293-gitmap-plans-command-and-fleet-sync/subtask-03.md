# Subtask 03: JSON & YAML Machine-Readable Serializers

## Status: Completed
- **Parent Plan:** [293-gitmap-plans-command-and-fleet-sync.md](../../pending/293-gitmap-plans-command-and-fleet-sync.md)
- **Target Files:**
  - `cli/cmdplans/plans_format.go`

---

## 1. Context & Objectives
Implement serialization to `os.Stdout` for `--json` and `--yaml`/`--yml` flags. Ensure zero ANSI styling and clean pipeable streams.

---

## 2. Implementation Checklist
- [ ] Create `cli/cmdplans/plans_format.go`.
- [ ] Implement `RenderJSON(report PlansReport) error`.
- [ ] Implement `RenderYAML(report PlansReport) error`.
- [ ] Bubble up errors with leaf error sovereignty.
