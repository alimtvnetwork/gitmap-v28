---
name: inventory-pending-tasks
description: Discover, catalog, and sequence all pending plans, subtasks, and unresolved issues across .lovable/.
---

# Inventory Pending Tasks

Audits and catalogs pending work across `.lovable/plans/pending/`, `.lovable/plans/subtasks/`, `.lovable/issues/`, and `.lovable/cicd-issues/`.

## Workflow
1. Scan `.lovable/plans/01-index.md` and `plans/pending/`.
2. Verify subtask batch linkages in `plans/subtasks/`.
3. Check unresolved issues in `.lovable/issues/` and `.lovable/cicd-issues/`.
4. Present aligned sequence and execution readiness status.
