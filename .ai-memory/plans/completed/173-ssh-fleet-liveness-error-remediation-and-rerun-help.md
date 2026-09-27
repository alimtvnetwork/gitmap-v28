# Consolidated Completed Plan: 173 — SSH Fleet Liveness, Error Remediation, Resource-Aware Concurrency, and Rerun Help

## Header & Provenance
- **Canonical Spec Reference:** [02-spec/21-app/173-ssh-fleet-liveness-error-remediation-and-rerun-help.md](../../../02-spec/21-app/173-ssh-fleet-liveness-error-remediation-and-rerun-help.md)
- **Canonical RCA Reference:** [02-spec/22-app-issues/49-ssh-fleet-parallel-pull-machine-hangs-and-concurrency-multiplication-rca.md](../../../02-spec/22-app-issues/49-ssh-fleet-parallel-pull-machine-hangs-and-concurrency-multiplication-rca.md)
- **Status:** `completed`
- **Total Execution Steps / Loops:** 1 continuous dual-subagent iteration (A=2, H=2, 0 blockers)

---

## Consolidated Subtasks & Implemented Deliverables

### Subtask 01: Resource-Aware Concurrency Resolver & Half-Priority SSH Delegation
- **Traceability ID:** Task-01, Task-02
- **Target Files:**
  - `cli/cloneconcurrency/ssh_priority.go`
  - `cli/cloneconcurrency/resolve.go`
  - `cli/cmdssh/ssh_pull_fleet.go`
  - `cli/cmdpull/pullparallel.go`
  - `cli/cmdpull/pull_efficient.go`
- **Accomplished Changes:**
  - Implemented `IsSSHSession() bool` to detect SSH connections (`SSH_CLIENT`, `SSH_CONNECTION`, `GITMAP_SSH_DELEGATED`).
  - Implemented `ResolveWithPriority(n int, isSSHRequest bool) (int, bool)` supporting half priority (`NumCPU() / 2`, clamped to `[1, 2]`) when `isSSHRequest || IsSSHSession()`, falling back to standard CPU resolution otherwise.
  - In `cli/cmdpull/pullparallel.go`: updated `resolveParallelLimit` with `clampSSHLimit(limit)` capping workers to 2 under SSH sessions.
  - In `cli/cmdssh/ssh_pull_fleet.go`:
    - Dispatched `executeLocalVMPull` with `--parallel 2` to prevent host CPU starvation during fleet orchestration.
    - Dispatched `runRemotePullJSON` with `gitmap pa --json --parallel 2`.
    - Handled stale pending task self-healing via `recoverRemotePendingTask` and `extractPendingTaskID`.

### Subtask 02: `gitmap rerun help` Routing, Substring Guard & Embedded Help
- **Traceability ID:** Task-03, Task-04
- **Target Files:**
  - `cli/cmdagy/agy_rerun_help.go`
  - `cli/cmdagy/agy_rerun.go`
  - `cli/cmdagy/agy_rerun_project_resolve.go`
  - `cli/helptext/rerun.md`
- **Accomplished Changes:**
  - Implemented `renderAgyRerunHelp()` in `cli/cmdagy/agy_rerun_help.go` with formatted terminal guide for rerun synopsis, commands, flags, and examples.
  - Defined `agyRerunHelpCmd` (`Aliases: []string{"man", "info"}`) and registered to `agyRerunCmd.AddCommand(agyRerunHelpCmd)`.
  - Implemented `isRerunHelpRequested(args []string) bool` checking `help`, `-h`, `--help`, `-help` (case-insensitive) across all argument positions.
  - Added early help interception in `RunRerunTopLevelCLI` and `runAgyRerun` to render help and return cleanly without executing prompt replay.
  - In `cli/cmdagy/agy_rerun_project_resolve.go`: guarded `findProjectByFlexibleTarget` with `isHelpKeyword(target)` to block `"help"`, `"info"`, and `"man"` from ever fuzzy-matching project names (e.g. `strhelper`).
  - Authored comprehensive embedded markdown help in `cli/helptext/rerun.md`.
