# Prompt Architect: Canonical AI Prompts Library

This directory hosts the canonical, production-grade prompts library (V4 architecture) for the Prompt Architect meta-repository and all connected repositories.

> **Historical Archive Note:** Legacy versioned tiers (`execute/`, `v1/`, `v2/`, and `v3/`) have been archived in [`06-archive/`](../06-archive/) inside the `coding-guidelines` meta-repository and are never synchronized to downstream repositories.

---

## Core Architecture & Capabilities

1. **Antigravity Slash Command Links + GitMap AUM Engine:**
   - Interactive slash command links (`[/goal](slashCommand;goal)`, `[/learn](slashCommand;learn)`) combined with **GitMap AUM Engine** (`gitmap` CLI) as primary and Python scripts (`03-ai-scripts/`) as fallback.

2. **High-Speed File & Content Discovery (TOTAL BAN ON `Select-String` & `git grep`):**
   - Multi-Core Streaming Live Search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` (alias: `gitmap aum grep`) — streaming live disk text/regex scanner (replaces `Select-String`, `git grep`)
   - Instant Indexed Symbol Search: `gitmap search "<query>" [--limit <n>]` — cached SQLite symbol & keyword search
   - Universal Wildcard Search: `gitmap find "<pattern>" [-ext <ext>]` (<10ms across 10,000+ files)
   - Zero-Disk Streaming: `gitmap cat <filepath>`
   - Substring Filename Search: `gitmap find-files-any "<str>"` (alias: `gitmap ffa`)
   - Directory Indexing: `gitmap list-files [pattern]` (alias: `gitmap lf`)

3. **Pipeline AI & Zero-Credit-Waste Dynamic Waiting:**
   - Telemetry Query: `gitmap pipeline-ai status --json`
   - Adaptive Sleep: `gitmap pipeline-ai status -t <etaSeconds>` (replaces expensive busy-polling loops)
   - Failure Log Snippets: `gitmap pipeline errors` (alias: `gitmap pe`)
   - Runner Target Diagnostics: `gitmap pipeline details` (alias: `gitmap pd`)

4. **Autonomous Automated Release:**
   - Single-Command Release: `gitmap release --bump <patch|minor> -y` (alias: `gitmap r -y`)
   - Release History: `gitmap changelog` (alias: `gitmap cl`), `gitmap list-versions` (alias: `gitmap lv`)

5. **Instruction Precedence Mandates:**
   - **Standard Execute & Audit Prompts:** Enforce **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence)** — instructions provided ABOVE the prompt take highest priority.
   - **Below-Steps Prompts (`*-in-below-steps.md`, `32-cg-follow-other-prompts.md`):** Enforce **Bottom-Instruction Priority Mandate (Below Precedence)** — instructions appended after the trailing `--` and `## 🚨 High Priority Instructions Below` header take highest priority.

6. **Canonical V6 Parent Task N-Steps Continuous Loop (`14-execute/02-execute-parent-task-with-n-steps.md`):**
   - Pure parameterization driven entirely by header variables (`N = 300`, `A = 2`, `H = 2`, `C = 30`, `PHASE_1_BUDGET = 150`, `PHASE_2_BUDGET = 150`, `WAVES`) with zero hardcoded literal step or agent counts in the body.
   - **Mandatory Subagent Spawning Gate (`invoke_subagent`, Zero Solo Execution):** Requires spawning `A = 2` concurrent subagents (`H = 2` disjoint tasks per subagent, `TypeName: "self"` in Phase 2) across both Phase 1 discovery/spec generation and Phase 2 code execution.
   - 100% GitMap commit primacy via atomic `gitmap cpf` / `gitmap cpb` (eliminating manual `git add` and `git commit`), upstream `.gitignore` hygiene gate (R8) with automatic untracking of ignored files (`git rm --cached`), SQLite task tracking, self-contained worker briefs with language-specific rules, secrets gate, and <= 3,200-word footprint.
   - Historical versions (V2, V3, V4, V5, and `excute-parent-old.md`) are preserved in [`06-archive/execute/`](../06-archive/execute/).

---

## Directory Index

All 25 canonical prompt categories reside directly at the root of `01-prompts/`:

```text
01-prompts/
├── 00-folder-structure/
├── 01-prompt-library-setup/
├── 02-core-workflow/
├── 03-read-write/
├── 04-coding-standards/
├── 05-coding-guidelines/
├── 06-testing-and-qa/
├── 07-bug-fix/
├── 08-dry-code/
├── 09-commit-and-multi-agent-code-fix/
├── 10-ui-and-design/
├── 11-content-and-seo/
├── 12-old-plan-prompts/
├── 13-plan-audit/
├── 14-execute/
├── 15-cg-execute/
├── 16-ci-cd/
├── 17-release-management/
├── 18-insults/
├── 19-old-execute-prompts/
├── 20-ai-fix-script-prompts/
├── 21-temp-end-to-end-tests/
├── 22-letterly/
├── 23-cursor-prompts/
├── 24-sync/
│   ├── 01-sync.md
│   ├── 02-sync-other-codebase.md
│   └── readme.md
└── 25-ai-verification/
```

---

## Table of Categories

| # | Category | Key Prompts | Description |
| :---: | :--- | :--- | :--- |
| **00** | `00-folder-structure/` | `01-canonical-folder-structure.md` | Canonical repository layout and folder structure rules |
| **01** | `01-prompt-library-setup/` | `01-prompt-library-setup.md` | Library setup and prompt architecture scaffolding specification |
| **02** | `02-core-workflow/` | `01-next-steps.md`, `02-pending-tasks.md`, `03-unified-ai-prompt-v4.md` | Core planning, next steps, and unified autonomous execution protocol |
| **03** | `03-read-write/` | `01-write-antigravity.md`, `03-read-memory-latest.md`, `04-write-memory.md` | Agent configuration, memory retrieval, persistence, and proofreading |
| **04** | `04-coding-standards/` | `01-coding-guidelines.md`, `02-theming-guidelines.md` | Grounded multi-language coding standards and theming rules |
| **05** | `05-coding-guidelines/` | `01-plan-coding-guideline-audit.md`, `02-execute-coding-guideline-fix.md` | Coding guideline audits, plans, and automated fixes |
| **06** | `06-testing-and-qa/` | `01-autonomous-qa-and-testing-v4.md` | Autonomous test suite execution and quality verification |
| **07** | `07-bug-fix/` | `01-fix-with-rca.md` | Grounded 4-part root cause analysis and regression fix |
| **08** | `08-dry-code/` | `01-python-dry-architecture-and-caching.md` | DRY code architecture, script enums, and caching |
| **09** | `09-commit-and-multi-agent-code-fix/` | `01-boolean-improvements.md`, `03-commit-fix-v2.md`, `08-git-reconcile-and-resolve-conflict.md` | Atomic commits, boolean refactoring, and conflict resolution |
| **10** | `10-ui-and-design/` | `01-logo-create.md`, `02-react-ui-fixes-update.md`, `08-create-slide-deck.md` | Visual design, logos, slide decks, and React UI components |
| **11** | `11-content-and-seo/` | `01-jokes-ideas-generate.md`, `03-seo-optimization.md`, `05-update-readme.md` | Content copywriting, SEO optimization, and documentation sync |
| **12** | `12-old-plan-prompts/` | `01-plan-maximum-enforcement.md`, `05-plan-spec-steps.md` | Archived legacy planning protocols and audit workflows |
| **13** | `13-plan-audit/` | `01-inventory-pending-tasks.md`, `02-plan-spec-steps-v2.md`, `03-audit-app-spec.md` | Task inventory, spec planning V2, and blind-AI audits |
| **14** | `14-execute/` | `02-execute-parent-task-with-n-steps.md`, `09-parent-task-in-below-steps.md`, `14-run.md` | V6 autonomous parent task execution, continuous loops, and run orchestration |
| **15** | `15-cg-execute/` | `01-execute-coding-guideline-fix.md` through `20-*` | Granular coding guideline rule enforcement suite |
| **16** | `16-ci-cd/` | `01-ci-cd-fix-with-release.md`, `02-cicd-pipeline-create.md` | CI/CD pipeline diagnosis, creation, and release integration |
| **17** | `17-release-management/` | `01-release.md`, `02-patch-bump.md`, `03-minor-bump.md`, `04-major-bump.md` | Semantic release ceremonies, version bumps, and tags |
| **18** | `18-insults/` | `01-raw-insults.md` | Anti-carelessness discipline and quality enforcement |
| **19** | `19-old-execute-prompts/` | `01-execute-pending-tasks-old.md` | Archived legacy execute prompt templates |
| **20** | `20-ai-fix-script-prompts/` | `01-python-file-manipulator.md` | Dedicated AI utility script prompt definitions |
| **21** | `21-temp-end-to-end-tests/` | `01-temp-end-to-end-test.md` | Isolated temporary end-to-end tests with skip-by-default guards |
| **22** | `22-letterly/` | `01-mobile-letterly.md` through `10-execute-with-release-letterly.md` | Voice dictation prompt formatters for mobile, desktop, execution, and releases |
| **23** | `23-cursor-prompts/` | `01-mobile-letterly-cursor.md` through `10-execute-with-release-letterly-cursor.md` | Cursor IDE prompt formatters targeting `.cursor/skills/` resolution |
| **24** | `24-sync/` | `01-sync.md`, `02-sync-other-codebase.md` | Multi-repository synchronization engine: pull base branch, safety backup, asset mirroring with 5 boundary protections, atomic commit, and release tagging |
| **25** | `25-ai-verification/` | `01-retrospective-ai-verification.md` | Retrospective quality audit, guideline verification, and CI validation |
