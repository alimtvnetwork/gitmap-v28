# RCA-091: ListTaskHistory Undefined on *store.DB in cmd/tasks_list.go

**Date:** 2026-09-30
**Status:** Resolved
**Severity:** Critical (CI/CD Pipeline Failure on Build & Race Detector)
**Affected Workflow:** Release (`v6.429.0`), `race-detector`, `Cross-Platform Build`

---

## 1. Root Cause Analysis
During compilation on CI/CD (GitHub Actions), `go build` and `go test` failed with:
```text
cmd/tasks_list.go:24:26: db.ListTaskHistory undefined (type *store.DB has no field or method ListTaskHistory)
```
In `cli/cmd/tasks_list.go`, `openTasksDB()` returns `(*store.DB, error)`. `runTasksList()` called `db.ListTaskHistory(...)`.
However, `ListTaskHistory` and `InsertTaskHistory` were originally defined with receiver `(db *TasksSplitDB)` in `cli/store/tasks_split_db.go`.
Because `TasksSplitDB` embeds `*store.DB`, methods defined only on `*TasksSplitDB` are NOT accessible on instances of `*store.DB`.

## 2. Why Not Caught Earlier?
The pre-release step was run with `--skip-tests` during the automated release execution, delegating the full compilation and race-detection verification directly to the remote CI/CD matrix. The call site in `tasks_list.go` was introduced when wiring the unified task audit view without running local compilation on `cli/cmd`.

## 3. Remediation
1. Changed receiver in `cli/store/tasks_split_db.go`:
   - `func (db *TasksSplitDB) InsertTaskHistory(...)` -> `func (db *DB) InsertTaskHistory(...)`
   - `func (db *TasksSplitDB) ListTaskHistory(...)` -> `func (db *DB) ListTaskHistory(...)`
   - `func (db *TasksSplitDB) queryTaskHistoryRows(...)` -> `func (db *DB) queryTaskHistoryRows(...)`
2. Since `TasksSplitDB` embeds `*DB`, Go method promotion ensures that both `*store.DB` and `*store.TasksSplitDB` callers (`tasks_list.go`, `task_history_cmd.go`, and tests) have full access to `ListTaskHistory` and `InsertTaskHistory`.
3. Verified locally with `go build ./...` and `go test -count=1 ./cmd/... ./cmdtask/... ./store/...`, passing 100%.

## 4. Prevention
- Always ensure that methods required on the generic database handle returned by `openTasksDB()` are attached to `*store.DB` or access the embedded split DB correctly.
- Add local `go build ./...` verification before tagging releases.
