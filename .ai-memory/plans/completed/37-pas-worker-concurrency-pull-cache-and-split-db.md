# Milestone 37: PAS Worker Concurrency, Pull Cache, and Split-DB Architecture

- **Slug:** `pas-worker-concurrency-pull-cache-and-split-db`
- **Milestone Index:** `37`
- **Status:** `COMPLETED`
- **Source Plans Merged:** Plans 57, 58, 60, 61, 62, 64, 66, 193, 195, 197, 198, 199, 201
- **Folded Subtask Folders:**
  - `181-gitmap-ignore-and-cache-engine` (7 files)
  - `201-pas-formula-ignore-suite-cpar-and-cache` (4 files)
  - `61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry` (5 files)
  - `62-devtools-cache-discovery-tree-and-split-db` (4 files)
- **Target Subsystems:** `cli/cmdpas/`, `cli/store/`, `cli/cpar/`, `cli/cmddevtools/`, `repodb/`

---

## 1. Domain Context & Architectural Problem

As repository counts scaled across developer workstations, GitMap's mass pull operations (`gitmap pa`, `gitmap pat`) and status audits encountered severe bottlenecks:
1. **PAS Formula Inaccuracies & Deadlocks:** The Pull-All Status (PAS) formula exhibited race conditions and SQLite database write deadlocks when high concurrency worker pools attempted simultaneous updates to a single centralized SQLite database.
2. **Missing Ignore Hierarchy:** Mass operations lacked directory-level ignore awareness, leading to recursive scanning of bloated folders (`node_modules`, `vendor`, `.venv`, `.cargo`) and repeated pull failures on read-only submodules or third-party checkouts.
3. **Cache Invalidation & DevTools Overhead:** Developer caches (npm stores, Go build caches, cargo artifacts, temporary files) grew unbounded without dynamic discovery or persistent size tracking.
4. **Pull-Error Isolation Deficiencies:** Network interruptions and authentication failures during pull-all routines polluted general console logs and lacked structured error isolation for root cause diagnosis.

---

## 2. Synthesized Architectural Outcomes

### 2.1 GitMap PAS Formula & Concurrency Throttling
- **Worker Concurrency Engine:** Re-architected worker pools in `cli/cmdpas/pas_runner.go` with dynamic core scaling (`runtime.NumCPU() * 2`, bounded to safe max workers).
- **Deadlock-Free State Machine:** Introduced atomic in-memory state tracking for active pull jobs, decoupling worker execution from SQLite write contention. Worker completion events are queued via Go channels and written in micro-batches using WAL (Write-Ahead Logging) mode.
- **PAS Formula Calibration:** Standardized status computation metrics:
  - Repositories scanned vs. eligible for pull.
  - Fast-forwardable repositories (`isFastForwardable`).
  - Diverged / merge conflict branches requiring manual intervention (`hasDivergedBranch`).
  - Stale / unreachable remotes (`isRemoteReachable`).

### 2.2 Ignore Engine & CPAR (Commit-Pull-Array-Runner)
- **`.gitmapignore` Specification:** Implemented hierarchical ignore parsing supporting glob patterns, negation rules, and built-in presets (Node, Python, Go, Rust, Java).
- **CPAR Suite:** Implemented `cli/cpar/` facilitating composite commit and pull pipelines across arrays of repositories. CPAR sequences pre-flight checks, dirty-state validation, staged commits, and fetch/rebase cycles with fail-fast options.
- **Repository-Level Cache Store:** Added `repodb/` containing localized SQLite cache files for individual repository metadata, keeping `gitmap.db` lightweight and isolated from transient git tree hashes.

### 2.3 DevTools Dynamic Cache Discovery & Tree View
- **Multi-Category Cache Detection:** Implemented `cli/cmddevtools/cache_discovery.go` to scan and calculate storage consumption across:
  - Golang build/test caches (`GOCACHE`, `GOPATH/pkg/mod`).
  - Node.js stores (`~/.npm`, `~/.pnpm-store`, `yarn/cache`).
  - Python wheel caches and pip temp directories.
  - Rust cargo registry and target artifacts.
  - OS temporary folders and stale gitmap release buffers.
- **Terminal Tree Rendering:** Implemented double-line boxed tree visualization displaying directory size hierarchies, last modified timestamps, and safe purging recommendations (`gitmap devtool clear`).

### 2.4 Pull-Error Isolated Database Logging & Machine Telemetry
- **Split-DB Error Logging (`pull_errors.db`):** Partitioned pull failure events into a dedicated SQLite database table (`PullErrors`), capturing timestamp, repository path, branch name, raw git error, and network round-trip time.
- **Machine Telemetry Collection:** Added workstation telemetry reporting available memory, CPU load during mass pulls, disk I/O wait times, and network throughput to optimize worker count recommendations dynamically.

---

## 3. Go Type Contracts & Architecture

```go
// PASResult aggregates the execution outcome of a mass pull operation
type PASResult struct {
    TotalRepos     int           `json:"totalRepos"`
    UpdatedCount   int           `json:"updatedCount"`
    AlreadyUpToDate int          `json:"alreadyUpToDate"`
    FailedCount    int           `json:"failedCount"`
    SkippedCount   int           `json:"skippedCount"`
    Duration       time.Duration `json:"duration"`
    Errors         []PullError   `json:"errors"`
}

// PullError captures isolated pull failure diagnostics
type PullError struct {
    RepoSlug    string    `json:"repoSlug"`
    Branch      string    `json:"branch"`
    ErrorType   string    `json:"errorType"` // "AUTH", "NETWORK", "CONFLICT", "LOCKED"
    RawMessage  string    `json:"rawMessage"`
    RecordedAt  time.Time `json:"recordedAt"`
}

// CacheCategory encapsulates discovered developer cache metrics
type CacheCategory struct {
    Name        string `json:"name"`
    RootPath    string `json:"rootPath"`
    SizeBytes   int64  `json:"sizeBytes"`
    FileCount   int    `json:"fileCount"`
    IsPurgeSafe bool   `json:"isPurgeSafe"`
}
```

All repository queries use prepared statements with transaction rollbacks via `*appfault.AppError` and structured `Result[T]` containers.

---

## 4. Subtask Verification Ledger

| Folded Subtask Directory | Source Subtask Files | Verified Criteria & Delivered Artifacts |
| :--- | :--- | :--- |
| `181-gitmap-ignore-and-cache-engine` | 7 subtask files | `.gitmapignore` engine, hierarchical glob matcher, disk cache lifecycle, ignore presets. |
| `201-pas-formula-ignore-suite-cpar-and-cache` | 4 subtask files | PAS formula calculation, CPAR runner orchestration, repository cache commands. |
| `61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry` | 5 subtask files | High-concurrency worker pool, pull error split-DB isolation, workstation telemetry collector. |
| `62-devtools-cache-discovery-tree-and-split-db` | 4 subtask files | DevTools dynamic discovery scanner, tree view layout renderer, safe cache purging. |

---

## 5. Quality & Coding Guideline Compliance

- **Positive Booleans Only:** Enforced `isFastForwardable`, `isRemoteReachable`, `isPurgeSafe`, `hasDivergedBranch`.
- **Database Hygiene:** Connection pools enforce reentrant mutex locks (`sync.Mutex`) preventing database busy errors under concurrent read/write operations.
- **Path Portability:** All file path operations utilize normalized slash separators (`filepath.ToSlash`).
- **Memory Efficiency:** Streaming iterators utilized for large repository sets to keep memory footprints below 50MB during mass pull audits.
