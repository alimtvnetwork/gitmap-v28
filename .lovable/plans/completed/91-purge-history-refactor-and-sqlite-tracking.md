# 91-purge-history-refactor-and-sqlite-tracking

## Goal Description
Review, refactor, and harden the newly introduced purge CLI command and history tracking mechanism:
1. Normalize database primary keys and struct identifiers from all-caps ID to PascalCase Id (PurgeHistoryLogId canonical, Id alias) and camelCase id parameters to strictly adhere to repository naming guidelines and prevent linter rejections.
2. Adopt the centralized database pattern (ExecWrapper, QueryRowWrapper, QueryResult[T]) from gitmap/store/wrapper.go in gitmap/store/purge_history.go so all query executions have structured boolean states (isSuccess, isFailure) and standard stderr diagnostics.
3. Eliminate all 20+ instances of swallowed errors (_ ignores, unchecked os.MkdirAll, unchecked copyPurgeFile, unchecked sendToRecycleBin, unchecked git commands, unchecked file closes, and unhandled db.MarkPurgeHistoryRestored).
4. Decompose gitmap/cmd/purge.go (currently 290 lines with doPurge at 134 lines) into modular files <= 200 lines each with functions <= 15 lines (purge.go, purge_engine.go, purge_restore.go, purge_recycle_windows.go, purge_recycle_other.go).
5. Ensure cross-platform compatibility for temp directory backups (os.TempDir() -> filepath.Join(os.TempDir(), fmt.Sprintf("gitmap_purge_%d", timestamp))), path normalization with filepath.ToSlash() for git filter-repo, and build-tag isolated Recycle Bin calls.
6. Verify and harden the Python companion script 03-ai-scripts/30-purge-history.py against line budget (>200 lines) and non-Windows ctypes crashes.

---

## Custom Constraints & Rules (Task-Specific)
1. Strict ID Casing (Rule 1): All-caps ID is completely prohibited in Go structs, methods, parameters, and SQLite queries. Must use PurgeHistoryLogId as canonical PK, Id as struct alias, and id as parameter.
2. Zero Swallowed Errors (Rule 2): Every single returned error (file ops, git commands, db calls, IO copies, and JSON marshaling) must be checked immediately and wrapped via apperror.Wrap or returned up the stack.
3. Store Wrapper Pattern Reuse (Rule 3): All SQL queries in gitmap/store/purge_history.go must be executed via store.ExecWrapper and store.QueryRowWrapper.
4. Function & File Sizing Caps (Rule 4): Every source file must remain <= 200 lines (target <= 100 lines). Every function must remain <= 15 lines (target 8-15 lines).
5. Cross-Platform Recycle Isolation (Rule 5): Windows shell32.dll SHFileOperationW must reside in purge_recycle_windows.go guarded by //go:build windows, with a portable fallback in purge_recycle_other.go (//go:build !windows).

---

## Subtasks Decomposition
- 01-store-purge-history-id-and-wrapper.md: Refactor gitmap/store/purge_history.go to use PurgeHistoryLogId, Id, ExecWrapper, QueryRowWrapper, IsRestored, and wire into store.go:Migrate().
- 02-cmd-purge-split-and-error-handling.md: Refactor gitmap/cmd/purge.go and extract gitmap/cmd/purge_engine.go with zero swallowed errors and functions <= 15 lines.
- 03-purge-restore-and-cross-platform-recycle.md: Extract gitmap/cmd/purge_restore.go, gitmap/cmd/purge_recycle_windows.go, and gitmap/cmd/purge_recycle_other.go.
- 04-python-purge-history-hardening.md: Refactor 03-ai-scripts/30-purge-history.py to stay under 200 lines with cross-platform recycle and argument list subprocess execution.
- 05-unit-tests-and-ci-verification.md: Add unit and integration tests for purge and purge_history store, run go vet, tests, and CI local runner.
