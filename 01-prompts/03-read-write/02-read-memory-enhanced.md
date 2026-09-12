# Memory Retrieval, Git Commit History & Project Context Ingestion — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Load the project's identity, specifications, conventions, active plans, recent git commit history, and Root Cause Analysis (RCA) records into your context before starting any task. Never repeat a logged past failure.

/learn Ingest and internalize all past learnings, user corrections, patterns, coding rules, error philosophies, RCA logs, and project specifications from `.lovable/memory/learned/`, `.lovable/memory/`, `.lovable/issues/`, `.lovable/cicd-issues/`, and `.lovable/strictly-avoid.md` so Antigravity operates with zero hallucination.

The specs, `.lovable/` folder, `what-to-read.md`, recent git commit history, root `readme.md`, and the codebase as a whole are the single source of truth. Your training data is not. If the two disagree, the repo wins, every time.

Autonomously self-loop and read:

- /learn the last 10 git commits (via `git log -n 10 --stat` and diffs) to understand what changes were made across the files so you understand the latest changes and can work with them seamlessly.
- /learn `.lovable/what-to-read.md` first as the authoritative reading priority list before exploring the rest of the codebase.
- /learn the entire codebase as a whole to create memory.
- /learn the root `readme.md` to create memory.
- /learn the entire `.lovable/` folder (especially `what-to-read.md`, `.lovable/coding-guidelines.md` and all files they reference) to create memory.
- /learn every single folder, subfolder, and nested markdown file in the `spec/` directory (specifically `spec/02-coding-guidelines/`, `spec/03-error-manage/`, enum fixes, database conventions) to create memory.
- /learn all recent Root Cause Analysis (RCA) records, retrospectives, and past failure post-mortems in `.lovable/issues/01-<slug>.md`, `.lovable/cicd-issues/01-<slug>.md`, and `spec/03-error-manage/01-error-resolution/03-retrospectives/` to ensure past mistakes are never repeated.
- /learn all hard prohibitions in `.lovable/strictly-avoid.md`.
- Read every pending task across `.lovable/plans/pending/01-<slug>.md`, `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md`, `.lovable/issues/`, and `.lovable/cicd-issues/`, listing them out in full.

Note on spec folder naming: Spec folders follow the hyphenated pattern `spec/<NN>-<slug>/` where `<NN>` is a sequence prefix and `<slug>` is the descriptive name. These numbers and folder placements are not rigidly fixed and may switch or be reorganized between projects. This canonical layout represents the general architecture the AI must dynamically discover, inspect, and read in full.

You are done reading when you can, without guessing:

- Summarize the last 10 git commits, what files were changed, and their architectural intent.
- Confirm `.lovable/what-to-read.md` was read first and followed in full.
- Name the CODE RED rules.
- Name the naming, error-handling, and DB conventions.
- List what is currently in `.lovable/plans/pending/` (sequenced as `01-`, `02-`) and every active subtask.
- Point at the exact file that justifies any rule you enforce.
- Explain the whole codebase structure, active DB schemas, API route contracts, app features (`spec/21-app/`), coding guidelines (`spec/02-coding-guidelines/`), and error management philosophy (`spec/03-error-manage/`).
- List out all pending tasks and unresolved issues with accurate step counts.
- Confirm that every nested markdown file in `spec/` has been inspected and broken links identified.
- Confirm runtime toolchain and package dependency compatibility.
- Confirm that the root readme is strictly lowercase `readme.md` (and auto-fixed/committed/pushed if it was not).

If you cannot do that, keep reading. Do not start work.

## Reading Strategy: Fast Cached Tools & Autonomous Looping

The `.lovable/`, `spec/`, and knowledge base folders can be massive. To process this information with zero blind spots and sub-second latency:

> [!TIP]
> **FAST KNOWLEDGE BASE & FILE READING (CRITICAL MANDATE):**
> Whenever you need to read the knowledge base, scan directories, find something, or read file contents, you MUST use the high-performance Python reader script `03-ai-scripts/17-fast-file-reader.py` instead of generic tools.
>
> **How it works:**
> - It is a standalone, dependency-free Python script designed for AI agents that provides sub-millisecond file reading and directory exploration.
> - It leverages `tmp/cache/` (or `__pycache__` style internal memoization) for instant lookups (<1ms) and automatically falls back to live disk scanning.
>
> **If the script is missing (How to recreate it):**
> If `03-ai-scripts/17-fast-file-reader.py` is not found, you must immediately recreate it using only Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`).
> - **Requirements:** Import any shared constants from `03-ai-scripts/02-shared-engine.py` if present. Ensure strict UTF-8 encoding (`sys.stdout.reconfigure(encoding="utf-8")`). Implement 2-tier caching for repeated calls.
> - **Core CLI Flags to implement:**
>   - `--list-folder <folder_path> [--ext .md,.ts]`: Recursively lists files matching extensions.
>   - `--read-file <file_path> [--max-bytes N]`: Fast content extraction with bounded limits.
>   - `--search-pattern "<regex_term>" [--path <dir>]`: Rapid content searching using compiled regex.

1. Fast Cached Directory Discovery:
   - Use `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> [--ext <extensions>]` to enumerate repository structures in <15ms via `tmp/cache/`.
   - Use `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<text>" [--path <dir>] [--ext <extensions>]` for parallel multi-threaded content grepping using pre-compiled regexes.

2. Pre-Flight Script Authoring Checklist:
   - Inspect `03-ai-scripts/02-shared-engine.py` for centralized constants, `RegexPatternType` Enums with PascalCase members, and lazy regex memoization.
   - Inspect `03-ai-scripts/01-index.md` for tool inventory and performance baselines.

3. Autonomous looping enforcement:
   - The AI agent MUST take at least 50 distinct steps to autonomously self-loop through all directories and files across `spec/`, `.lovable/`, and application source trees.
   - Do not stop after one high-level glance. You must deeply explore and learn.
   - Systematically iterate through each directory layer.
   - **Crucial Rule:** You MUST author or update an Antigravity skill (`.agents/skills/<slug>/skill.md`) for every major spec/prompt you read to persist it as an agent capability.

2. Deep recursive spec traversal:
   - The AI must recursively inspect every subfolder and all nested `.md` files in `spec/` (`00-overview.md`, numbered specs `01-*.md`, `99-consistency-report.md`, `spec-index.md`, subdirectories).

3. Anti-hallucination and clarifying questions:
   - Scan internal markdown references across `.lovable/` and `spec/`.
   - If a referenced spec or issue file is missing on disk, or if requirements are ambiguous, the AI MUST NOT guess or assume.
   - It must automatically log an open question in `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md` or ask the user directly to ensure alignment.

4. Active schema and API contract mapping:
   - Ingest and maintain an in-memory map of active DB tables/schemas (`spec/04-database-conventions/`, `spec/23-app-db/`), API endpoints, and global state stores so downstream tasks have zero assumptions on field names or parameter shapes.

5. Tooling and runtime compatibility check:
   - Inspect package manifests (`package.json`, `tsconfig.json`, build configs) to catalog runtime targets, linter rules, and banned packages before completing onboarding.

6. Parallel sub-agents for deep reading:
   - You are allowed and strongly encouraged to spawn dedicated sub-agents to read items and synthesize memory in parallel.
   - When spawning a sub-agent for reading, give it a highly specific title reflecting exactly what it is reading (e.g., `Reading Auth Specs in spec/21-app`, `Scanning Error Management in spec/03-error-manage`). Do not use generic names. If an agent switches tasks, its title must change.
   - Assign sub-agents small, granular folders/files to read rather than asking one agent to read the entire codebase in a single pass.

7. Root `readme.md` lowercase self-healing exception:
   - If the root readme is uppercase `README.md` or incorrectly cased, immediately rename it to `readme.md`, commit, and push to git without asking.

8. Memory persistence:
   - You are allowed to write to the `.lovable/` directory to enhance project memory after reading.
   - Write summaries of what you learned into `.lovable/memory/learned/01-<slug>.md` (or `.lovable/memory/01-<slug>.md`), including file counts, to maintain context.
   - Update `.lovable/what-to-read.md` based on your progress to guide future reading workflows.
   - Document any discovered bugs into `.lovable/issues/01-<slug>.md` or `.lovable/suggestions.md`.
   - Capture open ambiguities or update execution plans.

9. Missing spec file protocol:
   - If a spec folder contains only `.gitkeep` or missing reference files, check the full names in `01-prompts/04-coding-standards/01-coding-guidelines.md` or existing plans.
   - Use available guidelines in the prompt library.
   - If critical information is absent, explicitly ask the user for the file.

10. CRITICAL read-only enforcement:
    - Other than fixing the root `readme.md` lowercase naming if needed, you MUST NOT refactor, edit, or write any application source code.
    - This is a strictly read and analysis phase.

---

## Root Cause Analysis (RCA) & Failure Memory Architecture (Non-Negotiable)

To guarantee institutional memory and prevent regressions across all workflows (Plan Mode, CI/CD Fix, Coding Guidelines Audit, and Execution), all Root Cause Analyses (RCAs) follow strict canonical paths, file naming formats, and anti-hallucination guardrails:

1. General Issue / Bug RCAs:
   - Path: `.lovable/issues/01-<slug>.md` (sequenced as `01-`, `02-`, etc.)
   - Index: Registered in `.lovable/01-index.md` (or `.lovable/plans/01-index.md`)
   - Mandatory Structure: Error description, exact file/line location, Root Cause Analysis (one-sentence root cause + deep analysis), fix strategy, and prevention checklist.

2. CI/CD Failure RCAs:
   - Path: `.lovable/cicd-issues/01-<slug>.md` (sequenced as `01-`, `02-`, etc.)
   - Index: Registered in `.lovable/cicd-index.md` (or `.lovable/cicd-index.md`)
   - Mandatory Structure: Raw pipeline error snippet, Root Cause Analysis, resolution applied, and "What NOT to Repeat" rules.

3. Retrospectives & Architectural Failure Learnings:
   - Path: `.lovable/memory/learned/01-<slug>.md` and `spec/03-error-manage/01-error-resolution/03-retrospectives/01-<slug>.md`
   - Index: Registered in `.lovable/memory/01-index.md`

4. Hard Avoidances / CODE RED Prohibitions:
   - Path: `.lovable/strictly-avoid.md`
   - Hard Rule: Append-only. Never overwrite or truncate existing entries. Any RCA that uncovers a forbidden anti-pattern must append a one-line rule here.

5. Plan Mode Integration:
   - Path: `.lovable/plans/pending/01-<slug>.md` & `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md`
   - Context Requirement: Must cite previous relevant RCAs and failure records so the new plan explicitly avoids repeating past errors.

6. Anti-Hallucination Guard for RCAs:
   - If an RCA file, referenced spec, or issue file is missing on disk, the AI MUST NOT guess or assume its contents. Stop and file an open ambiguity in `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md` or ask the user directly before proceeding.

---

## Phase 1: Load the Project & Recent Git History

### 1.0 Read the Last 10 Git Commits & File Changes (MANDATORY)

Before touching any code, drafting plans, or executing tasks, the AI MUST inspect the recent git history to ground its context in the latest repository state:

1. Run `git log -n 10 --stat` to view the last 10 commits, their commit messages, and the exact files modified, added, or deleted.
2. Understand what changes were made to the files:
   - Identify which packages, specifications, scripts, or documentation were recently edited or introduced.
   - For any complex, architectural, or ambiguous changes, run `git show <commit-sha> --stat` or inspect specific diffs (`git diff <commit-sha>~1 <commit-sha> -- <file>`) to understand the underlying implementation logic.
   - Ensure you do NOT revert recently added features, re-introduce anti-patterns that were just removed, or contradict recently established conventions.
3. Formulate an internal summary of:
   - What features, refactors, or bug fixes were completed across the last 10 commits.
   - Which files and packages are currently in active development.
   - Why those changes were made (architectural motivations and rationale).

### 1.1 Read `what-to-read.md` First (Authoritative Priority List)

1. Read `.lovable/what-to-read.md` in full before reading other files.
2. Follow the exact prioritized reading order defined in `.lovable/what-to-read.md` (e.g. `version.json`, `.lovable/memory/01-index.md`, learned memories, coding guidelines, active plans).
3. The reading sequence in `what-to-read.md` overrides any generic assumptions or ad-hoc file exploration.

### 1.2 Confirm Root `readme.md` Lowercase (Auto-Fix & Commit)

1. Root `readme.md` lowercase verification and auto-fix:
   - Verify that the root readme file is strictly named lowercase `readme.md`.
   - If an uppercase `README.md` exists or the casing is incorrect on disk or in git, immediately rename it to `readme.md`, remove the stale uppercase file, commit the change (`fix: ensure root readme is strictly lowercase readme.md`), and push to git without asking or second-guessing.
   - Read the root `readme.md` file for architecture, casing rules, repository layout, and AI entry points.

### 1.3 Read the Whole `.lovable/` Folder & Flag Markdown Files

Walk `.lovable/` recursively. Every file matters. Missing files are noted, not silently skipped.

> [!IMPORTANT]
> **MANDATORY `.lovable/*.md` FLAGGING & AUDIT RULE:**
> Every single `.md` file discovered directly in `.lovable/` (e.g. `.lovable/folder-structure.md`, `.lovable/prompts.md`, `.lovable/strictly-avoid.md`, `version.json`) or inside any nested subdirectory MUST be actively read, cataloged, and flagged in the memory log. If an unindexed or orphan markdown file is found, immediately raise a flag in the context log and register it into `.lovable/memory/01-index.md` or `.lovable/plans/01-index.md`.

| #   | Path                                                  | What you get                                                                                                                                |
| --- | ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `.lovable/what-to-read.md`                     | Authoritative reading order for this project. Read it first and follow all referenced files.                                            |
| 2   | `.lovable/folder-structure.md`                | Canonical architectural map of all `.lovable/` folders, scripts, and naming conventions.                                                  |
| 3   | `.lovable/strictly-avoid.md`                          | Hard prohibitions (CODE RED) — append-only, never truncate.                                                                               |
| 4   | `.lovable/coding-guidelines.md`     | Master single source of truth for cross-language coding guidelines, `Result[T]` envelopes, and `is`/`has` booleans.                       |
| 5   | `03-ai-scripts/`                            | Automated Python utilities (`01-file-manipulator.py`, `06-cicd-local-runner.py`, `08-fast-file-scanner.py`, `09-fast-cached-grep.py`, `01-index.md`). Pre-warms `tmp/cache/` file caches for rapid zero-overhead file discovery across folders. |
| 6   | `.lovable/prompts.md` + `01-prompts/`           | Canonical prompt index and mirrored execution prompts (`cg-execute/`, `execute/`, `ci-cd/`).                                              |
| 7   | `.lovable/memory/01-index.md`                         | Master index of institutional knowledge. Then read every file it references recursively.                                                 |
| 8   | `.lovable/plans/01-index.md`                             | Master index of all plans (pending + completed + subtasks). Read this before touching individual plan files.                               |
| 9   | `.lovable/plans/pending/`                             | Active plans, `01-<slug>.md` — read all and list out each pending item.                                                                   |
| 10  | `.lovable/plans/subtasks/`                            | Granular 5–8 file batch files (`batch-01.md`, etc.) linked from parent plans.                                                              |
| 11  | `.lovable/plan.md`                      | Instant recovery state from the most recent run failure.                                                                                   |
| 12  | `.lovable/plans/completed/`                           | Historical completed plans and execution logs.                                                                                              |
| 13  | `.lovable/issues/`                                    | General bugs and regressions — read and list out pending bugs.                                                                             |
| 14  | `.lovable/cicd-issues/`                               | CI/CD-specific failures — read and list out pending CI/CD issues.                                                                           |
| 15  | `.lovable/ambiguous-questions/01-new-ambiguity/`      | Open questions currently blocking work. Surface them immediately; do NOT guess past them.                                                  |
| 16  | `.lovable/ambiguous-questions/02-ambiguity-resolved/`  | Answered questions with applied solutions. Binding decisions; do not re-litigate.                                                          |
| 17  | `.agents/skills/`                                     | Progressive disclosure runbooks for Antigravity agents (`<slug>/skill.md`).                                                                |
| 18  | `.agents/rules/`                                      | Systemic behavioral and coding rules auto-enforced by Antigravity.                                                                         |
| 19  | `linter-scripts/`                                     | Repository linters (`check-newline-styling.py`, `check-markdown-header-spacing.py`, `check-relative-paths.py`).                            |
| 20  | Anything else under `.lovable/`                       | Read it. If the folder exists, it exists for a reason.                                                                                      |

---

### 1.2 Mandatory Skill & Rule Set Auto-Generation on Ingestion

Whenever the AI agent reads prompts or coding guidelines during memory ingestion:

1. **Auto-Generate Antigravity Skills for Prompts:**
   - For every prompt ingested from `01-prompts-category/` or `01-prompts/`, create or update a dedicated Antigravity skill in `.agents/skills/<slug>/skill.md` with standard YAML frontmatter (`name`, `description`).
   - Ensures any agent in the ecosystem can activate the skill on demand.

2. **Auto-Generate Antigravity Rules for Coding Guidelines:**
   - For all coding guidelines ingested from `.lovable/coding-guidelines.md` or `spec/02-coding-guidelines/`, synthesize and write authoritative agent rules into `.agents/rules/<slug>.md` and inject essential constraints into `AGENTS.md`.
   - Core rules enforced:
     - **Strict Boolean Standard:** `is and has only (can, should, was, etc. are banned)`.
     - **No Bare Void in Go:** Functions must return `Result[T]` or `*appfault.AppError`.
     - **Parameter Structs:** Banned loose >2-3 parameters; use `*Params` structs.
     - **Vertical Line Gaps:** Mandatory blank lines before `if`, after `}`, before `return`, and around multiline struct calls.
     - **5–8 Files Micro-Batching:** All refactors broken into bounded subtasks.

---

### 1.3 Loop Through the Entire `spec/` Directory, Subfolders, and Nested Files

Systematically loop through the `spec/` folder, dynamically matching canonical hyphenated names (numbers may vary between projects):

- `spec/01-spec-authoring-guide/`: Spec authoring conventions, required files, format requirements.
- `spec/02-coding-guidelines/` (or `01-prompts/04-coding-standards/01-coding-guidelines.md`): Zero-tolerance coding standards, function size caps (8 lines preferred, 15 max), boolean naming (`is*`, `has*`, positive framing), immutable patterns, DRY priority 1.
- `spec/03-error-manage/`: Error management philosophy — never swallow errors, log operation name and key inputs on every catch, wrap errors without losing cause, typed errors only, universal response envelopes (`{ data, errors[], meta }`).
- `spec/04-database-conventions/`: Database schema, table naming (PascalCase), columns (camelCase), primary keys (`{Table}Id`), SQLite/ORM rules, ERD requirements.
- `spec/05-split-db-architecture/` through `spec/19-main-worker-service/`: Architectural specs for config, design system, docs viewer, code blocks, CLI, workflows, and release.
- `spec/21-app/`: App specification, domain architecture, core capabilities, routes, and business rules.
- `spec/22-app-issues/`, `spec/23-app-db/`, `spec/24-app-ui-design-system/`: App-specific issues, schemas, and design systems.
- Recursively read all nested markdown files (`*.md`), overview documents (`00-overview.md`), consistency reports (`99-consistency-report.md`), and `spec-index.md`. If a folder contains only `.gitkeep`, fallback to the prompt library guideline or ask the user.

### 1.4 Loop Through the Entire Codebase as a Whole

Autonomously survey the codebase structure end-to-end:

- Root configuration files, package manifests, build scripts, tsconfig, linter configs.
- Application directory (`src/` or app root), entry points, routing tree, components, state management stores, and utility modules.
- Database schemas, models, and migrations (`db/`, `prisma/`, `drizzle/`, SQLite tables).
- Asset directories (`assets/`).
- Verify how data flows from input to state, backend/storage, and UI presentation.

### 1.5 The Two Index Files

Two indexes decide what you read next. Treat them as required entry points, not as summaries:

- `.lovable/memory/01-index.md` lists every institutional-knowledge file. If it points at 12 files, you read 12 files.
- `.lovable/plans/01-index.md` lists every plan (pending, completed, subtasks) with its slug, status, and one-line intent. Use it to pick which plan files to open in full. If it is missing, create it as part of the next code change.

### 1.6 Self-Check (Internal, Before Phase 2)

- CODE RED rules?
- Naming conventions (files, folders, DB columns, variables)?
- Root readme strictly lowercase `readme.md`?
- Error-handling philosophy (`spec/03-error-manage/`)?
- What is in `.lovable/plans/pending/` (sequenced as `01-`, `02-`) and `plans/subtasks/` right now (exact list)?
- Active DB schemas, table columns, and API contracts?
- App specs and domain architecture (`spec/21-app/`)?
- Whole codebase layout and component flow?
- Top forbidden patterns?

If any answer is fuzzy, go back and reread by looping through the files again. Do not proceed.

---

## Phase 2: Consolidated Guidelines

Read `spec/17-consolidated-guidelines/` (or `spec/17-consolidated-guidelines/`) in numeric order (`01-*.md` through `18-*.md`). Each file is a self-contained policy document. Missing folder: note it and continue.

---

## Phase 3: Spec Authoring Rules

Read `spec/01-spec-authoring-guide/` in numeric order. You should come out knowing:

- File and folder naming conventions (`<NN>-<slug>/`).
- Required files per spec folder (`00-overview.md`, `99-consistency-report.md`).
- The `.lovable/` layout (see Phase 1.1).
- The linter infrastructure.

---

## Phase 4: Task-Driven Deep Dives

Only open a spec folder when the current task needs it.

| Task involves...                          | Read                                    |
| ---------------------------------------- | --------------------------------------- |
| Writing or reviewing code                | `spec/02-coding-guidelines/`            |
| Error handling                           | `spec/03-error-manage/`                 |
| Database schema or queries               | `spec/04-database-conventions/`         |
| SQLite / multi-DB architecture           | `spec/05-split-db-architecture/`        |
| Config systems                           | `spec/06-seedable-config-architecture/` |
| UI theming, CSS variables, design tokens | `spec/07-design-system/`                |
| Documentation viewer features            | `spec/08-docs-viewer-ui/`               |
| Code block rendering                     | `spec/09-code-block-system/`            |
| PowerShell scripts                       | `spec/11-powershell-integration/`       |
| CI/CD pipelines                          | `spec/12-cicd-pipeline-workflows/`      |
| CLI self-update                          | `spec/14-update/`                       |
| WordPress plugins                        | `spec/18-wp-plugin-how-to/`             |
| App-specific features                    | `spec/21-app/`                          |
| Known app bugs                           | `spec/22-app-issues/`                   |
| App-specific DB schema                   | `spec/23-app-db/`                       |
| App-specific UI + design system          | `spec/24-app-ui-design-system/`         |

Inside each folder: `00-overview.md` -> numbered files -> `99-consistency-report.md`.

Fallbacks when the canonical numbered folder is absent: `.lovable/coding-guidelines.md`, `spec/02-coding-guidelines/`, `coding-guidelines/`, `spec/03-error-manage/`, `01-prompts/04-coding-standards/01-coding-guidelines.md`. Numbered folder wins on conflict; call the conflict out in the plan's Context.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination Contract

1. If the specs are silent on a rule, that rule does not exist. Do not invent one.
2. Specs beat training data. Always.
3. Cite the file and section when you enforce a rule.
4. When a spec is ambiguous or missing, ask questions. Do not "use best judgement".
5. Do not blend this project's conventions with conventions from other projects you have seen.
6. No filler. No "hope this helps", no "let me know".

---

## Memory Update Protocol

```
New info discovered
├─ Institutional knowledge (pattern / convention / decision)?
│   YES → .lovable/memory/01-<slug>.md  +  update .lovable/memory/01-index.md
├─ Must never happen again?
│   YES → .lovable/strictly-avoid.md
├─ Idea, not yet approved?
│   YES → .lovable/suggestions.md
├─ New user command / convention?
│   YES → .lovable/spec/commands/01-<slug>.md
├─ Bug / regression?
│   YES → .lovable/issues/01-<slug>.md   (or .lovable/cicd-issues/ if CI/CD)
├─ New or changed plan?
│   YES → .lovable/plans/pending/01-<slug>.md  +  update .lovable/plans/01-index.md
├─ Ambiguity / unclear requirement blocking progress?
│   YES → .lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md
├─ User just answered a previously-open ambiguity?
│   YES → mv the file to .lovable/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md,
│         append `## Resolution` (answer + applied solution), flip Status: resolved
└─ None of the above → do not persist.
```

Hard rules:

- Folder is `.lovable/memory/`, never `memories/`.
- Adding a memory file always updates `.lovable/memory/01-index.md`.
- Adding, moving, or completing a plan always updates `.lovable/plans/01-index.md`.
- Ambiguity folders: `01-new-ambiguity/` for open, `02-ambiguity-resolved/` for answered. On answer, MOVE the file (never copy) so it exists in exactly one place. Every resolved file carries a `## Resolution` section.
- Never guess past an open ambiguity. If one exists and is relevant to the current task, stop and surface it before doing work.
- Editing existing memory or index files preserves unrelated content. No silent truncation.
- Any code-base change bumps the minor version.

---

## Completion Confirmation

After Phases 1-3, reply exactly:

```
Onboarding complete.

- Recent git commits inspected: [10] (from git log -n 10 --stat)
- What-to-read followed: [yes] (from .lovable/what-to-read.md)
- Memory files read: [X]
- Consolidated guidelines read: [Y]
- Spec authoring files read: [Z]
- Pending plans: [N]  (from .lovable/plans/01-index.md)
- CI/CD issues absorbed: [M]  (from .lovable/cicd-issues/)
- Open ambiguities: [K]  (from .lovable/ambiguous-questions/01-new-ambiguity/)
- Resolved ambiguities on file: [R]  (from .lovable/ambiguous-questions/02-ambiguity-resolved/)

I understand:
- Latest commit changes: [brief summary of recent file changes and architectural intent from the last 10 commits]
- CODE RED rules: [top 3-5]
- Naming conventions: [brief]
- Error handling: [one sentence]
- Active DB schemas & contracts: [key models / tables]
- Active plans & pending tasks: [slugs from .lovable/plans/pending/ and subtasks]
- Strict avoidances: [top 3-5]
- Blocking ambiguities: [slugs, or "none"]

Ready for tasks.
```

Then stop. No next-step suggestions, no exploratory questions.

---

## Pre-Reply Checklist (All Must Be True)

/goal Complete the checklist properly until done can do self-looping.

1. [ ] Inspected the last 10 git commits (messages, file change stats, and diffs via `git log -n 10 --stat`) to understand recent modifications and file evolutions.
2. [ ] /learn `.lovable/what-to-read.md` first and followed its prioritized reading sequence in full.
3. [ ] Confirmed root readme is strictly lowercase `readme.md` (auto-fixed, committed, and pushed if uppercase or missing).
4. [ ] /learn the root `readme.md` file (casing rules, architecture, entry points).
5. [ ] Walked `.lovable/` recursively, no folder or file skipped silently, and flagged all `.lovable/*.md` files.
6. [ ] /learn `.lovable/memory/01-index.md` and every file it points at.
7. [ ] /learn `.lovable/plans/01-index.md`, every file in `pending/` (sequenced as `01-`, `02-`), and all active subtasks.
8. [ ] Skimmed `.lovable/plans/completed/` for recent history.
9. [ ] /learn every file in `.lovable/spec/commands/`.
10. [ ] /learn every file in `.lovable/issues/` and `.lovable/cicd-issues/`.
11. [ ] /learn every file in `.lovable/ambiguous-questions/01-new-ambiguity/` and `02-ambiguity-resolved/`.
12. [ ] Scanned for broken links or missing docs and surfaced them under open ambiguities.
13. [ ] Ingested active schema models, DB column conventions, and API route shapes.
14. [ ] Verified runtime dependencies and package compatibility.
15. [ ] Recursively traversed and read every subfolder, nested markdown file (`*.md`), overview, and consistency report within `spec/` (e.g. `spec/01-spec-authoring-guide/`, `spec/02-coding-guidelines/`, `spec/03-error-manage/`, `spec/04-database-conventions/`, `spec/21-app/`, etc.).
16. [ ] Autonomously surveyed and looped through the entire codebase as a whole (all application code, entry points, routes, components, state stores, utilities, and configuration files).
17. [ ] /learn `spec/17-consolidated-guidelines/` (or `spec/17-consolidated-guidelines/`) in numeric order (or noted missing).
18. [ ] /learn `spec/01-spec-authoring-guide/` in numeric order (or noted missing).
19. [ ] Can name CODE RED rules, naming conventions, error-handling philosophy without guessing.
20. [ ] Can list every pending plan slug and subtask from memory.
21. [ ] Checked whether the repo contains explicit tone, strictly-avoid, or prior-stupidity instructions and applied them without softening.
22. [ ] Did not replace hard user wording with polite generic language.
23. [ ] Emitted the Completion Confirmation block verbatim, then stopped.
24. [ ] Confirmed that reading remained strictly read-only regarding the codebase (no source code refactored, only memory and lowercase readme auto-fix updated).

## Actionable Items & Checklist

/goal Complete the checklist properly until done can do self-looping.

1. [ ] Inspect the last 10 git commits via `git log -n 10 --stat` to understand recent file changes and ensure seamless continuity with the latest code state.
2. [ ] /learn `.lovable/what-to-read.md` and follow all prioritized reading entries before exploring the codebase.
3. [ ] /learn the coding guidelines in: `.lovable/coding-guidelines.md` and create memory.
4. [ ] /learn the condition extraction in: `spec/02-coding-guidelines/01-cross-language/01-index.md` and create memory.
5. [ ] /learn the formatting and braces in: `spec/02-coding-guidelines/01-cross-language/01-index.md` and create memory.
6. [ ] /learn the multi-line formatting in: `spec/02-coding-guidelines/01-cross-language/01-index.md` and create memory.
7. [ ] /learn the boolean guidelines in: `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md` and create memory.
8. [ ] /learn the anti-hallucination rules in: `spec/02-coding-guidelines/01-cross-language/01-index.md` and create memory.
9. [ ] /learn the error management architecture in: `spec/03-error-manage/01-index.md` (and related error manage files) and create memory.
10. [ ] /learn all recent Root Cause Analysis (RCA) files in `.lovable/issues/`, `.lovable/cicd-issues/`, and `spec/03-error-manage/01-error-resolution/03-retrospectives/` to prevent recurring errors.
11. [ ] /learn all hard prohibitions in `.lovable/strictly-avoid.md` and verify zero violations.
12. [ ] /learn the enum standards and fixes in: `spec/17-consolidated-guidelines/07-enum-standards.md` and `spec/17-consolidated-guidelines/07-enum-standards.md` and create memory.
13. [ ] /learn ALL other single-file specs in `spec/02-coding-guidelines/` and create memory.
14. [ ] /learn the overarching main task plan.
15. [ ] Ensure the git repository starts completely clean.
16. [ ] /goal Complete all work on the current branch only.
17. [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
18. [ ] Group all completed work into a single logical commit.
19. [ ] Push the commit to the remote repository.
20. [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: read-memory-enhanced
- status: active
