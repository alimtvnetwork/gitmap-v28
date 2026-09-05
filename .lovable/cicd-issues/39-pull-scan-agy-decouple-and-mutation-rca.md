# CI/CD Issue 39: Antigravity Project Mutation and Command Decoupling in Pull and Scan

- **Job**: Gitmap CLI Execution / Workspaces Integrity / Clean Separation of Concerns
- **Type**: BUGFIX / ARCHITECTURAL INTEGRITY
- **Detected**: 2026-09-05
- **Status**: resolved

## Error

During `gitmap scan` and `gitmap pull`, operations unintendedly mutated Antigravity project workspaces or presented Antigravity-specific commands:

1. `gitmap scan` unconditionally modified `~/.gemini/config/projects/<uuid>.json` for every discovered repository without explicit user authorization.
2. `gitmap pull` reported dirty repository remediation using the terminology `gitmap reconcile`, creating a command name collision with the Antigravity management command `gitmap agy reconcile`.

## 1. Why It Happened

In release v6.115.0 (commit `3787743a`), Antigravity workspace integration features were introduced. During this addition:
- A loop invoking `workspacesync.SyncAntigravity` was mistakenly placed inside `syncClonedReposToVSCodePM` in `gitmap/cmd/clonepmsync.go`. Because `syncRecordsToVSCodePM` in `gitmap/cmd/vscodepmscan.go` delegates to `syncClonedReposToVSCodePM`, every routine `gitmap scan` execution created and mutated JSON configurations in the user's Antigravity project folder (`~/.gemini/config/projects/`).
- Concurrently, dirty repository remediation following `gitmap pull` was labeled as "reconciliation" (`gitmap reconcile`), even though Antigravity already owned `gitmap agy reconcile` for repairing missing Antigravity workspace project links. This caused confusion, where standard git pull operations appeared to require Antigravity-specific reconciliation commands.

## 2. How It Happened

The invocation flow occurred through two decoupled paths:

### Scan Mutation Flow
1. User runs `gitmap scan`.
2. `gitmap/cmd/scan.go` calls `syncRecordsToVSCodePM(records, quiet)`.
3. `gitmap/cmd/vscodepmscan.go` forwards to `syncClonedReposToVSCodePM(pairs, false)`.
4. `gitmap/cmd/clonepmsync.go` previously contained an iteration over all pairs calling `workspacesync.SyncAntigravity(pair.RootPath, pair.Name)`.
5. `workspacesync.SyncAntigravity` inspected `~/.gemini/config/projects/` and wrote a new `<uuid>.json` project file for every scanned repository.

### Pull Command Collision Flow
1. User runs `gitmap pull` on a set of repositories where local changes exist.
2. `gitmap/cmd/pull.go` identifies dirty repositories and calls `PrintRemediationSummary(remItems)` in `gitmap/cmd/remediation_box.go`.
3. The printed instructions stated:
   ```text
   To reconcile later:
     gitmap reconcile               (interactive walkthrough)
     gitmap reconcile <repo> [1|2|3] (target specific repo)
     gitmap reconcile --all stash    (apply stash to all)
   Reconcile these repositories now? [Y/n]:
   ```
4. This gave the erroneous impression that `gitmap pull` was tied to or dependent on `agy` commands, colliding directly with `gitmap agy reconcile` (`gitmap/cmd/agy_reconcile.go`).

## 3. Root Cause

1. **Leaky Abstraction in VS Code PM Sync**: `syncClonedReposToVSCodePM` was intended exclusively for VS Code Project Manager sync (`alefragnani.project-manager`). Injecting `workspacesync.SyncAntigravity` into a function meant for VS Code PM coupled Antigravity project mutation to both clone and scan without user consent.
2. **Command Vocabulary Collision**: The git dirty-working-tree remediation command was ambiguously named `reconcile` alongside the existing `gitmap fix` command, conflicting with `gitmap agy reconcile`. Remediation of dirty git repositories is a git workspace maintenance task (`gitmap fix`), whereas `agy reconcile` is an Antigravity IDE configuration task.
3. **Scan Database Phase Naming**: In `gitmap/cmd/scan.go`, the phase for removing stale SQLite records was titled `scan.reconcileDB` and logged `reconciling missing/stale db entries...`, unnecessarily overloading the term "reconcile".

## 4. Code Fix

1. **Eliminated Antigravity Sync from Clone/Scan (`gitmap/cmd/clonepmsync.go`)**:
   - Removed the `workspacesync.SyncAntigravity` loop from `syncClonedReposToVSCodePM`.
   - Removed the unused `workspacesync` package import.
   - Refactored `syncClonedReposToVSCodePM` and extracted `executeVSCodePMSync` to adhere strictly to the <= 15 lines per function coding guideline.
   - Antigravity syncing is now strictly opt-in via explicit user-invoked commands (`gitmap agy sync` or `gitmap clone-sync`).

2. **Decoupled Pull Remediation from AGY Vocabulary (`gitmap/cmd/remediation_box.go` & `gitmap/cmd/reconcile_prompt.go`)**:
   - Updated `printRemediationCLIHelp` to guide users to `gitmap fix`, `gitmap fix <repo> [1|2|3]`, and `gitmap fix --all stash`.
   - Updated interactive prompt text to `Remediate these repositories now? [Y/n]:`.
   - Created `runInteractiveRemediation` and extracted `promptForRemediation` with clean <= 15 line boundaries.
   - Retained `runInteractiveReconciliation` as a backward-compatible forwarder.

3. **Renamed Scan Database Pruning Phase (`gitmap/cmd/scan.go` & `gitmap/cmd/reconcile_db.go`)**:
   - Renamed benchmark phase `scan.reconcileDB` to `scan.pruneStaleDB`.
   - Updated logs in `gitmap/cmd/reconcile_db.go` from `reconciling missing/stale db entries...` to `pruning missing/stale db entries...`.
   - Added `runPruneStaleDB` as primary entrypoint with modular sub-functions (`pruneStaleRecords`, `deleteStaleEntries`) all <= 15 lines.
