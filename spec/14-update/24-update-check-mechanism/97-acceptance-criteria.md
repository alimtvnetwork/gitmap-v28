# Acceptance Criteria

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

An implementation is **conformant** iff every item below is verifiable
in code, tests, and runtime behavior.

---

## A. Discovery

| # | Criterion | Verified by |
|---|-----------|-------------|
| 1 | All six probes (V + V+1..V+5) fire in parallel — measured wall time of full discovery ≤ 6s on a healthy network | Integration test with mocked HTTP |
| 2 | Discovery NEVER walks past V+5, even if all five lookahead probes succeed | Unit test: assert `len(Candidates) == 5` always |
| 3 | A 404 on any probe is treated as not-found (no retry, no error escalation) | Unit test |
| 4 | A 200 with malformed JSON is logged and treated as not-found | Unit test |
| 5 | The highest version returning HTTP 200 + valid JSON wins | Unit test with mixed 200/404 matrix |
| 6 | Per-probe HTTP timeout = 5s; total discovery deadline = 10s | Code review + timeout test |

## B. Persistence

| # | Criterion | Verified by |
|---|-----------|-------------|
| 7 | `UpdateChecker` table stores BOTH `RawJson` AND each parsed column | Schema inspection + insert/select test |
| 8 | `UpdateStatus` is implemented as a code enum AND a DB lookup table with TINYINT PK | Schema inspection + enum test |
| 9 | All DB writes wrapped in try/catch (or language equivalent) | Static check / code review |
| 10 | A failed re-check does NOT clear `LatestVersion`/`HasUpdate` from a prior successful check | Unit test |
| 11 | When no DB exists, JSON fallback at `~/.<CliName>/data/UpdateChecker.json` works with the same interface | Backend-swap integration test |
| 12 | JSON fallback writes are atomic (tmp + rename) | Crash-injection test |

## C. CLI Behavior

| # | Criterion | Verified by |
|---|-----------|-------------|
| 13 | Plain `update-check` runs synchronously, prints formatted output, AND persists | E2E test |
| 14 | `update-check --async` returns in < 200 ms and the child writes to the store | Process-timing test |
| 15 | `update-check --force` bypasses the interval gate | Unit test |
| 16 | `do-update` runs unattended (no interactive prompt) | E2E test |
| 17 | `do-update` failure leaves `HasUpdate = true` so the warning persists | Unit test |
| 18 | Exit codes match [06 §3](./06-cli-commands.md#3-exit-codes) | E2E test |

## D. Pre-Command Hook

| # | Criterion | Verified by |
|---|-----------|-------------|
| 19 | Pre-hook returns in < 50 ms on the happy path | Benchmark |
| 20 | Pre-hook NEVER blocks the user's command, even on subsystem error | Fault-injection test |
| 21 | Pre-hook does NOT recurse on `update-check` / `do-update` | Unit test |
| 22 | `BackgroundUpdateCheckEnabled = false` disables the spawn | Config test |
| 23 | `PendingUpdateWarningEnabled = false` disables the warning | Config test |
| 24 | Trailing warning prints to **stderr**, not stdout | Pipe-redirection test |

## E. Logging & Errors

| # | Criterion | Verified by |
|---|-----------|-------------|
| 25 | Every error is written to `~/.<CliName>/Logs/UpdateChecker.log` with file + line metadata | Log-format test |
| 26 | Subsystem errors set `ErrorMessage`, `ErrorAt`, and `UpdateStatusId = Failed` | Unit test |
| 27 | Errors are NEVER silently swallowed except at the named pre-hook boundary | Static check (CODE RED P1) |
| 28 | Log file is capped at 1 MiB and rotated to `.log.1` | File-size test |

## F. Naming & Standards

| # | Criterion | Verified by |
|---|-----------|-------------|
| 29 | All JSON keys, JSON string-enum values, table names, column names, and code identifiers use PascalCase | `validate-guidelines.py` |
| 30 | All free-text columns are nullable with no DEFAULT (Schema Rule 12) | SQL linter |
| 31 | `UpdateChecker` includes `Description`, `Notes`, `ErrorMessage` (Schema Rule 11 — transactional) | Schema inspection |
| 32 | No nested `if` in service code (CODE RED P6) | `validate-guidelines.py` |

## G. Migration & Backwards-Compat

| # | Criterion | Verified by |
|---|-----------|-------------|
| 33 | A project upgrading from JSON to SQLite migrates the JSON contents into the new table on first init | Migration test |
| 34 | Status JSON `NewRepoUrl` non-null surfaces a migration banner via the warning channel | E2E test |

---

**Total: 34 criteria.** All must pass before tagging a release that
ships this subsystem.

---

*Acceptance Criteria — v1.0.0 — 2026-04-20*
