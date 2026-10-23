# Plan 299 Architecture Specification: Store SSH, History, Database Queries, Templates & Visibility Monadic Results

## User Request (Verbatim)
```text
@[01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md]
```

## 1. Executive Summary & Architectural Intent
This specification advances the repository-wide migration away from raw multi-value Go error tuples (`([]T, error)`, `(T, error)`, `(map[K]V, error)`), standardizing on single-return monadic containers (`result.Wrap[T]`, `result.WrapSlice[T]`, `result.WrapMap[K, V]`) with pointer-attached null safety and structured `*apperror.AppError` returns.

Plan 299 targets secondary database repositories in `cli/store/`, template discovery in `cli/templates/`, and pattern parsing in `cli/visibility/`:
1. `cli/store/ssh_repo.go`, `sshkey.go`, `ssh_known_hosts_repo.go`, `ssh_hist_repo.go`, `tasks_split_db.go`:
   - Modernize `ListHosts`, `ListHostsByRole`, `ListHostsByTarget`, `ListSSHHosts`, `ListSSHKeys`, `SSHKeyNames`, `ListSSHKnownHosts`, `ListSSHHistory`, `ListTaskHistory` to canonical monadic `*Wrap` methods.
2. `cli/store/alias.go`, `bookmark.go`, `release.go`, `workdir_list.go`:
   - Modernize `FindAliasesByRepoID`, `ListAliases`, `ListAliasesWithRepo`, `ListUnaliasedRepos`, `ListBookmarks`, `ListReleases`, `ListWorkDirs` to canonical monadic `*Wrap` methods.
   - Centralize all domain types in `cli/store/types.go`.
3. `cli/templates/` & `cli/visibility/`:
   - Create `cli/templates/types.go` centralizing `EntrySliceWrap` and `EntryWrap`.
   - Add `ListWrap()` in `cli/templates/list.go` delegating to existing `List()`.
   - Create `cli/visibility/types.go` centralizing `PatternSliceWrap` and `PatternWrap`.
   - Add `ParsePatternListWrap()` in `cli/visibility/pattern.go` delegating to existing `ParsePatternList()`.

## 2. Scope Boundaries
- **In-Scope Subsystems & Files:**
  - `cli/store/types.go`
  - `cli/store/ssh_repo.go`, `cli/store/sshkey.go`, `cli/store/ssh_known_hosts_repo.go`, `cli/store/ssh_hist_repo.go`, `cli/store/tasks_split_db.go`
  - `cli/store/alias.go`, `cli/store/bookmark.go`, `cli/store/release.go`, `cli/store/workdir_list.go`
  - `cli/templates/types.go`, `cli/templates/list.go`
  - `cli/visibility/types.go`, `cli/visibility/pattern.go`
- **Excluded / Deferred:**
  - `go build` and `go test` (strictly banned under Rule 1).
  - Version bumping or changelog edits (strictly banned under Rule 10).
  - Raw git commands inside subagents (strictly banned under Rule 7).

## 3. Core Technical Invariants
1. **Single Monadic Return:** Operations return `result.Wrap[T]`, `result.WrapSlice[T]`, or `result.WrapMap[K, V]`.
2. **Backward Compatibility Guarantee:** Existing public methods preserve their original function signatures by delegating directly to the newly introduced monadic `*Wrap` methods (`res := s.MethodWrap(); return res.Items(), res.Err`).
3. **DRY Types Centralization:** Types declared exclusively in `types.go` per package.
4. **AST Compliance:** Zero warnings from `python 03-ai-scripts/50-fastgate.py`.

## 4. Multi-Agent Partitioning & Ownership Matrix
- **Wave 1 (Tree Swarm Concurrency $A = 3$, $H = 2$):**
  - Worker 01 (`self`): Owns `cli/store/ssh_repo.go`, `sshkey.go`, `ssh_known_hosts_repo.go`, `ssh_hist_repo.go`, `tasks_split_db.go`.
  - Worker 02 (`self`): Owns `cli/store/types.go`, `alias.go`, `bookmark.go`, `release.go`, `workdir_list.go`.
  - Worker 03 (`self`): Owns `cli/templates/*`, `cli/visibility/*`.
- **Wave 2 (Verification & Push):**
  - Root Orchestrator: Validates AST via FastGate, verifies secrets gate, updates registry, and executes atomic GitMap commit.
