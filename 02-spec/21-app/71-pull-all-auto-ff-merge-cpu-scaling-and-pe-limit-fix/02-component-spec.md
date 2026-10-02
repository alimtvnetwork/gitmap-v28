# Component Specification: CLI Flag Disambiguation, Auto FF-Merge, Batch Fix & CPU Concurrency

## 1. Component Boundaries

| Component | Primary Files | Responsibilities |
|---|---|---|
| **Pipeline Flags & Logs** | `cli/cmdpipeline/pipeline_flags.go`<br>`cli/cmdpipeline/pipeline_logs.go` | Parse `-l`, `-limit`, `--limit`, `-n`, `--lines`; filter out consumed flag values; apply line limits to error logs & failure runs. |
| **Repository Identity & Mapper** | `cli/mapper/mapper.go`<br>`cli/store/repo_sanitize.go`<br>`cli/cmdpull/pull_target_resolver.go` | Fallback repo name to directory basename when remote is empty; sanitize existing DB entries with slug `"unknown"`. |
| **Pull Cloner & FF-Merge** | `cli/cloner/safe_pull.go`<br>`cli/cmdpull/pull_worker.go` | Broaden divergence detection tokens; execute `--no-rebase` merge fallback; skip repositories lacking remotes without 4x retry loops. |
| **Batch Fix & Interactive Prompt** | `cli/cmdpull/pull_db_sync.go`<br>`cli/cmdpull/pull_remediation.go`<br>`cli/cmdpull/pull.go`<br>`cli/cmd/rootutility.go` | Record failed pull items to DB; display interactive resolution prompt; provide `gitmap pull-fix` (`pf`). |
| **CPU Auto-Scaling** | `cli/cloneconcurrency/pull_concurrency.go`<br>`cli/cmdpull/pull_efficient.go`<br>`cli/helptext/pull-all.md` | Calculate adaptive worker counts; provide `--high-perf` and `--low-cpu` presets; document flags in help. |

---

## 2. API & Method Signatures

### 2.1 Pipeline Flags (`cli/cmdpipeline/pipeline_flags.go`)
```go
type PipelineErrorFlags struct {
    ...
    HasLimit             bool
    Limit                int
}

func parseLimitFlag(args []string) (int, bool)
func isValueFlag(flag string) bool
func isConsumedFlagValue(val string, flags *PipelineErrorFlags) bool
func hasSuppressOutputArg(args []string) bool
```

### 2.2 Mapper & Identity (`cli/mapper/mapper.go`)
```go
func deriveRepoNameFromPath(absPath, relPath string) string
func buildOneRecord(repo scanner.RepoInfo, opts BuildOptions) model.ScanRecord
```

### 2.3 Pull Cloner & Safe Pull (`cli/cloner/safe_pull.go`)
```go
func isDivergedOutput(output string) bool
func attemptAutoMergePull(dir, branch string, progress SafePullProgressFunc) (string, error)
```

### 2.4 Batch Remediation (`cli/cmdpull/pull_remediation.go`)
```go
func PromptInteractiveBatchFix(failedRepos []PullFailureItem) error
func RunBatchPullFix(args []string) error
```

### 2.5 CPU Concurrency (`cli/cloneconcurrency/pull_concurrency.go`)
```go
type ConcurrencyPreset string

const (
    PresetAutoScale    ConcurrencyPreset = "auto"
    PresetHighPerf     ConcurrencyPreset = "high-perf"
    PresetLowCPU       ConcurrencyPreset = "low-cpu"
)

func ResolveAdaptiveConcurrency(preset ConcurrencyPreset, userLimit int) int
```

---

## 3. Acceptance Criteria & Test Scenarios

1. **AC-1 (Pipeline Error Flag Parsing)**:
   - `gitmap pe -l 10` parses `Limit = 10` and does NOT search for a repository named `"10"`.
   - `gitmap pe -limit 10` and `gitmap pe --limit 10` parse `Limit = 10`.
   - `gitmap pe -n 10` parses `Limit = 10` and does NOT suppress output log.
   - `gitmap pe -n` without trailing number sets `HasSuppressOutputLog = true`.
2. **AC-2 (Local Repository Identity & Skip)**:
   - Repositories without remotes (e.g. `d:\work\repo-cache`) obtain `RepoName = "repo-cache"` and `Slug = "repo-cache"` instead of `"unknown"`.
   - `gitmap pa` checks if a repo has remotes before attempting pull; if local-only, records as `skipped` without 4x retry delay.
3. **AC-3 (Auto Fast-Forward Merge Fallback)**:
   - When pulling a divergent branch where git requires reconciling divergent branches, GitMap detects divergence and executes `git pull --progress --no-rebase --no-edit --autostash`.
   - If clean merge succeeds, repository status is marked clean.
   - If conflicts exist, runs `git merge --abort` and marks repository as failed due to conflicts.
4. **AC-4 (Batch Remediation & Interactive Prompt)**:
   - When `gitmap pa` finishes with failed repositories in an interactive terminal, displays interactive prompt with options `[1/a] All at once`, `[2/s] One by one`, `[q/n] Skip`.
   - Running `gitmap pull-fix` or `gitmap pf` directly triggers batch remediation of all failed repositories from the previous pull.
5. **AC-5 (CPU Auto-Scaling & Presets)**:
   - Adaptive concurrency scales based on core count (capped appropriately).
   - `--high-perf` uses up to 12 workers with zero delay.
   - `--low-cpu` uses conservative concurrency (1-2 workers) with throttling pauses.
   - `gitmap pa --help` documents all concurrency options.
