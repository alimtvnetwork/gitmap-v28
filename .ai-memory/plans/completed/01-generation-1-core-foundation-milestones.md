# Generation 1 Milestone: Core Architecture Foundation & Subsystem Baseline

**Milestone ID:** 01  
**Generational Scope:** Historical Milestones 01 through 15  
**Status:** Completed  
**Completion Date:** 2026-09-19  
**Target Scope:** Coding Guidelines, Error Management, Type Safety, CI/CD, Database, Git Engine, SSH Fleet, Terminal UI, Chrome Vault, Multi-OS Installers, History Tracer  
**Canonical Specifications:**
- [01-cli-architecture](../../../02-spec/21-app/01-cli-architecture/01-architecture-spec.md)
- [02-scanner-and-projects](../../../02-spec/21-app/02-scanner-and-projects/01-architecture-spec.md)
- [03-git-operations-and-pull](../../../02-spec/21-app/03-git-operations-and-pull/01-architecture-spec.md)
- [04-fleet-nodes-and-ssh](../../../02-spec/21-app/04-fleet-nodes-and-ssh/01-architecture-spec.md)
- [06-database-and-split-db](../../../02-spec/21-app/06-database-and-split-db/01-architecture-spec.md)
- [07-pipeline-and-diagnostics](../../../02-spec/21-app/07-pipeline-and-diagnostics/01-architecture-spec.md)
- [08-distribution-and-release](../../../02-spec/21-app/08-distribution-and-release/01-architecture-spec.md)

---

## 1. Executive Summary & Consolidated Generational Scope

Generation 1 represents the foundational architecture and subsystem maturity baseline of GitMap, synthesizing 15 initial milestones (covering 108 original implementation plans and 42 subtask directories). This generational document compacts raw historical execution chronicles while maintaining 100% fidelity on architectural decisions, Go type safety contracts, error propagation invariants, and verified operational capabilities.

### Historical Milestones Synthesized:
1. **Milestone 01: Coding Guidelines, Sizing, Booleans & Style Quality** (12 original plans): Enforced function body caps (<= 15 lines), file size caps (<= 100 lines), inverted guard clauses, affirmative boolean naming (`is*`, `has*`), and strict relative path hygiene.
2. **Milestone 02: Centralized Error Architecture & Cliexit Engine** (5 original plans): Unified error handling around `*appfault.AppError`, standard error code ranges (`E1001`–`E9000`), zero swallowed errors, and centralized exit handlers (`cliexit.Fail`, `cliexit.Exit`).
3. **Milestone 03: Type Safety, Signatures, Enums & React Architecture** (10 original plans): Replaced multi-value error tuples with `Result[T]` monadic wrappers, typed enums with `*Type` suffixes, parameter structs for functions with > 3 parameters, and modular React components.
4. **Milestone 04: CI/CD Pipelines, Multi-Worker Runners & Real-Time Streaming** (10 original plans): Multi-worker local CI/CD runner (`01-run-all-local.py`), real-time streaming output, smart test inventory caching, and automated PR trigger reconciliation.
5. **Milestone 05: Database Engine, SQLite Schema, Joins & Typed Scanners** (6 original plans): Standardized PascalCase SQLite schemas, singular table names, `{Entity}Id` primary keys, zero-swallow row scanners, and WAL mode single-writer concurrency.
6. **Milestone 06: Git Operations, Commit Engines & Interactive Remediation** (14 original plans): Flat commit command suite (`gitmap c`, `gitmap commit`), automatic git staging, commit-right path resolution, interactive dirty-repo remediation, and delta extraction.
7. **Milestone 07: SSH Nodes, Cluster Delegation & Remote Execution Engine** (2 original plans): Concurrent SSH execution engine (`gitmap se`), dynamic `~/.ssh/config` templating, and broadcast command delegation.
8. **Milestone 08: Terminal UI, Help Parity & Interactive Macro Builder** (15 original plans): Lipgloss ANSI terminal UI, column-aligned help tables, command AST parity, and interactive macro builder (`gitmap mb`).
9. **Milestone 09: Chrome Profile Management, Picker & Token Vault** (3 original plans): Chromium Local State 13-attribute registration, graphical profile picker integration, preflight validation flags (`--json`, `--fnf`), and reversible Base64 + Caesar cipher token vault.
10. **Milestone 10: Multi-OS Installers, Scripts & Web Stacks** (18 original plans): Unified installation scripts (`run.ps1`, `run.sh`, `local-install`), VMware shared folder persistence, scripts-fixer engine, and web stack installers (Nginx, WordPress, Laravel).
11. **Milestone 11: Completed Plans Historical Consolidation** (1 original plan): Pre-flight backup verification protocols and initial plan compaction methodology.
12. **Milestone 12: Git Deleted Files Tracer & Purger** (1 original plan): Deep Git history scanning and surgical packfile purge tool (`33-git-history-tracer-and-purger.py`) with restore and discard modes.
13. **Milestone 13: CI/CD Pipeline Logs, SQLite Error History & Incremental DB** (11 original plans): Repository-scoped `repodb/pipeline.db`, bounded error extraction (-1, -2, -3 offsets), compact OK-line suppression, and two-pass log parser.
14. **Milestone 14: Pipeline AGY Fix Injection Suite** (7 original plans): Automated AGY prompt generation from failing pipeline runs, deduplication flags (`--force`), and direct IDE queue injection.
15. **Milestone 15: SSH Multi-Command, Multi-Machine Join & Liveness Parity** (14 original plans): Concurrent port 22 TCP liveness scanning, space-delimited target resolution, host joining/recall, and aligned terminal help parity.

---

## 2. Preserved Go Type Contracts & Architecture Invariants

### 2.1 Core Error Envelope & Monadic Result Wrappers

```go
package appfault

// ErrorCode defines standardized classification for domain errors
type ErrorCode string

const (
    ErrValidation ErrorCode = "E1001"
    ErrNotFound   ErrorCode = "E2001"
    ErrDatabase   ErrorCode = "E3001"
    ErrGit        ErrorCode = "E4001"
    ErrSSH        ErrorCode = "E5001"
    ErrPipeline   ErrorCode = "E6001"
    ErrInternal   ErrorCode = "E9001"
)

// AppError is the universal structured error envelope
type AppError struct {
    Code       ErrorCode `json:"code"`
    Message    string    `json:"message"`
    Err        error     `json:"-"`
    Caller     string    `json:"caller"`
    StatusCode int       `json:"statusCode"`
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return e.Message + ": " + e.Err.Error()
    }
    return e.Message
}

// Result provides monadic single-value return with structured error
type Result[T any] struct {
    Value T
    Error *AppError
}

func Success[T any](v T) Result[T] {
    return Result[T]{Value: v, Error: nil}
}

func Failure[T any](err *AppError) Result[T] {
    return Result[T]{Error: err}
}
```

### 2.2 StreamWriter Pipeline Contract

```go
package stream

import "io"

// StreamWriter defines unbuffered, progressive streaming output
type StreamWriter interface {
    WriteLine(line string) error
    WriteChunk(p []byte) (int, error)
    Flush() error
}

// JSONEnvelopeV2 provides standard CLI machine-readable response payload
type JSONEnvelopeV2 struct {
    IsSuccess bool        `json:"isSuccess"`
    Data      interface{} `json:"data,omitempty"`
    Error     *AppError   `json:"error,omitempty"`
    Meta      MetaPayload `json:"meta"`
}

type MetaPayload struct {
    DurationMs int64  `json:"durationMs"`
    Command    string `json:"command"`
    Version    string `json:"version"`
}
```

### 2.3 Database Invariants & Tesla ID Standards
- **Singular PascalCase Tables:** All table definitions use PascalCase (`Repository`, `Project`, `SshHost`, `PipelineRun`).
- **Tesla ID Suffix:** Primary keys are strictly typed integer IDs named `{TableName}Id` (`RepositoryId`, `SshHostId`). Never bare `id` or `ID`.
- **Concurrency & WAL:** `SetMaxOpenConns(1)` for all SQLite database handles to avoid locking conflicts, with WAL mode enabled (`PRAGMA journal_mode=WAL`).
- **Binary-Anchored Paths:** All database root resolution uses `filepath.EvalSymlinks` to prevent symlink drift.

---

## 3. Dense Architectural Decisions Ledger

| Milestone | Domain | Architectural Decision | Technical Rationale & Invariant |
| :--- | :--- | :--- | :--- |
| **M01** | Guidelines | Inverted guard clauses & early returns | Flattened branching depth to <= 1; eliminated deeply nested if blocks. |
| **M01** | Guidelines | Positive boolean naming (`is*`, `has*`) | Eliminated double negatives (`!isNotReady`); positive polarity enforcement. |
| **M02** | Error Mgmt | Universal `*appfault.AppError` wrapping | Eliminated bare `error` propagation; every error captures code, caller, and context. |
| **M02** | Error Mgmt | Centralized CLI exit via `cliexit.Fail` | Terminated unhandled `os.Exit` calls scattered across subcommands. |
| **M03** | Type Safety | Monadic `Result[T]` returns | Replaced error tuples `(T, error)` across packages for single-return purity. |
| **M03** | Type Safety | Enums suffixed with `*Type` | Clarified type definitions in Go and TS; eliminated magic string literals. |
| **M04** | CI/CD | Multi-worker local test runner | Decoupled CI execution from slow cloud triggers; enabled sub-second local test runs. |
| **M04** | CI/CD | Test inventory manifest caching | Avoids scanning entire repository for tests on every run; incremental hashes. |
| **M05** | Database | PascalCase SQLite schema + joins | Standardized entity mapping; enabled universal join query builder. |
| **M05** | Database | Safe scanner generator | Eliminated raw column index scans; guaranteed zero swallowed SQL scan errors. |
| **M06** | Git Engine | Flat commit CLI (`gitmap c`, `cm`) | Replaced complex manual git commit flags with semantic flat commit flags. |
| **M06** | Git Engine | Automatic stage and dirty remediation | Automatically detects unstaged changes and offers surgical interactive staging. |
| **M07** | SSH Fleet | Concurrent SSH execution pool | Parallelized commands across remote nodes with non-blocking stdout/stderr streaming. |
| **M07** | SSH Fleet | Dynamic `~/.ssh/config` templating | Preserves existing user host configs while safely merging GitMap cluster nodes. |
| **M08** | Terminal UI | Lipgloss ANSI table alignment | Eliminated misaligned monospace ASCII dumps; auto-fit terminal width. |
| **M08** | Terminal UI | Interactive Macro Builder (`mb`) | Visual step recorder with validation before macro disk persistence. |
| **M09** | Chrome Vault | 13-attribute Chromium Local State | Prevents profile corruption on import by ensuring 100% schema completeness. |
| **M09** | Chrome Vault | Reversible Base64 + Caesar cipher vault | Secures sensitive tokens in profile stores while remaining lightweight. |
| **M10** | Installers | Dual runner scripts (`run.ps1`, `run.sh`) | Unified bootstrapping across Linux, macOS, and Windows with automatic dependency check. |
| **M10** | Installers | VMware shared folder crontab sync | Resolved VMware mount disappearance upon guest reboot via persistent systemd units. |
| **M12** | Git Tracer | Deep history purge tool | Scans Git packfiles for deleted artifacts and purges historical commit trees. |
| **M13** | Telemetry | Bounded stack trace extraction | Captures 5 leading + 20 trailing lines around failing test lines in pipeline DB. |
| **M13** | Telemetry | OK-line suppression in logs | Reduced raw pipeline logs by > 80% while keeping actionable error context. |
| **M14** | AGY Fix | Automatic pipeline fix prompt injection | Direct feed of failing CI frames into AGY prompt queue for autonomous remediation. |
| **M15** | SSH Liveness | Port 22 concurrent TCP dialer | Pre-flight node discovery filters unreachable hosts with a 1500ms timeout. |

---

## 4. Synthesized Verification Outcomes Ledger

| Domain | Scope | Verified Command / Test Gate | Concrete Outcome | Status |
| :--- | :--- | :--- | :--- | :---: |
| **Coding Style** | M01 | `python3 03-ai-scripts/04-lint-coding-guidelines.py` | 0 nesting violations, 100% positive booleans, all files <= 100 lines. | `VERIFIED` |
| **Error Handling** | M02 | `go test ./pkg/appfault/...` | All errors wrap with caller frame; zero swallowed errors. | `VERIFIED` |
| **Type Safety** | M03 | `go test ./pkg/types/...` | Monadic `Result[T]` unwraps cleanly; enums validate with `*Type`. | `VERIFIED` |
| **CI/CD Runner** | M04 | `python3 01-run-all-local.py --fast` | Fast worker pool runs target test suites under 12 seconds. | `VERIFIED` |
| **Database** | M05 | `go test ./pkg/db/...` | PascalCase tables validated; single-writer WAL mode verified. | `VERIFIED` |
| **Git Operations** | M06 | `gitmap commit -m "test"` | Staged changes committed with normalized commit author and trailer. | `VERIFIED` |
| **SSH Nodes** | M07 | `gitmap se "echo ready"` | Command dispatched across registered cluster nodes concurrently. | `VERIFIED` |
| **Terminal UI** | M08 | `gitmap help` | High-contrast Lipgloss tables aligned to active terminal width. | `VERIFIED` |
| **Chrome Picker** | M09 | `gitmap chrome import --json` | 13 Local State attributes verified; profile picker opens correctly. | `VERIFIED` |
| **Installers** | M10 | `bash run.sh --check` | Pre-flight dependency audit passes across clean Ubuntu and macOS. | `VERIFIED` |
| **Git Purger** | M12 | `python3 03-ai-scripts/33-git-history-tracer-and-purger.py --scan` | Successfully identifies deleted packfile objects without data loss. | `VERIFIED` |
| **Pipeline DB** | M13 | `gitmap pe -t` | Negative commit offset (-1, -2) loads bounded error stack traces. | `VERIFIED` |
| **AGY Injection** | M14 | `gitmap agy inject --last` | Injects failing pipeline frame directly into AGY IDE prompt queue. | `VERIFIED` |
| **SSH Liveness** | M15 | `gitmap nodes ping` | Port 22 TCP probes return healthy status for active fleet nodes. | `VERIFIED` |

---

## 5. Quality Gates & Preserved State Invariants

- **`isConsolidated == true`:** Historical Milestones 01–15 fully synthesized into this single Generation 1 reference.
- **`isAuthoritative == true`:** Supersedes all previous exploratory plans and subtasks in this domain.
- **`isPreserved == true`:** All Go contracts (`AppError`, `Result[T]`, `StreamWriter`) preserved with complete type fidelity.
- **`hasVerifiedOutcome == true`:** All listed verification commands and quality gates tested and recorded.
- **`isPendingIsolated == true`:** Zero modifications performed on pending or open plans.
