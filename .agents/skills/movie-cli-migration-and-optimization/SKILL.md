---
name: movie-cli-migration-and-optimization
description: Migration and DRY optimization plan orchestration for movie-cli-v8.
---

# Movie CLI Migration and Optimization

This skill encapsulates the workflow for migrating a legacy `.lovable` repository to the modern `.ai-memory` and `02-spec` folder architecture, and subsequently performing DRY code optimizations on the target codebase (specifically replacing legacy error wrappers with `pkg/appfault` and enforcing strict boolean standards).

## Instructions

1. **Phase 1 (Planning & Spec Generation)**:
   - Perform a `git pull` on the target repository.
   - Synchronize the `.ai-memory`, `02-spec`, `01-prompts`, `03-ai-scripts`, and `agents.md` from the central orchestration repo to the target repo.
   - Move existing `spec/08-app` and `spec/09-app-issues` to `02-spec/21-app/` and `02-spec/22-app-issues/` respectively.
   - Migrate `.lovable/memory/` into `.ai-memory/memory/`.
   - Delete legacy `.lovable` and `spec` folders.
   - Create a DRY optimization plan focusing on error standardization, strict boolean enforcement, and struct unification.

2. **Phase 2 (Execution)**:
   - Spawn execution subagents to safely execute the migration commands.
   - Spawn execution subagents to refactor the Go codebase (replace `apperror` with `appfault`, fix boolean flag names, extract embedded options structs).
   - Ensure line length formatting (<100 characters) is adhered to.
   - NO test running or build checking is permitted during this loop.

3. **Consolidation**:
   - Merge subtasks into a single `.ai-memory/plans/completed/` log.
   - Commit and push all changes in a single atomic git commit.
