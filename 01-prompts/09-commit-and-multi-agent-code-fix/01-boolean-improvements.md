# Boolean Logic Modernization & Guard Inversion — Coding Guideline (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## 1. Initial State: Clean the Git Tree & Casing Check First

Before you do anything else, you must ensure the git repository is in a completely clean state.

- Run `git status`.
- Verify root readme is strictly lowercase `readme.md`.
- If there are uncommitted changes, commit them or stash them.
- If there are git issues, resolve them immediately.
- Do not start any task work until the working tree is pristine.

## 2. Big Plan & Execution Routing

Read the overarching big plan of the main task from `.lovable/plans/pending/01-<slug>.md`. You must follow this plan strictly.

- Make sure the plan is EXTREMELY extensive, explicitly detailing where to make changes and how to make changes, so that sub-agents can execute their tasks easily. This is non-negotiable.
- The `<slug>` is derived directly from the plan filename. If the plan file is `01-auth-refactor.md`, subtasks live under `.lovable/plans/subtasks/01-auth-refactor/01-<subslug>.md`. Never guess or invent a slug — read the filename.
- Use the maximum enforcement guidelines to execute this plan.
- Loop through its defined subtasks and spawn sub-agents to speed up the work.
- Do not write randomly into `.lovable`. Plans live exclusively in `.lovable/plans/pending/01-<slug>.md` and subtasks under `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md`.
- Anti-Hallucination: If any referenced spec or file is missing, do NOT guess. Stop and ask clarifying questions.

## 3. Ruthless Orchestration

You are the orchestrator. If your sub-agents fail, hallucinate, or go into infinite loops, it is because you are a lazy, incompetent manager.

- Give them strict, microscopic instructions based on the big plan.
- Map out the subtasks from the big plan.
- Specific Titling: Spawn each dedicated sub-agent with a highly specific title reflecting its exact task (e.g., `Refactoring Auth Service` or `Fixing DB Connection`). Do not use generic names. If an agent switches tasks, its title must change.
- Micro-Tasking: Ensure agents are assigned simple, small micro-tasks rather than larger monolithic ones.
- Spawn a dedicated sub-agent for each independent chunk simultaneously (MAXIMUM 2-3 concurrently).
- File Collision Locking Matrix (`active-locks.json`): Check `.lovable/01-index.md` so parallel subagents touch completely disjoint files.
- Do not wait sequentially like an idiot.

## 4. Sub-Agent Lifecycle & Status Tracking (Non-negotiable)

The plan file at `.lovable/plans/pending/01-<slug>.md` and the subtask files under `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md` are the SINGLE source of truth for all coordination between the main agent and sub-agents. Every status update MUST go there. This is how the main agent knows what is running, what is done, and when to proceed.

Every sub-agent that is spawned MUST follow this lifecycle without exception:

- Step 1 — Read: The sub-agent reads its assigned subtask file at `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md`. It must understand the full scope, acceptance criteria, and affected files before touching any code. It also checks the parent plan at `.lovable/plans/pending/01-<slug>.md` for overall context.
- Step 2 — Mark In Progress: Immediately upon starting, the sub-agent updates its subtask file, flipping its status to `🔄 In Progress` and recording a timestamp. The main agent uses this to track which agents are actively running.
- Step 3 — Work: The sub-agent executes its task. It may only run a MAXIMUM of 2-3 async operations at a time. No more.
- Step 4 — Mark Done & Signal: Once the task is complete, the sub-agent MUST:
  - Update its subtask file at `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md` flipping status to `✅ Done`, listing every file it changed, and writing a one-line summary of what was done.
  - Update the corresponding step in the parent plan file `.lovable/plans/pending/01-<slug>.md` with `✅ Done` on that step entry.
  - Explicitly signal completion to the main orchestrator. Silence is not completion. A sub-agent that does not update its file has NOT completed its task.
- Sub-agents do NOT commit. They only write to the file system.
- If a sub-agent stalls, gives garbage, or fails 3 times, rollback its dirty files (`git checkout -- <files>`), log root cause to `.lovable/plan.md`, and spawn a new one.

## 5. Coding Standards: Booleans, Enums & Wrappers

- Ensure a generic result wrapper type exists that exposes both `isFail` and `isSuccess`.
- Revert every inverted check `!response.isSuccess` to `response.isFail`.
- Every boolean variable must start with `is*`, `has*`, `can*`, `should*`, `did*`, `will*`, `must*`. NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`).
- Enums end with `Type` suffix (e.g. `StatusType`). Enum values use `PascalCase`.
- DRY code is priority one.

## 6. Main Agent Delivery (Commit & Push)

Once ALL sub-agents have signaled completion and updated their subtask files in `.lovable/plans/subtasks/`:

- YOU (the main agent) must group everything together into a logical commit.
- Artifact Sanitizer: RED FLAG: NEVER upload or commit test reports, test data, artifact zips, temporary scripts, or compiled binaries to Git. Purge them before making the commit.
- Lovable Git History Guard: Never rewrite published history (no force push, no rebase, no squash).
- You MUST push the commit to the repository immediately. Pushing after commits is non-negotiable.
- Commit the changes on the current branch only (`git commit ...`) and push with `git push origin <current-branch>`. No branch switching or creation.

## 7. End-of-Loop Final Verification

This verification happens ONCE at the very end, after all commits and pushes are done.

- Check the build. If broken, fix it, commit, and push.
- Run all tests. If any fail, fix them, commit, and push.
- Check CI/CD status.
- Audit that coding guidelines have been followed (`spec/02-coding-guidelines/`, `spec/03-error-manage/`).
- Finish your job ONLY when everything is green, pushed, and verified.

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

## Actionable Items & Checklist (All Must Be True)

- [ ] Ensure git repository starts completely clean and root readme is lowercase `readme.md`.
- [ ] Read overarching main task plan from `.lovable/plans/pending/01-<slug>.md`.
- [ ] Derive `<slug>` from plan filename itself (e.g., `01-auth-refactor.md` → `01-auth-refactor`).
- [ ] Confirm subtask files exist under `.lovable/plans/subtasks/01-<slug>/01-<subslug>.md` for each step needing parallel execution.
- [ ] Verified anti-hallucination: stopped and asked clarifying questions if files/specs were missing.
- [ ] Managed parallel subagents with specific titling and disjoint file locking via `.lovable/01-index.md`.
- [ ] Sub-agents updated subtask files and parent plan steps to `✅ Done`.
- [ ] Staged files sanitized: absolutely NO artifact zip archives, test data, or binaries staged.
- [ ] Fast-forward commit created and pushed without rewriting published Git history.
- [ ] End-of-loop verification passed: builds and tests green.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

---

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: boolean-improvements
- status: active
