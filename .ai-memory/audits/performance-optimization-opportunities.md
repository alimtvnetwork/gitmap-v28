# GitMap Performance Optimization Opportunities Audit

> **Status**: Catalog & Diagnostic Audit (For discussion; no premature modifications applied)
> **Author**: Antigravity Orchestrator
> **Target Subsystems**: Pipeline Log Ingestion, SQLite Database, Git Subprocesses, Memory Allocations, and UI Rendering

---

## Executive Summary
This audit catalogs high-impact performance optimization opportunities across GitMap. The identified items focus on reducing subprocess invocation overhead, eliminating heap allocations in log processing loops, optimizing SQLite disk I/O, and parallelizing network downloads.

---

## 1. Pipeline Log Ingestion & Error Filtering

| Target Subsystem | Current Implementation | Proposed Optimization | Expected Impact |
| :--- | :--- | :--- | :--- |
| **Section & Job Log Downloads** (`cmdpipeline/`) | Sequential or single-thread HTTP requests per section | Worker pool bounded to 4–6 concurrent workers with HTTP keep-alive connection reuse | **3x–5x faster** download times for large multi-job CI runs |
| **Line-by-Line Log Filtering** (`cmdpipeline/pipeline_parallel_filter.go`) | `strings.ToLower(strings.TrimSpace(line))` allocated on every line | In-place byte scanning with zero-allocation ASCII case folding (`bytes.EqualFold`) | **60% reduction** in GC pressure and heap allocations during 20k+ line log runs |
| **Buffer Management** (`formatDbErrorRecords`, `formatAggregatedErrorLogs`) | Dynamically growing `strings.Builder` with frequent reallocations | `sync.Pool` of reusable `strings.Builder` with preallocated capacities | Eliminates repeated memory fragmentation |

---

## 2. SQLite Database & Storage Layer

| Target Subsystem | Current Implementation | Proposed Optimization | Expected Impact |
| :--- | :--- | :--- | :--- |
| **Split DB Fallback Queries** (`cli/pipelinedb/`) | Repeated `OpenPipelineSplitDb` and `Close` calls per fallback operation | Cached single-instance DB handle per repository with mutex protection | Eliminates file open/lock syscalls on Windows |
| **Query Prepared Statements** (`QueryDetailedErrorLogsByRunId`) | Ad-hoc SQL query parsing on each invocation | Prepared statements (`db.PrepareContext`) cached in struct fields | **40% faster** query execution across pipeline inspection loops |
| **Transaction Batching** (`writeCachedPipelineJobs`) | Individual inserts inside loops without explicit transactions | Wrap batch inserts in single `BEGIN TRANSACTION ... COMMIT` | **10x–20x faster** batch writes during log ingestion |

---

## 3. Git Subprocess & Process Lifecycle

| Target Subsystem | Current Implementation | Proposed Optimization | Expected Impact |
| :--- | :--- | :--- | :--- |
| **Windows Process Spawning** (`exec.Command("git", ...)`) | Spawning a new `git.exe` process for every repo status check in `gitmap status` | Batching operations or leveraging `git cat-file --batch` / `git status --porcelain=v2` | Windows process spawn is ~15ms per invocation; batching saves seconds on large mono-repos |
| **Clone & Pull Worker Scaling** (`cmdclone/`, `cmdpull/`) | Static worker concurrency | Dynamic CPU core & memory-aware scaling capped at available system thread pool | Prevents memory exhaustion (commit limit errno 1455) under heavy load |

---

## 4. Help Text & Command Routing

| Target Subsystem | Current Implementation | Proposed Optimization | Expected Impact |
| :--- | :--- | :--- | :--- |
| **Help Text Catalog Resolution** (`cli/helptext/catalog.go`) | Filesystem checks and raw file reading on every `--help` invocation | Static embedded map (`embed.FS` with precomputed byte lengths) | Instantaneous O(1) help rendering without disk read |
| **Terminal Table Formatting** (`cli/termtable/`) | Repeated string concatenation and padding during row renders | Preallocated byte buffer and direct ANSI byte writes | Smooth flicker-free table rendering |

---

## 5. Prioritized Discussion Roadmap

1. **Phase A (High Priority / Low Risk)**: Batch SQLite transactions in pipeline caching; reuse `sync.Pool` builders in error log parsing.
2. **Phase B (Medium Priority / Medium Risk)**: Parallelize section log downloads with bounded worker pool; embed help catalog into memory.
3. **Phase C (Future Architectural Evolution)**: Transition heavy batch Git operations to long-lived daemons or batch porcelain streams.
