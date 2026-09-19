# 128 — AUM Automation Suite & Polyglot Script Migration Roadmap

## Overview

**Module Number:** 128  
**Version:** 1.1.0  
**Updated:** 2026-09-19  
**Status:** Strategic Architecture & Roadmap  
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

The `03-ai-scripts/` directory houses 44 specialized repository automation algorithms. This specification organizes all 44 scripts into a cohesive 6-phase migration and integration roadmap.

```
03-ai-scripts/
├── [Foundation & Engine]
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
├── [Phase 2: Topology, Schema & Database Generation]
│   ├── 18-codebase-topology-discoverer➔ Graphs entrypoints, schemas, workflows, languages
│   ├── 30-db-struct-enum-generator.py ➔ SQLite schema to Go/TS struct & enum generator
│   ├── 31-db-migration-runner.py      ➔ Transactional SQLite schema migration runner
│   └── 34-schema-scanner.py           ➔ Split-DB schema validator against conventions
│
├── [Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers]
│   ├── 06-cicd-local-runner.py        ➔ Parallel multi-core checker orchestration
│   ├── 16-installer-smoke-tester.py   ➔ Cross-platform installer dry-run tester
│   ├── 28-go-preflight-ci.py          ➔ Local preflight CI quality gates
│   ├── 33-test-inventory-generator.py ➔ Test inventory manifest and duration baselines
│   └── 34-purge-github-actions-art... ➔ Actions zero-storage purge workflow
│
├── [Phase 4: Release, SemVer & Version Synchronization]
│   ├── 14-version-sync-checker.py     ➔ Multi-project version synchronization
│   ├── 29-release-bumper.py           ➔ SemVer version bumping and tag creation
│   ├── 29-release-orchestrator.py     ➔ Full release ceremony orchestration
│   └── 38-milestone-consolidator.py   ➔ GitHub milestone and changelog consolidator
│
├── [Phase 5: Documentation, Spec Migration & Memory]
│   ├── 09-cli-help-auditor.py         ➔ CLI help AST simulation & documentation parity
│   ├── 20-plan-consolidator.py        ➔ Subtask and memory plan consolidator
│   ├── 22-doc-path-linter.py          ➔ Markdown relative link integrity auditor
│   ├── 23-coding-guideline-path-con...➔ Coding guideline link consolidator
│   ├── 24-spec-path-migrator.py       ➔ Spec numbering and reference migrator
│   ├── 25-repo-migrator.py            ➔ Repository structure reorganization
│   ├── 31-md-gap-fixer.py             ➔ Markdown spacing and table alignment fixer
│   └── 32-deep-consolidator.py        ➔ Multi-tier plan memory consolidator
│
└── [Phase 6: Git Hygiene, Cleanup & History Purging]
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

### Phase 2: Topology, Schema & Database Generation (Planned & Next)
*Target: Database automation, migrations, and code generation.*

1. **`gitmap aum topology` (from `18-codebase-topology-discoverer.py`):**
   - Discovers entry points, database schemas, workflow files, and language distributions, caching results into `.gitmap/data/automation/sql.db`.
2. **`gitmap aum db-generate` (from `30-db-struct-enum-generator.py`):**
   - Inspects SQLite tables and generates strongly typed Go structs, TypeScript interfaces, and enum definitions adhering to PascalCase rules.
3. **`gitmap aum db-migrate` (from `31-db-migration-runner.py`):**
   - Executes ordered SQL migration files with transaction rollback protection and migration history tracking.
4. **`gitmap aum schema-audit` (from `34-schema-scanner.py`):**
   - Audits SQLite split-db schemas against repository naming rules.

---

### Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers (Planned)
*Target: Eliminating CI feedback latency with parallel local verification.*

1. **`gitmap aum preflight` (from `28-go-preflight-ci.py` & `06-cicd-local-runner.py`):**
   - Executes parallel multi-core checks across CPU cores (gofmt, linters, relative paths, nested ifs, tests).
2. **`gitmap aum test-inventory` (from `33-test-inventory-generator.py`):**
   - Discovers all unit and integration tests across the repository, generating `.ai-memory/test-inventory.json` with duration baselines.
3. **`gitmap aum purge-actions` (from `34-purge-github-actions-artifacts.py`):**
   - Queries GitHub Actions API and purges obsolete artifacts to maintain 0.0 GB storage footprint (Rule R18).

---

### Phase 4: Release, SemVer & Version Synchronization (Planned)
*Target: Unbreakable automated release ceremony.*

1. **`gitmap aum version-sync` (from `14-version-sync-checker.py`):**
   - Verifies version alignment across `cli/constants/constants_version.go`, `package.json`, `Cargo.toml`, and docs.
2. **`gitmap aum release-bump` (from `29-release-bumper.py` & `29-release-orchestrator.py`):**
   - Increments SemVer, updates changelog, synchronizes package files, and triggers release pipeline.
3. **`gitmap aum milestones` (from `38-milestone-consolidator.py`):**
   - Consolidates GitHub milestone issues into release notes.

---

### Phase 5: Documentation, Spec Migration & Memory Consolidation (Planned)
*Target: Keeping 120+ specifications and AI memory compact and gapless.*

1. **`gitmap aum help-audit` (from `09-cli-help-auditor.py`):**
   - Simulates CLI `--help` across all subcommands and verifies documentation parity in `cli/helptext/`.
2. **`gitmap aum plan-consolidate` (from `20-plan-consolidator.py` & `32-deep-consolidator.py`):**
   - Clusters completed tasks and subtasks in `.ai-memory/plans/` into consolidated milestone summaries, reducing file count while preserving outcomes.
3. **`gitmap aum doc-links` (from `22-doc-path-linter.py`):**
   - Validates all relative markdown links across `02-spec/` and `.ai-memory/`.

---

### Phase 6: Git Hygiene, Cleanup & History Purging (Planned)
*Target: Fast cleanup of artifacts and git hygiene.*

1. **`gitmap aum clean-artifacts` (from `19-artifact-remover.py`):**
   - Safely deletes compiled binaries, `.pytest_cache`, `__pycache__`, and temporary test dumps without touching tracked sources.
2. **`gitmap aum changed-files` (from `27-git-changed-files.py`):**
   - Returns JSON list of modified and untracked files relative to `origin/main` for targeted linter passes.
3. **`gitmap aum purge-history` (from `30-purge-history.py` & `33-git-history-tracer-and-purge.py`):**
   - Safely cleans large historical git blobs without rewriting recent commit lineage.

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

# Phase 2 Database & Topology (Planned)
gitmap aum topology                     # Codebase entrypoints & schema graph
gitmap aum db-generate                  # SQLite schema to Go/TS struct generator
gitmap aum db-migrate                   # Ordered SQL migrations with rollback
gitmap aum schema-audit                 # Split-db schema auditor

# Phase 3 CI/CD & Multi-Core (Planned)
gitmap aum preflight                    # Local multi-core preflight runner
gitmap aum test-inventory               # Test inventory manifest generator
gitmap aum purge-actions                # GitHub Actions 0.0 GB storage purger

# Phase 4 Release & SemVer (Planned)
gitmap aum version-sync                 # Multi-project version synchronization
gitmap aum release-bump                 # SemVer bump & changelog synchronizer

# Phase 5 Documentation & Memory (Planned)
gitmap aum help-audit                   # CLI help AST parity auditor
gitmap aum plan-consolidate             # Memory plan & subtask consolidator
gitmap aum doc-links                    # Markdown relative link integrity auditor

# Phase 6 Git Hygiene & Purge (Planned)
gitmap aum clean-artifacts              # Temporary artifact and build purger
gitmap aum changed-files                # Git changed and staged file detector
gitmap aum purge-history                # Historical blob purger
```

---

## 5. Cross-References

- LLM Train & Chained Curriculum: [`./127-llm-train-and-chained-agent-curriculum.md`](./127-llm-train-and-chained-agent-curriculum.md)
- Polyglot Worker Runner: [`./124-polyglot-worker-orchestrator-and-automation-runner.md`](./124-polyglot-worker-orchestrator-and-automation-runner.md)
- LLM Orchestration Playbook: [`./125-automation-llm-orchestration-guide.md`](./125-automation-llm-orchestration-guide.md)
- Rule R19 File Size Guard: [`../../02-coding-guidelines/01-cross-language/31-file-size-guard-and-binary-exclusion.md`](../../02-coding-guidelines/01-cross-language/31-file-size-guard-and-binary-exclusion.md)
