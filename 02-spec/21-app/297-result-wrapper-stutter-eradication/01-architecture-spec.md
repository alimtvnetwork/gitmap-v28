# Architecture Specification 297: Result Wrapper Package Stutter Eradication & Codebase-Wide Monadic Modernization

## 1. Executive Summary

This architecture specification defines the second phase of Result wrapper modernization across GitMap. While Plan 296 established foundation Result containers (`result.Wrap[T]`, `result.WrapSlice[T]`, `result.WrapMap[K, V]`) in `store`, `gitutil`, and `cmdpipeline`, residual package stutter (`result.Result[...]`) and raw generic bracket leakage persist across `cmdagy`, `prdb`, `osuser`, `osclean`, `cmdautomation`, `committransfer`, `lazyregex`, `macro`, `netip`, `osfix`, and utility subsystems.

This specification enforces:
1. **Total Eradication of `result.Result[...]`:** All type alias definitions and function return signatures MUST use canonical `result.Wrap[...]`, `result.WrapSlice[...]`, or `result.WrapMap[...]` aliases.
2. **Centralization in `types.go`:** Packages declaring domain models or result envelopes MUST define them as single reusable named types in each package's `types.go`.
3. **Structured `*apperror.AppError` Returns:** Zero raw standard library `error` returns in leaf business logic; zero double-wrapping (`FailWrap` + `WrapSimple`).
4. **Zero-Build, Zero-Test Mandate:** Static AST verification via FastGate (`03-ai-scripts/50-fastgate.py`), guaranteeing zero broken callers, zero nested `if` statements, positive booleans, and 15-line function caps.

---

## 2. Core Architectural Pillars

### 2.1 The Stutter Anti-Pattern Elimination
Writing `result.Result[T]` repeats the package name (`result.Result`) and leaks raw generic type brackets into function signatures. All instances are modernized to:
- `result.Wrap[T]` for single entities.
- `result.WrapSlice[T]` for slice collections.
- `result.WrapMap[K, V]` for keyed lookups.

### 2.2 Domain Types Centralization (`types.go` Invariant)
Every package declares its Result container aliases in a dedicated `types.go` file:
```go
package mypackage

type (
	ItemWrap      = result.Wrap[Item]
	ItemSliceWrap = result.WrapSlice[Item]
	ItemMapWrap   = result.WrapMap[string, Item]
)
```
Implementation files and callers MUST consume these single named types rather than writing raw generic brackets inline.

---

## 3. Subsystem Boundaries

| Subsystem | Target Files | Key Deliverables |
| :--- | :--- | :--- |
| **`cmdagy`** | `cli/cmdagy/types.go`, `agy_agentapi_ops.go`, `agy_open_paths.go`, `agy_cmd.go`, `agy_ping.go` | Centralize `AgyProcessInfoWrap`, eliminate `result.Result[string]` and `result.Result[AgyProcessInfo]` |
| **`prdb`** | `cli/prdb/pr_split_conn.go`, `cli/prdb/pr_split_ops.go` | Migrate connection and CRUD methods from `result.Result[...]` to `PrSplitDbWrap`, `PullRequestWrap`, `PrBranchRecordWrap` |
| **`osuser` / `osclean`** | `cli/osuser/group_portable.go`, `portable.go`, `cli/osclean/dev_clean.go`, `terminal_clean.go` | Migrate portable imports/exports to `PortableGroupWrap`, `PortableUserWrap`, `DevCleanSummaryWrap` |
| **`cmdautomation`** | `cli/cmdautomation/*_types.go`, `runtime_db.go`, `worker_pool.go` | Eliminate `result.Result` across phase 1–6 monads and runtime database queries |
| **`committransfer`** | `cli/committransfer/pr_engine.go` | Migrate replay and PR execution functions to `PRCommitInfoWrap`, `PrSplitDbWrap`, and `result.StringResult` |
| **Utilities** | `cmddownload`, `cmdupdate`, `cmdzsh`, `lazyregex`, `macro`, `netip`, `osfix`, `movemerge`, `cmdschedule` | Modernize type aliases from `result.Result` to `result.Wrap` |

---

## 4. Verification Protocol

1. **Zero Builds / Zero Tests:** No `go build` or `go test` calls.
2. **FastGate Quality Gate:** `python 03-ai-scripts/50-fastgate.py` must pass with 0 errors across all modified files.
3. **Commit Ceremony:** Atomic commit via `gitmap cpf` with hyphen format.
