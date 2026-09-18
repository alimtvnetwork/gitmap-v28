# Generic CLI Creation Guidelines — Overview

> **Related specs:**
> - [02-project-structure.md](02-project-structure.md) — package layout and file organization
> - [03-subcommand-architecture.md](03-subcommand-architecture.md) — dispatch pattern and entry point
> - [13-checklist.md](13-checklist.md) — phased implementation plan referencing all specs
> - [20-terminal-output-design.md](20-terminal-output-design.md) — terminal rendering architecture

## Purpose

This specification is a **complete, self-contained blueprint** for
building production-quality CLI tools. Hand it to any AI assistant
or developer and they can implement a well-structured CLI from scratch.

These guidelines are language-agnostic in principle but use Go for
concrete examples. Adapt syntax to your target language.

---

## Design Philosophy

| Principle | Detail |
|-----------|--------|
| Consistency over cleverness | Predictable patterns across all commands |
| Convention over configuration | Sensible defaults; config is optional |
| Fail fast, fail clearly | Bad input → immediate error with actionable message |
| One responsibility per unit | Each file, function, and package does one thing |
| No magic strings | Every literal in a constants package |
| Self-documenting | Help text, version, and examples built into the binary |

---

## Document Index

| # | File | Topic |
|---|------|-------|
| 01 | [01-overview.md](01-overview.md) | This document — philosophy, scope, index |
| 02 | [02-project-structure.md](02-project-structure.md) | Package layout, file organization, naming |
| 03 | [03-subcommand-architecture.md](03-subcommand-architecture.md) | Routing, dispatch, handler pattern |
| 04 | [04-flag-parsing.md](04-flag-parsing.md) | Per-command flags, defaults, validation |
| 05 | [05-configuration.md](05-configuration.md) | Three-layer config (defaults → file → flags) |
| 06 | [06-output-formatting.md](06-output-formatting.md) | Terminal, CSV, JSON, Markdown, scripts |
| 07 | [07-error-handling.md](07-error-handling.md) | Exit codes, error messages, batch errors |
| 08 | [08-code-style.md](08-code-style.md) | Function length, file length, naming, conditionals |
| 09 | [09-help-system.md](09-help-system.md) | Embedded help files, `--help` interception |
| 10 | [10-database.md](10-database.md) | Local persistence, schema, upsert patterns |
| 11 | [11-build-deploy.md](11-build-deploy.md) | Build scripts, deploy, self-update |
| 12 | [12-testing.md](12-testing.md) | Test structure, conventions, coverage |
| 13 | [13-checklist.md](13-checklist.md) | Step-by-step implementation checklist for AI |
| 14 | [14-date-formatting.md](14-date-formatting.md) | Centralized date display format |
| 15 | [15-constants-reference.md](15-constants-reference.md) | Every constant category with naming patterns |
| 16 | [16-verbose-logging.md](16-verbose-logging.md) | Verbose/debug logging with `--verbose` flag |
| 17 | [17-progress-tracking.md](17-progress-tracking.md) | Progress reporting for batch operations |
| 18 | [18-batch-execution.md](18-batch-execution.md) | Exec command for running commands across repos |
| 19 | [19-shell-completion.md](19-shell-completion.md) | Tab-completion for PowerShell, Bash, Zsh |
| 20 | [20-terminal-output-design.md](20-terminal-output-design.md) | Rich terminal report formatting and color system |

---

## How to Use This Spec

1. **Start with `13-checklist.md`** — it gives a sequenced plan.
2. **Reference individual docs** as you implement each layer.
3. **Every code example is a pattern** — adapt names, not structure.
4. **All constraints are mandatory** unless explicitly marked optional.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
