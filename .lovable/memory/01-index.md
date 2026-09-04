# Memory: index.md

Updated: 2026-04-27

# Project Memory

## Core

- 🔴 CODE RED: Never swallow errors. Zero-nesting (no nested if). Max 2 operands. Positively named guard functions.
- 🔴 CODE RED: Strict metrics: functions 8-15 lines, files < 300 lines, React components < 100 lines.
- 🔴 CODE RED: Standalone scripts — NO `!important`, NO `as unknown` / `as any`, NO error-swallowing `return null` in catch, NO magic strings (use enums in `types.ts`), NO bare top-level functions (class-based, DI). Hide via class toggle + CSS transition — never nested `requestAnimationFrame`. Always blank line before `return`.
- 🔴 PRE-WRITE: Before producing any standalone-script code, read the target repo's coding-guideline spec AND the sibling `standalone-scripts/*/src/` files first. See mem://issues/payment-banner-hider-rca.
- 🔴 README install commands: ONE line per command, NEVER `#` comments inside install code fences. Platform shown by header above the block. Install section comes FIRST after badges, before TOC. Match on-site `InstallSection.tsx` order (4 full-repo + 7 named bundles).
- 🔴 Release & Migration UI: exactly TWO cards (Windows PowerShell + macOS/Linux Bash), one-liner each. NEVER add "skip latest probe" variants. User locked this 2026-04-24.
- 🔴 Canonical repo slug: `alimtvnetwork/coding-guidelines-v24`. Any v14/v15 reference is a bug — run repo-wide `grep -rn` after every rebrand.
- 🔴 NEVER sync 01-app, 02-app-issues, 03-general, 03-tasks, or 12-consolidated-guidelines from gitmap-v3. All maintained locally.
- 🔴 CODE RED: Never commit test results, test reports, `.test-report.*` files, compiled binaries, or temporary data to the repo. They MUST be excluded via `.gitignore`. Test outputs cause repository bloat and leak architectural/environment details. Always check `git status` before `git add .` to ensure no generated artifacts are included.
- 🔴 STRICTLY PROHIBITED #01 — readme.txt timestamp/date/time generator: TOTAL BAN. Never build, suggest, propose, spec, document, or ask about any feature that writes a timestamp, date, time, "Malaysia-formatted" content, prefix phrase, or auto-refresh into readme.txt — no scripts, no sub-commands, no npm tasks, no sync hooks, no docs, no instructions, no clarifying questions, no follow-up offers. Only acceptable readme.txt action: one-shot manual edit of the exact text the user types that turn. See mem://constraints/readme-txt-timestamp-total-ban and `.lovable/strictly-avoid.md`.
- Naming: PascalCase for all internal IDs, DB, JSON, Types. Exceptions: Rust uses snake_case identifiers.
- DB Schema: PascalCase naming. PKs are `{TableName}Id` (INTEGER PRIMARY KEY AUTOINCREMENT). No UUIDs.
- DB Schema Rules 10/11/12: Entity/ref tables need `Description TEXT NULL`; transactional need `Notes`+`Comments TEXT NULL`; all must be nullable, no DEFAULT. Join tables exempt.
- Workflow: Spec-First (`spec/`) and Issue-First (`03-issues/`).
- Global Namespace: Always derive module imports and repository namespaces dynamically from the repository's canonical root configuration or manifest (`package.json`, `go.mod`, `version.json`).
- Version sync: bump package.json → `node scripts/sync-version.mjs` → `node scripts/sync-spec-tree.mjs`.
- Execution: Break complex requests into discrete tasks. Wait for "next" prompt to continue.
- 🔴 NO-QUESTIONS MODE active (40-task batch from 2026-04-26): never call ask_questions; log ambiguities to `.lovable/question-and-ambiguity/xx-title.md` and update its README index. Resume on explicit "ask questions" signal.
- 🔴 Docs Viewer keybinds: `Cmd/Ctrl+K` = search dialog, `Cmd/Ctrl+J` = command palette ("Open Spec Overview" pinned). Do not rebind. See mem://sessions/2026-04-27-docs-viewer-quick-jump-and-sync-banner.
- 🔴 GitHub Sync Banner reads `version.json` (`git.shortSha`, `git.branch`, `updated`). Per-SHA dismiss via `localStorage["lovable.github-sync-banner.dismissed-sha"]` — never make it permanently dismissable.
- 🔴 Docs sidebar tree diagnostics: opt-in via `localStorage["lovable.tree-diagnostics.enabled"]="1"`. Logs prefixed `[tree:<category>]`. Off by default. Never enable by default. The Lovable IDE file tree is platform UI — NOT instrumentable from project code; do not conflate it with the in-app docs sidebar. See mem://sessions/2026-04-27-tree-diagnostics-logging.
- 🔴 Idiomatic `-er` Go interface naming: NEVER suffix Go interfaces with `Interface` (`Writer[T]`, `Streamer[T]`, `Compiler`).
- 🔴 Strict `Id` / `id` naming: TOTAL BAN on uppercase acronym `ID` in struct fields and variables (`UserId`, `OrderId`, `TraceId`).
- 🔴 Strict positive boolean prefixes: ALL boolean fields and variables MUST have positive prefixes (`isActive`, `isSuccess`, `isValid`).
- 🔴 Parallel CI/CD Local Runner: Worker group concurrency via `ThreadPoolExecutor` with selective log filtering (`--failed` suppresses noisy logs on pass; isolates failures).

- [Project Context & Learned Guidelines](learned/01-project-context-and-guidelines.md) — Comprehensive ingestion of repo identity, CODE RED rules, coding guidelines, error philosophy, and active plans.
- [Pluggable Logger & Uber Zap Architecture](learned/02-logger-swapping-and-uber-zap-architecture.md) — Architecture and code samples for dynamic Log Changer, Formatter/Writer pipeline, context.Context tracing, and Uber Zap integration.
- [Parallel CI/CD Runner & Selective Log Filtering](learned/03-parallel-cicd-runner-and-log-filtering.md) — Concurrency workgroups via ThreadPoolExecutor, duration metrics, and clean log suppression.
- [StreamWriter Contracts & Naming Standards](learned/04-streamwriter-contracts-and-naming-standards.md) — Idiomatic -er interfaces, ReentrantMutex, Bytes[T], JsonResult multi-source ingestion, boolean prefixes, and Id naming standard.

- [Fast File Indexing & Caching Strategy](standards/05-fast-file-indexing-and-caching.md) — Pre-computed file scanning and index caching in `tmp/` via `08-fast-file-scanner.py` for rapid multi-step lookups.
- [Prompt Synchronization Architecture](standards/04-prompt-synchronization-architecture.md) — All prompts authored in `01-prompts/` and compiled to flat `01-prompts/*.md` via `scripts/update-prompts.ps1` without external Git clones.
- [Version Source of Truth Standard](standards/version-source-of-truth.md) — Canonical `version.json` standard at repo root.
- [Release Ceremony](preferences/release-ceremony.md) — Triggers "release"/"bump version"/"bump version + add changelog + pin to root readme". Full ceremony: pick bump tier → sync all version pins (version.json, manifest.json, src/shared/constants.ts, readme.md, standalone-scripts SDK + instruction.ts files) → changelog.md entry → pin root readme → aggregate-prompts if prompt sources touched → verify no stale refs. No confirmation, no plan mode.
- [Install Command Formatting](mem://constraints/install-command-formatting) — One-line installs, no inline comments, per-platform headers, mirror InstallSection.tsx order.
- [readme.txt Timestamp TOTAL BAN](mem://constraints/readme-txt-timestamp-total-ban) — No script/02-spec/doc/suggestion of any timestamp generator for readme.txt. Hard prohibition.
- [No readme.txt Time Instructions](mem://constraints/no-readme-time-instructions) — Never produce instructions/docs/how-to-run for the readme.txt Malaysia date-time generator.
- [2026-04-24 Batch Cleanup + Rebrand](mem://sessions/2026-04-24-batch-cleanup-and-rebrand) — Slug rebrand to v16, Release & Migration UI lock, 11 plan items closed (B5/B6/B7/B8/B10/B11/09/10/12/B2 + UI).
- [Blank Line Between If Guards](mem://constraints/blank-line-between-if-guards) — Rule 5 applied to all markdown snippets and source code.
- [SQL Linter Rules](mem://sessions/2026-04-sql-linter-rules) — DB-FREETEXT-001 (presence) + MISSING-DESC-001 (presence+Rule 12+waivers), shared _lib, waiver syntax.
- [Axios Pinning](mem://constraints/axios-version-pinning) — Exact pinned versions only (1.14.0/0.30.3). Blocked versions: 1.14.1, 0.30.4.
- [Database Architecture](mem://architecture/database-schema) — PascalCase naming, no UUIDs, Vw prefixes for views.
- [Error Handling](mem://architecture/error-handling) — 'apperror' package, explicit file/path logging required.
- [PowerShell Style](mem://style/powershell-naming) — lowercase-kebab-case files, PascalCase Verb-Noun functions.
- [Development Workflow](mem://processes/development-workflow) — Spec-first workflow, linter enforcement, clean docs.
- [React ForwardRef Warning](mem://constraints/react-app-forwardref-warning) — Ignore lovable.js App.tsx ref console warning.
- [Code Red Guidelines](mem://standards/code-red-guidelines) — Full rules for zero-nesting, booleans, metrics.
- [Standards Enforcement](mem://processes/automated-standards-enforcement) — linter-scripts validation requirements.
- [Naming Conventions](mem://style/naming-conventions) — Zero-Underscore policy, full uppercase acronyms.
- [Caching Policy](mem://architecture/caching-policy) — Explicit TTL, deterministic keys, invalidate on mutation.
- [Nested Code Fences](mem://issues/nested-code-fence-data-corruption) — 4-backtick fences required for nested markdown blocks.
- [TypeScript Patterns](mem://standards/typescript-patterns) — Named interfaces for unions, TypedAction, explicit types.
- [Enum Standards](mem://standards/enum-standards) — Cross-language PascalCase enums, strict parsing methods.
- [Split Database](mem://architecture/split-database) — Root, App, Session hierarchical SQLite with WAL and Casbin.
- [Seedable Config](mem://architecture/seedable-configuration) — SemVer GORM merge of config.seed.json.
- [Self Update Arch](mem://features/self-update-architecture) — Rename-first deployment, atomicity with latest.json.
- [Doc Standards](mem://project/documentation-standards) — Mandatory numeric folders (01-20 Core, 21+ App), JSON tree syncing.
- [Author Attribution](mem://project/author-attribution) — Md. Alim Ul Karim, Riseup Asia LLC, SEO/footer requirements.
- [Avoid Gitmap Sync](mem://constraints/avoid-app-sync) — NEVER sync app, app-issues, general, tasks, or consolidated-guidelines from gitmap-v3.
- [Install Command Formatting](mem://constraints/install-command-formatting) — README top install area must mirror UI order; one-line commands only; bundles before ToC.
- [Standalone Script Standards](mem://constraints/standalone-script-standards) — Hard rules for browser/userscript files: no !important, no as-unknown, no error swallowing, class+DI, enums in types.ts, styles.ts, hide via class+transition.
- [Payment Banner Hider RCA](mem://issues/payment-banner-hider-rca) — Root cause for the macro-ahk-v23 regression and the mandatory pre-write checklist that prevents repeats.
- [No-Questions Mode](mem://workflow/no-questions-mode) — 40-task no-questions run; ambiguity logging template, sequencing, and resume trigger.
- [Avoid Time Suggestions in readme.txt](mem://avoid/02-no-time-suggestions-in-readme-txt) — Strictly Prohibited #01: zero time-related chatter, suggestions, or auto-updates around readme.txt; persist across sessions.
- [2026-04-27 Docs Viewer Quick-Jump + Sync Banner](mem://sessions/2026-04-27-docs-viewer-quick-jump-and-sync-banner) — Cmd/Ctrl+J command palette with Open-Spec-Overview pinned, GithubSyncBanner reading version.json, .gitmap/ removed, spec/ visibility incident diagnosed.
- [2026-04-27 Tree Diagnostics Logging](mem://sessions/2026-04-27-tree-diagnostics-logging) — Opt-in `[tree:*]` ring-buffer logging, floating TreeDiagnosticsPanel, instrumentation in useSpecData/DocsSidebar/SpecTreeNav. Clarified IDE-tree vs in-app-tree.

## Memories — Backfilled (orphans recovered 2026-04-27)

- [Avoid Per-Task Folders](mem://avoid/01-avoid-per-task-folders) — Folder-shape prohibitions for one-off task files.
- [Skip Stub Spec Folders](mem://constraints/skip-stub-spec-folders) — Do not generate placeholder/stub spec subfolders.
- [Coding-Guidelines Consolidation Plan (DONE)](mem://done/coding-guidelines-consolidation-plan) — Archived plan kept for historical reference.
- [Exclude-Paths Glob Support](mem://features/exclude-paths-glob-support) — `.codeguidelines.toml` exclusion globs (B9).
- [Release-Pinned Installer](mem://features/release-pinned-installer) — Pinned-version install flow.
- [Spec Link Checker](mem://features/spec-link-checker) — Cross-spec link validation rules.
- [Visual Rendering System](mem://features/visual-rendering-system) — Markdown viewer rendering pipeline.
- [Nested Code Fence Rendering](mem://issues/nested-code-fence-rendering) — Companion to nested-code-fence-data-corruption; renderer specifics.
- [README Bundle Installers](mem://project/03-readme-bundle-installers) — Bundle install matrix synced with `InstallSection.tsx`.
- [Naming Compliance Issues](mem://project/naming-compliance-issues) — Tracker for renames + compliance fixes.
- [Phase 2 Content-Overlap Audit](mem://project/phase2-content-overlap-audit) — Spec consolidation audit notes.
- [Phase 3 Consolidated Structure Design](mem://project/phase3-consolidated-structure-design) — Target consolidated structure.
- [v2.2 Error-Spec Changes](mem://project/v2.2-error-spec-changes) — Error spec migration notes.
- [2026-04-19 Distribution Runner + Slides](mem://sessions/2026-04-19-distribution-runner-slides) — Session log.
- [2026-04-19 Perf + Boolean Naming + Schema](mem://sessions/2026-04-19-perf-boolean-naming-schema) — Session log.
- [2026-04-23 Quickstart + Zero Violations](mem://sessions/2026-04-23-quickstart-and-zero-violations) — Session log.
- [2026-04-23 README Code-Red Walkthrough](mem://sessions/2026-04-23-readme-code-red-walkthrough) — Session log.
- [Installer Behavior Standards](mem://standards/installer-behavior) — Canonical installer UX rules.
- [Suggestions Tracker (memory mirror)](mem://suggestions/01-suggestions-tracker) — Mirrors `.lovable/suggestions.md` highlights.
- [Plan Tracker (workflow)](mem://workflow/01-plan-tracker) — Workflow-state mirror of `.lovable/29-plan.md`.
- [Session 2026-08-09 Code Red Refactor](mem://sessions/08-2026-08-09-code-red-refactor) — Enum enforcement and boolean fix.
- [Absolute Paths Ban](mem://specs/01-absolute-paths-ban) — Direct user mandate to use standalone relative paths.
- [Avoid Absolute File System Paths](mem://avoid/03-absolute-file-system-paths) — Do not use file:/// absolute paths.
- [Version Source of Truth](mem://standards/version-source-of-truth) — `version.json` at repo root is the single canonical source of truth for versions across all languages and tools. Self-explaining metadata (`_purpose`, `_instructions`) embedded directly in the JSON. Changing version is done solely in `version.json` and propagated via `npm run sync`.
- [Release Architecture Map](mem://standards/release-architecture-map) — Architectural map of version propagation, sync pipeline (`sync-version`, `sync-spec-tree`, `sync-health-score`, `sync-readme-stats`, `sync-guidelines`), and release ceremony.
- [Prompt Architect Version Tracking](mem://standards/prompt-architect-version-tracking) — Installers inject a `promptArchitectByRiseupAsia` block into target repo's `version.json` when copying `.lovable/prompts`. Block contains author (Md. Alim Ul Karim, Chief Software Engineer), sourceRepository, installedAt, version, lastCommit, and fileMapping array. Template at `prompt-version.template.json`.
- 🔴 Install scripts dynamically inject `codingGuideline` and `promptArchitectByRiseupAsia` sections into target repository's `version.json` without overwriting the file. The `.lovable/memory` directory is packaged in bundles and transferred to target repositories during install and update phases so all downstream repos inherit memory standards.
