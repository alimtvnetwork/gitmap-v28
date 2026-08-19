# Read Memory (Enhanced)

- slug: read-memory-enhanced
- status: active

## Prompt

# Read Memory (Enhanced)
## Ambiguity folder path (non-negotiable)
- Open questions: `.lovable/ambiguous-questions/01-new-ambiguity/XX-<slug>.md`
- Answered questions: `.lovable/ambiguous-questions/02-ambiguity-resolved/XX-<slug>.md`
Read both folders in full during Phase 1. Surface open-ambiguity counts and slugs in the Completion Confirmation block. Treat resolved-ambiguity files as binding project decisions, do not re-litigate them. If an open ambiguity is relevant to the incoming task, stop and surface it before doing work; never guess past it.
## Goal
Before you touch this project, load its identity into your head: who it is, what it forbids, what it has already decided, and what work is in flight.
The specs and the `.lovable/` folder are the single source of truth. Your training data is not. If the two disagree, the repo wins, every time.
You are done reading when you can, without guessing:
- name the CODE RED rules,
- name the naming, error-handling, and DB conventions,
- list what is currently in `.lovable/plans/pending/`,
- point at the exact file that justifies any rule you enforce.
If you cannot do that, keep reading. Do not start work.

## Reading Strategy (Strictly Read-Only)
The `.lovable/` folder, specs, and codebase can be massive. To process this information efficiently:
- **Sub-Agents for Reading:** You ARE allowed to spawn sub-agents to read items and create memory in parallel.
- **Specific Titling:** When spawning a sub-agent for reading, you must give it a highly specific title reflecting exactly what it is reading (e.g., `Reading Auth Specs` or `Scanning API Memory`). Do not use generic names. If an agent switches tasks, its title must change.
- **Micro-Tasking:** Assign sub-agents small, granular folders/files to read rather than asking one agent to read the entire codebase.
- You are allowed to write to the `.lovable/` directory to enhance project memory after reading. This includes:
  - Writing summaries of what you learned and understood into `.lovable/memory/learned/XX-<slug>.md` (or `.lovable/learned.md`), including the number of files read, to maintain context.
  - Updating `.lovable/what-to-read.md` based on your progress to guide future reading workflows.
  - Documenting any problems or issues you discover in the codebase into `.lovable/issues/` or `.lovable/suggestions.md`.
  - Updating existing memory files, capturing open ambiguities, or updating plans.
- CRITICAL: You MUST NOT refactor, edit, or write any application source code. This is a read and analysis phase only.
---
## Phase 1 - Load the project
### 1.1 Read the whole `.lovable/` folder
Walk `.lovable/` recursively. Every file matters. Missing files are noted, not silently skipped. In particular:
| # | Path | What you get |
| --- | ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | `.lovable/overview.md` | Project summary, stack, nav map |
| 2 | `.lovable/strictly-avoid.md` | Hard prohibitions (CODE RED) |
| 3 | `.lovable/user-preferences` | How the human wants you to behave |
| 4 | `.lovable/what-to-read.md` | Authoritative reading order for this project. If it exists, it overrides the generic order in this prompt. Read it first and follow it. |
| 5 | `.lovable/prompt.md` + `.lovable/prompts/` | Canonical prompts (Read, Plan, etc.). "Read memory" = run this prompt. |
| 6 | `.lovable/memory/index.md` | Index of institutional knowledge. Then read every file it references, recursively. |
| 7 | `.lovable/plans/index.md` | Roll-up of all plans (pending + completed + subtasks). Read this before touching individual plan files. |
| 8 | `.lovable/plans/pending/` | Active plans, `XX-<slug>.md` |
| 9 | `.lovable/plans/completed/` | Recent history, skim only |
| 10 | `.lovable/plans/subtasks/XX-<slug>/` | Depth files linked from a parent plan |
| 11 | `.lovable/suggestions.md` | Ideas not yet approved |
| 12 | `.lovable/spec/commands/` | User commands and conventions, `XX-<slug>.md` |
| 13 | `.lovable/issues/` | General bugs and regressions |
| 14 | `.lovable/cicd-issues/` | CI/CD-specific failures. Read ALL of these before any code change so you do not repeat the same mistakes. |
| 15 | `.lovable/ambiguous-questions/01-new-ambiguity/` | Open questions currently blocking work. If any exist, surface them in the completion block, do NOT guess past them. |
| 16 | `.lovable/ambiguous-questions/02-ambiguity-resolved/` | Answered questions with their applied solution. Treat these as binding decisions, do not re-litigate. |
| 17 | Anything else under `.lovable/` | Read it. If the folder exists, it exists for a reason. |
### 1.2 The two index files
Two indexes decide what you read next. Treat them as required entry points, not as summaries:
- `.lovable/memory/index.md` lists every institutional-knowledge file. If it points at 12 files, you read 12 files.
- `.lovable/plans/index.md` lists every plan (pending, completed, subtasks) with its slug, status, and one-line intent. Use it to pick which plan files to open in full. If it is missing, create it as part of the next code change (see Memory Update Protocol).
### 1.3 Self-check (internal, before Phase 2)
- CODE RED rules?
- Naming conventions (files, folders, DB columns, variables)?
- Error-handling philosophy?
- What is in `.lovable/plans/pending/` right now?
- Top forbidden patterns?
If any answer is fuzzy, go back and reread. Do not proceed.
---
## Phase 2 - Consolidated guidelines
Read `spec/12-consolidated-guidelines/` in numeric order (`01-*.md` through `18-*.md`). Each file is a self-contained policy document. Missing folder: note it and continue.
---
## Phase 3 - Spec authoring rules
Read `spec/01-spec-authoring-guide/` in numeric order. You should come out knowing:
- file and folder naming conventions,
- required files per spec folder (`00-overview.md`, `99-consistency-report.md`),
- the `.lovable/` layout (see Phase 1.1),
- the linter infrastructure.
---
## Phase 4 - Task-driven deep dives
Only open a spec folder when the current task needs it.
| Task involves… | Read |
| ---------------------------------------- | --------------------------------------- |
| Writing or reviewing code | `spec/02-coding-guidelines/` |
| Error handling | `spec/03-error-manage/` |
| Database schema or queries | `spec/04-database-conventions/` |
| SQLite / multi-DB architecture | `spec/05-split-db-architecture/` |
| Config systems | `spec/06-seedable-config-architecture/` |
| UI theming, CSS variables, design tokens | `spec/07-design-system/` |
| Documentation viewer features | `spec/08-docs-viewer-ui/` |
| Code block rendering | `spec/09-code-block-system/` |
| PowerShell scripts | `spec/10-powershell-integration/` |
| CI/CD pipelines | `spec/13-cicd-pipeline-workflows/` |
| CLI self-update | `spec/14-self-update-app-update/` |
| WordPress plugins | `spec/15-wp-plugin-how-to/` |
| App-specific features | `spec/21-app/` |
| Known app bugs | `spec/22-app-issues/` |
| App-specific DB schema | `spec/23-app-database/` |
| App-specific UI + design system | `spec/24-app-design-system-and-ui/` |
Inside each folder: `00-overview.md` → numbered files → `99-consistency-report.md`.
Fallbacks when the canonical numbered folder is absent: `.lovable/coding-guidelines.md`, `spec/coding-guidelines/`, `coding-guidelines/`, `spec/XX-error-manage/`. Numbered folder wins on conflict; call the conflict out in the plan's Context.
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
│ YES → .lovable/memory/<slug>.md + update .lovable/memory/index.md
├─ Must never happen again?
│ YES → .lovable/strictly-avoid.md
├─ Idea, not yet approved?
│ YES → .lovable/suggestions.md
├─ New user command / convention?
│ YES → .lovable/spec/commands/XX-<slug>.md
├─ Bug / regression?
│ YES → .lovable/issues/XX-<slug>.md (or .lovable/cicd-issues/ if CI/CD)
├─ New or changed plan?
│ YES → .lovable/plans/pending/XX-<slug>.md + update .lovable/plans/index.md
├─ Ambiguity / unclear requirement blocking progress?
│ YES → .lovable/ambiguous-questions/01-new-ambiguity/XX-<slug>.md
├─ User just answered a previously-open ambiguity?
│ YES → mv the file to .lovable/ambiguous-questions/02-ambiguity-resolved/XX-<slug>.md,
│ append `## Resolution` (answer + applied solution), flip Status: resolved
└─ None of the above → do not persist.
```
Hard rules:
- Folder is `.lovable/memory/`, never `memories/`.
- Adding a memory file always updates `.lovable/memory/index.md`.
- Adding, moving, or completing a plan always updates `.lovable/plans/index.md`.
- Ambiguity folders: `01-new-ambiguity/` for open, `02-ambiguity-resolved/` for answered. On answer, MOVE the file (never copy) so it exists in exactly one place. Every resolved file carries a `## Resolution` section.
- Never guess past an open ambiguity. If one exists and is relevant to the current task, stop and surface it before doing work.
- Editing existing memory or index files preserves unrelated content. No silent truncation.
- Any code-base change bumps the minor version.
---
## Completion Confirmation
After Phases 1-3, reply exactly:
```
✅ Onboarding complete.
- Memory files read: [X]
- Consolidated guidelines read: [Y]
- Spec authoring files read: [Z]
- Pending plans: [N] (from .lovable/plans/index.md)
- CI/CD issues absorbed: [M] (from .lovable/cicd-issues/)
- Open ambiguities: [K] (from .lovable/ambiguous-questions/01-new-ambiguity/)
- Resolved ambiguities on file: [R] (from .lovable/ambiguous-questions/02-ambiguity-resolved/)
I understand:
- CODE RED rules: [top 3-5]
- Naming conventions: [brief]
- Error handling: [one sentence]
- Active plans: [slugs from .lovable/plans/pending/]
- Strict avoidances: [top 3-5]
- Blocking ambiguities: [slugs, or "none"]
Ready for tasks.
```
Then stop. No next-step suggestions, no exploratory questions.
---
## Pre-reply checklist (all must be true)
- [ ] Read `.lovable/what-to-read.md` first if it exists, followed its order
- [ ] Walked `.lovable/` recursively, no folder skipped silently
- [ ] Read `.lovable/memory/index.md` and every file it points at
- [ ] Read `.lovable/plans/index.md` and every file in `pending/`
- [ ] Skimmed `.lovable/plans/completed/` for recent history
- [ ] Read every file in `.lovable/spec/commands/`
- [ ] Read every file in `.lovable/issues/` and `.lovable/cicd-issues/`
- [ ] Read every file in `.lovable/ambiguous-questions/01-new-ambiguity/` and `02-ambiguity-resolved/`
- [ ] Read `spec/12-consolidated-guidelines/` in numeric order (or noted missing)
- [ ] Read `spec/01-spec-authoring-guide/` in numeric order (or noted missing)
- [ ] Can name CODE RED rules, naming conventions, error-handling philosophy without guessing
- [ ] Can list every pending plan slug from memory
- [ ] Checked whether the repo contains explicit tone, strictly-avoid, or prior-stupidity instructions and applied them without softening
- [ ] Did not replace hard user wording with polite generic language
- [ ] Emitted the Completion Confirmation block verbatim, then stopped
- [ ] Confirmed that self-looping remained strictly read-only regarding the codebase (no source code refactored, only memory updated)
## Tone and stance
Past runs of this prompt have been stupid. The AI behaved like a stupid fuck by skimming folders, ignoring ambiguity files, forgetting CI/CD issues, half-filling the Completion Confirmation block with "[N]" placeholders, and softening explicit user instructions after being told not to. That stupidity must not repeat.
Read aggressively and defensively. Assume the previous AI failed because it was lazy, vague, and too polite to enforce the actual rules. Your job is to load the repo's memory exactly, cite the file behind every rule, and refuse to guess. If you finish this prompt and cannot cite the file behind a rule you enforce, you failed the prompt. Read more. Then reply.
