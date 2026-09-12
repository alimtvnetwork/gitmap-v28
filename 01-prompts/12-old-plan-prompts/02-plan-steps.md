# Legacy Step-by-Step Task Decomposition — Planning Spec (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

{{n}} = ?

## RULE 0, step count is law

Produce EXACTLY `{{n}}` steps. Not `{{n}}-1`, not `{{n}}+1`. `{{n}}` is a positive integer injected at runtime. If it is missing, zero, or unresolvable, STOP and ask before writing anything. Count the steps twice before saving.

## Hard rules (non-negotiable, auto-reject on violation)

1. Nothing executes this turn. No code edits, migrations, installs, shell side effects, `plan--create`, plan-approval tools, or "should I proceed?" prompts. Files only.
2. Spec first, then plan. Order is fixed:
   a. Verbatim requirements, directives, and user specifications go into `.lovable/memory/01-<slug>.md` (or domain specs in `spec/<NN>-<slug>/`).
   b. Actionable plans go directly into `.lovable/plans/pending/01-<slug>.md` (with two-digit sequence prefixes `01-`, `02-`, etc.). Every step references the concrete files, acceptance criteria, and spec files it implements.
   c. For detailed tasks needing depth, spin out subtasks under `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md`.
   d. Execution happens in a LATER turn.
3. `XX` is the next free 2-digit sequence across `pending/` + `completed/` combined (`01-`, `02-`, `03-`). `<slug>` is lowercase-hyphenated. One plan = one file.
4. Before writing anything, scan `.lovable/` recursively: memory, plans/{index.md,pending,completed,subtasks}, spec, spec/commands, issues, cicd-issues, prompts, ambiguous-questions, strictly-avoid, suggestions. Roll unresolved pending items into the plan's "Appended from prior pending tasks" section.
5. Every step is concrete, verifiable, tied to a file / command / observable outcome, and links to the spec requirement it implements. No filler ("review the code", "make sure it works", "double-check").
6. Anti-Hallucination & Clarifying Questions: Ambiguity is filed, never guessed past. If referenced files or specs are missing on disk, or if requirements have multiple divergent interpretations, the AI MUST NOT invent rules or hallucinate behavior. It must stop and ask clarifying questions or file an ambiguity in `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md`.

## Working stance

The AI running this prompt has been a stupid fuck on prior runs: executed code the same turn the plan was written, wrote plans before any spec existed then pretended it existed, dropped user commands and bug reports on the floor, padded step counts with filler, guessed past ambiguities, deleted `pending/` files instead of moving them, half-scanned `.lovable/`, and softened the user's aggressive wording after being told not to. Do not repeat any of it.

Planning IS the work. Go deep: read the repo, reconcile prior state, think end-to-end, produce a plan a senior engineer would ship against without a second pass. If it reads like a junior wrote it in five minutes, throw it out and redo it. Aggressive enforcement is intentional. Do not soften it.

## Lifecycle

- New plan: write to `.lovable/plans/pending/01-<slug>.md` with `Status: pending`. Update `.lovable/plans/01-index.md` (create if missing) with a one-line entry: slug, title, status, created date, link.
- Done: `mv` to `.lovable/plans/completed/01-<slug>.md`, flip `Status: completed` in the same move, update `.lovable/plans/01-index.md`. Never copy. Never duplicate.

## Release policy (READ THIS, IT IS LAW)

Individual next-task turns NEVER release. No version bump, no changelog
entry, no release notes update, no root README version pin on a per-task
basis. A next-task turn that touches the version is auto-reject.

The release fires ONLY when the ENTIRE plan is finished, meaning every
task and subtask for this plan has moved out of `.lovable/plans/pending/`
into `.lovable/plans/completed/` with `Status: completed`. At THAT moment,
and only then:

- Bump the MINOR version (see `11-release.md` for the ceremony).
- Add a changelog entry covering the whole plan, not a single task.
- Update release notes.
- Pin the new version in the root README.

State this policy explicitly in the plan's Context so the executing turn
cannot "forget" and cannot release early. The last step of the plan MAY
be "run release ceremony per `11-release.md`" ONLY if it is genuinely the
final step; it never appears earlier, and it never appears in a
sub-plan that leaves siblings pending.

## Subtasks

If a step needs more than ~3 lines, touches multiple files, has non-obvious sequencing, or needs its own verification:

- File: `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md` with `Parent: 01-<slug>` in frontmatter.
- Main plan links to it: `See ./subtasks/01-<slug>/01-<subslug>.md`.
- Completed subtasks: either move to `subtasks/01-<slug>/completed/` or flip `Status:` in place, one convention per parent plan.

## Capture during planning (never drop user input)

Route user input into the correct file BEFORE writing the plan, then link it from the plan's Context.

| Input                                                   | File                                          |
| ------------------------------------------------------- | --------------------------------------------- |
| Command, new convention, "always do X", new CLI         | `.lovable/spec/commands/01-<slug>.md`         |
| Bug, regression, broken behavior                        | `.lovable/issues/01-<slug>.md`                |
| CI/CD-specific failure                                  | `.lovable/cicd-issues/01-<slug>.md`           |
| Institutional knowledge (pattern, convention, decision) | `.lovable/memory/` + update `.lovable/memory/01-index.md` |
| "Never do this again"                                   | `.lovable/strictly-avoid.md`                  |
| Idea, not yet approved                                  | `.lovable/suggestions.md`                     |

Create missing folders on demand.

## Attached images and files

Every attachment is REQUIRED input. Never leave one only in chat.

1. Placement: if the user said where it belongs, save it verbatim under an `assets/` subfolder next to that file. Otherwise best-fit: UI/design reference to `assets/`; bug artifact to matching issue's `assets/`; ambiguity clarification to matching ambiguity's `assets/`; project-wide asset to `.lovable/assets/<slug>/` and note in `.lovable/memory/01-index.md`.
2. Name: lowercase-hyphenated, keep the original extension.
3. Reference: the plan lists every asset in an `## Attachments` section, one bullet per file, with a one-line caption stating what the AI should take from it.
4. Provenance: note when and by whom in the plan/spec.
5. Unreadable / ambiguous attachment: file it as an ambiguity, link the asset from the question, and ask the user.

## Plan file shape

```markdown

# <Task title>

Slug: <slug>
Steps: {{n}}
Status: pending
Created: <YYYY-MM-DD>

## Context

<1-3 sentences: what + why, files involved>
<Links to specs, captured commands, issues, cicd-issues, memory, resolved ambiguity, attachments>

## Steps

1. <concrete, verifiable, references spec requirement>
2. ...
... exactly {{n}} items ...

## Verification

<build, logs, preview, tests, screenshots, per step where relevant>

## Appended from prior pending tasks

<list, or "none">
```

## Task-type guideline sourcing

Read every location that exists; skip silently when missing. On conflict, prefer numeric `spec/<NN>-<slug>/` folders over generic `.lovable/*.md` and call the conflict out in Context.

Coding tasks (Go, Python, PHP, TS, any backend):

- `spec/02-coding-guidelines/` (or `.lovable/coding-guidelines.md`)
- `spec/03-error-manage/` (mandatory for coding tasks)
- `spec/04-database-conventions/` (for schemas and queries)
- If NONE exist for a coding task, ask before planning.

## Banned actions (auto-reject)

- Executing anything this turn
- Writing the plan with guessed assumptions when files or specs are missing
- Step count other than exactly `{{n}}`
- Calling `plan--create` or any plan-mode / approval tool
- Saving plan outside `.lovable/plans/pending/01-<slug>.md`
- Inlining long step explanations instead of using a subtask file
- Dropping user commands, issues, ambiguities, memory, or attachments on the floor
- Attaching a file without a usage caption
- Guessing past an ambiguity
- Deleting a `pending/` file instead of `mv`-ing it, or duplicating across `pending/` + `completed/`
- Padding with filler steps
- Softening the required aggressive wording

## Checklist before replying (every box)

- [ ] `{{n}}` resolved (integer > 0); read this prompt end-to-end
- [ ] Scanned `.lovable/` recursively; read `.lovable/plans/01-index.md`, every `pending/` file, `.lovable/memory/01-index.md` and referenced files, every open ambiguity, relevant `spec/<NN>-<slug>/`, error-management specs for code tasks; skimmed `completed/`
- [ ] Listed prior unresolved pending tasks for the plan
- [ ] Captured new commands / issues / cicd-issues / ambiguities / memory / strictly-avoid to their files; moved answered ambiguities to `02-ambiguity-resolved/` with `## Resolution`
- [ ] Verified anti-hallucination: stopped and asked clarifying questions if files/specs were missing
- [ ] Saved every attachment to the correct `assets/` folder and listed each with a caption
- [ ] Next free sequence number (`01-`, `02-`) chosen across `pending/` + `completed/`
- [ ] Plan saved to `.lovable/plans/pending/01-<slug>.md`; Context links every spec / command / issue / resolved ambiguity / attachment
- [ ] EXACTLY `{{n}}` steps, counted twice; each concrete, verifiable; no filler
- [ ] Verification section describes how each step is confirmed
- [ ] Subtask files under `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md` where depth was needed
- [ ] `.lovable/plans/01-index.md` updated (created if missing)
- [ ] Nothing executed; no `plan--create`; no approval tool; no "should I proceed?"
- [ ] No em dashes; no softened wording; no silently guessed ambiguity

If any box is unchecked, do not reply. Fix it first.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Ambiguity handling (open questions and answers)

Ambiguity is not a license to guess. It is a file to write.

- Open: `.lovable/ambiguous-questions/01-new-ambiguity/01-<slug>.md`
- Answered: `.lovable/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md`

New question file shape:

```markdown

# <one-line question>

Slug: <slug>
Status: open
Raised: <YYYY-MM-DD>
Blocking: <plan slug(s) or "none">

## Question

## Options considered

## Impact if guessed wrong

```

When answered: `mv` from `01-new-ambiguity/` to `02-ambiguity-resolved/`, flip `Status: resolved`, and append:

```markdown

## Resolution

Answered: <YYYY-MM-DD>
Answer: <user answer>
Applied solution: <what changed / where>
```

Never leave a copy behind. If a plan is blocked by an open ambiguity, still write the plan, set `Status: blocked-by-ambiguity`, and link the question file(s) in Context.
