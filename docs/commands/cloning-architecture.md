# GitMap Cloning Architecture Specification

> **Component:** Multi-Engine Git Cloning Subsystem  
> **Author:** MD ALIM UL KARIM (Chief Software Engineer)  
> **Sponsor:** RISEUP ASIA LLC  
> **Target Package Scope:** `cli/cloner/`, `cli/clonefrom/`, `cli/clonenow/`, `cli/clonepick/`, `cli/clonenext/`, `cli/cloneconcurrency/`  
> **Canonical Version:** `v6.497.0`  
> **Specification References:** [02-spec/21-app/](../../02-spec/21-app/), [02-spec/02-coding-guidelines/](../../02-spec/02-coding-guidelines/)  

---

## 1. Executive Summary & Architectural Motivation

GitMap provides five specialized cloning engines orchestrated through a centralized, cycle-free concurrency coordinator. Rather than collapsing cloning into an unmaintainable monolithic God-object, GitMap strictly segregates responsibilities across modular packages:

1. **`cli/cloner`**: The scan-pipeline cloner operating on in-memory and database-backed `model.ScanRecord` collections.
2. **`cli/clonefrom`**: The plan-driven cloner consuming user-authored arbitrary JSON/CSV files with customized branches, depths, and destinations.
3. **`cli/clonenow`**: The round-trip artifact cloner reading `gitmap scan` output files, faithfully preserving relative folder hierarchies and offering seamless HTTPS/SSH transport toggles.
4. **`cli/clonepick`**: The sparse-checkout and partial clone engine providing interactive path selection, SQLite selection persistence, and `--replay` capabilities.
5. **`cli/clonenext`**: The version-incrementing engine (`-vN`, `v++`, `v+1`) resolving remote repositories, validating via the GitHub API, and processing directory batches.
6. **`cli/cloneconcurrency`**: The centralized, zero-dependency leaf package providing bounded worker-pool allocation, adaptive CPU scaling, and SSH priority throttling.

---

## 2. Formal Defense Against Monolithic Single-File Collapse

A frequent anti-pattern in developer tooling is attempting to combine distinct operational models into a single "universal" cloner. GitMap explicitly rejects this collapse for concrete architectural reasons:

- **Divergent Schema Contracts:** `cli/cloner` consumes internal scanning records (`model.ScanRecord`), `cli/clonefrom` accepts arbitrary user manifests with optional columns (`clonefrom.Plan`), and `cli/clonenow` requires strict scan-artifact parity (`clonenow.Plan`). Forcing them together creates leaky abstractions, bloated union structs, and runtime type assertions.
- **Dependency Cycle Elimination:** `cli/cloner` depends on `model` and database layers. `cli/clonefrom` and `cli/clonenow` must remain completely decoupled from database models so their dry-run renderers and parsers can be tested hermetically in environments without Git or SQLite.
- **Orthogonal Execution Semantics:** `cli/clonepick` manages git sparse-checkout trees and interactive terminal UI buffers (`bubbletea` / raw ANSI). Mixing sparse-checkout plumbing with multi-repo batch parallel workers introduces severe accidental complexity.
- **Independent Test Surfaces:** Each engine possesses dedicated golden file tests, mock git executors, and malformed input test suites that run independently in parallel during CI/CD.

```
                              ┌─────────────────────────┐
                              │     User CLI Input      │
                              └────────────┬────────────┘
                                           │
         ┌───────────────────┬─────────────┼─────────────┬───────────────────┐
         │                   │             │             │                   │
         ▼                   ▼             ▼             ▼                   ▼
   gitmap clone       gitmap clone-from  gitmap clone-now gitmap clone-pick gitmap clone-next
         │                   │             │             │                   │
         ▼                   ▼             ▼             ▼                   ▼
   ┌───────────┐       ┌───────────┐ ┌───────────┐ ┌───────────┐       ┌───────────┐
   │ cli/cloner│       │clonefrom  │ │ clonenow  │ │ clonepick │       │ clonenext │
   └─────┬─────┘       └─────┬─────┘ └─────┬─────┘ └─────┬─────┘       └─────┬─────┘
         │                   │             │             │                   │
         └───────────────────┴──────┬──────┴─────────────┴───────────────────┘
                                    │
                                    ▼
                      ┌───────────────────────────┐
                      │    cli/cloneconcurrency   │
                      │  Adaptive Worker Resolver │
                      └─────────────┬─────────────┘
                                    │
                                    ▼
                      ┌───────────────────────────┐
                      │    OS / Subprocess Pool   │
                      │     (git clone workers)   │
                      └───────────────────────────┘
```

---

## 3. Engine Specifications & Data Models

### 3.1 `cli/cloner` — Scan-Driven Pipeline Cloner

- **CLI Invocations:** `gitmap clone <file>`, `gitmap clone-all`
- **Data Model:** `model.ScanRecord`
  ```go
  type ScanRecord struct {
      RepoName     string
      RelativePath string
      HTTPSUrl     string
      SSHUrl       string
      Branch       string
      BranchSource string
      RemoteName   string
      HeadCommit   string
  }
  ```
- **Primary Entrypoints:**
  - `CloneFromFile(sourcePath, targetDir string, isSafePull bool) (model.CloneSummary, error)`
  - `CloneFromFileWithOptions(sourcePath, targetDir string, opts CloneOptions) (model.CloneSummary, error)`
  - `CloneAll(records []model.ScanRecord, targetDir string, opts CloneOptions) model.CloneSummary`
- **Key Capabilities:**
  - Direct integration with `gitmap scan` database records.
  - Safe pull verification (`isSafePull`) to avoid clobbering uncommitted work.
  - LFS failure detection and automatic retry (`lfs_retry.go`).
  - Terminal batch progress reporting with elapsed time and per-repo status indicators.

---

### 3.2 `cli/clonefrom` — Plan-Driven Manifest Cloner

- **CLI Invocations:** `gitmap clone-from <manifest.json|csv> [--execute] [--workers N]`
- **Data Model:** `clonefrom.Plan` and `clonefrom.Row`
  ```go
  type Plan struct {
      Source string
      Format string // "json" | "csv"
      Rows   []Row
  }

  type Row struct {
      URL      string // Validated Git repository source
      Dest     string // Target directory relative to cwd
      Branch   string // Optional branch name override
      Depth    int    // Optional shallow clone depth
      Checkout string // Mode: "" | "auto" | "skip" | "force"
  }
  ```
- **Primary Entrypoints:**
  - `ParseFile(path string) (Plan, error)`
  - `Render(plan Plan, out io.Writer)` (Hermetic dry-run preview)
  - `Execute(plan Plan, cwd string, workers int, progress io.Writer) ([]Result, error)`
- **Key Capabilities:**
  - Permissive CSV parsing supporting UTF-8 BOM, CRLF/bare-CR normalization, and dynamic column headers (`repo`, `path`, `url`).
  - Safe-by-default dry-run preview: requires explicit `--execute` flag to perform actions on disk.
  - Pre-clone hooks (`BeforeRowHook`) enabling custom setup before git execution.

---

### 3.3 `cli/clonenow` — Scan Artifact Round-Trip Cloner

- **CLI Invocations:** `gitmap clone-now <scan-artifact.json|csv|txt> [--mode ssh|https] [--on-exists skip|update|force]`
- **Data Model:** `clonenow.Plan` and `clonenow.Row`
  ```go
  type Plan struct {
      Source     string
      Format     string // "json" | "csv" | "text"
      Mode       string // "https" | "ssh"
      OnExists   string // "skip" | "update" | "force"
      Rows       []Row
      CoerceURL  func(string) string
      PersistURL func(string)
  }

  type Row struct {
      RepoName     string
      HTTPSUrl     string
      SSHUrl       string
      Branch       string
      RelativePath string // Verbatim destination path from original scan
  }
  ```
- **Primary Entrypoints:**
  - `ParseFile(path, mode, onExists string) (Plan, error)`
  - `Render(plan Plan, out io.Writer)`
  - `Execute(plan Plan, cwd string, workers int, progress io.Writer) ([]Result, error)`
- **Key Capabilities:**
  - Re-establishes exact multi-tier directory layouts across developer machines and build agents.
  - Runtime transport switching between HTTPS and SSH keys without manual file editing.
  - Idempotent existing-repo handling via `OnExists` policies (`skip`, `update`, `force`).

---

### 3.4 `cli/clonepick` — Sparse-Checkout & Interactive Partial Cloner

- **CLI Invocations:** `gitmap clone-pick <repo-url> <path...>`, `gitmap clone-pick --replay <id|name>`, `gitmap clone-pick <repo-url> --ask`
- **Data Model:** `clonepick.Plan`
  ```go
  type Plan struct {
      Name            string
      RepoCanonicalId string
      RepoUrl         string
      Mode            string
      Branch          string
      Depth           int
      Cone            bool
      KeepGit         bool
      DestDir         string
      Paths           []string
      UsedAsk         bool
      DryRun          bool
      Quiet           bool
      Force           bool
      PreClonedSrc    string
  }
  ```
- **Primary Entrypoints:**
  - `ParseArgs(args []string, flags Flags) (Plan, error)`
  - `Execute(plan Plan, cwd string) error`
  - `PickerTUI(repoUrl string, branch string) (selectedPaths []string, err error)`
  - `LoadFromDB(db *sql.DB, key string) (Plan, error)`
- **Key Capabilities:**
  - Automatic cone mode detection: activates efficient cone sparse-checkout for folders, falling back to full pattern matching for files/globs.
  - Two-stage clone promotion: `--ask` uses a metadata-only clone (`--filter=blob:none --no-checkout`) to populate the interactive picker, promoting it into the destination upon confirmation.
  - SQLite persistence in `CloneInteractiveSelection` table to enable deterministic `--replay`.

---

### 3.5 `cli/clonenext` — Version Bumping & Remote Batch Cloner

- **CLI Invocations:** `gitmap clone-next <version-spec>`, `gitmap cn <version-spec>`, `gitmap cn --batch-csv <file.csv>`
- **Data Model:** `clonenext.ParsedRepo`
  ```go
  type ParsedRepo struct {
      BaseName       string
      CurrentVersion int
      HasVersion     bool
  }
  ```
- **Primary Entrypoints:**
  - `ParseRepoName(name string) ParsedRepo`
  - `ResolveTarget(parsed ParsedRepo, arg string) (int, error)`
  - `TargetRepoName(baseName string, version int) string`
  - `ReplaceRepoInURL(remoteURL, currentRepo, targetRepo string) string`
  - `LoadBatchFromCSV(path string) ([]string, error)`
  - `WalkBatchFromDir(dir string) ([]string, error)`
- **Key Capabilities:**
  - Version specifiers supported: `v++`, `v+1`, and explicit `vN` (e.g. `v15`).
  - Remote repository existence validation via GitHub API before disk allocation.
  - Batch ingestion for multi-repository version synchronizations.

---

## 4. Centralized Concurrency: `cli/cloneconcurrency`

Concurrency across all five cloning engines is governed exclusively by `cli/cloneconcurrency`. By isolating concurrency logic in a dedicated leaf package, circular dependencies are eliminated.

### 4.1 Concurrency Resolution API

```go
// Resolve translates user CLI flags into a bounded worker count.
func Resolve(n int) (int, bool)

// ResolveWorkerHands resolves two-dimensional worker and hand allocations.
func ResolveWorkerHands(workerCount, handCount int, isWWOH bool, isSSH bool) (int, int)

// ResolveAdaptivePullConcurrency calculates dynamic worker count based on system load.
func ResolveAdaptivePullConcurrency(profile CPUProfile, userSpecified int) int
```

### 4.2 Resolution Rules & Invariants

1. **Negative Guard:** When `n < 0`, `Resolve` returns `0, false`. Callers treat this as an immediate validation error (`exit 1`).
2. **Auto-Scaling Default:** When `n == 0` (flag omitted), concurrency defaults to `max(1, runtime.NumCPU())`.
3. **Explicit Capping:** When `n > 0`, the specified worker count is returned verbatim without silent down-clamping, ensuring user configurations are respected.
4. **SSH Session Throttling:** When `IsSSHSession()` is active, concurrency is throttled to prevent resource exhaustion and network buffer bloat on remote nodes.
5. **Adaptive CPU Profiles:**
   - `CPUProfileAuto`: Allocates 1 to 6 workers dynamically based on CPU core brackets.
   - `CPUProfileHighPerf`: Allocates up to 12 workers for dedicated high-throughput systems.
   - `CPUProfileLowCPU`: Throttles workers to `cores / 4` (max 2) for resource-constrained environments.

---

## 5. Architectural Comparison Matrix

| Property | `cli/cloner` | `cli/clonefrom` | `cli/clonenow` | `cli/clonepick` | `cli/clonenext` |
|---|---|---|---|---|---|
| **Input Source** | `model.ScanRecord` | User manifest file | `gitmap scan` output | Remote Git URL | Local repo or CSV |
| **Manifest Format** | SQLite / memory | JSON / CSV | JSON / CSV / Text | CLI arguments | CLI / CSV / directory |
| **Path Policy** | Recorded relative | User-specified dest | Preserved relative | Cwd or `--dest` | Base + `-vN` suffix |
| **Execution Mode** | Direct / Safe-pull | Dry-run by default | Dry-run by default | Single-repo sparse | Direct / Batch |
| **Transport Toggle**| In-memory switch | Inherits URL | `Mode` (ssh/https) | `Mode` (ssh/https) | URL string rewrite |
| **Concurrency** | Bounded pool | Worker channel | Worker channel | Sequential / Pre-clone | Bounded batch pool |
| **State Persistence**| `repodb` SQLite | None (stateless) | None (stateless) | `CloneInteractiveSelection` | Local state tracking |

---

## 6. Verification & Quality Gates

All five engines must satisfy these non-negotiable verification gates:

- **Hermetic Parsing:** Parsers must execute without spawned subprocesses or external network calls.
- **Dry-Run Parity:** Dry-run renderers must compute identical target paths to the execution engines.
- **Clean Whitespace & LF:** All documentation and Go source files must adhere strictly to LF line endings and standard formatting (`gofmt`).
- **Relative Path Adherence:** Zero absolute file system paths are permitted in configuration or documentation files.
