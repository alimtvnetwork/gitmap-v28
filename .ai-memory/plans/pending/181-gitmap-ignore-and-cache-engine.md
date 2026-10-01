# Plan: 181-gitmap-ignore-and-cache-engine

Spec Reference: [02-spec/21-app/181-gitmap-ignore-and-cache-engine/01-overview.md](../../../02-spec/21-app/181-gitmap-ignore-and-cache-engine/01-overview.md)

## Architectural Context & Blast Radius
This task implements cross-repository features for Gitmap (pull all, commit all, ignore all), split-DB caching engine, and SSH node execution (PAS formula).
- **Blast Radius**: Modifies `cli/` command routers and handlers. Introduces new split-DB schema in `internal/cache/` or similar. Adjusts task enqueue logic for history tracking.

## Deliverables & Subtasks
1. **Task-01: Gitmap Pull All & SSH Optimization**
   - Subtask 1: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/01-pull-all-optimization.md`
2. **Task-02: Gitmap Fix Ignore All (Local & SSH)**
   - Subtask 2: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/02-fix-ignore-all.md`
3. **Task-03: Gitmap Commit Push All Repos (cpar)**
   - Subtask 3: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/03-commit-push-all.md`
4. **Task-04: Gitmap Ignore Group Engine**
   - Subtask 4: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/04-ignore-engine.md`
5. **Task-05: Gitmap Cache Engine Split-DB**
   - Subtask 5: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/05-cache-engine.md`
6. **Task-06: Gitmap "See" Commands**
   - Subtask 6: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/06-see-commands.md`
7. **Task-07: Fix Existing Bug (Screenshot)**
   - Subtask 7: `.ai-memory/plans/subtasks/181-gitmap-ignore-and-cache-engine/07-fix-screenshot-bug.md`
