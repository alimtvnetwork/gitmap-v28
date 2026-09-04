# What to Read

> Canonical map of what the AI must read before working on this project.
> Last updated: 2026-09-05T01:46:30Z

## Changelog

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
- `.lovable/memory/standards/version-source-of-truth.md`, why: mandatory standard for version.json single source of truth, 'inherit' keyword for sub-packages, and release sync workflow
- `.lovable/coding-guidelines.md`, why: baseline rules and coding standards
- `.lovable/plans/01-index.md`, why: active roadmap and pending tasks
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
- `.lovable/plans/pending/01-macro-step-open-chrome-failure.md`
- `.lovable/plans/completed/63-chrome-profile-picker-visibility.md`
