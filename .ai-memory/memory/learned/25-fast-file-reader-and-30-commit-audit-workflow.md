# Fast File Reader, 30-Commit Audit & Prompt 2.2.0 Protocols

- Slug: fast-file-reader-and-30-commit-audit-workflow
- Date: 2026-09-13
- Category: learned
- Status: permanent

## Overview

In response to the workflow upgrade to Prompt Version 2.2.0 ("Memory Persistence & Issue Logging — Workflow"), this memory documents the required architectural protocols for rapid repository exploration, deep commit history auditing, and structured 20-task tracking across AI loop turns.

The previous session achieved successful resolution of the v6.227.0 release pipeline failure (GitHub Actions run #34746417661 passed 100% green across all four jobs: Build and Release, Installer Dry-Run Windows, Installer Smoke, and Installer Smoke Windows). This turn formalizes the expanded memory and exploration standards introduced in version 2.2.0.

## Key Architectural Protocols

### 1. Fast Cached Exploration via `17-fast-file-reader.py` (Rule 16)
- AI agents must avoid slow, high-overhead recursive shell calls (`Get-ChildItem -Recurse`, `dir /s`, or nested bash walks) when inspecting repository structure.
- Canonical tool: `python 03-ai-scripts/17-fast-file-reader.py`
  - Folder listing: `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> [--ext .md,.ts]`
  - File reading: `python 03-ai-scripts/17-fast-file-reader.py --read-file <path> [--max-bytes N]`
  - Content pattern search: `python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<term>" [--path <dir>]`
- Performance characteristics: Uses pre-computed cache files in `tmp/cache/` for sub-millisecond lookups (<1ms), automatically falling back to live disk walking if cache is unavailable.
- Strict encoding: Reconfigures stdout to UTF-8 (`sys.stdout.reconfigure(encoding="utf-8")`) to prevent Windows CP1252 character mapping exceptions.

### 2. Mandatory 30-Commit Git History Audit (Rule 17)
- Memory authoring must never rely on chat memory or assumptions alone.
- Before writing memory, the AI must execute `git log -n 30 --oneline` (and `git log -n 30 --stat` where required) to inspect the last 30 commits.
- This extracts:
  - Recent architectural directions and refactoring trajectories (e.g. types.go extraction, parameter struct consolidation, Result wrapper unification).
  - Applied user directives and constraints (e.g. skipping tests during release, affirmative boolean naming).
  - Bug fixes and root cause remediations (e.g. installer seed URL fixes, CI dry-run path synchronization).
  - Release increments and changelog entries.

### 3. Recent 20-Task Tracking & Compact Task Register (Rule 18)
- The AI must maintain an explicit register of the last 20 tasks/plans in `.ai-memory/plans/01-index.md` under `## Recent Completed Tasks Register (Last 20 Tasks)`.
- This register must be referenced in `.ai-memory/what-to-read.md`.
- Newly authored memories must link back to these recent tasks so that subsequent agent iterations maintain continuous contextual awareness without re-reading the entire git history from scratch.

### 4. Prohibition of Em Dashes in Responses
- Per checklist item 19 and repo standards, completion confirmation blocks and responses must not contain Unicode em dashes. Hyphens (`-`) or colons (`:`) must be used instead.

### 5. Detailed Spec Preservation Policy (Rule 14)
- Simple routine tasks may be consolidated into session summaries to avoid file bloat.
- Detailed specifications, architectural designs, non-negotiable rules, domain specifications under `02-spec/`, and complex requirement documents must NEVER be consolidated, summarized, or truncated. Full fidelity and exact wording must be preserved.

## Verified Commit Audit Summary (Last 30 Commits)

- `fe162756`: docs(memory): record v6.227.0 release orchestration, installer seed URL fix, and release dryrun path
- `9109ee90`: release: v6.227.0 fix installer seed data URLs and release dryrun path
- `5522fe9f`: fix: update installer seed data URLs and release dryrun path to cli/scripts
- `cc701409`: release: v6.226.0 fix CI pipeline e2e bash syntax error, guard typescript AST import, and auto-migrate ssh tables
- `9b78d059`: fix(ci): balance e2e bash if block, guard typescript AST import, and auto-migrate ssh tables
- `af79067b`: fix(release): add retry logic with backoff for push_release_artifacts
- `0f0fcddf`: release: v6.225.0 robust Windows read-only git pack artifact sweeping and storage hygiene
- `e91bd8bb`: fix(hygiene): use robust_rmtree to unlock and sweep Windows read-only git pack artifacts
- `d30d398b`: chore(inventory): sync test timings after v6.224.0 verification
- `cc968591`: release: v6.224.0 CPU freeness scaling and automated pre-build/pre-test disk storage hygiene
- `26ced7e2`: fix(linter): enable --allow-serial-runners for concurrent golangci-lint invocations
- `83b93927`: fix(hygiene): enforce automatic pre-build and test artifact purging and disk space reclamation
- `1c216e27`: chore: update test inventory timestamps after release v6.223.0
- `0fff9c9e`: release: v6.223.0 unify 32-worker test pool, eliminate long-tail stalls, and accelerate coverage gate to 2s
- `a3343383`: perf: unify unit test priority worker pool across 32 threads, mock network queries, and optimize coverage gate from 193s to 2s
- `76fee8db`: chore: update test inventory manifest timestamps
- `1d5e2a04`: release: v6.222.0 route agy to Desktop IDE, add agm manager command, and optimize CPU freeness worker scaling
- `cdb0100c`: feat: route agy to Desktop IDE, add agm manager command, and optimize CPU freeness worker scaling to 32 threads
- `f76f875a`: feat(install): default agy install to Antigravity IDE and route agm install to Antigravity Manager / Tools
- `4620060c`: feat(install): route agy alias to Antigravity Desktop IDE and detect installed path
- `2a1385c4`: docs(prompts,skills,spec): enforce IsDefined over !isEmpty and clarify map lookup canonical naming
- `1b9c5bbb`: docs(prompts): clarify IsDefined replacement for !isEmpty and map lookup naming
- `703a857c`: feat(ci,workflows,prompts,skills): enforce zero-storage actions mandate, standardize IsDefined over !isEmpty, and sync prompts/skills
- `e11e1615`: refactor(cg): standardize IsDefined affirmative naming, ban isExists, and forbid compound negatives
- `9ba7650c`: refactor: centralize types.go & remove installation tests from E2E suites (Plan 148)
- `f6dd75a9`: refactor(params): enforce parameter structs, affirmative boolean fields, and AppError returns across cluster, clone, and visibility
- `f665d89b`: refactor(db,cluster): migrate cluster and ssh slice returns to ResultSlice and centralize types.go
- `c3b2e564`: refactor(result): centralize single reusable Result aliases into types.go and enforce affirmative boolean naming
- `1e4cdff2`: refactor(result): enforce pointer-attached null safety, core predicates, and single return envelopes
- `4654b0e8`: refactor(params): enforce parameter structs, affirmative boolean fields, and AppError returns
