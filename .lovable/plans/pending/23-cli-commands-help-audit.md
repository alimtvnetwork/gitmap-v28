# Plan 23: CLI Commands, Help Text Parity & Help UI Architecture Audit

## Overview

Comprehensive audit and verification of all CLI entry points, commands, subcommands, flags, options, and help text descriptions across the repository, ensuring 100% discoverability, standard UI formatting, and concrete usage examples.

---

## Phase 1: Command-to-Help Parity Ledger

| Command / Script | Implemented Subcommands | Registered in Help UI? | Flag Coverage % | Missing Help Text / Examples | Planned Fix | Status |
|---|---|:---:|:---:|---|---|:---:|
| `gitmap scan` | `scan`, `s` | ✅ YES | 100% | Complete | Verified helptext/scan.md | DONE |
| `gitmap clone` | `clone`, `cl` | ✅ YES | 100% | Complete | Verified helptext/clone.md | DONE |
| `gitmap push` | `push`, `ps` | ✅ YES | 100% | Complete | Verified helptext/push.md | DONE |
| `gitmap pull` | `pull`, `pl` | ✅ YES | 100% | Complete | Verified helptext/pull.md | DONE |
| `gitmap commit-in` | `commit-in`, `cin` | ✅ YES | 100% | Complete | Verified helptext/commit-in.md | DONE |
| `gitmap ssh` | `create`, `list`, `status`, `copy`, `cat`, `delete`, `config`, `join`, `login`, `alias`, `exec` | ✅ YES | 100% | Complete | Documented in root usage & ssh.md | DONE |
| `gitmap vscode` | `add`, `rm`, `ls`, `pap`, `plugins` | ✅ YES | 100% | Complete | Verified helptext/vscode.md | DONE |
| `gitmap agy` | `add`, `rm`, `ls`, `stats`, `update` | ✅ YES | 100% | Complete | Verified helptext/agy.md | DONE |
| `gitmap llm` | `llm` (`--url` support) | ✅ YES | 100% | Complete | Verified helptext/llm.md & raw URL | DONE |
| `gitmap doctor` | `doctor`, `doc` | ✅ YES | 100% | Complete | Verified helptext/doctor.md | DONE |
| `gitmap schedule` | `schedule`, `sch` | ✅ YES | 100% | Complete | Verified helptext/schedule.md | DONE |
| `gitmap tree` | `tree`, `tr` | ✅ YES | 100% | Complete | Verified helptext/tree.md | DONE |

---

## Subtasks Breakdown

- [x] [01-help-auditor-and-ledger.md](.lovable/plans/subtasks/23-cli-commands-help/01-help-auditor-and-ledger.md) — Create and verify `06-cli-help-auditor.py` to scan all primary commands.
- [x] [02-root-help-formatting.md](.lovable/plans/subtasks/23-cli-commands-help/02-root-help-formatting.md) — Compact column formatting (column 30) and 4-space multiline indent in `rootusage.go`.
- [x] [03-subcommand-coverage.md](.lovable/plans/subtasks/23-cli-commands-help/03-subcommand-coverage.md) — Verify full subcommand dispatching for SSH, VS Code, AGY, and Cluster.
- [x] [04-verification.md](.lovable/plans/subtasks/23-cli-commands-help/04-verification.md) — Run all CLI parity test suites, linters, and quality gates.
