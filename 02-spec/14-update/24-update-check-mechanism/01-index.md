# Update Check Mechanism — Overview

. **CRITICAL AI INSTRUCTION:** This `01-index.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.


> **Version:** 1.0.0
> **Created:** 2026-04-20
> **Status:** Active
> **AI Confidence:** Production-Ready
> **Ambiguity:** None
> **Parent:** [02-spec/14-update/01-index.md](../01-index.md)

---

## Purpose

Defines the **non-blocking, parallel, status-script-driven** update-check
mechanism for every CLI in the Riseup Asia stack. This is the **detection**
half of the update lifecycle. The **application** half (rename-first deploy,
handoff, cleanup) is already covered by `02-self-update-overview.md` through
`08-console-safe-handoff.md` in the parent folder.

This is **not a cron job**. There is no daemon, no scheduled task, no
background service. Every check is initiated by the CLI itself, gated by
a configurable interval (default 12h), and runs as a fire-and-forget
detached process so the user is never blocked.

---

## Keywords

`update-check` · `status-script` · `parallel-discovery` · `version-probe` ·
`pre-command-hook` · `do-update` · `update-checker-table` · `pascal-case` ·
`async-detached` · `seedable-config`

---

## Defining Properties

1. **Status script lives on GitHub, not on the user machine.** The CLI
   fetches `Status.ps1` (Windows) or `Status.sh` (Unix) from
   `raw.githubusercontent.com/{owner}/{repo}/main/` and executes it
   server-side semantics — the script's job is to emit a JSON document.
2. **Parallel V → V+5 discovery, no walking.** Six probes fire at once.
   Even if all five lookahead probes succeed, discovery never continues
   past V+5.
3. **Newer-repo detection by suffix.** Repo names encode the major
   version (`repo-v15`, `repo-v16`, …). The discoverer probes
   `repo-v{N+1}` … `repo-v{N+5}` against the same owner.
4. **Interval-gated, not time-driven.** The CLI checks the stored
   `LastCheckedAt`. If `Now - LastCheckedAt < CheckIntervalHours`, it
   skips. Otherwise it spawns `update-check --async` and continues.
5. **Persistence in BOTH modes.** Sync `update-check` prints AND writes.
   `--async` writes silently. There is no `--persist` flag.
6. **PascalCase everywhere.** JSON keys, JSON string-enum values,
   table names, column names, and code identifiers all use PascalCase.

---

## File Inventory

| # | File | Description |
|---|------|-------------|
| 00 | [01-index.md](./01-index.md) | This file — index and defining properties |
| 01 | [02-fundamentals.md](./02-fundamentals.md) | Discovery algorithm, version-probe URLs, owner/repo conventions |
| 02 | [02-status-script-json.md](./03-status-script-json.md) | `Status.ps1` / `Status.sh` output JSON schema (PascalCase) |
| 03 | [03-combined-json.md](./04-combined-json.md) | Combined JSON after parallel discovery — fed to consumers |
| 04 | [04-database-schema.md](./05-database-schema.md) | `UpdateChecker` table, `UpdateStatus` enum + lookup |
| 05 | [05-update-checker-service.md](./06-update-checker-service.md) | Reusable `UpdateCheckerService` class contract |
| 06 | [06-cli-commands.md](./07-cli-commands.md) | `update-check` (flags) and `do-update` commands |
| 07 | [07-pre-command-hook.md](./08-pre-command-hook.md) | Pre-command hook, interval gate, trailing warning line |
| 08 | [08-error-handling.md](./09-error-handling.md) | Try/catch policy, file-system log, `ErrorMessage` column |
| 09 | [09-json-fallback-store.md](./10-json-fallback-store.md) | JSON-file storage when no SQLite database exists |
| 97 | [97-acceptance-criteria.md](./97-acceptance-criteria.md) | 16-point acceptance criteria |
| 97b | [97-changelog.md](./98-changelog.md) | Changelog |
| 99 | [99-consistency-report.md](./99-consistency-report.md) | Consistency report |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Self-update (apply phase) | [../02-self-update-overview.md](../02-self-update-overview.md) |
| Install script contract | [../19-install-scripts.md](../19-install-scripts.md) |
| Version-pinned installers | [../../16-generic-release/09-version-pinned-release-installers.md](../../16-generic-release/09-version-pinned-release-installers.md) |
| Seedable config (where settings live) | [../../06-seedable-config-architecture/01-index.md](../../06-seedable-config-architecture/01-index.md) |
| Local SQLite store | [../../13-generic-cli/10-database.md](../../13-generic-cli/10-database.md) |

---

## Resolved Decisions (no ambiguity remains)

| # | Question | Decision |
|---|---|---|
| 1 | Walk past V+5 if all five succeed? | **No.** All six probes fire in parallel at once; no continuation past V+5. |
| 2 | Does `do-update` ask for confirmation? | **No.** The trailing warning line on prior commands serves as notice; invocation is consent. |
| 3 | JSON path when no DB exists? | `~/.<CliName>/data/UpdateChecker.json` |
| 4 | `--force` flag on `update-check`? | **Yes.** Bypasses the interval gate. |
| 5 | Pre-command hook opt-out? | **Yes**, via `BackgroundUpdateCheckEnabled` (default `true`) in `06-Seedable-Config`. |
| 6 | How to detect a newer repo? | Probe `repo-v{N+1}…N+5` suffixes; status JSON also carries optional `NewRepoUrl`. |
| 7 | `--persist` flag? | **Removed.** Both sync and `--async` modes persist. |

---

*Update Check Mechanism Overview — v1.0.0 — 2026-04-20*
