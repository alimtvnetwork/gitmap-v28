# Consistency Report & Verification Matrix: GitMap Heatmap Discovery, URL Split-DB Caching, Segmented Downloader & LLM Training

## Document Metadata
- **Document ID**: `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/99-consistency-report.md`
- **Task Identifier**: Task 314 (`gitmap-heatmap-and-url-handling-instructions`)
- **Status**: `APPROVED / ARCHITECTURAL SPECIFICATION`
- **Version**: `1.0.0`
- **Classification**: Cross-Specification Audit & Verification Protocol
- **Disjoint Owners**: Lead Orchestrator & Multi-Agent Swarm

---

## Cross-Specification Reference Topology

This specification family forms an integrated, non-overlapping architectural contract governing Task 314:

```text
02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/
├── 00-overview.md             (Executive Overview, Verbatim Request, Topology, Traceability Matrix)
├── 01-architecture-spec.md     (Heatmap Multi-Tier Discovery, Alias Joins, Recent Tasks, AI Suggestions)
├── 02-component-spec.md        (URL Split-DB repodb/url_cache.db, 10-Split Segmented Downloader & aria2c)
└── 99-consistency-report.md     (Consistency Audit, Compliance Scorecard, Acceptance Criteria, Test Matrix)
```

### Cross-Document Linkages
1. **00-overview.md -> 01-architecture-spec.md**: Task-02 and Task-03 trace directly to the 3-tier database scanner, alias resolution via `db.ListAliasesWithRepo()`, and multi-repository task aggregation in `cli/heatmap/`.
2. **00-overview.md -> 02-component-spec.md**: Task-04, Task-05, and Task-06 trace to the URL Split-DB engine (`repodb/url_cache.db`), raw unbuffered streaming in `cli/cmdurl/`, and 10-connection parallel Range / `aria2c` downloading in `cli/downloader/`.
3. **00-overview.md -> 99-consistency-report.md**: Validates requirement satisfaction across all 10 tasks and defines acceptance gates and automated test suites.
4. **01-architecture-spec.md & 02-component-spec.md -> 06-subtask-llm-train-heatmap-linking.md**: Ensures the training curriculum linker consumes both local companion databases (`repo-cache/.heatmap/`) and remote REST endpoints via the cached URL reader.
5. **02-component-spec.md -> 07-subtask-coding-guidelines-and-skills.md**: Establishes repository standards for segmented URL caching and documents CLI usage in `.agents/skills/gitmap/skill.md`.

---

## Coding Guidelines Compliance Scorecard

All modules, types, and functions designed and implemented under Task 314 strictly conform to repository coding standards:

| Guideline Family | Core Requirement | Design & Architectural Enforcement | Compliance Status |
| :--- | :--- | :--- | :--- |
| **Boolean Naming & Polarity** | All booleans must start with `is*` or `has*` only. Zero negative booleans (`noCache`, `disabled` banned). | `hasCached`, `hasAria2`, `isExpired`, `isValid`, `isAvailable`, `hasRangeSupport`. CLI `--no-cache` flag maps immediately to internal positive struct field `hasForceRefresh: true`. | `100% COMPLIANT` |
| **Function Length Limits** | Max 15 lines per function; $\le 8$ lines preferred. | Complex operations decomposed into small single-purpose helper functions (e.g. `probeContentLength`, `spawnAria2Worker`, `writeChunkToBuffer`, `formatRepoAliasLine`). | `100% COMPLIANT` |
| **Zero Nested Conditionals** | No nested `if` statements. Guard clauses and early returns mandatory. | Every conditional block inverted into immediate guard returns; branching logic partitioned across discrete helper functions. | `100% COMPLIANT` |
| **Error Return Sovereignty** | Leaf functions returning `error` must return wrapped errors directly (`*apperror.AppError`). Never call `cliexit.HandleError` in leaf packages. | `cli/store/`, `cli/downloader/`, `cli/heatmap/`, and `cli/cmdurl/` leaf functions wrap errors using `apperror.WrapSimple(err, "<scope>")` and bubble to top-level command dispatchers in `cli/root.go`. | `100% COMPLIANT` |
| **Vertical Spacing & Return Style** | Blank line after `}` and blank line before `return`. | Formatted with consistent vertical rhythm across all new and refactored Go files. | `100% COMPLIANT` |
| **Strict Relative Paths** | Zero hardcoded drive letters or absolute file URI schemes. | All database paths resolved relative to `EffectiveDataDir()` or Git root using `filepath.Join` and POSIX-normalized relative strings. | `100% COMPLIANT` |
| **Worker Git Command Ban** | Subagents never run git commands to prevent `.git/index.lock` collisions. | Zero git invocations within subagents. Lead orchestrator executes atomic `gitmap cpf` at conclusion. | `100% COMPLIANT` |
| **Zero Builds & Tests in Planning** | No intermediate full test runs or full builds during planning and coding micro-batches. | Build and full-suite test execution strictly deferred to decoupled validation stages. | `100% COMPLIANT` |

---

## Full Acceptance Criteria Matrix

### Area 1: Heatmap Discovery & Help Architecture
- [ ] **AC-1.1 (CLI Invocation)**: Running `gitmap heatmap` or `gitmap rc heatmap` without arguments outputs a self-documenting discovery overview rather than failing with argument syntax errors.
- [ ] **AC-1.2 (Flag Routing)**: Passing `--repo <name|alias>` properly routes to repository-filtered heatmap display without injecting erroneous subcommand tokens into CLI dispatch.
- [ ] **AC-1.3 (Multi-Tier DB Inventory)**: Discovery output enumerates active databases across all 3 tiers:
  - Root DB: `EffectiveDataDir()/heatmap.db`
  - Per-Repo DBs: `EffectiveDataDir()/repodb/<repo>.db`
  - Multi-AI Repo-Cache DBs: `repo-cache/.heatmap/<seq>-<ai>/<repo>/heatmap.db`
- [ ] **AC-1.4 (Alias Cross-Referencing)**: Active repositories are displayed alongside their registered short aliases resolved dynamically from `db.ListAliasesWithRepo()`.
- [ ] **AC-1.5 (Recent Task Aggregation)**: Terminal output presents the top 5–10 recent heatmap tasks executed across 2–5 repositories (minimum 5 total tasks), displaying Task ID, Repo Name, Action Count, Success Rate, and Timestamp.
- [ ] **AC-1.6 (AI Terminal Suggestions)**: The output renders actionable command suggestion pills (e.g. `gitmap heatmap --repo <alias>`, `gitmap rc heatmap --task <slug>`, `gitmap heatmap --export json`).

### Area 2: Dedicated URL Split-DB Caching Engine
- [ ] **AC-2.1 (Split-DB Creation)**: GitMap manages a dedicated database file at `EffectiveDataDir()/repodb/url_cache.db`, registered in `SplitDatabaseRegistry`.
- [ ] **AC-2.2 (High-Performance PRAGMAs)**: Database connection initializes with `PRAGMA journal_mode = WAL`, `PRAGMA synchronous = NORMAL`, and `PRAGMA mmap_size = 268435456` (256 MB memory mapping).
- [ ] **AC-2.3 (Single-Writer Governance)**: Write operations (`Set`, `UpdateMetadata`, `PurgeExpired`) are protected by an in-process sync mutex to guarantee single-writer safety across concurrent swarm subagents.
- [ ] **AC-2.4 (Deduplication & Revalidation)**: URLs are indexed by SHA-256 hash. Re-requests verify HTTP ETag / Last-Modified timestamps using lightweight `If-None-Match` / `If-Modified-Since` conditional headers.
- [ ] **AC-2.5 (Sub-Millisecond Read)**: Valid cached entries stream directly from SQLite to stdout in $< 0.5\text{ ms}$, bypassing all network socket creation.
- [ ] **AC-2.6 (CLI Command)**: `gitmap url read <url>` outputs raw unbuffered bytes via `output.Raw()`. Supports `--refresh` (force network revalidation), `--no-cache` (bypass cache write), and `--json` (emit cache metadata and headers).

### Area 3: 10-Split Segmented & Aria2 Downloader
- [ ] **AC-3.1 (Size Threshold Detection)**: When fetching URLs or files exceeding 1–2 MB (configurable threshold `downloadSegmentMinBytes = 1 << 20`), the engine triggers parallel segmentation.
- [ ] **AC-3.2 (Aria2c Auto-Detection)**: The engine inspects `exec.LookPath("aria2c")`. If present on host, GitMap invokes `aria2c -s 10 -x 10 -k 5K -d <dest_dir> -o <filename> <url>` for accelerated parallel transfer.
- [ ] **AC-3.3 (Native Go 10-Worker Fallback)**: If `aria2c` is not installed or returns non-zero, the engine automatically falls back to native Go parallel HTTP Range requests across 10 concurrent worker routines with 5KB chunk alignment.
- [ ] **AC-3.4 (URL Cache Integration)**: Content fetched via segmented downloading is automatically committed to `repodb/url_cache.db` so subsequent reads execute locally.

### Area 4: LLM Training Integration & Repository Standards
- [ ] **AC-4.1 (Heatmap Training Flag)**: `gitmap llm train` and `gitmap llm train muse` support `--heatmap <path|url>` to ingest decision traces and action logs into training datasets.
- [ ] **AC-4.2 (Remote & Local Telemetry Ingestion)**: If `--heatmap` points to a URL, GitMap reads through the cached URL engine; if it points to a file, it reads directly from `repo-cache/.heatmap/<seq>-<ai>/<repo>/heatmap.db`.
- [ ] **AC-4.3 (Coding Guidelines Specification)**: Authored `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md` and `11-heatmap-telemetry-and-url-caching.md` specifying telemetry formats and URL caching rules.
- [ ] **AC-4.4 (GitMap Skill Synchronization)**: Updated `.agents/skills/gitmap/skill.md` with complete CLI reference for heatmap discovery, URL caching, and segmented downloading.
- [ ] **AC-4.5 (Release Ceremony)**: Executed SemVer minor version bump (`gitmap release` / `minor-bump`) and atomic commit via `gitmap cpf`.

---

## Verification Test Matrices

### 1. Unit Test Matrix

| Test Suite File | Target Function / Method | Verification Objective | Pass Criteria |
| :--- | :--- | :--- | :--- |
| `cli/heatmap/discovery_test.go` | `DiscoverHeatmapDatabases()` | Scan root, per-repo, and repo-cache directories with mock files. | Returns populated slice with correct database tiers and relative paths. |
| `cli/heatmap/recent_tasks_test.go` | `AggregateRecentTasks()` | Query mock SQLite databases across multiple repositories. | Returns top 5–10 recent tasks across 2–5 repositories ordered descending by timestamp. |
| `cli/heatmap/suggestions_test.go` | `BuildTerminalSuggestions()` | Generate suggestions from repository list and aliases. | Generates non-empty suggestion list including alias-based command pills. |
| `cli/store/url_cache_test.go` | `Get()`, `Set()` | Insert and retrieve mock HTTP payloads in `url_cache.db`. | Retrieved payload matches inserted bytes; SHA-256 hash verified; latency $< 1\text{ ms}$. |
| `cli/store/url_cache_test.go` | `HasExpired()` | Evaluate TTL and conditional validation flags. | Correct boolean return based on expiration timestamp and ETag presence. |
| `cli/downloader/segmented_test.go` | `CalculateSegments()` | Split 5 MB resource into 10 byte ranges. | Exactly 10 contiguous segments covering range `0` to `totalSize - 1`. |
| `cli/downloader/aria2_test.go` | `DetectAria2()` | Probe system environment for `aria2c` binary. | Returns positive boolean `hasAria2` without panicking when executable missing. |
| `cli/cmd/llm/llm_heatmap_test.go` | `ParseHeatmapSource()` | Parse local relative file path vs remote HTTP URL. | Correctly classifies source type and invokes appropriate ingestion driver. |

### 2. End-to-End Integration Test Matrix (Isolated)

| E2E Test Identifier | Execution Command | Target Validation Scope | Verification Invariant |
| :--- | :--- | :--- | :--- |
| **E2E-HM-01** | `gitmap heatmap` | Discovery overview, database inventory, alias mapping. | Clean table rendered; exit code 0; zero SQL errors. |
| **E2E-HM-02** | `gitmap rc heatmap --repo <alias>` | Specific repository filtering via short alias. | Filters exclusively to targeted repo; displays action logs and task statistics. |
| **E2E-URL-01** | `gitmap url read <local-test-url>` | First fetch from local mock HTTP test server. | Payload saved to `repodb/url_cache.db`; raw bytes streamed to stdout. |
| **E2E-URL-02** | `gitmap url read <local-test-url>` | Second fetch for same URL. | Served directly from SQLite cache in $< 0.5\text{ ms}$; zero server hits recorded. |
| **E2E-URL-03** | `gitmap url read --json <test-url>` | Cache inspection mode. | Outputs JSON document containing URL, SHA-256, byte size, ETag, and hit count. |
| **E2E-DL-01** | Segmented download test with mock Range server | 10 parallel HTTP Range requests for 2 MB payload. | Reassembles file into bit-for-bit identical SHA-256 hash. |
| **E2E-LLM-01** | `gitmap llm train --heatmap <test-db-path>` | Ingestion of decision traces into training dataset. | Extracts action logs, formats prompt/completion pairs, logs ingested record count. |

### 3. Static Analysis & Linter Verification Matrix

| Linter / Check Script | Target Scope | Enforcement Rule | Required Outcome |
| :--- | :--- | :--- | :--- |
| `03-ai-scripts/07-relative-path-fixer.py` | All `.go`, `.md`, `.json` files | Zero absolute paths or drive letters. | `0 violations found` |
| `03-ai-scripts/check-enum-and-boolean.py` | New and modified Go packages | Positive booleans (`is*`, `has*` only); enum `Type` suffixes. | `0 violations found` |
| AST Function Limit Inspector | `cli/heatmap/`, `cli/store/`, `cli/downloader/`, `cli/cmdurl/` | Function bodies $\le 15$ lines. | `0 violations found` |
| AST Nested If Inspector | New and modified Go source files | Zero nested `if` statements. | `0 violations found` |

---

## Rollback & Safety Gates

1. **Database Isolation**: The URL caching Split-DB (`repodb/url_cache.db`) is completely isolated from `gitmap.db` and repository index files. In the event of SQLite database corruption, deleting `repodb/url_cache.db` causes the engine to recreate an empty database on the next execution without data loss.
2. **Network Fallback Chain**: If `aria2c` execution encounters errors or terminates with non-zero exit status, the downloader falls back immediately to the native Go 10-connection parallel Range downloader. If HTTP Range headers are rejected with HTTP 200, it falls back to standard single-stream streaming.
3. **Heatmap Discovery Non-Interference**: Heatmap discovery commands perform strictly read-only queries against `heatmap.db` and `.heatmap/` directories, ensuring existing tracking sessions are never locked or mutated.
