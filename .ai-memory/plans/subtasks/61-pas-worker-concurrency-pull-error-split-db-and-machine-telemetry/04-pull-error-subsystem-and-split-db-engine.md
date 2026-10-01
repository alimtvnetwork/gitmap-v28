# Subtask 04: Pull-Error Subsystem & Split-DB Engine

> **Parent Plan:** [61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry](../../pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Tracking Spec:** [198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Primary File Targets:** `cli/cmdpullerror/`, `cli/store/split_db_pull_error.go`, `cli/store/pull_split_db_errors.go`, `cli/cmd/rootcore.go`  

---

## 1. Objective

Implement the `gitmap pull-error` (`pulle`, `pull-e`, `pull-errors`) subsystem backed by a persistent Split-DB table:
1. Records all failed repository pull operations, error messages, stack traces, timestamps, and actionable remediation commands.
2. Context-aware CLI inspection:
   - When inside a repo (`cwd` contains `.git`), inspects that specific repository's latest pull error.
   - When outside any repo, lists all recent pull failures across all repositories.
3. Supports `--ssh` to aggregate remote fleet pull errors.

---

## 2. Implementation Scope

- **`cli/store/split_db_pull_error.go` & `cli/store/pull_split_db_errors.go`:**
  - Create table `pull_errors` in Split-DB (`data/gitmap-pull.db`).
  - Methods: `RecordPullError`, `ListRecentPullErrors`, `GetPullErrorByRepo`.
- **`cli/cmdpullerror/`:**
  - Implement CLI command handler for `gitmap pull-error [target] [--json] [--ssh]`.
  - Format diagnostic error cards: Repo name, Timestamp, Stage, Error message, Remediation hint.
- **`cli/cmd/rootcore.go`:**
  - Register root command aliases `pull-error`, `pull-errors`, `pulle`, `pull-e`.

---

## 3. Verification

- Run `gitmap pull-error` in terminal and verify clean output.
- Run `gitmap pull-error --json` and verify structured output.
- Verify compilation with `go build ./...` in `cli/`.
