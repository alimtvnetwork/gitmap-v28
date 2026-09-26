# Memory Retrieval & Project Context Ingestion — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Load the project's identity, specifications, conventions, active plans, and recent Root Cause Analysis (RCA) records into your context before starting any task. Never repeat a logged past failure.

/learn Ingest and internalize all past learnings, user corrections, patterns, coding rules, error philosophies, RCA logs, and project specifications from `.ai-memory/memory/learned/`, `.ai-memory/memory/`, `.ai-memory/issues/`, `.ai-memory/cicd-issues/`, and `.ai-memory/strictly-avoid.md` so Antigravity operates with zero hallucination.

The specs, `.ai-memory/` folder, `what-to-read.md`, root `readme.md`, and the codebase as a whole are the single source of truth. Your training data is not. If the two disagree, the repo wins, every time.

Autonomously self-loop and read:

- /learn the entire codebase as a whole to create memory.
- /learn the root `readme.md` to create memory.
- /learn the entire `.ai-memory/` folder (especially `what-to-read.md`, `.ai-memory/coding-guidelines.md` and all files they reference) to create memory.
- /learn every single folder, subfolder, and nested markdown file in the `02-spec/` directory (specifically `02-spec/02-coding-guidelines/`, `02-02-spec/03-error-manage/`, enum fixes, database conventions) to create memory.
- /learn all recent Root Cause Analysis (RCA) records, retrospectives, and past failure post-mortems in `.ai-memory/issues/01-<slug>.md`, `.ai-memory/cicd-issues/01-<slug>.md`, and `02-spec/03-error-manage/01-error-resolution/03-retrospectives/` to ensure past mistakes are never repeated.
- /learn all hard prohibitions in `.ai-memory/strictly-avoid.md`.
- Read every pending task across `.ai-memory/plans/pending/01-<slug>.md`, `.ai-memory/plans/subtasks/01-<slug>/01-<subslug>.md`, `.ai-memory/issues/`, and `.ai-memory/cicd-issues/`, listing them out in full.

Note on spec folder naming: Spec folders follow the hyphenated pattern `02-spec/<NN>-<slug>/` where `<NN>` is a sequence prefix and `<slug>` is the descriptive name. These numbers and folder placements are not rigidly fixed and may switch or be reorganized between projects. This canonical layout represents the general architecture the AI must dynamically discover, inspect, and read in full.

You are done reading when you can, without guessing:

- Name the CODE RED rules.
- Name the naming, error-handling, and DB conventions.
- List what is currently in `.ai-memory/plans/pending/` (sequenced as `01-`, `02-`) and every active subtask.
- Point at the exact file that justifies any rule you enforce.
- Explain the whole codebase structure, active DB schemas, API route contracts, app features (`02-02-spec/21-app/`), coding guidelines (`02-spec/02-coding-guidelines/`), and error management philosophy (`02-02-spec/03-error-manage/`).
- List out all pending tasks and unresolved issues with accurate step counts.
- Confirm that every nested markdown file in `02-spec/` has been inspected and broken links identified.
- Confirm runtime toolchain and package dependency compatibility.
- Confirm that the root readme is strictly lowercase `readme.md` (and auto-fixed/committed/pushed if it was not).

If you cannot do that, keep reading. Do not start work.

## STRICTLY FORBIDDEN: Writing Repository Files or Modifying Folder Structure (TOTAL BAN)

> [!CAUTION]
> **READ-ONLY MANDATE (TOTAL BAN ON REPOSITORY WRITES):**
> - The reading phase is strictly, 100% read-only for the repository workspace.
> - **NEVER** modify, create, rename, or delete any files or directories inside the repository during this workflow.
> - **NEVER** write or update memory files (e.g. `.ai-memory/memory/`, `.ai-memory/memory/learned/`, `.ai-memory/readme.md`).
> - **NEVER** auto-generate or write skills/rules (`.agents/skills/`, `.agents/rules/`).
> - **NEVER** modify `readme.md`, `version.json`, or any repository configurations.
> - **NEVER** execute `git add`, `git commit`, `git push`, or modify git working tree state during reading.

> [!IMPORTANT]
> **ISOLATED USER TEMP DIRECTORY FOR AGENT COMMUNICATION (ONLY PERMITTED WRITE LOCATION):**
> - If sub-agents or the AI MUST write intermediate scratchpad data, caching, or messages to communicate with each other during reading (e.g. for sub-agents to communicate, coordinate, or record intermediate reading observations):
>   - It is **STRICTLY AND ONLY** permitted to write in the system **user temp directory** (`$TEMP`, `%TEMP%`, or `os.path.join(tempfile.gettempdir(), ...)`).
>   - Inside the user temp directory, it MUST create a dedicated folder named after the repository (e.g. `%TEMP%/coding-guidelines/` or `$TEMP/coding-guidelines/`).
>   - Inside `%TEMP%/<repo-name>/`, replicate the necessary mirror structure for scratch data or agent communication.
>   - This is the **ONLY** allowable location for any writing during the reading workflow. Zero repository files may be created or changed.
>   - If agent communication via files is not strictly required, **DO NOT WRITE ANYTHING AT ALL**.

## Reading Strategy: Mandatory Autonomous Looping & Parallel Subagents

> [!TIP]
> **Tool Hierarchy & 3-Tier Fallback Protocol (Strict Zero-Write Mandate):**
> 1. **Tier 1 (Fast Cached Python Tools):**
>    - Use `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> [--ext <extensions>]` and `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<text>" [--path <dir>]` for sub-millisecond cached exploration when available.
>    - **DO NOT** recreate or write Python scripts if they are missing. The read workflow is strictly read-only.
> 2. **Tier 2 (GitMap Acceleration):**
>    - If Python scripts are absent or cannot run, check if GitMap is installed (`gitmap`) and leverage GitMap CLI verbs to rapidly inspect, search, and read repository files with zero disk writes:
>      - **Directory & File Discovery:**
>        - `gitmap list-files` (alias `lf`) `[pattern] [-ext <extensions>]`: List tracked repository files with optional extension filtering.
>          - Example: `gitmap list-files "*"` (list all repository files)
>          - Example: `gitmap lf "*" -ext "md"` (list all markdown files across repository)
>          - Example: `gitmap list-files "02-spec/*"` (list files in spec directory)
>        - `gitmap find` (alias `f`) `<wildcard*>`: High-speed wildcard/glob file search.
>          - Example: `gitmap find "01*" -ext "md"` (find all markdown files starting with 01)
>          - Example: `gitmap find "*index*"` (find all index files)
>        - `gitmap find-files` (alias `ff`) `<exact-name>`: Exact filename lookup.
>          - Example: `gitmap ff "what-to-read.md"` (locate authoritative reading list)
>          - Example: `gitmap find-files "readme.md"` (locate repository readme)
>        - `gitmap find-files-any` (alias `ffa`) `<substring>`: Substring filename search.
>          - Example: `gitmap ffa "guide" -ext "md"` (find all guides)
>        - `gitmap find-files-startswith` (alias `ffs`) `<prefix>`: Prefix filename search.
>          - Example: `gitmap ffs "what-"` (find files starting with prefix)
>        - `gitmap find-files-endswith` (alias `ffe`) `<suffix>`: Suffix filename search.
>          - Example: `gitmap ffe "avoid.md"` (find strictly-avoid files)
>      - **File Reading & Content Inspection:**
>        - `gitmap cat <filepath>`: Stream raw file content directly to stdout without touching disk.
>          - Example: `gitmap cat .ai-memory/what-to-read.md` (read authoritative reading sequence)
>          - Example: `gitmap cat 02-spec/readme.md` (read spec index directly)
>          - Example: `gitmap cat readme.md` (read project identity and guidelines)
>      - **Repository Status & Changelog Context:**
>        - `gitmap status` (alias `st`): Display branch state, clean/dirty working tree, and ahead/behind counts.
>          - Example: `gitmap status`
>        - `gitmap changelog` (alias `cl`) `[--latest]`: Read concise release notes and version history.
>          - Example: `gitmap changelog --latest` (inspect most recent release notes)
>        - `gitmap list-versions` (alias `lv`): List tagged versions in descending order.
>          - Example: `gitmap list-versions --limit 5`
>      - **CI/CD Pipeline & Failure Diagnostics (RCA Context):**
>        - `gitmap pipeline status` (alias `pl status`): Check live CI/CD pipeline state and remaining ETA.
>          - Example: `gitmap pipeline status`
>        - `gitmap pipeline history` (alias `pl history`): Inspect recent commits pipeline execution tree and failure status.
>          - Example: `gitmap pipeline history`
>        - `gitmap pipeline errors` (alias `pl errors`): Fetch and inspect failed step error logs for Root Cause Analysis.
>          - Example: `gitmap pipeline errors` (view latest failed workflow logs)
> 3. **Tier 3 (Native Agent Process Fallback):**
>    - If GitMap is also not installed or available, smoothly fall back to the AI agent's native built-in process and tools (`view_file`, `list_dir`, `grep_search`, `find_by_name`, or standard shell commands like `cat`, `ls`, `grep`).
>    - The AI must seamlessly proceed with reading using its native capabilities without halting, complaining, or writing any files to the repository.

The `.ai-memory/` folder, specs, and entire codebase can be massive. To process this information with zero blind spots:

1. Autonomous looping enforcement:
   - The AI agent MUST autonomously loop through all directories and files across `02-spec/`, `.ai-memory/`, and application source trees.
   - Do not stop after one high-level glance.
   - Systematically iterate through each directory layer.

2. Deep recursive spec traversal:
   - The AI must recursively inspect every subfolder and all nested `.md` files in `02-spec/` (`00-overview.md`, numbered specs `01-*.md`, `99-consistency-report.md`, `spec-index.md`, subdirectories).

3. Anti-hallucination and clarifying questions:
   - Scan internal markdown references across `.ai-memory/` and `02-spec/`.
   - If a referenced spec or issue file is missing on disk, or if requirements are ambiguous, the AI MUST NOT guess or assume.
   - It must automatically log an open question in `.ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md` or ask the user directly to ensure alignment.

4. Active schema and API contract mapping:
   - Ingest and maintain an in-memory map of active DB tables/schemas (`02-02-spec/04-database-conventions/`, `02-02-spec/23-app-db/`), API endpoints, and global state stores so downstream tasks have zero assumptions on field names or parameter shapes.

5. Tooling and runtime compatibility check:
   - Inspect package manifests (`package.json`, `tsconfig.json`, build configs) to catalog runtime targets, linter rules, and banned packages before completing onboarding.

6. Parallel sub-agents for deep reading:
   - You are allowed and strongly encouraged to spawn dedicated sub-agents to read items and synthesize memory in parallel.
   - When spawning a sub-agent for reading, give it a highly specific title reflecting exactly what it is reading (e.g., `Reading Auth Specs in 02-spec/21-app`, `Scanning Error Management in 02-spec/03-error-manage`). Do not use generic names. If an agent switches tasks, its title must change.
   - Assign sub-agents small, granular folders/files to read rather than asking one agent to read the entire codebase in a single pass.

7. Root `readme.md` lowercase self-healing exception:
   - If the root readme is uppercase `README.md` or incorrectly cased, immediately rename it to `readme.md`, commit, and push to git without asking.

8. Memory persistence deferred (Zero repository writes during reading):
   - You are STRICTLY FORBIDDEN from writing to `.ai-memory/` or the repository during the reading workflow.
   - Do NOT write summaries into `.ai-memory/memory/learned/` or update `.ai-memory/readme.md`.
   - Do NOT update `.ai-memory/what-to-read.md` or `.ai-memory/plans/`.
   - Ingest knowledge into memory context and internal state only.
   - If sub-agents need to communicate or record scratch notes during reading, write exclusively to the user temp directory: `%TEMP%/<repo-name>/` (or `$TEMP/<repo-name>/`).

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
   - Path: `.ai-memory/issues/01-<slug>.md` (sequenced as `01-`, `02-`, etc.)
   - Index: Registered in `.ai-memory/readme.md` (or `.ai-memory/plans/readme.md`)
   - Mandatory Structure: Error description, exact file/line location, Root Cause Analysis (one-sentence root cause + deep analysis), fix strategy, and prevention checklist.

2. CI/CD Failure RCAs:
   - Path: `.ai-memory/cicd-issues/01-<slug>.md` (sequenced as `01-`, `02-`, etc.)
   - Index: Registered in `.ai-memory/cicd-index.md` (or `.ai-memory/cicd-index.md`)
   - Mandatory Structure: Raw pipeline error snippet, Root Cause Analysis, resolution applied, and "What NOT to Repeat" rules.

3. Retrospectives & Architectural Failure Learnings:
   - Path: `.ai-memory/memory/learned/01-<slug>.md` and `02-spec/03-error-manage/01-error-resolution/03-retrospectives/01-<slug>.md`
   - Index: Registered in `.ai-memory/memory/readme.md`

4. Hard Avoidances / CODE RED Prohibitions:
   - Path: `.ai-memory/strictly-avoid.md`
   - Hard Rule: Append-only. Never overwrite or truncate existing entries. Any RCA that uncovers a forbidden anti-pattern must append a one-line rule here.

5. Plan Mode Integration:
   - Path: `.ai-memory/plans/pending/01-<slug>.md` & `.ai-memory/plans/subtasks/01-<slug>/01-<subslug>.md`
   - Context Requirement: Must cite previous relevant RCAs and failure records so the new plan explicitly avoids repeating past errors.

6. Anti-Hallucination Guard for RCAs:
   - If an RCA file, referenced spec, or issue file is missing on disk, the AI MUST NOT guess or assume its contents. Stop and file an open ambiguity in `.ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md` or ask the user directly before proceeding.

---

## Phase 1: Load the Project

### 1.0 Read `what-to-read.md` and Confirm Root `readme.md` Lowercase (Auto-Fix & Commit)

1. Read `.ai-memory/what-to-read.md` (or `.ai-memory/what-to-read.md`). This is the authoritative reading order for the project and overrides any generic order. Follow every file and order it specifies.
2. Root `readme.md` lowercase verification and auto-fix:
   - Verify that the root readme file is strictly named lowercase `readme.md`.
   - If an uppercase `README.md` exists or the casing is incorrect on disk or in git, immediately rename it to `readme.md`, remove the stale uppercase file, commit the change (`fix: ensure root readme is strictly lowercase readme.md`), and push to git without asking or second-guessing.
   - Read the root `readme.md` file for architecture, casing rules, repository layout, and AI entry points.

### 1.1 Read the Whole `.ai-memory/` Folder & Flag Markdown Files

Walk `.ai-memory/` recursively. Every file matters. Missing files are noted, not silently skipped.

> [!IMPORTANT]
> **MANDATORY `.ai-memory/*.md` FLAGGING & AUDIT RULE:**
> Every single `.md` file discovered directly in `.ai-memory/` (e.g. `.ai-memory/folder-structure.md`, `.ai-memory/prompts.md`, `.ai-memory/strictly-avoid.md`, `version.json`) or inside any nested subdirectory MUST be actively read, cataloged, and flagged in the memory log. If an unindexed or orphan markdown file is found, immediately raise a flag in the context log and register it into `.ai-memory/memory/readme.md` or `.ai-memory/plans/readme.md`.

| #   | Path                                                  | What you get                                                                                                                                |
| --- | ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `.ai-memory/what-to-read.md`                     | Authoritative reading order for this project. Read it first and follow all referenced files.                                            |
| 2   | `.ai-memory/folder-structure.md`                | Canonical architectural map of all `.ai-memory/` folders, scripts, and naming conventions.                                                  |
| 3   | `.ai-memory/strictly-avoid.md`                          | Hard prohibitions (CODE RED) — append-only, never truncate.                                                                               |
| 4   | `.ai-memory/coding-guidelines.md`     | Master single source of truth for cross-language coding guidelines, `Result[T]` envelopes, and `is`/`has` booleans.                       |
| 5   | `03-ai-scripts/`                            | Automated Python scripts (`01-file-manipulator.py`, `02-guideline-autofixer.py`, `03-cicd-local-runner.py`, `index.md`).                 |
| 6   | `.ai-memory/prompts.md` + `01-prompts/`           | Canonical prompt index and mirrored execution prompts (`cg-execute/`, `execute/`, `ci-cd/`).                                              |
| 7   | `.ai-memory/memory/readme.md`                         | Master index of institutional knowledge. Then read every file it references recursively.                                                 |
| 8   | `.ai-memory/plans/readme.md`                             | Master index of all plans (pending + completed + subtasks). Read this before touching individual plan files.                               |
| 9   | `.ai-memory/plans/pending/`                             | Active plans, `01-<slug>.md` — read all and list out each pending item.                                                                   |
| 10  | `.ai-memory/plans/subtasks/`                            | Granular 5–8 file batch files (`batch-01.md`, etc.) linked from parent plans.                                                              |
| 11  | `.ai-memory/plan.md`                      | Instant recovery state from the most recent run failure.                                                                                   |
| 12  | `.ai-memory/plans/completed/`                           | Historical completed plans and execution logs.                                                                                              |
| 13  | `.ai-memory/issues/`                                    | General bugs and regressions — read and list out pending bugs.                                                                             |
| 14  | `.ai-memory/cicd-issues/`                               | CI/CD-specific failures — read and list out pending CI/CD issues.                                                                           |
| 15  | `.ai-memory/ambiguous-questions/01-new-ambiguity/`      | Open questions currently blocking work. Surface them immediately; do NOT guess past them.                                                  |
| 16  | `.ai-memory/ambiguous-questions/02-ambiguity-resolved/`  | Answered questions with applied solutions. Binding decisions; do not re-litigate.                                                          |
| 17  | `.agents/skills/`                                     | Progressive disclosure runbooks for Antigravity agents (`<slug>/skill.md`).                                                                |
| 18  | `.agents/rules/`                                      | Systemic behavioral and coding rules auto-enforced by Antigravity.                                                                         |
| 19  | `linter-scripts/`                                     | Repository linters (`check-newline-styling.py`, `check-markdown-header-spacing.py`, `check-relative-paths.py`).                            |
| 20  | Anything else under `.ai-memory/`                       | Read it. If the folder exists, it exists for a reason.                                                                                      |

---

### 1.2 Mandatory Skill & Rule Set Auto-Generation on Ingestion

Whenever the AI agent reads prompts or coding guidelines during memory ingestion:

1. **Auto-Generate Antigravity Skills for Prompts:**
   - For every prompt ingested from `01-prompts/` or `01-prompts/`, create or update a dedicated Antigravity skill in `.agents/skills/<slug>/skill.md` with standard YAML frontmatter (`name`, `description`).
   - Ensures any agent in the ecosystem can activate the skill on demand.

2. **Auto-Generate Antigravity Rules for Coding Guidelines:**
   - For all coding guidelines ingested from `.ai-memory/coding-guidelines.md` or `02-spec/02-coding-guidelines/`, synthesize and write authoritative agent rules into `.agents/rules/<slug>.md` and inject essential constraints into `AGENTS.md`.
   - Core rules enforced:
     - **Strict Boolean Standard:** `is, has as prefix is only acceptable and nothing else acceptable including but not limited to can, should etc`.
     - **No Bare Void in Go:** Functions must return `Result[T]` or `*apperror.AppError`.
     - **Parameter Structs:** Banned loose >2-3 parameters; use `*Params` structs.
     - **Vertical Line Gaps:** Mandatory blank lines before `if`, after `}`, before `return`, and around multiline struct calls.
     - **5–8 Files Micro-Batching:** All refactors broken into bounded subtasks.

---

### 1.3 Loop Through the Entire `02-spec/` Directory, Subfolders, and Nested Files

Systematically loop through the `02-spec/` folder, dynamically matching canonical hyphenated names (numbers may vary between projects):

- `02-02-spec/01-spec-authoring-guide/`: Spec authoring conventions, required files, format requirements.
- `02-spec/02-coding-guidelines/` (or `01-prompts/04-coding-standards/01-coding-guidelines.md`): Zero-tolerance coding standards, function size caps (8 lines preferred, 15 max), boolean naming (`is*`, `has*`, positive framing), immutable patterns, DRY priority 1.
- `02-02-spec/03-error-manage/`: Error management philosophy — never swallow errors, log operation name and key inputs on every catch, wrap errors without losing cause, typed errors only, universal response envelopes (`{ data, errors[], meta }`).
- `02-02-spec/04-database-conventions/`: Database schema, table naming (PascalCase), columns (camelCase), primary keys (`{Table}Id`), SQLite/ORM rules, ERD requirements.
- `02-02-spec/05-split-db-architecture/` through `02-02-spec/19-main-worker-service/`: Architectural specs for config, design system, docs viewer, code blocks, CLI, workflows, and release.
- `02-02-spec/21-app/`: App specification, domain architecture, core capabilities, routes, and business rules.
- `02-02-spec/22-app-issues/`, `02-02-spec/23-app-db/`, `02-02-spec/24-app-ui-design-system/`: App-specific issues, schemas, and design systems.
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

- `.ai-memory/memory/readme.md` lists every institutional-knowledge file. If it points at 12 files, you read 12 files.
- `.ai-memory/plans/readme.md` lists every plan (pending, completed, subtasks) with its slug, status, and one-line intent. Use it to pick which plan files to open in full. If it is missing, create it as part of the next code change.

### 1.6 Self-Check (Internal, Before Phase 2)

- CODE RED rules?
- Naming conventions (files, folders, DB columns, variables)?
- Root readme strictly lowercase `readme.md`?
- Error-handling philosophy (`02-02-spec/03-error-manage/`)?
- What is in `.ai-memory/plans/pending/` (sequenced as `01-`, `02-`) and `plans/subtasks/` right now (exact list)?
- Active DB schemas, table columns, and API contracts?
- App specs and domain architecture (`02-02-spec/21-app/`)?
- Whole codebase layout and component flow?
- Top forbidden patterns?

If any answer is fuzzy, go back and reread by looping through the files again. Do not proceed.

---

## Phase 2: Consolidated Guidelines

Read `02-02-spec/17-consolidated-guidelines/` (or `02-02-spec/17-consolidated-guidelines/`) in numeric order (`01-*.md` through `18-*.md`). Each file is a self-contained policy document. Missing folder: note it and continue.

---

## Phase 3: Spec Authoring Rules

Read `02-02-spec/01-spec-authoring-guide/` in numeric order. You should come out knowing:

- File and folder naming conventions (`<NN>-<slug>/`).
- Required files per spec folder (`00-overview.md`, `99-consistency-report.md`).
- The `.ai-memory/` layout (see Phase 1.1).
- The linter infrastructure.

---

## Phase 4: Task-Driven Deep Dives

Only open a spec folder when the current task needs it.

| Task involves...                          | Read                                    |
| ---------------------------------------- | --------------------------------------- |
| Writing or reviewing code                | `02-spec/02-coding-guidelines/`            |
| Error handling                           | `02-02-spec/03-error-manage/`                 |
| Database schema or queries               | `02-02-spec/04-database-conventions/`         |
| SQLite / multi-DB architecture           | `02-02-spec/05-split-db-architecture/`        |
| Config systems                           | `02-02-spec/06-seedable-config-architecture/` |
| UI theming, CSS variables, design tokens | `02-02-spec/07-design-system/`                |
| Documentation viewer features            | `02-02-spec/08-docs-viewer-ui/`               |
| Code block rendering                     | `02-02-spec/09-code-block-system/`            |
| PowerShell scripts                       | `02-02-spec/11-powershell-integration/`       |
| CI/CD pipelines                          | `02-02-spec/12-cicd-pipeline-workflows/`      |
| CLI self-update                          | `02-02-spec/14-update/`                       |
| WordPress plugins                        | `02-02-spec/18-wp-plugin-how-to/`             |
| App-specific features                    | `02-02-spec/21-app/`                          |
| Known app bugs                           | `02-02-spec/22-app-issues/`                   |
| App-specific DB schema                   | `02-02-spec/23-app-db/`                       |
| App-specific UI + design system          | `02-02-spec/24-app-ui-design-system/`         |

Inside each folder: `00-overview.md` -> numbered files -> `99-consistency-report.md`.

Fallbacks when the canonical numbered folder is absent: `.ai-memory/coding-guidelines.md`, `02-spec/02-coding-guidelines/`, `coding-guidelines/`, `02-02-spec/03-error-manage/`, `01-prompts/04-coding-standards/01-coding-guidelines.md`. Numbered folder wins on conflict; call the conflict out in the plan's Context.

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
│   YES → .ai-memory/memory/01-<slug>.md  +  update .ai-memory/memory/readme.md
├─ Must never happen again?
│   YES → .ai-memory/strictly-avoid.md
├─ Idea, not yet approved?
│   YES → .ai-memory/suggestions.md
├─ New user command / convention?
│   YES → .ai-memory/spec/commands/01-<slug>.md
├─ Bug / regression?
│   YES → .ai-memory/issues/01-<slug>.md   (or .ai-memory/cicd-issues/ if CI/CD)
├─ New or changed plan?
│   YES → .ai-memory/plans/pending/01-<slug>.md  +  update .ai-memory/plans/readme.md
├─ Ambiguity / unclear requirement blocking progress?
│   YES → .ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md
├─ User just answered a previously-open ambiguity?
│   YES → mv the file to .ai-memory/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md,
│         append `## Resolution` (answer + applied solution), flip Status: resolved
└─ None of the above → do not persist.
```

Hard rules:

- Folder is `.ai-memory/memory/`, never `memories/`.
- Adding a memory file always updates `.ai-memory/memory/readme.md`.
- Adding, moving, or completing a plan always updates `.ai-memory/plans/readme.md`.
- Ambiguity folders: `01-new-ambiguity/` for open, `02-ambiguity-resolved/` for answered. On answer, MOVE the file (never copy) so it exists in exactly one place. Every resolved file carries a `## Resolution` section.
- Never guess past an open ambiguity. If one exists and is relevant to the current task, stop and surface it before doing work.
- Editing existing memory or index files preserves unrelated content. No silent truncation.
- Any code-base change bumps the minor version.

---

## Completion Confirmation

After Phases 1-3, reply exactly:

```
Onboarding complete.

- Memory files read: [X]
- Consolidated guidelines read: [Y]
- Spec authoring files read: [Z]
- Pending plans: [N]  (from .ai-memory/plans/readme.md)
- CI/CD issues absorbed: [M]  (from .ai-memory/cicd-issues/)
- Open ambiguities: [K]  (from .ai-memory/ambiguous-questions/01-new-ambiguity/)
- Resolved ambiguities on file: [R]  (from .ai-memory/ambiguous-questions/02-ambiguity-resolved/)

I understand:
- CODE RED rules: [top 3-5]
- Naming conventions: [brief]
- Error handling: [one sentence]
- Active DB schemas & contracts: [key models / tables]
- Active plans & pending tasks: [slugs from .ai-memory/plans/pending/ and subtasks]
- Strict avoidances: [top 3-5]
- Blocking ambiguities: [slugs, or "none"]

Ready for tasks.
```

Then stop. No next-step suggestions, no exploratory questions.

---

## Pre-Reply Checklist (All Must Be True)

/goal Complete the checklist properly until done can do self-looping.

1. [ ] /learn `.ai-memory/what-to-read.md` (or `.ai-memory/what-to-read.md`) first and followed its order in full.
2. [ ] Confirmed root readme is strictly lowercase `readme.md` (auto-fixed, committed, and pushed if uppercase or missing).
3. [ ] /learn the root `readme.md` file (casing rules, architecture, entry points).
4. [ ] Walked `.ai-memory/` recursively, no folder or file skipped silently, and flagged all `.ai-memory/*.md` files.
5. [ ] /learn `.ai-memory/memory/readme.md` and every file it points at.
6. [ ] /learn `.ai-memory/plans/readme.md`, every file in `pending/` (sequenced as `01-`, `02-`), and all active subtasks.
7. [ ] Skimmed `.ai-memory/plans/completed/` for recent history.
8. [ ] /learn every file in `.ai-memory/spec/commands/`.
9. [ ] /learn every file in `.ai-memory/issues/` and `.ai-memory/cicd-issues/`.
10. [ ] /learn every file in `.ai-memory/ambiguous-questions/01-new-ambiguity/` and `02-ambiguity-resolved/`.
11. [ ] Scanned for broken links or missing docs and surfaced them under open ambiguities.
12. [ ] Ingested active schema models, DB column conventions, and API route shapes.
13. [ ] Verified runtime dependencies and package compatibility.
14. [ ] Recursively traversed and read every subfolder, nested markdown file (`*.md`), overview, and consistency report within `02-spec/` (e.g. `02-02-spec/01-spec-authoring-guide/`, `02-spec/02-coding-guidelines/`, `02-02-spec/03-error-manage/`, `02-02-spec/04-database-conventions/`, `02-02-spec/21-app/`, etc.).
15. [ ] Autonomously surveyed and looped through the entire codebase as a whole (all application code, entry points, routes, components, state stores, utilities, and configuration files).
16. [ ] /learn `02-02-spec/17-consolidated-guidelines/` (or `02-02-spec/17-consolidated-guidelines/`) in numeric order (or noted missing).
17. [ ] /learn `02-02-spec/01-spec-authoring-guide/` in numeric order (or noted missing).
18. [ ] Can name CODE RED rules, naming conventions, error-handling philosophy without guessing.
19. [ ] Can list every pending plan slug and subtask from memory.
20. [ ] Checked whether the repo contains explicit tone, strictly-avoid, or prior-stupidity instructions and applied them without softening.
21. [ ] Did not replace hard user wording with polite generic language.
22. [ ] Emitted the Completion Confirmation block verbatim, then stopped.
23. [ ] Confirmed that reading remained strictly read-only regarding the codebase (no source code refactored, only memory and lowercase readme auto-fix updated).

## Actionable Items & Checklist

/goal Complete the checklist properly until done can do self-looping.

1. [ ] /learn the coding guidelines in: `.ai-memory/coding-guidelines.md` and create memory.
2. [ ] /learn the condition extraction in: `02-spec/02-coding-guidelines/01-cross-language/readme.md` and create memory.
3. [ ] /learn the formatting and braces in: `02-spec/02-coding-guidelines/01-cross-language/readme.md` and create memory.
4. [ ] /learn the multi-line formatting in: `02-spec/02-coding-guidelines/01-cross-language/readme.md` and create memory.
5. [ ] /learn the boolean guidelines in: `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/readme.md` and create memory.
6. [ ] /learn the anti-hallucination rules in: `02-spec/02-coding-guidelines/01-cross-language/readme.md` and create memory.
7. [ ] /learn the error management architecture in: `02-spec/03-error-manage/readme.md` (and related error manage files) and create memory.
8. [ ] /learn all recent Root Cause Analysis (RCA) files in `.ai-memory/issues/`, `.ai-memory/cicd-issues/`, and `02-spec/03-error-manage/01-error-resolution/03-retrospectives/` to prevent recurring errors.
9. [ ] /learn all hard prohibitions in `.ai-memory/strictly-avoid.md` and verify zero violations.
10. [ ] /learn the enum standards and fixes in: `02-spec/17-consolidated-guidelines/07-enum-standards.md` and `02-spec/17-consolidated-guidelines/07-enum-standards.md` and create memory.
11. [ ] /learn ALL other single-file specs in `02-spec/02-coding-guidelines/` and create memory.
12. [ ] /learn the overarching main task plan.
13. [ ] Ensure the git repository starts completely clean.
14. [ ] /goal Complete all work on the current branch only.
15. [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
16. [ ] Confirmed that reading remained 100% read-only for the repository (zero files modified, zero folder structure changes).
17. [ ] Confirmed zero git commits or pushes triggered during reading.
18. [ ] Confirmed that if any temporary agent communication was required, it was written strictly to the user temp directory (%TEMP%/<repo-name>/ or $TEMP/<repo-name>/) without polluting the repo.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: read-memory-old
- status: active
