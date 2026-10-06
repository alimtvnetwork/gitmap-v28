# What to Read

> Canonical map of what the AI must read before working on this project.
> Last updated: 2026-10-06T22:20:00Z

## Changelog

- 2026-10-06T22:15:00Z, Memory write: Plan 235 / Spec 235, Application Specifications & Completed Plans Consolidation & Memory Reduction, synthesized 8 Canonical Domain Clusters in `02-spec/21-app/`, consolidated 163 completed plans into Milestones 28–42 in `.ai-memory/plans/completed/`, folded 177 completed subtasks across 32 directories, compacted `.ai-memory/memory/` into 36 dense reference files, and isolated all pending/active files.
- 2026-10-06T18:00:00Z, Memory write: Plan 230 / Spec 230, Security Token Purge, Installer Navigation ($work, $def), Pull Auto-Remediation, AGM Linux Update quiet mode and failure resolver, Multi-Instance API (JSON/REST), and Settings UI Modernization.
- 2026-10-02T10:25:00Z, Memory write: Plan 65 completed, Spec 65, GitMap nodes table reordering, update-all-zip SCP distribution, search clean cache invalidation, prompt search primacy, repo deduplication with dual solutions, and gitignore policy.
- 2026-10-02T10:15:00Z, Memory write: Plan 66 completed, Issue 62 (RCA), Pull-All duplicate repository elimination, database case-insensitive collation (`COLLATE NOCASE`), and pull pipeline in-memory deduplication.
- 2026-10-01T15:45:00Z, Memory write: RCA-26 / Spec 60, Avoid 06, Non-idempotent PowerShell macro removal prevention, platform-adaptive shims (Test-Path), and multi-user binary synchronization.
- 2026-10-01T15:35:00Z, Memory write: Plan 60 completed with all 7 subtasks verified, CPAR suite, Split-DB cache engine, and ignore grouping suite.
- 2026-10-01T15:30:00Z, Memory write: Spec 201, Plan 64 (completed), Pull Concurrency Throttling, Split-DB GitIgnore Cache Engine, Subprocess Deadlock Prevention, and Single-Hand SSH Delegation.
- 2026-10-01T14:45:00Z, Memory write: Spec 200, Command 13, Flat Git Commit Workflow (`gitmap commit`, `gitmap cm`), Auto-Staging, `--push` and `--dry-run` Flags, and Test Suite Validation.
- 2026-10-01T14:10:00Z, Memory write: Spec 199, Plan 62 (completed), Devtools Cache Dynamic Discovery, Tree View Rendering & Split-DB Persistence.
- 2026-10-01T13:45:00Z, Memory write: Spec 198, Plan 61 (completed), PAS Worker Concurrency, Pull-Error Split-DB Subsystem, Machine Telemetry, Heartbeat Progress, and Wincredman Remediation.
- 2026-09-19T09:30:00Z, Memory write: Consolidated 142 completed plans and 72 subtasks down to 27 dense milestones under backup branch `backup/plans-consolidation-20260919-093052`.

## Before any task (always)

- `version.json`, why: single source of truth for the repository version, backend/frontend sections, and sub-package version tracks. All codebases must import this file for version information.
- `.ai-memory/what-to-read.md`, why: canonical map of project state, invariants, and navigation instructions.
- `.ai-memory/memory/learned/01-project-context-and-standards.md`, why: canonical learned memory of repo identity, CODE RED rules, coding guidelines, error philosophy, and active plans.
- `.ai-memory/memory/learned/02-cli-contracts-and-help-architecture.md`, why: Cobra CLI hierarchy, command conventions, middle-ellipsized tables, and JSON envelope standards.
- `.ai-memory/memory/learned/03-scanner-deduplication-and-path-collation.md`, why: scanner deduplication, path collation, `.gitmapignore` caching, and OS sensitivity.
- `.ai-memory/memory/learned/04-git-operations-flat-commit-and-pull.md`, why: pull-all concurrency, semantic flat commit suite, push recovery, and staging safeguards.
- `.ai-memory/memory/learned/05-ssh-fleet-ping-and-remote-exec.md`, why: unified fleet nodes, reachability probing, remote SSH delegation, and multi-node execution.
- `.ai-memory/memory/learned/06-rsa-credential-vault-and-auth.md`, why: RSA-OAEP salt credential encryption vault, askpass mechanics, and password interception safety.
- `.ai-memory/memory/learned/07-antigravity-integration-and-ide-sync.md`, why: Antigravity SDK workflows, multi-conversation prompts, IDE synchronization, and profile migrations.
- `.ai-memory/memory/learned/08-sqlite-split-db-and-cache-lifecycle.md`, why: three-tier SQLite Split-DB engine (`gitmap.db`, `installation.db`, `repodb/pipeline.db`), WAL mode, and connection pooling.
- `.ai-memory/memory/learned/09-pipeline-telemetry-and-rca-remediation.md`, why: pipeline diagnostics (`pe`/`pea`), traceback extractors, and automated 4-part RCA engine.
- `.ai-memory/memory/learned/10-cross-platform-installers-and-runners.md`, why: NSIS Windows installers, Linux archives, cross-platform runners (`run.ps1`/`run.sh`), and release ceremony.
- `.ai-memory/memory/learned/14-coding-guidelines-and-positive-booleans.md`, why: positive boolean conventions (`is*`, `has*`), 8-to-15 line functions, and structured errors.
- `.ai-memory/memory/avoid/01-test-and-build-execution-guards.md`, why: strict ban on heavy build or full test executions during routine agent turns.
- `.ai-memory/memory/avoid/04-boolean-naming-and-error-swallowing.md`, why: prohibition against double negatives, negative boolean flags, and swallowed errors.
- `.ai-memory/memory/avoid/05-hardcoded-paths-and-secrets-exposure.md`, why: total ban on absolute filesystem paths, drive letters, and committed secrets or refresh tokens.
- `.ai-memory/coding-guidelines.md`, why: baseline rules and coding standards.
- `.ai-memory/plans/readme.md`, why: active roadmap, pending tasks, and the 42 consolidated milestone registries.

## Before writing code

- `02-spec/21-app/readme.md`, why: explore authoritative domain clusters:
  - `02-spec/21-app/01-cli-architecture/` (CLI & Shell UX)
  - `02-spec/21-app/02-scanner-and-projects/` (Scanner & AUM)
  - `02-spec/21-app/03-git-operations-and-pull/` (Git & Pull)
  - `02-spec/21-app/04-fleet-nodes-and-ssh/` (Fleet & SSH Vault)
  - `02-spec/21-app/05-antigravity-and-ide/` (AGY & IDEs)
  - `02-spec/21-app/06-database-and-split-db/` (Split-DB SQLite)
  - `02-spec/21-app/07-pipeline-and-diagnostics/` (CI/CD & RCA)
  - `02-spec/21-app/08-distribution-and-release/` (Installers & Release)

## Before adding a feature

- `02-spec/`, why: ensure it fits within existing specs

## Before writing a spec

- `02-spec/01-spec-authoring-guide/`, why: follow authoring format

## Before adding a unit test

- `02-spec/02-coding-guidelines/`, why: testing conventions

## See also

- Root `readme.md` (must stay in sync with this file)
- `docs/benchmarks/benchmark.md`, why: Native AUM search vs Go walk vs Python grep benchmark report
- `02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md`, why: Master consolidation architecture and reduction ledger
