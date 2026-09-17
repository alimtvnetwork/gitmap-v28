# What to Read

> Canonical map of what the AI must read before working on this project.
> Last updated: 2026-09-13T16:15:00Z

## Changelog

- 2026-09-17T22:05:00Z, Memory write: Git commit history audit (last 10 commits), Coding Guideline 24 OS test isolation, hermetic mock runner decoupling, and full repository context ingestion.
- 2026-09-13T20:25:00Z, Memory write: Linux archive package installation (gitmap install tar <archive>), intelligent package strategy detection (binary, script, source, single gz), step-by-step progress display, and dual-database uninstaller integration.
- 2026-09-13T20:20:00Z, Memory write: Official Google Cloud Storage Antigravity artifacts, Linux desktop launcher and broken symlink cleanup, universal uninstallation engine, and stack trace preservation.
- 2026-09-13T16:15:00Z, Memory write: Fast cached exploration via 17-fast-file-reader.py, mandatory 30-commit git history audit, recent 20-task register in plans/01-index.md, and prompt 2.2.0 standards.
- 2026-09-13T16:05:00Z, Memory write: Fixed installer dry-run and smoke test paths (migrated legacy gitmap/ to cli/), added local filesystem fallbacks in install scripts, bumped minor version to v6.227.0, and cut release.
- 2026-09-09T05:07:00Z, Memory write: Interactive macro builder PWD header, in-builder ls/find/search/replace helpers, offline mock installer server, and parallel runner multi-agent section concurrency.
- 2026-09-05T01:46:30Z, Memory write: Macro execution failure RCA, open command cross-platform shim requirement, and Chrome profile reconcile sync.
- 2026-09-05T01:43:00Z, Recorded pending task: Macro step execution failure on 'open chrome' (E9000:EXECUTION).
- 2026-09-04T17:41:00Z, Memory write: Chrome profile picker registration, Local State 13-attribute schema, process concurrency protection, and orphan profile reconciliation.
- 2026-09-04T17:39:00Z, Memory write: parallel multi-worker CI/CD local runner, selective log filtering, streamwriter contracts, and naming standards.
- 2026-08-09T18:21:37Z, Memory write: code red refactor and strict absolute path avoidance.

## Before any task (always)

- `version.json`, why: single source of truth for the repository version, backend/frontend sections, and sub-package version tracks. All codebases must import this file for version information.
- `.lovable/memory/01-index.md`, why: core memory index
- `.lovable/memory/learned/01-project-context-and-guidelines.md`, why: canonical learned memory of repo identity, CODE RED rules, coding guidelines, error philosophy, and active plans
- `.lovable/memory/learned/03-parallel-cicd-runner-and-log-filtering.md`, why: parallel local runner concurrency, duration tracking, and log suppression standard
- `.lovable/memory/learned/04-streamwriter-contracts-and-naming-standards.md`, why: streamwriter contracts, reentrant locker, monadic Bytes[T], JsonResult multi-source ingestion, boolean prefixes, and Id naming standard
- `.lovable/memory/learned/05-chrome-profile-picker-registration-and-concurrency.md`, why: Chrome profile picker visibility contracts, Local State 13-attribute schema, Preferences sanitization, and reconcile engine
- `.lovable/memory/learned/06-macro-step-execution-and-shell-open.md`, why: macro execution engine, Windows 'open' failure analysis, and cross-platform command shims
- `.lovable/memory/learned/10-interactive-macro-builder-pwd-ls-commands.md`, why: interactive macro builder PWD header display, in-builder ls/find/search/replace helpers, and offline mock release testing architecture
- `.lovable/memory/learned/24-installer-paths-seed-urls-and-v6-227-0-release.md`, why: elimination of legacy gitmap/ paths, installer seed fallbacks, and release dryrun contracts
- `.lovable/memory/learned/25-fast-file-reader-and-30-commit-audit-workflow.md`, why: fast cached exploration via 17-fast-file-reader.py, mandatory 30-commit audit, and prompt 2.2.0 standards
- `.lovable/memory/learned/26-official-antigravity-artifacts-and-universal-uninstall.md`, why: official GCS Antigravity artifacts, Linux desktop cleanup, universal uninstall engine, and stack traces
- `.lovable/memory/learned/27-linux-archive-installer-intelligent-strategy.md`, why: Linux archive package installer, intelligent strategy detection, step-by-step progress, and uninstaller integration
- `.lovable/memory/learned/28-git-commit-history-os-isolation-and-context-ingestion.md`, why: recent 10 commits, Coding Guideline 24 OS test isolation, hermetic mock runner decoupling, and context ingestion


- `03-ai-scripts/01-index.md`, why: local automation tools and CI/CD parallel runner specifications
- `.lovable/memory/standards/version-source-of-truth.md`, why: mandatory standard for version.json single source of truth, 'inherit' keyword for sub-packages, and release sync workflow
- `.lovable/coding-guidelines.md`, why: baseline rules and coding standards
- `.lovable/plans/01-index.md`, why: active roadmap, pending tasks, and the Recent Completed Tasks Register (last 20 tasks)
- `.lovable/strictly-avoid.md`, why: hard constraints and anti-patterns
- `.lovable/ambiguous-questions/01-new-ambiguity/`, why: open questions

## Before writing code

- `spec/`, why: understand feature specifications

## Before adding a feature

- `spec/`, why: ensure it fits within existing specs

## Before writing a spec

- `02-spec/01-spec-authoring-guide/`, why: follow authoring format

## Before adding a unit test

- `02-spec/02-coding-guidelines/`, why: testing conventions

## See also

- Root `readme.md` (must stay in sync with this file)
- `.lovable/plans/completed/09-chrome-profile-management-picker-and-token-vault.md`
