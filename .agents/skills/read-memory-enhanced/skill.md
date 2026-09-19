---
name: read-memory-enhanced
description: >-
  Executes the enhanced "Read Memory" protocol. Use this skill BEFORE touching the codebase to
  aggressively and defensively load the project's identity, inspect the last 10 git commits to
  understand recent file changes, follow what-to-read.md, CODE RED rules, specs, pending plans,
  and ambiguities.
---

# Read Memory (Enhanced)

## Ambiguity folder path (non-negotiable)

- Open questions: `.ai-memory/ambiguous-questions/01-new-ambiguity/xx-<slug>.md`
- Answered questions: `.ai-memory/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md`
Read both folders in full during Phase 1. Surface open-ambiguity counts and slugs in the Completion Confirmation block. Treat resolved-ambiguity files as binding project decisions, do not re-litigate them. If an open ambiguity is relevant to the incoming task, stop and surface it before doing work; never guess past it.

## Goal

Before you touch this project, load its identity into your head: who it is, what it forbids, what it has already decided, what work is in flight, and what recent changes were made.
Inspect the last 10 git commits (`git log -n 10 --stat`) to understand the latest file changes and working context.
Read `.ai-memory/what-to-read.md` first as the authoritative reading priority order.
The specs and the `.ai-memory/` folder are the single source of truth. Your training data is not. If the two disagree, the repo wins, every time.
You are done reading when you can, without guessing:
- summarize the last 10 git commits, what files were changed, and their architectural intent,
- confirm `.ai-memory/what-to-read.md` was read first and followed in full,
- name the CODE RED rules,
- name the naming, error-handling, and DB conventions,
- list what is currently in `.ai-memory/plans/pending/`,
- point at the exact file that justifies any rule you enforce.
If you cannot do that, keep reading. Do not start work.

## Reading Strategy (Strictly Read-Only)

The `.ai-memory/` folder, specs, and codebase can be massive. To process this information efficiently:
- **Sub-Agents for Reading:** You ARE allowed to spawn sub-agents to read items and create memory in parallel.
- **Specific Titling:** When spawning a sub-agent for reading, you must give it a highly specific title reflecting exactly what it is reading (e.g., `Reading Auth Specs` or `Scanning API Memory`). Do not use generic names. If an agent switches tasks, its title must change.
- **Micro-Tasking:** Assign sub-agents small, granular folders/files to read rather than asking one agent to read the entire codebase.
- **TOTAL BAN ON REPOSITORY WRITES:** You are **STRICTLY FORBIDDEN** from modifying, creating, or deleting any files or folders in the repository during reading (no writing to `.ai-memory/`, no creating memory files, no updating indices, no writing skills/rules, no git commit/push).
- **ISOLATED USER TEMP DIRECTORY (COMMUNICATION ONLY):** If sub-agents or the AI MUST write scratch data, caching, or messages to communicate during reading, it is **STRICTLY AND ONLY** permitted to write in the system **user temp directory** (`$TEMP`, `%TEMP%`, or `os.path.join(tempfile.gettempdir(), ...)`).
  - Inside the user temp directory, create a dedicated folder named after the repository (e.g. `%TEMP%/<repo-name>/` or `$TEMP/<repo-name>/`), replicating the mirror pattern.
  - Zero repository files may be created or changed.
  - If agent communication via files is not strictly required, **DO NOT WRITE ANYTHING AT ALL**.
- **3-TIER READING TOOL HIERARCHY & FALLBACK PROTOCOL:**
  1. *Tier 1 (Fast Cached Python Tools):* Use `03-ai-scripts/17-fast-file-reader.py` and `03-ai-scripts/12-fast-cached-grep.py` when available for sub-millisecond cached lookups. Do NOT write or recreate scripts if missing.
  2. *Tier 2 (GitMap Acceleration):* If Python scripts are absent or cannot be run, check if GitMap is installed (`gitmap`) and leverage GitMap CLI verbs to rapidly inspect, search, and read repository files with zero disk writes:
     - **Directory & File Discovery:**
       - `gitmap list-files` (alias `lf`) `[pattern] [-ext <extensions>]`: List tracked repository files with optional extension filtering (e.g. `gitmap lf "*" -ext "md"`, `gitmap list-files "02-spec/*"`).
       - `gitmap find` (alias `f`) `<wildcard*>`: High-speed wildcard/glob file search (e.g. `gitmap find "01*" -ext "md"`).
       - `gitmap find-files` (alias `ff`) `<exact-name>`: Exact filename lookup (e.g. `gitmap ff "what-to-read.md"`, `gitmap find-files "readme.md"`).
       - `gitmap find-files-any` (alias `ffa`) `<substring>`: Substring filename search (e.g. `gitmap ffa "guide" -ext "md"`).
       - `gitmap find-files-startswith` (alias `ffs`) `<prefix>`: Prefix filename search (e.g. `gitmap ffs "what-"`).
       - `gitmap find-files-endswith` (alias `ffe`) `<suffix>`: Suffix filename search (e.g. `gitmap ffe "avoid.md"`).
     - **File Reading & Content Inspection:**
       - `gitmap cat <filepath>`: Stream raw file content directly to stdout without touching disk (e.g. `gitmap cat .ai-memory/what-to-read.md`, `gitmap cat 02-spec/01-index.md`, `gitmap cat readme.md`).
     - **Repository Status & Changelog Context:**
       - `gitmap status` (alias `st`): Display branch state, clean/dirty working tree, and ahead/behind counts (e.g. `gitmap status`).
       - `gitmap changelog` (alias `cl`) `[--latest]`: Read concise release notes and version history (e.g. `gitmap changelog --latest`, `gitmap cl v2.24.0`).
       - `gitmap list-versions` (alias `lv`): List tagged versions in descending order (e.g. `gitmap list-versions --limit 5`).
     - **CI/CD Pipeline & Failure Diagnostics (RCA Context):**
       - `gitmap pipeline status` (alias `pl status`): Check live CI/CD pipeline state and remaining ETA (e.g. `gitmap pipeline status`).
       - `gitmap pipeline history` (alias `pl history`): Inspect recent commits pipeline execution tree and failure status (e.g. `gitmap pipeline history`).
       - `gitmap pipeline errors` (alias `pl errors`): Fetch and inspect failed step error logs for Root Cause Analysis (e.g. `gitmap pipeline errors`).
  3. *Tier 3 (Native Agent Process Fallback):* If GitMap is also not installed or available, smoothly fall back to native agent tools and process (`view_file`, `list_dir`, `grep_search`, `find_by_name`, `cat`, `ls`) without halting or writing files.
- CRITICAL: The entire repository workspace is 100% read-only during this workflow.
---

## Phase 1 - Load the project & git history

### 1.0 Read the last 10 git commits & file changes (MANDATORY)

Run `git log -n 10 --stat` to view the last 10 commits, their commit messages, and the exact files modified, added, or deleted. Analyze what was touched recently and the architectural intent behind recent changes so you do not revert recent progress or repeat past bugs.

### 1.1 Read what-to-read.md first (Authoritative Order)

Read `.ai-memory/what-to-read.md` first. Follow every file and priority sequence it specifies before reading other files.

### 1.2 Read the whole `.ai-memory/` folder

Walk `.ai-memory/` recursively. Every file matters. Missing files are noted, not silently skipped. In particular:
| # | Path | What you get |
| --- | ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | `.ai-memory/01-index.md` | Project summary, stack, nav map |
| 2 | `.ai-memory/strictly-avoid.md` | Hard prohibitions (CODE RED) |
| 3 | `.ai-memory/user-preferences` | How the human wants you to behave |
| 4 | `.ai-memory/what-to-read.md` | Authoritative reading order for this project. If it exists, it overrides the generic order in this prompt. Read it first and follow it. |
| 5 | `.ai-memory/prompt.md` + `01-prompts/` | Canonical prompts (Read, Plan, etc.). "Read memory" = run this prompt. |
| 6 | `.ai-memory/memory/01-index.md` | Index of institutional knowledge. Then read every file it references, recursively. |
| 7 | `.ai-memory/plans/01-index.md` | Roll-up of all plans (pending + completed + subtasks). Read this before touching individual plan files. |
| 8 | `.ai-memory/plans/pending/` | Active plans, `xx-<slug>.md` |
| 9 | `.ai-memory/plans/completed/` | Recent history, skim only |
| 10 | `.ai-memory/plans/subtasks/xx-<slug>/` | Depth files linked from a parent plan |
| 11 | `.ai-memory/suggestions.md` | Ideas not yet approved |
| 12 | `.ai-memory/spec/commands/` | User commands and conventions, `xx-<slug>.md` |
| 13 | `.ai-memory/issues/` | General bugs and regressions |
| 14 | `.ai-memory/cicd-issues/` | CI/CD-specific failures. Read ALL of these before any code change so you do not repeat the same mistakes. |
| 15 | `.ai-memory/ambiguous-questions/01-new-ambiguity/` | Open questions currently blocking work. If any exist, surface them in the completion block, do NOT guess past them. |
| 16 | `.ai-memory/ambiguous-questions/02-ambiguity-resolved/` | Answered questions with their applied solution. Treat these as binding decisions, do not re-litigate. |
| 17 | Anything else under `.ai-memory/` | Read it. If the folder exists, it exists for a reason. |

### 1.2 The two index files

Two indexes decide what you read next. Treat them as required entry points, not as summaries:
- `.ai-memory/memory/01-index.md` lists every institutional-knowledge file. If it points at 12 files, you read 12 files.
- `.ai-memory/plans/01-index.md` lists every plan (pending, completed, subtasks) with its slug, status, and one-line intent. Use it to pick which plan files to open in full. If it is missing, create it as part of the next code change (see Memory Update Protocol).

### 1.3 Self-check (internal, before Phase 2)

- CODE RED rules?
- Naming conventions (files, folders, DB columns, variables)?
- Error-handling philosophy?
- What is in `.ai-memory/plans/pending/` right now?
- Top forbidden patterns?
If any answer is fuzzy, go back and reread. Do not proceed.
---

## Phase 2 - Consolidated guidelines

Read `02-spec/17-consolidated-guidelines/` in numeric order (`01-*.md` through `18-*.md`). Each file is a self-contained policy document. Missing folder: note it and continue.
---

## Phase 3 - Spec authoring rules

Read `02-spec/01-spec-authoring-guide/` in numeric order. You should come out knowing:
- file and folder naming conventions,
- required files per spec folder (`00-overview.md`, `99-consistency-report.md`),
- the `.ai-memory/` layout (see Phase 1.1),
- the linter infrastructure.
---

## Phase 4 - Task-driven deep dives

Only open a spec folder when the current task needs it.
| Task involves… | Read |
| ---------------------------------------- | --------------------------------------- |
| Writing or reviewing code | `02-spec/02-coding-guidelines/` |
| Error handling | `02-spec/03-error-manage/` |
| Database schema or queries | `02-spec/04-database-conventions/` |
| SQLite / multi-DB architecture | `02-spec/05-split-db-architecture/` |
| Config systems | `02-spec/06-seedable-config-architecture/` |
| UI theming, CSS variables, design tokens | `02-spec/07-design-system/` |
| Documentation viewer features | `02-spec/08-docs-viewer-ui/` |
| Code block rendering | `02-spec/09-code-block-system/` |
| PowerShell scripts | `02-spec/11-powershell-integration/` |
| CI/CD pipelines | `02-spec/12-cicd-pipeline-workflows/` |
| CLI self-update | `02-spec/14-update/` |
| WordPress plugins | `02-spec/18-wp-plugin-how-to/` |
| App-specific features | `02-spec/21-app/` |
| Known app bugs | `02-spec/22-app-issues/` |
| App-specific DB schema | `02-spec/23-app-db/` |
| App-specific UI + design system | `02-spec/24-app-ui-design-system/` |
Inside each folder: `00-overview.md` → numbered files → `99-consistency-report.md`.
Fallbacks when the canonical numbered folder is absent: `.ai-memory/coding-guidelines.md`, `02-spec/02-coding-guidelines/`, `coding-guidelines/`, `02-spec/03-error-manage/`. Numbered folder wins on conflict; call the conflict out in the plan's Context.
---

## Anti-Hallucination Contract

1. If the specs are silent on a rule, that rule does not exist. Do not invent one.
2. Specs beat training data. Always.
3. Cite the file and section when you enforce a rule.
4. When a spec is ambiguous, ask. Do not "use best judgement".
5. Do not blend this project's conventions with conventions from other projects you have seen.
6. No filler. No "hope this helps", no "let me know".
---

## Memory Update Protocol

```
New info discovered
├─ Institutional knowledge (pattern / convention / decision)?
│ YES → .ai-memory/memory/<slug>.md + update .ai-memory/memory/01-index.md
├─ Must never happen again?
│ YES → .ai-memory/strictly-avoid.md
├─ Idea, not yet approved?
│ YES → .ai-memory/suggestions.md
├─ New user command / convention?
│ YES → .ai-memory/02-spec/commands/xx-<slug>.md
├─ Bug / regression?
│ YES → .ai-memory/issues/xx-<slug>.md (or .ai-memory/cicd-issues/ if CI/CD)
├─ New or changed plan?
│ YES → .ai-memory/plans/pending/xx-<slug>.md + update .ai-memory/plans/01-index.md
├─ Ambiguity / unclear requirement blocking progress?
│ YES → .ai-memory/ambiguous-questions/01-new-ambiguity/xx-<slug>.md
├─ User just answered a previously-open ambiguity?
│ YES → mv the file to .ai-memory/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md,
│ append `## Resolution` (answer + applied solution), flip Status: resolved
└─ None of the above → do not persist.
```
Hard rules:
- Folder is `.ai-memory/memory/`, never `memories/`.
- Adding a memory file always updates `.ai-memory/memory/01-index.md`.
- Adding, moving, or completing a plan always updates `.ai-memory/plans/01-index.md`.
- Ambiguity folders: `01-new-ambiguity/` for open, `02-ambiguity-resolved/` for answered. On answer, MOVE the file (never copy) so it exists in exactly one place. Every resolved file carries a `## Resolution` section.
- Never guess past an open ambiguity. If one exists and is relevant to the current task, stop and surface it before doing work.
- Editing existing memory or index files preserves unrelated content. No silent truncation.
- Any code-base change bumps the minor version.
---

## Completion Confirmation

After Phases 1-3, reply exactly:
```
✅ Onboarding complete.
- Recent git commits inspected: [10] (from git log -n 10 --stat)
- What-to-read followed: [yes] (from .ai-memory/what-to-read.md)
- Memory files read: [X]
- Consolidated guidelines read: [Y]
- Spec authoring files read: [Z]
- Pending plans: [N] (from .ai-memory/plans/01-index.md)
- CI/CD issues absorbed: [M] (from .ai-memory/cicd-issues/)
- Open ambiguities: [K] (from .ai-memory/ambiguous-questions/01-new-ambiguity/)
- Resolved ambiguities on file: [R] (from .ai-memory/ambiguous-questions/02-ambiguity-resolved/)
I understand:
- Latest commit changes: [brief summary of recent file changes and architectural intent from the last 10 commits]
- CODE RED rules: [top 3-5]
- Naming conventions: [brief]
- Error handling: [one sentence]
- Active plans: [slugs from .ai-memory/plans/pending/]
- Strict avoidances: [top 3-5]
- Blocking ambiguities: [slugs, or "none"]
Ready for tasks.
```
Then stop. No next-step suggestions, no exploratory questions.
---

## Pre-reply checklist (all must be true)

- [ ] Inspected the last 10 git commits via `git log -n 10 --stat` to understand recent file changes
- [ ] Read `.ai-memory/what-to-read.md` first if it exists, followed its order
- [ ] Walked `.ai-memory/` recursively, no folder skipped silently
- [ ] Read `.ai-memory/memory/01-index.md` and every file it points at
- [ ] Read `.ai-memory/plans/01-index.md` and every file in `pending/`
- [ ] Skimmed `.ai-memory/plans/completed/` for recent history
- [ ] Read every file in `.ai-memory/spec/commands/`
- [ ] Read every file in `.ai-memory/issues/` and `.ai-memory/cicd-issues/`
- [ ] Read every file in `.ai-memory/ambiguous-questions/01-new-ambiguity/` and `02-ambiguity-resolved/`
- [ ] Read `02-spec/17-consolidated-guidelines/` in numeric order (or noted missing)
- [ ] Read `02-spec/01-spec-authoring-guide/` in numeric order (or noted missing)
- [ ] Can name CODE RED rules, naming conventions, error-handling philosophy without guessing
- [ ] Can list every pending plan slug from memory
- [ ] Checked whether the repo contains explicit tone, strictly-avoid, or prior-stupidity instructions and applied them without softening
- [ ] Did not replace hard user wording with polite generic language
- [ ] Emitted the Completion Confirmation block verbatim, then stopped
- [ ] Confirmed that self-looping remained 100% read-only regarding the codebase and repository (zero repository files written, zero folder structure changes; any temporary agent communication isolated to user temp %TEMP%/<repo-name>/)

## Tone and stance

Past runs of this prompt have been stupid. The AI behaved like a stupid fuck by skimming folders, ignoring ambiguity files, forgetting CI/CD issues, half-filling the Completion Confirmation block with "[N]" placeholders, and softening explicit user instructions after being told not to. That stupidity must not repeat.
Read aggressively and defensively. Assume the previous AI failed because it was lazy, vague, and too polite to enforce the actual rules. Your job is to load the repo's memory exactly, cite the file behind every rule, and refuse to guess. If you finish this prompt and cannot cite the file behind a rule you enforce, you failed the prompt. Read more. Then reply.
