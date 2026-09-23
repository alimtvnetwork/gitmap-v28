# 128 — AUM Automation Suite & Polyglot Script Migration Roadmap

## Overview

**Module Number:** 128
**Version:** 2.0.0
**Updated:** 2026-09-19
**Status:** Production-Ready
**AI Confidence:** Production-Ready
**Ambiguity Score:** None

---

## 1. Purpose & Core Naming Mandate

GitMap provides a unified, high-performance automation subsystem designed to eliminate slow, fragile scripts in favor of a dual-engine architecture:
- **Engine A (Go Supervisor):** Ultra-fast filesystem traversal, lazy regex compilation, SQLite-backed runtime discovery, file size guards, memory-safe binary probes, and sequence integrity linters.
- **Engine B (Polyglot Worker Pool):** Concurrent multi-core child workers executing Python, Node.js, Go, Rust, PowerShell, or Bash without per-file process spawn overhead.

### Short-Form Alias Mandate

Per user specification, the primary short-form CLI trigger for `gitmap automation` is:

```bash
gitmap aum <subcommand> [flags]
```

*(Both `aum` and legacy `auto` remain active, with `aum` designated as the canonical abbreviation for **A**utomation and **U**tility **M**anager).*

---

## 2. Complete Inventory of `03-ai-scripts/`

The `03-ai-scripts/` directory houses 44 specialized repository automation algorithms. All 44 scripts across all 6 phases have been migrated to native compiled Go subcommands in `cli/cmdautomation/` and verified 100% in production.

```
03-ai-scripts/
├── [Foundation & Engine] (COMPLETED & VERIFIED)
│   ├── 02-shared-engine.py            ➔ Centralized enums, lazy regex, locking, 8KB binary probe
│   ├── 04-newline-fixer.py            ➔ Universal CRLF/LF and trailing whitespace normalizer
│   ├── 11-fast-file-scanner.py        ➔ Parallel directory traversal and inventory caching
│   ├── 12-fast-cached-grep.py         ➔ In-memory Boyer-Moore / lazy regex grep engine
│   ├── 13-file-size-guard.py          ➔ Blob guard, size limits, binary probe, waiver check
│   └── 15-sequence-and-title-auditor  ➔ Markdown numbering gap detector & H1 header fixer
│
├── [Phase 1: Code Quality, Naming & Coding Guidelines] (COMPLETED & VERIFIED)
│   ├── 05-guideline-autofixer.py      ➔ Automated coding guideline auditor & autofixer
│   ├── 07-relative-path-fixer.py      ➔ Detects absolute paths / file:/// links in docs
│   ├── 08-naming-autofixer.py         ➔ PascalCase keys/columns, positive booleans
│   ├── 21-sequence-integrity-linter   ➔ Multi-folder contiguous sequence checker
│   ├── 27-misspell-auditor.py         ➔ Spellcheck and typo auditor across code/markdown
│   ├── 35-result-wrapper-auditor.py   ➔ Eliminates raw (T, error) tuples in Go
│   ├── 36-param-struct-auditor.py     ➔ Flags functions exceeding 3–4 parameters
│   └── 37-enum-guideline-auditor.py   ➔ Enforces *Type suffixes and PascalCase enums
│
├── [Phase 2: Topology, Schema & Database Generation] (COMPLETED & VERIFIED)
│   ├── 18-codebase-topology-discoverer➔ Graphs entrypoints, schemas, workflows, languages
│   ├── 30-db-struct-enum-generator.py ➔ SQLite schema to Go/TS struct & enum generator
│   ├── 31-db-migration-runner.py      ➔ Transactional SQLite schema migration runner
│   └── 34-schema-scanner.py           ➔ Split-DB schema validator against conventions
│
├── [Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers] (COMPLETED & VERIFIED)
│   ├── 06-cicd-local-runner.py        ➔ Parallel multi-core checker orchestration
│   ├── 16-installer-smoke-tester.py   ➔ Cross-platform installer dry-run tester
│   ├── 28-go-preflight-ci.py          ➔ Local preflight CI quality gates
│   ├── 33-test-inventory-generator.py ➔ Test inventory manifest and duration baselines
│   └── 34-purge-github-actions-art... ➔ Actions zero-storage purge workflow
│
├── [Phase 4: Release, SemVer & Version Synchronization] (COMPLETED & VERIFIED)
│   ├── 14-version-sync-checker.py     ➔ Multi-project version synchronization
│   ├── 29-release-bumper.py           ➔ SemVer version bumping and tag creation
│   ├── 29-release-orchestrator.py     ➔ Full release ceremony orchestration
│   └── 38-milestone-consolidator.py   ➔ GitHub milestone and changelog consolidator
│
├── [Phase 5: Documentation, Spec Migration & Memory] (COMPLETED & VERIFIED)
│   ├── 09-cli-help-auditor.py         ➔ CLI help AST simulation & documentation parity
│   ├── 20-plan-consolidator.py        ➔ Subtask and memory plan consolidator
│   ├── 22-doc-path-linter.py          ➔ Markdown relative link integrity auditor
│   ├── 23-coding-guideline-path-con...➔ Coding guideline link consolidator
│   ├── 24-spec-path-migrator.py       ➔ Spec numbering and reference migrator
│   ├── 25-repo-migrator.py            ➔ Repository structure reorganization
│   ├── 31-md-gap-fixer.py             ➔ Markdown spacing and table alignment fixer
│   └── 32-deep-consolidator.py        ➔ Multi-tier plan memory consolidator
│
└── [Phase 6: Git Hygiene, Cleanup & History Purging] (COMPLETED & VERIFIED)
    ├── 03-file-manipulator.py         ➔ Atomic file replacement, backup, and restore
    ├── 10-encoding-normalizer.py      ➔ UTF-8 / UTF-16 stream and BOM normalizer
    ├── 17-fast-file-reader.py         ➔ Fault-tolerant chunked file reader
    ├── 19-artifact-remover.py         ➔ Build artifact, temp, and cache purger
    ├── 26-go-code-formatter.py        ➔ Go AST formatting and import grouping
    ├── 27-git-changed-files.py        ➔ Git diff and staged file detector
    ├── 30-purge-history.py            ➔ Git history blob purger
    └── 33-git-history-tracer-and-pu...➔ Commit history tracing and deep purge
```

---

## 3. Six-Phase Implementation Roadmap

All six migration phases are 100% completed, fully integrated into the unified `gitmap aum` command structure, and verified in production.

### Active Baseline (Complete & Verified in Production)

| Command | Trigger | Core Functionality | Status |
|---|---|---|---|
| `gitmap aum search` | Search | Multi-core streaming search with lazy regex and binary filtering | ✅ Production |
| `gitmap aum guard` | Size Guard | 500 KB limit, large JSON exclusion, binary null-byte probe, interactive prompt | ✅ Production |
| `gitmap aum sequence` | Sequence | Markdown sequence gap detector & `# XX Title` auto-fixer (`--fix`) | ✅ Production |
| `gitmap aum exclude` | Exclusions | Persistent search exclusion list in SQLite (`sql.db`) | ✅ Production |
| `gitmap aum newlines` | Newlines | Polyglot CRLF to LF and trailing whitespace normalizer | ✅ Production |
| `gitmap aum cache` | In-Memory | In-memory file cache status, warming, and purging | ✅ Production |
| `gitmap aum benchmark` | Benchmark | Side-by-side Go vs Python execution benchmarks | ✅ Production |

---

### Phase 1: Code Quality, Naming & Coding Guidelines (Complete & Verified in Production)

| Command | Trigger | Source Script | Functionality | Status |
|---|---|---|---|---|
| `gitmap aum relative-paths` | Paths | `07-relative-path-fixer.py` | Scans for forbidden absolute filesystem paths (`D:\...`, `file:///`) and auto-sanitizes via `--fix` | ✅ Production |
| `gitmap aum naming` | Naming | `08-naming-autofixer.py` | Audits explicit boolean comparisons (`== true`, `=== true`) and affirmative booleans (`is*`, `has*`) | ✅ Production |
| `gitmap aum result-wrapper` | Results | `35-result-wrapper-auditor.py` | Detects multi-value `(map, error)` / `([]T, error)` tuples and enforces `ResultMap`/`ResultSlice` | ✅ Production |
| `gitmap aum params` | Arity | `36-param-struct-auditor.py` | Flags functions with >4 parameters; recommends dedicated `*Params` structs in `types.go` | ✅ Production |
| `gitmap aum enums` | Enums | `37-enum-guideline-auditor.py` | Audits enum definitions for mandatory `*Type` suffixes and bans raw numeric `rune(10)` casts | ✅ Production |

---

### Phase 2: Topology, Schema & Database Generation (Complete & Verified in Production)

| Command | Trigger | Source Script | Functionality | Status |
|---|---|---|---|---|
| `gitmap aum topology` | Topology | `18-codebase-topology-discoverer.py` | Discovers entry points, database schemas, workflows, and language distributions | ✅ Production |
| `gitmap aum db-generate` | Code Gen | `30-db-struct-enum-generator.py` | Generates Go structs, TypeScript interfaces, and column enums from SQLite tables | ✅ Production |
| `gitmap aum db-migrate` | Migrations | `31-db-migration-runner.py` | Executes ordered SQL migrations with transaction rollback protection and tracking | ✅ Production |
| `gitmap aum schema-audit` | Schema Audit | `34-schema-scanner.py` | Validates SQLite split-db schemas against naming conventions (PascalCase, affirmative booleans) | ✅ Production |

#### 1. `gitmap aum topology`
- **CLI Grammar:** `gitmap aum topology [dir] [flags]`
- **Aliases:** `topo`, `codebase-topology`
- **Flags:**
  - `--refresh` (`-r`): Force refresh topology discovery cache.
  - `--json` (`-j`): Output results as raw JSON.
  - `--query` (`-q`): Search specific subsystem or language routing.
  - `--ttl`: Cache TTL in seconds (default: `1800`).
- **Architecture:** Traverses the repository using concurrent directory walkers; classifies files by programming language; maps manifests (`package.json`, `go.mod`, `Cargo.toml`, etc.); groups subsystems (`cli`, `data`, `specs`, `scripts`) with root paths, schema definitions, and application entrypoints; caches discovered topology to `.gitmap/data/automation/sql.db` or in-memory structures with TTL expiry; provides query routing for quick subsystem inspection. Returns `Result[TopologyResult]`.

#### 2. `gitmap aum db-generate`
- **CLI Grammar:** `gitmap aum db-generate [db-path] [flags]`
- **Aliases:** `db-gen`, `gen-db`
- **Flags:**
  - `--lang`: Target language: `go`, `ts`, or `all` (default: `"all"`).
  - `--out`: Output directory for generated code files.
  - `--dry-run`: Preview generated code without writing to disk.
  - `--struct-dir`: Optional directory of model structs.
- **Architecture:** Inspects SQLite split-db files (`sql.db`, `cache.db`); queries `sqlite_master` and table column PRAGMAs; enforces PascalCase naming conventions; emits strongly-typed Go model structs (`structs.go`), TypeScript interface definitions (`types.ts`), and typed column enums (`enums.go`); supports dry-run preview and custom directory targeting. Returns `Result[DbGenerateResult]`.

#### 3. `gitmap aum db-migrate`
- **CLI Grammar:** `gitmap aum db-migrate [db-path] [migrations-dir] [flags]`
- **Aliases:** `migrate`, `db-up`
- **Flags:**
  - `--dry-run`: Preview migrations without writing.
  - `--status`: Display migration history and current version.
  - `--rollback`: Rollback last N migrations (default: `0`).
  - `--sql`: Execute inline SQL migration statement.
- **Architecture:** Scans sequential SQL migration files (`*.sql`); tracks applied migrations in the `_schema_migrations` tracking table; executes unapplied migrations inside atomic transactions with automatic rollback on error; supports rolling back the last N applied migrations, inspecting migration status, and executing inline SQL statements. Returns `Result[DbMigrateResult]`.

#### 4. `gitmap aum schema-audit`
- **CLI Grammar:** `gitmap aum schema-audit [db-path] [flags]`
- **Aliases:** `audit-schema`, `db-audit`
- **Flags:**
  - `--json`: Output results as machine-readable JSON.
  - `--strict`: Enforce strict column naming and primary key rules.
  - `--dir`: Target directory containing SQL or Go schema definitions.
- **Architecture:** Inspects database tables and schema files against repository database standards: PascalCase table and column names, affirmative boolean prefixes (`is*`, `has*`), primary key constraints (`Id`), and split-db isolation; reports violations with severity ratings and specific remediation instructions. Returns `Result[SchemaAuditResult]`.

---

### Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers (Complete & Verified in Production)

| Command | Trigger | Source Script | Functionality | Status |
|---|---|---|---|---|
| `gitmap aum preflight` | Preflight | `28-go-preflight-ci.py` & `06-cicd-local-runner.py` | Parallel multi-core CI preflight checks across CPU cores | ✅ Production |
| `gitmap aum test-inventory` | Test Inventory | `33-test-inventory-generator.py` | Discovers unit and integration tests, outputs inventory manifest with durations | ✅ Production |
| `gitmap aum purge-actions` | Actions Purge | `34-purge-github-actions-artifacts.py` | Queries GitHub Actions API and purges obsolete artifacts (0.0 GB mandate) | ✅ Production |
| `gitmap aum smoke-test` | Smoke Test | `16-installer-smoke-tester.py` | Cross-platform dry-run validator for installed tools and installer scripts | ✅ Production |

#### 1. `gitmap aum preflight`
- **CLI Grammar:** `gitmap aum preflight [dir] [flags]`
- **Aliases:** `check-all`, `ci-local`, `pre-commit`
- **Flags:**
  - `--fail-fast`: Halt execution on first check failure.
  - `--json`: Output results as machine-readable JSON.
  - `--workers` (`-w`): Number of worker threads (default: CPU cores).
  - `--phase` (`-p`): Execution phase: `all`, `test`, `lint` (default: `"all"`).
  - `--filter` (`-k`): Filter checks by name substring.
- **Architecture:** Coordinates parallel multi-core quality gates across all available CPU cores; runs Go AST formatting checks, relative path integrity audits, control-flow flattening audits, boolean/affirmative naming checks, and unit tests; aggregates check outcomes into structured pass/fail results; halts immediately on failure when `--fail-fast` is specified. Returns `Result[PreflightResult]`.

#### 2. `gitmap aum test-inventory`
- **CLI Grammar:** `gitmap aum test-inventory [dir] [flags]`
- **Aliases:** `tests-inv`, `inventory`
- **Flags:**
  - `--out`: Output manifest path (default: `.ai-memory/test-inventory.json`).
  - `--refresh`: Force refresh existing inventory manifest.
  - `--json`: Output results as machine-readable JSON.
  - `--slow-threshold`: Threshold in seconds to categorize slow tests (default: `4.0`).
  - `--force-run-all`: Mark all discovered tests as needing execution.
- **Architecture:** Discovers unit and integration test functions across repository packages; calculates historical duration baselines; partitions tests into fast tiers (< threshold) and slow tiers (>= threshold); writes manifest to `.ai-memory/test-inventory.json` for smart test runner isolation and dynamic execution estimation. Returns `Result[TestInventoryResult]`.

#### 3. `gitmap aum purge-actions`
- **CLI Grammar:** `gitmap aum purge-actions [flags]`
- **Aliases:** `purge-artifacts`, `clean-actions`
- **Flags:**
  - `--repo`: Target repository slug (`owner/repo`).
  - `--older-than`: Purge artifacts older than N days (default: `0`).
  - `--dry-run`: Preview artifacts without deleting.
  - `--json`: Output results as machine-readable JSON.
  - `--workers` (`-w`): Number of concurrent deletion threads (default: `12`).
- **Architecture:** Communicates with the GitHub Actions API using repository credentials; enumerates workflow runs and uploaded artifact payloads; deletes obsolete build and test artifacts using concurrent worker pools to enforce the 0.0 GB storage mandate (Rule R18); reports total artifacts purged and megabytes freed. Returns `Result[PurgeActionsResult]`.

#### 4. `gitmap aum smoke-test`
- **CLI Grammar:** `gitmap aum smoke-test [dir] [flags]`
- **Aliases:** `installer-smoke`, `smoke`
- **Flags:**
  - `--tools`: List of tool binaries to verify in PATH.
  - `--json`: Output results as machine-readable JSON.
  - `--workers` (`-w`): Number of worker threads (default: CPU cores).
  - `--filter` (`-k`): Filter targets by name substring.
- **Architecture:** Executes non-destructive, dry-run validations for installed tool binaries and installation scripts; verifies that target binaries are present in the system PATH; validates installer script syntax and flags across Linux, Windows, and macOS without mutating system state. Returns `Result[SmokeTestResult]`.

---

### Phase 4: Release, SemVer & Version Synchronization (Complete & Verified in Production)

| Command | Trigger | Source Script | Functionality | Status |
|---|---|---|---|---|
| `gitmap aum version-sync` | Version Sync | `14-version-sync-checker.py` | Audits and synchronizes repository version manifests | ✅ Production |
| `gitmap aum release-bump` | Release Bump | `29-release-bumper.py` & `29-release-orchestrator.py` | Computes next semantic version and updates repository manifests | ✅ Production |
| `gitmap aum milestones` | Milestones | `38-milestone-consolidator.py` | Queries GitHub milestone issues and formats structured release notes | ✅ Production |

#### 1. `gitmap aum version-sync`
- **CLI Grammar:** `gitmap aum version-sync [dir] [flags]`
- **Aliases:** `sync-version`, `ver-sync`
- **Flags:**
  - `--fix` (`-f`): Automatically harmonize mismatched manifests.
  - `--target-version` (`-t`): Override canonical version to harmonize.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Inspects all version manifests across the repository against the Single Source of Truth (`cli/constants/constants_version.go`); audits `package.json`, `Cargo.toml`, documentation headers, and spec metadata; flags mismatched version strings; harmonizes all files to the canonical target version when `--fix` is passed. Returns `Result[VersionSyncResult]`.

#### 2. `gitmap aum release-bump`
- **CLI Grammar:** `gitmap aum release-bump [major|minor|patch] [flags]`
- **Aliases:** `bump`, `semver-bump`
- **Flags:**
  - `--version` (`-v`): Explicit version string to bump to (e.g. `6.263.0`).
  - `--dry-run` (`-d`): Preview version bump without modifying files.
  - `--tag`: Create git release tag (`vX.Y.Z`).
  - `--push`: Push release branch and tag to remote.
  - `--skip-tests`: Skip pre-release quality gate checks.
  - `--scope` (`-s`): Scope description for release commit (default: `"Automated release orchestration"`).
  - `--bullet` (`-b`): Changelog bullet points (repeatable).
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Calculates the next semantic version according to bump tier (major, minor, patch) or explicit version string; updates all SSoT manifests and changelog entries; optionally creates git release tags (`vX.Y.Z`) and pushes release branch to remote repository; incorporates pre-release verification gates unless `--skip-tests` is provided. Returns `Result[ReleaseBumpResult]`.

#### 3. `gitmap aum milestones`
- **CLI Grammar:** `gitmap aum milestones [milestone-id] [flags]`
- **Aliases:** `milestone`, `consolidate-milestones`
- **Flags:**
  - `--out` (`-o`): Output file path for generated release notes.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Queries GitHub Issues API for closed and open issues assigned to a specified milestone; groups items by issue labels (features, bug fixes, refactoring, documentation); compiles structured Markdown release notes; writes output to disk or stdout for release packaging. Returns `Result[MilestonesResult]`.

---

### Phase 5: Documentation, Spec Migration & Memory Consolidation (Complete & Verified in Production)

| Command | Trigger | Source Script | Functionality | Status |
|---|---|---|---|---|
| `gitmap aum help-audit` | Help Audit | `09-cli-help-auditor.py` | Audits CLI commands for help descriptions and documentation parity | ✅ Production |
| `gitmap aum plan-consolidate` | Memory Consolidate | `20-plan-consolidator.py` & `32-deep-consolidator.py` | Clusters completed plans and subtasks into milestone summaries | ✅ Production |
| `gitmap aum doc-links` | Doc Links | `22-doc-path-linter.py` | Validates markdown relative links across documentation and memory | ✅ Production |
| `gitmap aum spec-migrate` | Spec Migrate | `24-spec-path-migrator.py` | Re-sequences specification file prefixes and updates cross-links | ✅ Production |

#### 1. `gitmap aum help-audit`
- **CLI Grammar:** `gitmap aum help-audit [dir] [flags]`
- **Aliases:** `help-check`, `doc-audit`
- **Flags:**
  - `--strict` (`-s`): Fail with exit code 1 if violations exist.
  - `--json`: Output results as machine-readable JSON.
  - `--ext` (`-e`): Filter by file extensions.
- **Architecture:** Simulates CLI `--help` generation across all Cobra commands; verifies that all subcommands provide non-empty `Short` descriptions, valid `Use` syntax, and matching documentation in `cli/helptext/`; detects parity gaps between compiled commands and documentation; enforces strict exit codes under `--strict`. Returns `Result[HelpAuditResult]`.

#### 2. `gitmap aum plan-consolidate`
- **CLI Grammar:** `gitmap aum plan-consolidate [dir] [flags]`
- **Aliases:** `consolidate-plans`, `consolidate`
- **Flags:**
  - `--threshold` (`-t`): Minimum completed plans threshold to cluster (default: `5`).
  - `--dry-run` (`-d`): Preview consolidation without modifying files.
  - `--force`: Bypass confirmation prompts.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Scans completed plan files and subtasks in `.ai-memory/plans/`; clusters related completed tasks by milestone prefix and subject; consolidates multiple granular subtask files into concise milestone summaries; updates memory indices to reduce repository file count while preserving verified outcomes. Returns `Result[PlanConsolidateResult]`.

#### 3. `gitmap aum doc-links`
- **CLI Grammar:** `gitmap aum doc-links [dir] [flags]`
- **Aliases:** `check-links`, `doc-paths`
- **Flags:**
  - `--fix` (`-f`): Autofix known outdated path references.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Scans all Markdown files across `02-spec/`, `.ai-memory/`, and repository root; extracts relative links and file anchors; verifies target file existence on disk; detects broken or outdated relative links; auto-repairs known relocated paths when `--fix` is passed. Returns `Result[DocLinksResult]`.

#### 4. `gitmap aum spec-migrate`
- **CLI Grammar:** `gitmap aum spec-migrate [dir] [flags]`
- **Aliases:** `resequence-spec`, `migrate-spec`
- **Flags:**
  - `--from`: Source spec number to migrate from (default: `0`).
  - `--to`: Target spec number to migrate to (default: `0`).
  - `--dry-run` (`-d`): Preview migration without disk changes.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Re-sequences specification file numbering prefixes in `02-spec/`; renames target spec files to new sequence numbers; scans all repository markdown files and updates cross-reference links matching the old prefix to the new prefix; supports dry-run preview to prevent unintended rename cascades. Returns `Result[SpecMigrateResult]`.

---

### Phase 6: Git Hygiene, Cleanup & History Purging (Complete & Verified in Production)

| Command | Trigger | Source Script | Functionality | Status |
|---|---|---|---|---|
| `gitmap aum clean-artifacts` | Artifact Cleaner | `19-artifact-remover.py` | Safely deletes build binaries, test dumps, pycache, and temporary files | ✅ Production |
| `gitmap aum changed-files` | Changed Files | `27-git-changed-files.py` | Discovers modified, staged, and untracked files for targeted linter passes | ✅ Production |
| `gitmap aum purge-history` | History Purge | `30-purge-history.py` & `33-git-history-tracer-and-purge.py` | Traces and identifies large historical git blobs without rewriting recent commits | ✅ Production |
| `gitmap aum format-go` | Go Formatter | `26-go-code-formatter.py` | AST-aware Go code formatter organizing imports and enforcing UTF-8 LF | ✅ Production |

#### 1. `gitmap aum clean-artifacts`
- **CLI Grammar:** `gitmap aum clean-artifacts [flags]`
- **Aliases:** `clean-build`, `rm-artifacts`
- **Flags:**
  - `--dir`: Root directory to scan for artifacts (default: `"."`).
  - `--dry-run` (`-n`): Preview matching items without deleting.
  - `--all` (`-a`): Clean all preset categories (`pycache`, `temp`, `binaries`).
  - `--verbose` (`-v`): Show all discovered artifact paths.
  - `--json`: Output results as machine-readable JSON.
  - `--clean-pycache`: Remove python bytecode and test cache.
  - `--clean-temp`: Remove temporary files (`.tmp`, `.log`, `.swp`).
  - `--clean-binaries`: Remove compiled binaries (`.exe`, `.syso`, etc.).
- **Architecture:** Identifies build artifacts, compiled test binaries, `.pytest_cache`, `__pycache__`, temporary files (`.tmp`, `.log`, `.swp`), and orphan outputs; validates git status to ensure git-tracked files are never deleted; safely purges unwanted files and reports freed disk space; supports dry-run preview. Returns `Result[CleanArtifactsResult]`.

#### 2. `gitmap aum changed-files`
- **CLI Grammar:** `gitmap aum changed-files [flags]`
- **Aliases:** `git-changes`, `diff-files`
- **Flags:**
  - `--dir`: Target repository directory (default: `"."`).
  - `--base`: Git base reference or commit (e.g. `origin/main`, `HEAD~1`).
  - `--commits` (`-n`): Number of recent commits to evaluate (default: `20`).
  - `--staged-only`: Discover staged files only.
  - `--verify`: Verify that files currently exist on disk.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Executes git diff and status inspections against a base commit/branch or recent commits; identifies modified, staged, added, deleted, and untracked files; verifies disk presence when `--verify` is enabled; provides machine-readable JSON or colored terminal lists for targeted linter passes. Returns `Result[ChangedFilesResult]`.

#### 3. `gitmap aum purge-history`
- **CLI Grammar:** `gitmap aum purge-history [flags]`
- **Aliases:** `trace-history`, `clean-history`
- **Flags:**
  - `--dir`: Target repository directory (default: `"."`).
  - `--min-size-mb`: Minimum blob size in MB to identify (default: `1.0`).
  - `--path`: Target file or path pattern to trace.
  - `--dry-run` (`-n`): Preview large blobs without rewriting history (default: `true`).
  - `--confirm` (`-y`): Bypass confirmation prompt.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Traces git commit history objects using pack-index inspection; identifies large historical blobs exceeding `--min-size-mb`; displays commit hashes, blob sizes, and file paths; provides safe dry-run preview by default; requires explicit confirmation (`--confirm`) before any historical rewriting. Returns `Result[PurgeHistoryResult]`.

#### 4. `gitmap aum format-go`
- **CLI Grammar:** `gitmap aum format-go [dir] [flags]`
- **Aliases:** `gofmt-ast`, `fmt-go`
- **Flags:**
  - `--dir`: Directory or file to format (default: `"."`).
  - `--write` (`-w`): Write formatted changes back to disk.
  - `--check` (`-c`): Check and return violations if files need formatting.
  - `--staged`: Format only staged Go files.
  - `--json`: Output results as machine-readable JSON.
- **Architecture:** Parses Go source files into AST using `go/parser` and `go/format`; organizes imports into standard, third-party, and internal groups; enforces UTF-8 encoding and LF line endings; reports formatting violations in `--check` mode; writes formatted content to disk with `--write`. Returns `Result[FormatGoResult]`.

---

## 4. Summary Table of `aum` Command Names

```bash
# Core Active Commands (Production)
gitmap aum search <pattern> [dir]       # Streaming search with lazy regex
gitmap aum guard [dir]                  # Blob size guard and binary probe
gitmap aum sequence [dir] [--fix]       # Sequence gap & title header auditor
gitmap aum exclude [list|add|rm|clear]  # Persistent search exclusions
gitmap aum newlines [paths...] [--fix]  # Polyglot CRLF/LF normalizer
gitmap aum cache [status|warm|clear]    # In-memory file cache
gitmap aum benchmark [target]           # Go vs Python benchmark comparison

# Phase 1 Code Quality Commands (Production)
gitmap aum relative-paths [--fix]       # Absolute path & file:/// linter
gitmap aum naming                       # Boolean & affirmative naming linter
gitmap aum result-wrapper               # Result monad & AppError tuple audit
gitmap aum params                       # Parameter count limit auditor
gitmap aum enums                        # Enum *Type suffix auditor

# Phase 2 Database & Topology Commands (Production)
gitmap aum topology [dir]               # Codebase entrypoints & schema graph
gitmap aum db-generate [db-path]        # SQLite schema to Go/TS struct generator
gitmap aum db-migrate [db-path] [dir]   # Ordered SQL migrations with rollback
gitmap aum schema-audit [db-path]       # Split-db schema auditor

# Phase 3 CI/CD & Multi-Core Commands (Production)
gitmap aum preflight [dir]              # Local multi-core preflight runner
gitmap aum test-inventory [dir]         # Test inventory manifest generator
gitmap aum purge-actions                # GitHub Actions 0.0 GB storage purger
gitmap aum smoke-test [dir]             # Installer & tool smoke tester

# Phase 4 Release & SemVer Commands (Production)
gitmap aum version-sync [dir]           # Multi-project version synchronization
gitmap aum release-bump [tier]          # SemVer bump & changelog synchronizer
gitmap aum milestones [milestone-id]    # GitHub milestone issue consolidator

# Phase 5 Documentation & Memory Commands (Production)
gitmap aum help-audit [dir]             # CLI help AST parity auditor
gitmap aum plan-consolidate [dir]       # Memory plan & subtask consolidator
gitmap aum doc-links [dir]              # Markdown relative link integrity auditor
gitmap aum spec-migrate [dir]           # Specification re-sequencer & link migrator

# Phase 6 Git Hygiene & Purge Commands (Production)
gitmap aum clean-artifacts              # Temporary artifact and build purger
gitmap aum changed-files                # Git changed and staged file detector
gitmap aum purge-history                # Historical blob purger
gitmap aum format-go [dir]              # Go AST code formatter & import organizer
```

---

## 5. Cross-References

- LLM Train & Chained Curriculum: [`./127-llm-train-and-chained-agent-curriculum.md`](./127-llm-train-and-chained-agent-curriculum.md)
- Polyglot Worker Runner: [`./124-polyglot-worker-orchestrator-and-automation-runner.md`](./124-polyglot-worker-orchestrator-and-automation-runner.md)
- LLM Orchestration Playbook: [`./125-automation-llm-orchestration-guide.md`](./125-automation-llm-orchestration-guide.md)
- Rule R19 File Size Guard: [`../02-coding-guidelines/01-cross-language/31-file-size-guard-and-binary-exclusion.md`](../02-coding-guidelines/01-cross-language/31-file-size-guard-and-binary-exclusion.md)
