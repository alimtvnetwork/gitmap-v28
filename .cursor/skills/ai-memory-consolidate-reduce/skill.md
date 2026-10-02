---
name: ai-memory-consolidate-reduce
description: Autonomously cluster, consolidate, and re-sequence completed plan files and subtasks in .ai-memory/plans/ to drastically reduce file count while preserving all verified outcomes.
---

# AI Memory Consolidation, Reduction & Safety Backup

Trigger Keywords & Aliases: `ai-memory-consolidate-reduce`, `consolidate-plans`, `consolidate completed plans`, `clean completed plans`, `resequence completed plans`, `merge plans`, `archive completed plans`, `cleanup plans completed`, `memory consolidation`, `backup and consolidate plans`, `compact plans`, `reduce plan file count`, `compact completed plans`

## Purpose
Autonomously create a timestamped backup branch, scan, analyze, cluster, aggressively consolidate, and re-sequence all completed plan files and subtasks within `.ai-memory/plans/` into minimal, hyper-compact milestone summaries, combining 2, 3, or more related tasks and common checklists into single files to drastically reduce total file count while strictly preserving 100% of core architectural concepts, verified task outcomes, error contracts, and decision logs.

## Step 0: Mandatory Safety Backup Protocol
Before modifying or deleting plans:
```bash
git pull origin $(git rev-parse --abbrev-ref HEAD)
BACKUP_BRANCH="backup/plans-consolidation-$(date +%Y%m%d-%H%M%S)"
git branch "$BACKUP_BRANCH"
git push origin "$BACKUP_BRANCH"
```

## Core Compaction Rules

1. **Cluster Tasks:** Combine 2, 3, or more related tasks into single milestone files.
2. **Prune Guideline Noise:** Remove pure formatting or routine guideline-fix micro-plans that contain zero business logic.
3. **Checklist Compaction:** Consolidate repeated checklists into a single pointer to `.ai-memory/coding-guidelines.md`.
4. **Monotonic Sequencing:** Re-sequence `.ai-memory/plans/completed/` continuously (`01-`, `02-`, `03-`, ...) with lowercase filenames and zero sequence gaps.
5. **Index Synchronization:** Update `.ai-memory/plans/readme.md` and `.ai-memory/what-to-read.md`.
