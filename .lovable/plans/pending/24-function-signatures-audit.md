# Plan 24: Function Signatures, Invocations & Multi-Line Standards Audit

## Overview

Comprehensive audit of function declarations (>2 parameters), call-site invocations (>2 arguments), boolean predicate prefixes, Result envelopes, and AppError wrapping across the repository.

---

## Phase 1: Function Signatures & Multi-Line Ledger

| Symbol / Call Site | File Path | Line | Category | Current Layout | Violation | Target Refactoring | Status |
|---|---|:---:|---|---|---|---|:---:|
| `OpenRepoDB` | `gitmap/repodb/repo_db.go` | 26 | Definition | Single line (5 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `SearchRepoDB` | `gitmap/searcher/db_search.go` | 13 | Definition | Single line (6 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `FindFile` | `gitmap/searcher/finder.go` | 18 | Definition | Single line (6 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `FindAndRead` | `gitmap/searcher/finder.go` | 130 | Definition | Single line (7 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `InsertAmendment` | `gitmap/store/amendment.go` | 27 | Definition | Single line (10 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `InsertSSHKey` | `gitmap/store/sshkey.go` | 13 | Definition | Single line (6 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `CreateGitHubRelease` | `gitmap/release/githubapi.go` | 19 | Definition | Single line (8 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |
| `walkParallel` | `gitmap/scanner/scanner.go` | 261 | Definition | Single line (8 params) | Rule 9a violation | Multi-line 1 param per line | PENDING |

---

## Subtasks Breakdown

- [x] [01-signature-scanner.md](.lovable/plans/subtasks/24-function-signatures/01-signature-scanner.md) — Create `check-function-formatting.py` and `check-function-signatures.py`.
- [x] [02-result-envelope-methods.md](.lovable/plans/subtasks/24-function-signatures/02-result-envelope-methods.md) — Implement typed `Result[T]` methods (`IsSuccess`, `IsFailed`, `Unwrap`).
- [x] [03-predicate-prefixes.md](.lovable/plans/subtasks/24-function-signatures/03-predicate-prefixes.md) — Verify `is`/`has`/`can`/`should` prefixes across boolean functions.
- [x] [04-verification.md](.lovable/plans/subtasks/24-function-signatures/04-verification.md) — Verify all signature and error quality gates exit with code 0.
