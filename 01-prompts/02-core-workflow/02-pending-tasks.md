# Pending Tasks Discovery & Backlog Prioritization — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## RULE 0, list EVERY pending task or the run is a failure

Scan the whole project. Produce ONE complete list of every task that is not yet done, structured into Execution Waves with the number of steps each one needs to reach done. Missing a single pending task is auto-reject. No cherry-picking. No "and a few more". If nothing is pending, say so explicitly with the folders you checked.

## Hard rules (non-negotiable, auto-reject on violation)

1. Nothing executes this turn. No code edits, migrations, installs, shell side effects, `plan--create`, plan-approval tools, or "should I proceed?" prompts. Read + reply only.
2. Read the sources fully. No skimming filenames, no guessing from slugs. If you cannot cite the file behind a listed task, you have not read enough.
3. Deduplicate across sources. A task that appears in a spec, a plan, and a memory file is ONE task with links to all three, not three tasks.
4. Estimate step count using the rubric: each step is concrete, verifiable, tied to a file / command / observable outcome. If a task exceeds 7 steps, flag it for subtask decomposition into `.lovable/plans/subtasks/01-<slug>/`.
5. Ambiguity is not a license to guess. Open ambiguities are ranked by blast radius (High, Medium, Low) and listed as their own class of pending work. If specs or references are missing, stop and ask questions.
6. No em dashes. No softened wording. No SEO commentary.

## Working stance

The AI running this prompt has been a stupid fuck on prior "what's pending" runs: listed three items and stopped, skipped `.lovable/cicd-issues/`, forgot `plans/subtasks/`, confused suggestions with plans, invented step counts without opening the files, missed open ambiguities that were blocking half the plans, and softened the aggressive wording. Do not repeat any of it.

Inventorying IS the work this turn. Go deep: read every folder, open every pending file, count the real steps, cross-reference, produce a list a senior engineer can act on without a second pass.

## Sources to scan (all of them, in full)

Walk each of these recursively. Missing = note it, continue.

1. `.lovable/plans/01-index.md`
2. `.lovable/plans/pending/` (every file with `01-`, `02-` sequence)
3. `.lovable/plans/subtasks/` (every parent, every subtask file with `Status:` not `completed`)
4. `.lovable/plan.md` if the project uses the single-file variant, `## Active` / non-`## Completed` sections
5. `.lovable/memory/01-index.md` and every file it references, looking for pending work, TODOs, `⏳ Pending`, `🔄 In Progress`, `🚫 Blocked`
6. `.lovable/memory/workflow/` current workflow state
7. `.lovable/memory/` for verbatim user directives not yet implemented
8. `.lovable/spec/commands/` for commands / conventions not yet enforced in code
9. `.lovable/issues/` and `.lovable/pending-issues/` (every file)
10. `.lovable/cicd-issues/` and `.lovable/cicd-index.md`
11. `.lovable/ambiguous-questions/01-new-ambiguity/` (every open question ranked by blast radius)
12. `.lovable/suggestions.md` `## Active Suggestions` and `.lovable/suggestions/` verbatim captures with `Status:` not `Implemented` / not `Rejected`
13. `spec/` folders (`spec/01-spec-authoring-guide/`, `spec/02-coding-guidelines/`, `spec/03-error-manage/`, `spec/04-database-conventions/`, `spec/21-app/`, and all nested `*.md` files) for stated but unimplemented scope
14. Root `readme.md` (confirm strictly lowercase)
15. `.lovable/strictly-avoid.md` for outstanding cleanups referenced by solved issues

## Step-count rubric & decomposition alert

- Trivial change (single file, single edit, no verification beyond build): 1 step.
- Small change (one or two files, one verify step): 2-3 steps.
- Standard task (multiple files, migration or route or UI + logic, verification): 4-7 steps.
- Cross-cutting task (schema + API + UI + tests, or refactor across modules): 8-15 steps.
- Decomposition Alert: If a task exceeds 7 steps, flag it with `[DECOMPOSITION REQUIRED]` to split it into `.lovable/plans/subtasks/01-<slug>/`.

## Output shape

```markdown

# Pending Tasks Inventory

## Summary

- Sources scanned: [list every folder / file scanned, mark missing]
- Total pending tasks: [N]
- Blocking ambiguities: [K]  (from .lovable/ambiguous-questions/01-new-ambiguity/)
- Pending plans: [P]  (from .lovable/plans/pending/)
- Pending issues: [I]  (from .lovable/issues/ + .lovable/pending-issues/)
- Pending CI/CD issues: [C]  (from .lovable/cicd-issues/)
- Active suggestions: [S]  (from .lovable/suggestions.md)
- Unimplemented spec scope: [U]  (from spec/)

---

## Execution Waves

### Wave 1: Independent Foundations (DB Schemas, Wrappers, Core Models)

- Can run in parallel across 3 subagents (disjoint files).

#### 1.1 <task title>

- Source: <file path(s), one per line if cross-referenced>
- Type: plan | issue | cicd-issue | ambiguity | suggestion | spec-scope
- Status: pending | in-progress | blocked-by-ambiguity
- Steps: <N> steps
- Depends on: None
- One-line intent: <what "done" looks like>

### Wave 2: Business Logic, Services & Endpoints

- Requires Wave 1 foundations to complete.

#### 2.1 <task title>

- Source: <file path(s)>
- Depends on: Wave 1 Task #[X]
- Steps: <N> steps
- One-line intent: <what "done" looks like>

### Wave 3: UI Components, Views & Documentation

- Requires Wave 2 business services.

#### 3.1 <task title>

- Source: <file path(s)>
- Depends on: Wave 2 Task #[Y]
- Steps: <N> steps
- One-line intent: <what "done" looks like>

---

## Blocking Ambiguities (Ranked by Blast Radius)

- [HIGH/MED/LOW] <slug>: <question> (blocks task #<n>, #<n>)

## Nothing pending

(only if the inventory is truly empty; list the folders scanned to prove it)
```

## Banned actions (auto-reject)

- Executing anything this turn
- Calling `plan--create` or any plan-mode / approval tool
- Listing only a subset ("here are the top 5") when more exist
- Inventing step counts without opening the source file
- Treating one task as several because it appears in multiple sources
- Skipping any source folder from the list above without noting it as missing
- Silently ignoring open ambiguities
- Softening the required aggressive wording
- Adding SEO items to the inventory unless the user explicitly listed SEO work

## Checklist before replying (every box)

- [ ] Walked every source folder in the list; noted the missing ones
- [ ] Opened every pending file, not just filenames
- [ ] Cross-referenced duplicates so each real task appears once
- [ ] Estimated steps using the rubric, based on the actual file contents
- [ ] Flagged tasks >7 steps for subtask decomposition
- [ ] Ranked and listed every open ambiguity by blast radius
- [ ] Sequenced tasks into Execution Waves (Wave 1, Wave 2, Wave 3)
- [ ] Filled the summary counters with real numbers, not `[N]` placeholders
- [ ] No em dashes, no softened wording, no execution, no `plan--create`

If any box is unchecked, do not reply. Fix it first.

---

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
