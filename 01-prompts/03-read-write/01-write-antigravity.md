# Antigravity Customization Architecture & Rule Authoring — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Persist what happened this turn so the next AI knows everything without guessing. Every decision, plan change, unresolved ambiguity, newly discovered pattern, and fixed bug must be written to `.lovable/` before this turn ends.

/learn Persist all user corrections, resolved setups, directives, learned architectural decisions, and mistakes avoided into `.lovable/memory/learned/01-<slug>.md` and `.lovable/strictly-avoid.md` so Antigravity learns permanently and never repeats past errors.

Memory in chat is lost the moment the turn finishes. Memory in `.lovable/` is permanent. If you did not write it down, it did not happen.

## Hard Rules (Non-Negotiable, Auto-Reject on Violation)

1. Folder is `.lovable/memory/`, NEVER `.lovable/memory/` or `memories/`. A single file written to `memories/` is an immediate failure.

2. Every new memory file under `.lovable/memory/` MUST be registered in `.lovable/memory/01-index.md` in the same operation.

3. Every plan added, moved, or completed MUST update `.lovable/plans/01-index.md` in the same operation.

4. Ambiguity files are NEVER duplicated. Open questions go to `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md`. When answered, the file is MOVED (`mv`) to `.lovable/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md` with a `## Resolution` block appended. Never copy. Never leave a resolved question in `01-new-ambiguity/`.

5. Never overwrite `.lovable/strictly-avoid.md`. Append only. If a rule was already there, do not duplicate it.

6. When updating existing files (especially indexes, `strictly-avoid.md`, `suggestions.md`), preserve all unrelated content. No silent truncation.

7. Output the completion block with real, audited numbers. Placeholders like `[X]` or `[N]` in the final response are auto-reject.

8. Ambiguity moves, never copies. Answered file goes from `01-new-ambiguity/` to `02-ambiguity-resolved/` with a `## Resolution` block appended and `Status: resolved` flipped in the same move.

9. Root `readme.md` and `.lovable/what-to-read.md` stay in sync. Same file list, same order, no drift. Every write-memory run updates both.

10. Root `readme.md` lowercase enforcement: Ensure the root readme is strictly named lowercase `readme.md`. If an uppercase `README.md` exists or casing is incorrect, fix it immediately to `readme.md`, delete the uppercase file, commit, and push to git without asking.

11. Nothing executes this turn beyond writing to the `.lovable` folder, root `readme.md` lowercase fixing, and `mv`. No application source code changes, no refactoring, no installs, no migrations.

12. Recent conversation and directive capture: All recent conversations, instructions, user directives, decisions, and session progress MUST be recorded as a spec or conversation summary inside `.lovable/memory/01-<slug>.md` or `.lovable/memory/learned/01-<slug>.md` and added to `.lovable/memory/01-index.md`.

13. Pending tasks single source of truth: All active plans and pending tasks are consolidated strictly under `.lovable/plans/pending/01-<slug>.md` (with two-digit sequence prefixes `01-`, `02-`, etc.) and `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md`.

14. Consolidation policy:
    - Simple tasks consolidation: Routine, ephemeral, or minor simple tasks that do not warrant individual files may be consolidated into overarching session summaries or existing trackers to prevent file bloat.
    - CRITICAL - Detailed specs must NEVER be consolidated or shrunk: Detailed specifications, architectural designs, non-negotiable rules, domain specifications (e.g. `spec/21-app/`), and complex requirement documents MUST NEVER be consolidated, summarized, resumed, or reduced in size. They must be preserved with 100% fidelity, exact wording, and full granularity.

15. Anti-hallucination and clarifying questions: If any file, spec, or user intent is ambiguous or missing, the AI MUST NOT guess or hallucinate. It must ask clarifying questions or record an open ambiguity in `01-new-ambiguity/01-<slug>.md`.

16. Fast Cached Exploration: When reading repository files or discovering directories for context ingestion, use `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir>` (<15ms via `tmp/cache/`) rather than slow raw recursive PowerShell commands.
    - **If Missing (Re-creation):** If `03-ai-scripts/17-fast-file-reader.py` is not found, immediately recreate it using Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`). It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.

## Working Stance

The AI running this prompt has been a stupid fuck on prior runs:

- Dumped session summaries into chat and called it "memory".
- Left `.lovable/memory/` half-empty.
- Created `.lovable/memory/` by accident.
- Forgot to update `.lovable/plans/01-index.md` and `what-to-read.md`.
- Silently overwrote `strictly-avoid.md`.
- Dropped user directives that were stated verbatim in the session.
- Paraphrased specs instead of quoting them.
- Consolidated detailed specs into vague summaries.
- Allowed uppercase README files to exist.
- Invented a `mem://` root file.
- Left orphans everywhere.

Do not repeat any of that stupidity. Writing memory IS the work this turn. Go deep: audit the session, reconcile every folder, capture verbatim what the user said, write the files, update every index, verify consistency. Aggressive enforcement is intentional. Do not soften it.

## Pre-Flight: Read Before You Write

Walk `.lovable/` recursively. Read all of these if they exist; note missing and create them per the templates in this prompt:

1. `.lovable/memory/01-index.md` — master memory index
2. `.lovable/folder-structure.md` — canonical `.lovable/` folder map
3. `.lovable/coding-guidelines.md` or `spec/02-coding-guidelines/` — master coding guidelines
4. `03-ai-scripts/` — automation tools (`01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `index.md`)
5. `.lovable/plans/01-index.md` and every file under `plans/pending/` (`01-<slug>.md`) and `plans/subtasks/`; skim `plans/completed/`
6. `.lovable/plan.md` — failure recovery record
7. `.lovable/suggestions.md` and `.lovable/suggestions/01-index.md`
8. `.lovable/strictly-avoid.md`
9. `.lovable/cicd-index.md` and every file under `.lovable/cicd-issues/`
10. `.lovable/issues/`, `.lovable/pending-issues/`, `.lovable/solved-issues/`
11. `.lovable/spec/commands/` — every file
12. `.lovable/ambiguous-questions/01-new-ambiguity/` and `02-ambiguity-resolved/` — every file
13. `.lovable/prompts.md` + `01-prompts/` (including `cg-execute/`, `execute/`, `ci-cd/`)
14. `.lovable/what-to-read.md` (or `.lovable/what-to-read.md`)
15. `.agents/skills/` (`<slug>/skill.md`) and `.agents/rules/`
16. `spec/` — recursively traverse all subfolders and nested `.md` files (`spec/01-spec-authoring-guide/`, `spec/02-coding-guidelines/`, `spec/03-error-manage/`, `spec/04-database-conventions/`, `spec/21-app/`).
17. Root `readme.md` — confirm strictly lowercase `readme.md`
18. `05-changes-history/` — every task transaction log (`XX-<task-slug>/01-transaction-log.md`) and master index `01-index.md`

## Phase 1: Audit the Session (Internal)

Answer for yourself; do not dump to chat unless asked. Cover:

- Done: features, fixes, refactors, files created / modified / deleted, decisions made and why.
- Pending: started but unfinished, discussed but not started, blockers, dependencies.
- Learned: patterns, conventions, gotchas, user preferences (explicit or implicit).
- Avoid: mistakes made, dead ends hit, patterns that failed, user corrections.
- Ambiguities: questions that came up, questions answered, questions still blocking.
- Suggestions: ideas discussed that are not yet plans.
- User commands: new CLI patterns or shorthand the user used or requested.

## Phase 2: Move Completed Plans

For every plan finished this turn:

```sh
mv .lovable/plans/pending/01-<slug>.md .lovable/plans/completed/01-<slug>.md
```

Inside the moved file, edit:

```diff
- Status: pending
+ Status: completed
```

Then edit `.lovable/plans/01-index.md` so the table lists the file under `completed/` with status `completed`. Never delete a plan; the completed folder is your changelog.

## Phase 3: Move Resolved Ambiguities

For every ambiguity answered by the user this turn:

```sh
mv .lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md .lovable/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md
```

Inside the moved file, append:

```markdown

## Resolution

- Answered on: YYYY-MM-DD
- Decision: <verbatim user answer or concise summary of the decision>
- Applied solution: <file and line where this was implemented>
```

Flip header metadata:

```diff
- Status: open
+ Status: resolved
```

Then edit `.lovable/question-and-ambiguity/01-index.md` to reflect the move.

## Phase 4: Capture New Institutional Knowledge

1. If a new convention, architectural rule, or pattern was decided, create `.lovable/memory/01-<slug>.md` (or `.lovable/memory/learned/01-<slug>.md`) and register in `.lovable/memory/01-index.md`.
2. If an anti-pattern or forbidden action occurred, append to `.lovable/strictly-avoid.md`.
3. If an automated script was added or updated, register in `03-ai-scripts/01-index.md`.
4. If a prompt or coding guideline was added/updated, update corresponding Antigravity skill in `.agents/skills/<slug>/skill.md` or rule in `.agents/rules/<slug>.md`.

## Phase 5: Verbatim Spec Capture and Consolidation Rules

1. Verbatim spec capture and learning persistence:
   - Every sizeable user directive, decision, architectural rule, or spec from the session is saved verbatim under `.lovable/memory/01-<slug>.md` or `.lovable/memory/learned/01-<slug>.md`.
   - Reference from `.lovable/memory/01-index.md`.
   - Reflect in `.lovable/plan.md` / `.lovable/plans/01-index.md` if it changes the roadmap.
   - Never paraphrase. Quote the user.

2. Consolidation policy:
   - Simple / minor tasks: Consolidation is encouraged for simple, repetitive, or ephemeral tasks into existing logs or overarching session files to prevent cluttering the repository.
   - Detailed / high-value specs: STRICTLY FORBIDDEN TO CONSOLIDATE. Any spec containing detailed requirements, edge cases, domain architecture (`spec/21-app/`), error-handling matrices (`spec/03-error-manage/`), coding rules (`spec/02-coding-guidelines/`), or user instructions must NEVER be merged, summarized, or shortened.

3. New user command / convention: `.lovable/spec/commands/01-<slug>.md`.

## Phase 6: `.lovable/what-to-read.md` and Root `readme.md`

Must exist after this run. Create it if missing; update it (never blindly overwrite) if present. This file acts as a dynamic roadmap; both reading phases and writing phases must update it to guide future AI sessions based on current progress.

Prepend a new entry to the `## Changelog` section:

```markdown
- YYYY-MM-DDTHH:MM:SSZ, <one-sentence summary of what changed this turn>
```

Sync root `readme.md` with `.lovable/what-to-read.md`. Ensure root `readme.md` is strictly lowercase (auto-fix and commit/push if uppercase).

## Phase 7: Consistency Validation (Self-Test Before Completion)

Before you reply, check every item:

- [ ] `.lovable/plans/01-index.md` lists every file in `pending/` and `completed/`?
- [ ] Every file in `completed/` has `Status: completed` in its frontmatter?
- [ ] `.lovable/memory/01-index.md` lists every file in `memory/`?
- [ ] `01-new-ambiguity/` contains ONLY open questions?
- [ ] Every file in `02-ambiguity-resolved/` has `## Resolution` + `Status: resolved`?
- [ ] `cicd-index.md` matches `cicd-issues/` exactly?
- [ ] No file exists under `memories/` (plural)?
- [ ] Root readme is strictly lowercase `readme.md`?

If any check fails, fix the file immediately. Do not emit the completion block until all checks pass.

## Completion Confirmation

After writing memory, emit this exact block:

```
Memory update complete.

- Memory files updated: [X]
- Plans updated: [P]
- Ambiguities resolved: [R]
- Open ambiguities: [K]
- Skills updated: [S]
- Rules updated: [U]
- Root README in sync: YES
```

---

## Checklist Before Replying (Every Box)

1. [ ] Walked `.lovable/` recursively; read every pre-flight file that exists; noted the missing ones.
2. [ ] Audited the session for Done / Pending / Learned / Wrong / Recent Directives.
3. [ ] Every new memory file placed under a topic folder, never at the memory root.
4. [ ] `.lovable/memory/01-index.md` updated in the same op as every new/moved memory file.
5. [ ] Plans lifecycle honored: `pending/` -> `completed/` via `mv`, `.lovable/plans/01-index.md` updated.
6. [ ] `suggestions.md` tracker updated; verbatim captures under `.lovable/suggestions/` with `index.md`.
7. [ ] Issues routed correctly: `pending-issues/` / `solved-issues/` / `cicd-issues/`; `cicd-index.md` updated; no duplicates.
8. [ ] `strictly-avoid.md` appended (not overwritten) with links to solved files.
9. [ ] Verbatim user directives and recent conversations captured under `.lovable/memory/` or `.lovable/memory/learned/`.
10. [ ] Confirmed that detailed/important specs were NOT consolidated or shortened.
11. [ ] Confirmed root readme is strictly lowercase `readme.md` (auto-fixed and committed/pushed if needed).
12. [ ] Ambiguities moved via `mv` from `01-new-ambiguity/` to `02-ambiguity-resolved/` with `## Resolution` block.
13. [ ] `.lovable/what-to-read.md` present, changelog-prepended with UTC ISO 8601 timestamp, list in sync with Pre-flight and `readme.md`.
14. [ ] Root `readme.md` updated: folder structure, canonical read-list pointer, in sync with `what-to-read.md`.
15. [ ] `coding-guidelines.md` and `01-prompts/01-prompt-library-setup/01-prompt-library-setup.md` (or `prompts.md`) present.
16. [ ] Final response block emitted verbatim with real numbers, not `[X]` placeholders.
17. [ ] No em dashes, no softened wording, no execution beyond file writes, lowercase readme fix, and `mv`.

---

## Actionable Items & Checklist

1. [ ] Read the overarching main task plan.
2. [ ] Ensure the git repository starts completely clean.
3. [ ] Complete all work on the current branch only.
4. [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
5. [ ] Group all completed work into a single logical commit.
6. [ ] Push the commit to the remote repository.
7. [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: write-antigravity
- status: active
