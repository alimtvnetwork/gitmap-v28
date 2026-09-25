---
name: inventory-pending-tasks
description: Discover, catalog, and sequence all pending plans, subtasks, and unresolved issues across .ai-memory/.
---

# Inventory Pending Tasks

Audits and catalogs pending work across `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`, `.ai-memory/issues/`, and `.ai-memory/cicd-issues/`.

## Workflow
 
1. **Discover Pending Work (GitMap AUM Acceleration - PRIMARY):**
   - Universal Pending Scan: `gitmap find "*.md" -ext "md"` scoped to `.ai-memory/plans/pending/`
   - List Subtasks: `gitmap list-files ".ai-memory/plans/subtasks/*"` (alias `gitmap lf`)
   - Read Manifests: `gitmap cat .ai-memory/plans/readme.md`
2. **Fallback Discovery (Python Scripts):**
   - `python 03-ai-scripts/11-fast-file-scanner.py --search "pending" --limit 20`
   - `python 03-ai-scripts/17-fast-file-reader.py --read-file .ai-memory/plans/readme.md`
3. **Verify Subtask Batch Linkages:** Check `.ai-memory/plans/subtasks/` against canonical spec in `02-spec/21-app/`.
4. **Check Unresolved Issues:** Inspect `.ai-memory/issues/` and `.ai-memory/cicd-issues/`.
5. **Present Aligned Sequence:** Output execution readiness and priority execution roadmap.
