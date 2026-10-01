# Parent Plan: 197-pas-fix
Spec Reference: [02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md](../../../02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md)

## Architectural Context
The user has requested sweeping changes to GitMap, consisting of 4 core pillars:
1. `gitmap pull all (pa)` async ignore-check refactoring and the SSH Fleet `pas` command formula.
2. The `gitmap ignore` group management system (`agtd`, `cgwp`, import/export).
3. The `gitmap fix ignore all (fia)`, `commit-push-all-repos (cpar)`, and `see (c)` commands.
4. The high-performance `gitmap cache` Split-DB system and multi-grep capabilities.

We also have outstanding bugfixes from a prior cycle:
- `nodes clone` alignment and `w3` liveness misreporting.
- `v6.436.0` Git Credential Store fix across fleet.

## Deliverables & Subtasks Map

- **Task-01**: Fix `nodes clone` Table Alignment & w3 Resilience.
  - Subtask: `.ai-memory/plans/subtasks/197-pas-fix/01-nodes-clone-fix.md`
- **Task-02**: Implement `gitmap pull all (pa)` Async Ignore Check Architecture.
  - Subtask: `.ai-memory/plans/subtasks/197-pas-fix/02-pa-async-refactor.md`
- **Task-03**: Create base `gitmap ignore` commands and group logic.
  - Subtask: `.ai-memory/plans/subtasks/197-pas-fix/03-ignore-group-logic.md`
- **Task-04**: Implement `gitmap cache create` and `cache search` Split-DB initialization.
  - Subtask: `.ai-memory/plans/subtasks/197-pas-fix/04-cache-split-db.md`
- **Task-05**: Implement `gitmap cpar` (commit-push-all-repos) and `gitmap fia` commands.
  - Subtask: `.ai-memory/plans/subtasks/197-pas-fix/05-cpar-fia-commands.md`

## Next Steps
Subagents A=2, H=2 will be spawned to tackle Subtask 01 and 02 in parallel.
