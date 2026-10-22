# Subtask 02: Plan File Discovery & Metadata Parser

## Status: Completed
- **Parent Plan:** [293-gitmap-plans-command-and-fleet-sync.md](../../pending/293-gitmap-plans-command-and-fleet-sync.md)
- **Target Files:**
  - `cli/cmdplans/plans_parser.go`

---

## 1. Context & Objectives
Traverse `.ai-memory/plans/` (root), `.ai-memory/plans/pending/`, and `.ai-memory/plans/completed/`. Extract sequence ID, slug, title from H1 headings, and normalized status from markdown headers. Exclude index/ledger files. Sort completed plans descending by ID, pending plans ascending by ID.

---

## 2. Implementation Checklist
- [ ] Create `cli/cmdplans/plans_parser.go`.
- [ ] Implement `DiscoverPlans(repoPath string) ([]PlanItem, []PlanItem, error)`.
- [ ] Implement `parsePlanFile(filePath, relPath, defaultStatus string) (*PlanItem, bool)`.
- [ ] Guard against nested ifs, adhere to $\le 15$ line limits.
