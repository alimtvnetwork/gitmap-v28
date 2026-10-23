# What to Read

> Canonical map of what the AI must read before working on this project.
> Last updated: 2026-10-08T16:00:00Z

## The Mindset — READ FIRST

How this repository is engineered. Internalize before writing any code:

1. **Staging repo — velocity over ceremony.** AI-driven, fast-forward development: versions ship fast, each a mix of many features. AI makes mistakes; we fix forward in the next version. Never slow the loop for process.
2. **Two-tier future.** This repo stays fast. When satisfied, code moves to a separate release repository with a stricter engineering process. Strictness belongs there, not here.
3. **Version churn is intentional.** Frequent versions are checkpoints, not milestones. Do not "fix" the cadence.
4. **Backup-first.** Significant work starts from a backup branch (`backup/<slug>`). `gitmap space backup-branch "<task>"` creates one from current code.
5. **Small files are engineering quality — ~300 lines max per file.** ~100-line median files, tiny functions, zero nesting. Split by concern into ~300-line files; never fragment into tiny shards (9 files of ~1000 lines → ~30 files, not 80+). Enforce the cap: no file crosses 300 lines without a recorded reason.
6. **Aliases are a feature, not sprawl.** Keep them — they make the CLI fast. Never prune for tidiness.
7. **DRY everything.** Help text, themes, colors: single source of truth, generated output. Never hand-maintain the same content in two formats.
8. **Enforcement beats documentation.** A rule on paper that isn't enforced is worse than no rule. Where enforcement lagged, close the gap.
9. **Interface + struct + binding.** Extensibility pattern: small interfaces, concrete structs, bound per the coding guidelines.

## Changelog

- 2026-10-23T21:40:00Z, Memory write: Plan 315 / Spec 21-app/315-gitmap-commit-tracking-and-database-integration, GitMap Commit Tracking, Universal SHA Resolution, External Commit Split-DB & AI Provenance Integration. Universal commit inspection (`gitmap commit info <sha>`, `gitmap ci <sha>`) resolving full, short, prefix, suffix, and substring SHAs across 4 data layers (Git repo, external Commit Split-DB `<EffectiveDataDir>/repodb/commit_history.db`, Heatmap DBs, and Agent task DBs). AI provenance commit flags (`--agent-id`, `--heatmap-id`, `--task`, `--spec`, `--ai`) across all commit commands (`gitmap c`, `cm`, `cpf`, `cpb`, `cpr`, `cpc`, `pcp`) with ambient fallbacks, rich ASCII tree view rendering, foreign AI attribution disclaimer, and V6 execution prompt & skill synchronization.

- 2026-10-23T20:45:00Z, Memory write: Plan 314 / Spec 21-app/314-gitmap-heatmap-and-url-handling-instructions, Heatmap Database Discovery & AI Help, URL Split-DB Caching, 10-Split Segmented Downloader & LLM Training Integration. Comprehensive 3-tier heatmap discovery (`<EffectiveDataDir>/heatmap.db`, per-repo split DBs, `repo-cache/.heatmap/<seq>-<ai>/`), alias enrichment (`db.ListAliasesWithRepo()`), `--repo` routing fix, multi-repo recent tasks aggregation across 2–5 repositories, dedicated URL caching Split-DB (`repodb/url_cache.db`), CLI `gitmap url read` with byte-faithful `output.Raw()` streaming, multi-tier 10-split segmented downloading (`aria2c -s 10 -x 10 -k 5K` + native Go parallel Range requests for payloads >= 1–2 MB), and LLM training integration linking heatmap companion files into training datasets.

- 2026-10-21T15:25:00Z, Memory write: Plan gitmap-enhancement-and-issue-resolution / Spec 21-app/gitmap-enhancement-and-issue-resolution, Database Read Introspection & 5-10s Metadata Cache, Visualizer UI Binary Cell Formatting & Drawer, Storytelling-v1 README Casing Remediation & Repo-Cache Recipe, Remote Node Clone GitHub Desktop & VS Code Registration, Remote SSH WinCredMan Resolution, and Machine Set Interactive Prompt.

- 2026-10-21T11:15:00Z, Memory write: Plan gitmap-command-and-skill-update / Spec 21-app/gitmap-command-and-skill-update, OS Dev-Clean Dry-Run & Cargo Caches: Always show ANSI dry-run preview and total reclaimable metrics before interactive [y/N] confirmation. Comprehensive Rust and Cargo cache sweeping (~/.cargo/registry/cache, registry/src, git/db, git/checkouts, .package-cache, target, CARGO_TARGET_DIR, sccache, ~/.rustup/downloads, ~/.rustup/tmp), tool clean invocations (cargo cache --autoclean, sccache --clear-cache, cargo clean --target-dir), and Split-DB devtools cache synchronization with missing ecosystem auto-discovery and post-clean zeroing.
- 2026-10-17T10:35:00Z, Memory write: Session programs 271–275. 271: nc auto-discovery reverted as broken (coordinator false claim caught), v6.552.0. 272: owner approved items 1–3 only (4–6 explicitly NOT approved — hard stop); deprecation policy removed, coverage floor gate, nc codegen; v6.553.0. 273: 4 design-system themes, `help example`, `mr` command, 20 proposals, Windows virus RCA. 274: CI 44/44 green, v6.558.0, 10+20 proposals, 2 E2E bugs documented (NOT fixed — awaiting owner). 275: 30 commands implemented, E2E 30/30, v6.559.0; rogue worker scope creep caught via git show --stat. Session record: `.ai-memory/memory/09-session-programs-271-275-2026-10-17.md`.
- 2026-10-17T10:30:00Z, Memory write: Muse Memory IO follow-ups F1–F5 / Spec 277. `gitmap llm train muse` now syncs the `alimtvnetwork/repo-cache` vault (repo-cache-first, `rc` special-repo key; repo-secrets as fallback), digest at `12-muse-ai/` (NN-sequence), last-read watermark tracked by gitmap in BinaryDataDir (0600 JSON, never an MD file), numbered DECISION TRACE in terminal output. Digest moved out of repo-secrets (secrets stay). New-muse onboarding: 9-step report-back checklist at `repo-cache/12-muse-ai/00-new-muse-onboarding.md` (install → update → login → `gitmap clone` → `gitmap pa` → master prompt → `llm train muse` → email → `gitmap scan` → final report).
- 2026-10-16T13:30:00Z, Memory write: Plan gitmap-command-enhancements / Spec 21-app/gitmap-command-enhancements, OS JSON Output (`gitmap os --json` / `-j`), Active Host Identity Footer, Fast Hardware Telemetry (`gitmap cpu`, `processor`, `ram`, `disk`) with 5-second cross-process cache protection, Network Adapter Metrics (`gitmap network`), and Fleet Nodes Task Queue & Delegation (`gitmap nodes tasks`).
- 2026-10-16T12:00:00Z, Memory write: Plan gitmap-enhancement-and-optimization / Spec 21-app/gitmap-enhancement-and-optimization, Resilient autostash pull-first execution with interactive post-pull remediation prompting, ANSI progress bar color fix on line clamp, PAE 48h commit recency window, multi-node CPU/RAM monitoring with dynamic pull concurrency boost, and semantic change grouping with markdown links for 'gitmap sends cp all'.
- 2026-10-10T03:20:00Z, Memory write: RCA-30 / Avoidance update, Multi-User Windows Binary Synchronization across all 4 target paths (`C:\Users\Alim\AppData\Local\gitmap\`, `AppData\Local\gitmap-cli\`, `bin\`), interactive macro step recording polish, and elimination of inverted boolean checks (`!isSuccess`) in favor of direct `.IsFail()` / `.is_fail` checks.
- 2026-10-08T16:00:00Z, Memory write: Added "The Mindset — READ FIRST": staging-repo velocity model, backup-first workflow, small-files enforcement, aliases-as-feature, DRY, enforcement-over-documentation, interface+struct+binding.

- 2026-10-06T22:55:00Z, Memory write: Plan 236 / Spec 236, Deep Spec Consolidation & Canonical Reduction, pruned 72 legacy folders and 179 loose files across `02-spec/21-app/` down to exactly 8 Canonical Clusters + active spec folder + readme (-96.2%), compacted completed plans into 18 authoritative milestones (-58.1% files, -81.1% lines), unified memory into 10 canonical domain references (-86.7%), and preserved all pending/active files.
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
- `.ai-memory/memory/00-project-governance-and-invariants.md`, why: canonical governance, positive booleans (`is*`, `has*`), relative paths, zero-build rules, and git mutation boundaries.
- `.ai-memory/memory/01-cli-architecture-and-contracts.md`, why: Cobra CLI hierarchy, command conventions, ANSI tables, and typed JSON Envelope V2 standards.
- `.ai-memory/memory/02-scanner-projects-and-deduplication.md`, why: scanner deduplication, path collation, `.gitmapignore` caching, and OS sensitivity.
- `.ai-memory/memory/03-git-operations-commit-and-pull.md`, why: pull-all concurrency ($\min(\text{cores}, 8)$), semantic flat commit suite (`gitmap c`), push recovery, and staging safeguards.
- `.ai-memory/memory/04-fleet-nodes-ssh-and-credentials.md`, why: unified fleet nodes, reachability probing, remote SSH delegation, and multi-node execution.
- `.ai-memory/memory/05-antigravity-and-ide-ecosystem.md`, why: Antigravity SDK workflows, multi-conversation prompts, IDE synchronization, and profile migrations.
- `.ai-memory/memory/06-database-engine-and-split-storage.md`, why: three-tier SQLite Split-DB engine (`gitmap.db`, `installation.db`, `repodb/pipeline.db`), WAL mode, and connection pooling.
- `.ai-memory/memory/07-pipeline-diagnostics-and-telemetry.md`, why: pipeline diagnostics (`pe`/`pea`), traceback extractors, failure heatmaps, and automated 4-part RCA engine.
- `.ai-memory/memory/08-distribution-installers-and-release.md`, why: NSIS Windows installers, Linux archives, cross-platform runners (`run.ps1`/`run.sh`), and release ceremony.
- `.ai-memory/coding-guidelines.md`, why: baseline rules and coding standards.
- `.ai-memory/plans/readme.md`, why: active roadmap, pending tasks, and the 18 consolidated milestone registries.

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
- `02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md`, why: Master consolidation architecture and reduction ledger
