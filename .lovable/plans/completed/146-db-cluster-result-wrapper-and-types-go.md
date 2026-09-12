# Plan 146: Result Wrapper Types, Collections & AppError Returns (DB, Cluster, and CmdPurge Types Centralization)

Trigger Keywords & Aliases: `cg-result-wrapper`, `cg-apperror-returns`, `cg-execute result-wrapper`, `audit result wrapper`, `fix map return error`, `fix slice return error`, `single return object audit`, `enforce apperror returns`, `enforce result map`, `fix multi-value returns`, `is-count-other-than`, `has-record`, `is-defined`, `result-wrapper-null-safety`, `pointer-null-safety`, `types-go-single-type`, `types-go-result-reuse`, `centralize-types-go`

> **Prompt Version:** 2.4.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)
> **Status:** COMPLETED

---

## 1. Executive Summary & Objective

Autonomously scan, discover, plan, refactor, and verify all Go functions returning multi-value error tuples (such as `(map[K]V, error)`, `([]T, error)`, or `(T, error)`), eliminating raw standard library error returns, centralizing all domain payload structs and Result type aliases into `types.go` within each package as single reusable types everywhere rather than scattering inline structs or raw generic Result declarations across implementation files, replacing multi-value returns with strongly-typed result wrappers (`ResultMap[K, V]`, `ResultSlice[T]`, `Result[T]`) and structured `*apperror.AppError` returns, guaranteeing a single return object, pointer-attached null safety (`*Result[T]`, `*ResultSlice[T]`, `*ResultMap[K, V]`) with methods attached to pointer receivers guarding against nil dereferences on line 1, implementing the standardized outer-layer inspection predicates (`IsSuccess`, `IsFailure`, `HasError`, `IsEmptyError`, `IsEmpty`, `HasRecord`, `IsDefined`, `IsCountOtherThan`, `Data`, `Items`, `AppError`, `Fault`, `Get`, `Has`, `Count`), and replacing clumsy compound checks (`err != nil || len(...) != N`, `IsFailure() || Count() != N`) with fluent `res.IsCountOtherThan(N)` across the entire codebase until 100% green without stopping.

Plan 146 focused on the `cli/db`, `cli/cluster`, and `cli/cmdpurge` subsystems across 4 milestones:
1. **DB Types Centralization & Slice Returns Migration (`cli/db/types.go`):** Created `cli/db/types.go` declaring `NodePathAlias`, `NodePathAliasSliceResult`, `ClusterNodeSliceResult`, `ClusterExecResultSliceResult`, `ClusterRunSliceResult`, `SSHConnectionSliceResult`. Refactored `ListClusterNodes`, `scanClusterNodeRows`, `SelectClusterExecResultsByRunId`, `scanClusterExecResultRows`, `ListClusterRuns`, `scanClusterRunRows`, `GetSSHConnections`, `scanSSHConnectionRows`, and `ListPathAliases` to return single `ResultSlice` envelopes.
2. **Cluster & CmdPurge Types Centralization (`cli/cluster/types.go`, `cli/cmdpurge/types.go`):** Created `cli/cluster/types.go` declaring `AliasEntry` and `AliasEntrySliceResult`. Created `cli/cmdpurge/types.go` declaring `TrackedLovableFilesMapResult`. Updated `cli/cluster/pathalias.go` and `cli/cmdpurge/purge_lovable.go`.
3. **Caller Site Modernization (`cli/cmd/cluster_ops.go`, `cli/cmdssh/`):** Modernized callers of `ListClusterNodes`, `SelectClusterExecResultsByRunId`, `ListClusterRuns`, and `GetSSHConnections` to check `res.IsFailure()`, read `res.Data`, and use fluent predicates (`IsCountOtherThan(5)`, `IsEmpty()`).
4. **Auditor & Quality Gates Verification:** Updated `03-ai-scripts/35-result-wrapper-auditor.py` to enforce `cli/db/` in `RESULT_SLICE_ENFORCED_PREFIXES` and `cli/db`, `cli/cluster`, `cli/cmdpurge` in `ENFORCED_TYPES_GO_PACKAGES`. Verified zero violations across all quality linters, locked modified files in test inventory, and pushed atomically via SSH.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Single return type mandate and micro-tasking.
- `spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md`: Types.go and single type definitions.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100–300 lines).
- `spec/03-error-manage/01-index.md`: Universal AppError wrapping and error envelopes.
- `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`: Error handling architecture and Result wrappers.
- `spec/03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/04-result-types.md`: Result[T], ResultSlice[T], and ResultMap[K, V] method specifications and pointer null-safety rules.
- `.agents/skills/cg-result-wrapper/skill.md`: Canonical skill for result wrappers and AppError returns.
- `.lovable/coding-guidelines.md`: Master consolidated coding guidelines.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Identifier | Issue | Target Refactoring | Status |
|:---|:---|:---:|:---|---|---|:---:|
| V-01 | `cli/db/types.go` | 1 | `types.go` missing | Package lacks central types.go | Create `cli/db/types.go` with domain payload models and Result aliases | Completed |
| V-02 | `cli/db/nodepath.go` | 9 | `NodePathAlias` | Struct declared inline | Move to `cli/db/types.go` and return `NodePathAliasSliceResult` | Completed |
| V-03 | `cli/db/clusternode.go` | 75 | `ListClusterNodes` | Returns legacy tuple `([]ClusterNode, *apperror.AppError)` | Return `ClusterNodeSliceResult` | Completed |
| V-04 | `cli/db/clusternode.go` | 86 | `scanClusterNodeRows` | Returns legacy tuple `([]ClusterNode, *apperror.AppError)` | Return `ClusterNodeSliceResult` | Completed |
| V-05 | `cli/db/clusterexecresult.go` | 111 | `SelectClusterExecResultsByRunId` | Returns legacy tuple `([]ClusterExecResult, *apperror.AppError)` | Return `ClusterExecResultSliceResult` | Completed |
| V-06 | `cli/db/clusterexecresult.go` | 126 | `scanClusterExecResultRows` | Returns legacy tuple `([]ClusterExecResult, *apperror.AppError)` | Return `ClusterExecResultSliceResult` | Completed |
| V-07 | `cli/db/clusterrun.go` | 109 | `ListClusterRuns` | Returns legacy tuple `([]ClusterRun, *apperror.AppError)` | Return `ClusterRunSliceResult` | Completed |
| V-08 | `cli/db/clusterrun.go` | 125 | `scanClusterRunRows` | Returns legacy tuple `([]ClusterRun, *apperror.AppError)` | Return `ClusterRunSliceResult` | Completed |
| V-09 | `cli/db/sshconnection.go` | 54 | `GetSSHConnections` | Returns legacy tuple `([]SSHConnection, *apperror.AppError)` | Return `SSHConnectionSliceResult` | Completed |
| V-10 | `cli/db/sshconnection.go` | 65 | `scanSSHConnectionRows` | Returns legacy tuple `([]SSHConnection, *apperror.AppError)` | Return `SSHConnectionSliceResult` | Completed |
| V-11 | `cli/cluster/types.go` | 1 | `types.go` missing | Package lacks central types.go | Create `cli/cluster/types.go` declaring `AliasEntry` and `AliasEntrySliceResult` | Completed |
| V-12 | `cli/cluster/pathalias.go` | 5 | `AliasEntry` | Struct declared inline; returns raw generic | Move to `cli/cluster/types.go` and return `AliasEntrySliceResult` | Completed |
| V-13 | `cli/cmdpurge/types.go` | 1 | `types.go` missing | Package lacks central types.go | Create `cli/cmdpurge/types.go` declaring `TrackedLovableFilesMapResult` | Completed |
| V-14 | `cli/cmdpurge/purge_lovable.go` | 29 | `getTrackedLovableFiles` | Returns raw generic `result.ResultMap[string, bool]` | Return `TrackedLovableFilesMapResult` | Completed |
| V-15 | `cli/cmd/cluster_ops.go` | 46 | `printClusterHistoryList` | Unpacks legacy tuple from `ListClusterRuns` | Modernize caller site to inspect `runsRes.IsFailure()` and `runsRes.Data` | Completed |
| V-16 | `cli/cmd/cluster_ops.go` | 91 | `printClusterRunDetails` | Unpacks legacy tuple from `SelectClusterExecResultsByRunId` | Modernize caller site to inspect `resultsRes.IsFailure()` and `resultsRes.Data` | Completed |
| V-17 | `cli/cmd/cluster_ops.go` | 161 | `runClusterExport` | Unpacks legacy tuple from `ListClusterNodes` | Modernize caller site to inspect `nodesRes.IsFailure()` and `nodesRes.Data` | Completed |
| V-18 | `cli/cmd/cluster_ops.go` | 275 | `getExistingClusterNodesMap` | Unpacks legacy tuple from `ListClusterNodes` | Modernize caller site to `existingRes.Data` and `existingRes.Count()` | Completed |
| V-19 | `cli/cmd/cluster_ops.go` | 396 | `runClusterList` | Unpacks legacy tuple from `ListClusterNodes` | Modernize caller site to inspect `nodesRes.IsFailure()` and `nodesRes.Data` | Completed |
| V-20 | `cli/cmdssh/sshjoin.go` | 84, 230 | `runSSHJoinList`, `runSSHJoinExport` | Unpacks legacy tuple from `GetSSHConnections` | Modernize caller site to inspect `connsRes.IsFailure()` and `connsRes.Data` | Completed |
| V-21 | `cli/cmdssh/sshexec.go` | 61 | `runSSHExec` | Unpacks legacy tuple from `GetSSHConnections` | Modernize caller site to inspect `connsRes.IsFailure()` and `connsRes.Data` | Completed |
| V-22 | `cli/db/clusterrun_test.go` | 118, 134 | `TestClusterRunAndExecResult` | Clumsy length checks and legacy tuple unpacking | Modernize to `IsFailure()`, `IsCountOtherThan(5)`, `IsEmpty()` | Completed |
| V-23 | `03-ai-scripts/35-result-wrapper-auditor.py` | 50, 60 | Enforced scopes | Missing `cli/db`, `cli/cluster`, `cli/cmdpurge` | Expand `RESULT_SLICE_ENFORCED_PREFIXES` and `ENFORCED_TYPES_GO_PACKAGES` | Completed |
