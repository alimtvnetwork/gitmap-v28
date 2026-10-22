# Component Specification 296: Monadic Result Containers & Types Sovereignty in Gitutil and Cmdpipeline

## Status: APPROVED
- **Spec ID:** 296-02
- **Slug:** `296-result-wrapper-and-apperror-returns`
- **Domain:** Component Design & Type Centralization

---

## 1. Executive Summary & Architectural Scope

In enterprise Go architecture, functions must never return multi-value error tuples such as `(T, error)`, `([]T, error)`, or `(int, error)`. Multi-value error returns violate single-return sovereignty, induce accidental `nil` dereferences, encourage swallowed errors, and bypass structured DB audit logging.

This component specification defines the architectural migration of low-level git operations and pipeline diagnostics to strongly-typed monadic Result containers (`result.Wrap[T]`, `result.WrapSlice[T]`). Specifically, it encompasses:
1. **Scalar and Slice Result Containers:** Migrating git inspection and execution functions to use canonical containers defined in `cli/result/types.go`:
   - `result.StringResult` (`Result[string]`)
   - `result.BoolResult` (`Result[bool]`)
   - `result.IntResult` (`Result[int]`)
   - `result.ByteSliceResult` (`Result[[]byte]`)
   - `result.StringSliceResult` (`ResultSlice[string]`)
2. **Domain Type Centralization in `types.go`:** Eliminating inline generic brackets in implementation signatures by establishing dedicated type aliases:
   - `cli/gitutil/types.go`: `CaseCollisionSliceWrap`, `RemoteBranchSliceWrap`
   - `cli/cmdpipeline/types.go`: `PEFormatProfileSliceWrap`, `PipelineSyncResultWrap`
3. **Affirmative Predicate Inspection:** Replacing legacy direct field access (`res.Value`, `res.Err`), negative checks (`!res.IsFailed()`), and obsolete helpers (`res.HasError()`) with affirmative predicates:
   - `res.IsFailure()`, `res.IsSuccess()`, `res.HasRecord()`, `res.IsEmpty()`
   - `res.Items()`, `res.UnwrapOr(defaultVal)`, `res.AppError()`

---

## 2. Component Boundaries & Target Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│ cli/gitutil/types.go (NEW)                                                   │
│ • CaseCollisionSlice, CaseCollisionSliceWrap                                 │
│ • RemoteBranchSlice, RemoteBranchSliceWrap                                   │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/gitutil/git_cmd_exec.go                                                  │
│ • RunGitInDir, RunGitCombinedInDir       -> result.ByteSliceResult           │
│ • GetRepoConfig, ReadRemoteOriginURL     -> result.StringResult              │
│ • GetRepoCurrentBranch                   -> result.StringResult              │
│ • IsRepoClean                            -> result.BoolResult                │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/gitutil/git_repo_query.go                                                │
│ • GetRepoHeadCommitSha, GetRepoRemoteURL -> result.StringResult              │
│ • GetRepoTrackingRemote                  -> result.StringResult              │
│ • HasUnpushedCommits, IsRepoDetached     -> result.BoolResult                │
│ • GetRepoDirtyFiles                      -> result.StringSliceResult         │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/gitutil/git_repo_commits.go                                              │
│ • GetRepoUnpushedCommitCount             -> result.IntResult                 │
│ • GetRepoMergeBase                       -> result.StringResult              │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/gitutil/gitutil.go                                                       │
│ • RemoteURL, RepoRoot, CurrentBranch     -> result.StringResult              │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/cmdpipeline/types.go (NEW)                                               │
│ • PEFormatProfileSlice, PEFormatProfileSliceWrap                             │
│ • PipelineSyncResultWrap                                                     │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/cmdpipeline/pipeline_history_render.go                                   │
│ • findLocalRepoDirForSlug                -> consumes repo.FindPathBySlug     │
│ • renderDetailedOrLegacyErrors           -> consumes detailRes.Items()       │
│ • renderCompactOrLegacyErrors            -> consumes compactRes.Items()      │
│ • queryCachedFailedRuns                  -> consumes runRes.Items()          │
├──────────────────────────────────────────────────────────────────────────────┤
│ cli/cmdpipeline/pipeline_fallback.go                                         │
│ • TryRetrievePreviousDbRun               -> consumes runRes affirmative      │
│ • queryDbRunErrorText                    -> consumes errRes/compactRes items │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. `cli/gitutil` Component Specification

### 3.1 Type Centralization in `cli/gitutil/types.go`

`cli/gitutil/types.go` centralizes domain collections and eliminates raw generic syntax across git operations.

```go
package gitutil

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// CaseCollisionSlice defines a collection of case collision records.
	CaseCollisionSlice = []CaseCollision

	// CaseCollisionSliceWrap is the canonical single reusable container for CaseCollision collections.
	CaseCollisionSliceWrap = result.Wrap[CaseCollisionSlice]

	// RemoteBranchSlice defines a collection of remote branch metadata.
	RemoteBranchSlice = []RemoteBranchInfo

	// RemoteBranchSliceWrap is the canonical single reusable container for remote branch collections.
	RemoteBranchSliceWrap = result.Wrap[RemoteBranchSlice]
)
```

### 3.2 Low-Level Process Execution in `cli/gitutil/git_cmd_exec.go`

Low-level process runners must return `result.ByteSliceResult` wrapping command standard output or structured `*apperror.AppError`.

#### `RunGitInDir`
```go
// RunGitInDir executes a git command with directory arguments and returns stdout output.
func RunGitInDir(repoPath string, args ...string) result.ByteSliceResult {
	cmd := createGitCmd(repoPath, args...)
	out, err := cmd.Output()
	if err != nil {
		return result.Fail[[]byte](apperror.WrapSimple(err, "RunGitInDir"))
	}

	return result.Ok(out)
}
```

#### `RunGitCombinedInDir`
```go
// RunGitCombinedInDir executes a git command with directory arguments and returns combined stdout and stderr.
func RunGitCombinedInDir(repoPath string, args ...string) result.ByteSliceResult {
	cmd := createGitCmd(repoPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return result.Fail[[]byte](apperror.WrapSimple(err, "RunGitCombinedInDir"))
	}

	return result.Ok(out)
}
```

#### `GetRepoConfig`
```go
// GetRepoConfig retrieves a git configuration value for key in repoPath.
func GetRepoConfig(repoPath, key string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitConfigCmd, constants.GitGetFlag, key)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
```

#### `ReadRemoteOriginURL`
```go
// ReadRemoteOriginURL retrieves the configured remote.origin.url for repoPath.
func ReadRemoteOriginURL(repoPath string) result.StringResult {
	return GetRepoConfig(repoPath, constants.GitRemoteOrigin)
}
```

#### `GetRepoCurrentBranch`
```go
// GetRepoCurrentBranch returns the active branch name of the repository at repoPath.
func GetRepoCurrentBranch(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
```

#### `IsRepoClean`
```go
// IsRepoClean reports whether the working tree at repoPath has no uncommitted changes.
func IsRepoClean(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, "status", "--porcelain")
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}

	isClean := len(strings.TrimSpace(string(res.Data))) == 0

	return result.Ok(isClean)
}
```

### 3.3 Query & Inspection Functions in `cli/gitutil/git_repo_query.go`

#### `GetRepoHeadCommitSha`
```go
// GetRepoHeadCommitSha retrieves the full HEAD commit SHA of the repository at repoPath.
func GetRepoHeadCommitSha(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
```

#### `GetRepoRemoteURL`
```go
// GetRepoRemoteURL retrieves the configured URL for the specified remoteName at repoPath.
func GetRepoRemoteURL(repoPath, remoteName string) result.StringResult {
	key := fmt.Sprintf("remote.%s.url", remoteName)

	return GetRepoConfig(repoPath, key)
}
```

#### `GetRepoTrackingRemote`
```go
// GetRepoTrackingRemote retrieves the upstream tracking reference for the active branch.
func GetRepoTrackingRemote(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, "@{u}")
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
```

#### `HasUnpushedCommits`
```go
// HasUnpushedCommits reports whether the repository has commits not yet pushed to upstream.
func HasUnpushedCommits(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, "cherry", "-v")
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}

	hasCommits := len(strings.TrimSpace(string(res.Data))) > 0

	return result.Ok(hasCommits)
}
```

#### `GetRepoDirtyFiles`
```go
// GetRepoDirtyFiles returns relative paths of uncommitted, modified, or untracked files in repoPath.
func GetRepoDirtyFiles(repoPath string) result.StringSliceResult {
	res := RunGitInDir(repoPath, constants.GitStatus, constants.GitPorcelainFlag)
	if res.IsFailure() {
		return result.FailSlice[string](res.AppError())
	}

	files := parsePorcelainDirtyLines(string(res.Data))

	return result.OkSlice(files)
}
```

#### `IsRepoDetached`
```go
// IsRepoDetached reports whether the repository is in a detached HEAD state.
func IsRepoDetached(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, constants.GitSymbolicRef, constants.GitQuietFlag, constants.GitHEAD)
	if res.IsFailure() {
		return result.Ok(true)
	}

	isDetached := len(strings.TrimSpace(string(res.Data))) == 0

	return result.Ok(isDetached)
}
```

### 3.4 Commit & Branch Queries in `cli/gitutil/git_repo_commits.go`

#### `GetRepoUnpushedCommitCount`
```go
// GetRepoUnpushedCommitCount returns the number of unpushed commits ahead of upstream tracking.
func GetRepoUnpushedCommitCount(repoPath string) result.IntResult {
	res := RunGitInDir(repoPath, "rev-list", "--count", "@{u}..HEAD")
	if res.IsFailure() {
		return result.Fail[int](res.AppError())
	}

	clean := strings.TrimSpace(string(res.Data))
	count, parseErr := strconv.Atoi(clean)
	if parseErr != nil {
		return result.Fail[int](apperror.WrapSimple(parseErr, "GetRepoUnpushedCommitCount"))
	}

	return result.Ok(count)
}
```

#### `GetRepoMergeBase`
```go
// GetRepoMergeBase finds the best common ancestor commit between two references.
func GetRepoMergeBase(repoPath, ref1, ref2 string) result.StringResult {
	res := RunGitInDir(repoPath, "merge-base", ref1, ref2)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
```

### 3.5 High-Level Git Utilities in `cli/gitutil/gitutil.go`

#### `RemoteURL`
```go
// RemoteURL returns the origin remote URL for a repo at the given path.
func RemoteURL(repoPath string) result.StringResult {
	return GetRepoConfig(repoPath, constants.GitRemoteOrigin)
}
```

#### `RepoRoot`
```go
// RepoRoot returns the absolute path of the working-tree root for the
// repo containing `startPath`.
func RepoRoot(startPath string) result.StringResult {
	out, err := runGit(startPath, constants.GitRevParse, "--show-toplevel")
	if err != nil {
		return result.Fail[string](apperror.WrapSimple(err, "RepoRoot"))
	}

	return result.Ok(strings.TrimSpace(out))
}
```

#### `CurrentBranch`
```go
// CurrentBranch returns the current branch name for a repo.
func CurrentBranch(repoPath string) result.StringResult {
	out, err := runGit(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {
		return result.Fail[string](apperror.WrapSimple(err, "CurrentBranch"))
	}

	return result.Ok(strings.TrimSpace(out))
}
```

---

## 4. `cli/cmdpipeline` Component Specification

### 4.1 Type Centralization in `cli/cmdpipeline/types.go`

```go
package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// PEFormatProfileSlice defines a collection of pipeline format profiles.
	PEFormatProfileSlice = []PEFormatProfile

	// PEFormatProfileSliceWrap is the canonical single reusable container for PEFormatProfile collections.
	PEFormatProfileSliceWrap = result.Wrap[PEFormatProfileSlice]

	// PipelineSyncResultWrap is the canonical single reusable container for PipelineSyncResult pointers.
	PipelineSyncResultWrap = result.Wrap[*PipelineSyncResult]
)
```

### 4.2 Caller Refactoring in `cli/cmdpipeline/pipeline_history_render.go`

In `pipeline_history_render.go`, replace obsolete `HasError()` checks with affirmative predicates:

```go
func findLocalRepoDirForSlug(slug string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	repo := store.NewRepoRepository(db.Conn())
	pathRes := repo.FindPathBySlug(slug)
	if pathRes.IsFailure() {
		logRepoLookupFailure(slug, pathRes.AppError())

		return ""
	}

	return pathRes.Data
}
```

In error log rendering functions, consume items via `detailRes.Items()` and `compactRes.Items()`:

```go
func renderDetailedOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	detailRes := db.QueryDetailedErrorLogsByRunId(runId)
	if detailRes.HasRecord() {
		renderCachedErrorsList(detailRes.Items())

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Items())
}

func renderCompactOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	compactRes := db.QueryCompactErrorLogsByRunId(runId)
	if compactRes.HasRecord() {
		renderCachedCompactErrorsList(compactRes.Items())

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Items())
}
```

In `queryCachedFailedRuns`, guard with `IsFailure()` and access items:

```go
func queryCachedFailedRuns(db *pipelinedb.PipelineSplitDb, count int) ([]pipelinedb.PipelineRunRecord, error) {
	limit := resolveMaxSyncLimit(count)
	runRes := db.QueryLastFailedRuns(limit)
	if runRes.IsFailure() {
		return nil, runRes.AppError()
	}

	return runRes.Items(), nil
}
```

### 4.3 Caller Refactoring in `cli/cmdpipeline/pipeline_fallback.go`

In `pipeline_fallback.go`, eliminate negative booleans (`!errRes.IsFailed()`, `!compactRes.IsFailed()`) and use affirmative predicates (`HasRecord()`, `IsFailure()`, `IsEmpty()`, `Items()`):

```go
func TryRetrievePreviousDbRun(repo string) (*pipelinedb.PipelineRunRecord, string, bool) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, "", false
	}
	defer db.Close()

	runRes := db.QueryLastFailedRuns(1)
	if runRes.IsFailure() || runRes.IsEmpty() {
		return nil, "", false
	}

	items := runRes.Items()
	rec := items[0]
	errLogs := queryDbRunErrorText(db, rec.RunId)

	return &rec, errLogs, true
}

func queryDbRunErrorText(db *pipelinedb.PipelineSplitDb, runId uint64) string {
	errRes := db.QueryDetailedErrorLogsByRunId(runId)
	if errRes.HasRecord() {
		return formatDbErrorRecords(errRes.Items())
	}

	compactRes := db.QueryCompactErrors(10)
	if compactRes.HasRecord() {
		return formatDbCompactRecords(compactRes.Items())
	}

	return ""
}
```

---

## 5. Caller Inspection Patterns & Predicate Sovereignty

| Operation | ❌ Anti-Pattern | ✅ Canonical Affirmative Pattern |
| :--- | :--- | :--- |
| **Failure Guard** | `if res.HasError()` / `if res.Err != nil` | `if res.IsFailure()` |
| **Success Branch** | `if !res.IsFailed()` / `if res.Err == nil` | `if res.IsSuccess()` |
| **Slice Extraction** | `for _, item := range res.Data` (raw) | `for _, item := range res.Items()` |
| **Slice Existence** | `if len(res.Data) > 0` | `if res.HasRecord()` |
| **Default Fallback** | `if res.IsFailure() { return "" } return res.Data` | `return res.UnwrapOr("")` |
| **Error Propagation** | `return res.Err` | `return res.AppError()` |

---

## 6. Coding Guidelines & Invariant Conformance

All refactored implementations must strictly satisfy:
1. **Error Sovereignty:** Leaf functions must return `Result[T]` wrapping an `*apperror.AppError`, never raw `error`, and must never call exit handlers or `os.Exit`.
2. **Boolean Naming:** Positive prefixes only (`is...`, `has...`). No negative booleans or inverted checks.
3. **15-Line Function Limit:** Flatten branches with early returns and guard clauses ($\le 15$ lines per function).
4. **Return Styling:** Mandatory blank line before `return` and blank line after closing brace `}`.
5. **No Blind Regex Mass Refactoring:** Surgical updates verified against function signatures and caller semantics.
6. **Zero Tests / Zero Compiler Builds:** Planning phase strictly forbids executing `go test` or `go build`.
