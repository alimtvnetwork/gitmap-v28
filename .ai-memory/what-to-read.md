# What to Read

> Canonical map of what the AI must read before working on this project.
> Last updated: 2026-09-13T16:15:00Z

## Changelog

- 2026-09-20T12:20:00Z, Memory write: Plan 46 Pipeline-AI live error streaming, early remediation command switching, AUM search extension-first optimization, vcvarsall fast tool locator (gitmap aum locate), and MD ALIM UL KARIM author attribution.
- 2026-09-20T11:55:00Z, Memory write: Plan 45 Composite Test Heatmap Engine, unified .ai-memory/cicd/cache/ architecture, repo_delta fast-path short-circuit, and fsync lock contention elimination.
- 2026-09-19T09:30:00Z, Memory write: Consolidated 142 completed plans and 72 subtasks down to 27 dense milestones (87.8% file count reduction) under backup branch backup/plans-consolidation-20260919-093052.
- 2026-09-19T02:00:00Z, Memory write: Plans 208–214 SSH multi-target command resolution and machine join, and CI/CD pipeline bounded stack traces (RCA 58).
- 2026-09-18T02:35:00Z, Memory write: Plan 194 parallel workflow/section log download for single commits, previous pipeline run/DB fallback, and two-pass non-mutating line execution and filtering.
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
- `.ai-memory/memory/01-index.md`, why: core memory index
- `.ai-memory/memory/learned/01-project-context-and-guidelines.md`, why: canonical learned memory of repo identity, CODE RED rules, coding guidelines, error philosophy, and active plans
- `.ai-memory/memory/learned/03-parallel-cicd-runner-and-log-filtering.md`, why: parallel local runner concurrency, duration tracking, and log suppression standard
- `.ai-memory/memory/learned/04-streamwriter-contracts-and-naming-standards.md`, why: streamwriter contracts, reentrant locker, monadic Bytes[T], JsonResult multi-source ingestion, boolean prefixes, and Id naming standard
- `.ai-memory/memory/learned/05-chrome-profile-picker-registration-and-concurrency.md`, why: Chrome profile picker visibility contracts, Local State 13-attribute schema, Preferences sanitization, and reconcile engine
- `.ai-memory/memory/learned/06-macro-step-execution-and-shell-open.md`, why: macro execution engine, Windows 'open' failure analysis, and cross-platform command shims
- `.ai-memory/memory/learned/10-interactive-macro-builder-pwd-ls-commands.md`, why: interactive macro builder PWD header display, in-builder ls/find/search/replace helpers, and offline mock release testing architecture
- `.ai-memory/memory/learned/24-installer-paths-seed-urls-and-v6-227-0-release.md`, why: elimination of legacy gitmap/ paths, installer seed fallbacks, and release dryrun contracts
- `.ai-memory/memory/learned/25-fast-file-reader-and-30-commit-audit-workflow.md`, why: fast cached exploration via 17-fast-file-reader.py, mandatory 30-commit audit, and prompt 2.2.0 standards
- `.ai-memory/memory/learned/26-official-antigravity-artifacts-and-universal-uninstall.md`, why: official GCS Antigravity artifacts, Linux desktop cleanup, universal uninstall engine, and stack traces
- `.ai-memory/memory/learned/27-linux-archive-installer-intelligent-strategy.md`, why: Linux archive package installer, intelligent strategy detection, step-by-step progress, and uninstaller integration
- `.ai-memory/memory/learned/28-git-commit-history-os-isolation-and-context-ingestion.md`, why: recent 10 commits, Coding Guideline 24 OS test isolation, hermetic mock runner decoupling, and context ingestion
- `.ai-memory/memory/learned/29-ssh-multi-target-pipeline-bounded-stacktrace-ingestion.md`, why: Plans 208–214 SSH multi-target command resolution and machine join, and CI/CD pipeline bounded stack traces (RCA 58)


- `03-ai-scripts/01-index.md`, why: local automation tools and CI/CD parallel runner specifications
- `.ai-memory/memory/standards/version-source-of-truth.md`, why: mandatory standard for version.json single source of truth, 'inherit' keyword for sub-packages, and release sync workflow
- `.ai-memory/coding-guidelines.md`, why: baseline rules and coding standards
- `.ai-memory/plans/01-index.md`, why: active roadmap, pending tasks, and the Recent Completed Tasks Register (last 20 tasks)
- `.ai-memory/strictly-avoid.md`, why: hard constraints and anti-patterns
- `.ai-memory/ambiguous-questions/01-new-ambiguity/`, why: open questions

## Before writing code

- `02-spec/`, why: understand feature specifications

## Before adding a feature

- `02-spec/`, why: ensure it fits within existing specs

## Before writing a spec

- `02-spec/01-spec-authoring-guide/`, why: follow authoring format

## Before adding a unit test

- `02-spec/02-coding-guidelines/`, why: testing conventions

## See also

- Root `readme.md` (must stay in sync with this file)
- `.ai-memory/plans/completed/09-chrome-profile-management-picker-and-token-vault.md`
- `02-spec/21-app/129-pr-commit-engines-and-sqlite-split-db.md`, why: PR commit replay engine, SQLite split-DB standardization, and final snapshot sync
- `cli/helptext/pr.md`, why: PR command family usage, HG help, and JSON examples
