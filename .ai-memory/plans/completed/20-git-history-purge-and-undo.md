# Milestone 20: Git History Purge & Undo Engine with Pre-Flight Graph Diff & SplitDB Journal

```text
Milestone ID: 20
Task Slug: 237-git-history-purge-and-undo
Specification: 02-spec/21-app/237-git-history-purge-and-undo/01-architecture-spec.md
Component Spec: 02-spec/21-app/237-git-history-purge-and-undo/02-component-and-cli-spec.md
Audit Ledger: 02-spec/21-app/237-git-history-purge-and-undo/00-master-audit-ledger.md
Date Completed: 2026-10-07
Status: COMPLETED (Verified 100%)
```

---

## 1. Executive Summary & Objective

Accidentally committed confidential secrets, API tokens, sensitive lead sheets, or credentials cannot be expunged via standard `git revert` or `git rm` commits, because Git preserves previous commit trees and binary blobs permanently in the object database.

Milestone 20 delivers an autonomous, high-speed Git history purge and undo subsystem inside GitMap, featuring:
1. **Target Modes**: Deep purge by folder path (`--folder`), relative file path (`--file`), single commit (`--commit`), multiple comma-separated commits, or pattern.
2. **Pre-Flight Visual Analysis**: Traverses commit history, isolates affected commits containing the targeted files, previews before/after Git graph states, displays affected branches and tags, and calculates blast radius (commits, blobs, bytes) before user confirmation.
3. **Interactive Confirmation with Flag Bypass**: Prompts `Proceed with history purge? [y/N]: ` with detailed impact warnings, auto-bypassed when `-y` or `--yes` is specified.
4. **Isolated Temp Blob Backup**: Copies all targeted blobs from every affected commit into `$TEMP/gitmap/history-backup/<repo-slug>/<operation-id>/<commit-sha>/<rel-path>` prior to any rewriting, generating SHA-256 integrity manifests.
5. **SplitDB Transaction Journal**: Persists atomic operation metadata in SQLite SplitDB (`PurgeHistoryLog` / `HistoryPurgeOperation`, `HistoryPurgeCommitMap`, `HistoryPurgeFile`, `HistoryUndoOperation`) recording original SHA, rewritten SHA, parent commit mapping, author metadata, and backup paths.
6. **Release Cleanup**: Detects and purges corresponding release assets or notes referencing the purged targets via GitHub REST API v3 and token discovery.
7. **Post-Purge Guidance & Undo Engine**: Emits the non-warranty restoration warning:
   `"You can undo this if you wanted to. We do not confirm this, but you can try: gitmap history undo <operation-id>"`
   When invoked, `gitmap history undo` reads the SplitDB journal, restores branch pointers from `refs/gitmap-backup/<operation-id>/<branch>`, replays/reconstructs commits from temp backup if necessary, and restores the working tree.

---

## 2. Subtask Execution & Verification Summary

### Subtask 01: SplitDB Schema & Temp Backup Vault Engine
- **Files Modified/Created:** `cli/store/purge_history_models.go`, `cli/store/purge_history.go`, `cli/tempdir/history_backup.go`.
- **Delivered:**
  - Singular PascalCase models: `HistoryPurgeOperation`, `HistoryPurgeCommitMap`, `HistoryPurgeFile`, `HistoryUndoOperation`.
  - Positive boolean columns: `IsDryRun`, `IsVerified`, `HasPushed`, `IsUndone`, `IsSuccess`.
  - SQLite WAL migration routines and CRUD methods: `InsertHistoryPurgeOperation`, `UpdateHistoryPurgeOperationStatus`, `InsertHistoryPurgeCommitMapBatch`, `InsertHistoryPurgeFileBatch`, `GetHistoryPurgeOperationById`, `GetLastHistoryPurgeOperation`, `InsertHistoryUndoOperation`.
  - Temp backup blob vault staging exact file blobs mapped by commit folder: `$TEMP/gitmap/history-backup/<repo-slug>/<operation-id>/<commit-sha>/<rel-path>`.
- **Status:** PASS (Exit 0).

### Subtask 02: Pre-Flight Analyzer & Graph Preview Engine
- **Files Created:** `cli/cmdpurge/preflight.go`, `cli/cmdpurge/graph_preview.go`.
- **Delivered:**
  - `RunPurgePreflight` scanning commits via `git log` / `git rev-list` to find all commits containing target folder, file, or commit hashes.
  - Remote tracking branch detection (`git branch -r --contains`) and working tree cleanliness validation.
  - Blast radius computation (total commits affected, total files, bytes).
  - Unicode/ASCII before-and-after commit topology diff card rendering.
  - Interactive prompt `Proceed with history purge? [y/N]: ` with `-y`/`--yes` bypass.
- **Status:** PASS (Exit 0).

### Subtask 03: Core Git History Purge & GitHub Release Pruning Engine
- **Files Modified/Created:** `cli/cmdpurge/purge_engine.go`, `cli/cmdpurge/release_prune.go`.
- **Delivered:**
  - Low-level Git plumbing tree rewriting (`git ls-tree`, `git mktree`, `git commit-tree`, `git update-ref`).
  - Safety backup ref creation (`refs/gitmap-backup/<operation-id>/<branch>`).
  - Aggressive garbage collection (`git reflog expire --expire=now --all`, `git prune`, `git gc --prune=now`).
  - GitHub release asset pruner discovering token via `ghtoken.DiscoverToken()`, deleting matching assets via GitHub REST API, and redacting mentions in release notes.
  - Post-purge advisory card output with exact undo instructions.
- **Status:** PASS (Exit 0).

### Subtask 04: Undo Engine & CLI Dispatch Wiring
- **Files Modified/Created:** `cli/cmdpurge/purge_undo.go`, `cli/cmd/history_purge_cmd.go`, `cli/cmd/roottooling.go`, `cli/cmd/root_cobra_completion.go`.
- **Delivered:**
  - `RunPurgeUndo` resolving operation ID from SplitDB, checking working tree cleanliness, restoring branch pointer from `refs/gitmap-backup/...`, and applying temp vault blob recovery fallback.
  - Cobra commands: `gitmap history purge` (aliases `hp`, `history clean`) and `gitmap history undo` (aliases `hu`, `history restore`, `undo-history`).
  - Wired into `toolingSystemEntries()` in `cli/cmd/roottooling.go`.
  - Registered dynamic tab-completion and flags in `cli/cmd/root_cobra_completion.go`.
- **Status:** PASS (Exit 0).

---

## 3. Command Usage & Syntax Reference

```bash
# 1. Purge a folder from all Git history (with interactive pre-flight preview & prompt)
gitmap history purge --folder secrets/confidential-leads

# 2. Purge a folder with automatic confirmation (-y)
gitmap history purge --folder secrets/confidential-leads -y

# 3. Purge a specific file from all Git history
gitmap history purge --file api/keys.json -y

# 4. Purge specific commits and prune GitHub release notes/assets
gitmap history purge --commit d4bba582 --release -y

# 5. Dry-run preflight inspection without modifying Git history
gitmap history purge --folder secrets --dry-run

# 6. Undo a history purge operation using SplitDB journal & temp backup vault
gitmap history undo <operation-id>
```
