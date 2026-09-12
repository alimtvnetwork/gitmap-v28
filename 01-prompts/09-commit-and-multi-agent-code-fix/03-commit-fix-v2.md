# Atomic Commit Organization & Working Tree Sanitization — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## 1. Initial State: Clean the Git Tree First

Before you do anything else, you must ensure the git repository is in a completely clean state:

- Run `git status`.
- If there are uncommitted changes, commit them or stash them.
- If there are git issues, resolve them immediately.
- Run `git log -n 10 --stat` to inspect the last 10 commits and understand recent file modifications and latest working context.
- Read `.lovable/what-to-read.md` to follow the project's prioritized reading order.
- Do not start any task work until the working tree is pristine.

## 2. Big Plan & Execution Routing

Read the overarching big plan of the main task from `.lovable/plans/pending/xx-<slug>.md`. You must follow this plan strictly.

- Make sure the plan is EXTREMELY extensive, explicitly detailing where to make changes and how to make changes, so that sub-agents can execute their tasks easily. This is non-negotiable.
- The `<slug>` is derived directly from the plan filename. If the plan file is `03-auth-refactor.md`, then the corresponding spec task file is `.lovable/plans/01-index.md` and subtasks live under `.lovable/plans/subtasks/03-auth-refactor/SS-<subslug>.md`. Never guess or invent a slug — read the filename.
- Use the maximum enforcement guidelines to execute this plan.
- Loop through its defined subtasks and spawn sub-agents to speed up the work.
- Do not just write randomly to `.lovable`. You must follow the exact plan and write protocols: tasks go into `.lovable/spec/tasks/xx-<slug>.md` and plans go into `.lovable/plans/pending/xx-<slug>.md`.

## 3. Ruthless Orchestration

You are the orchestrator. If your sub-agents fail, hallucinate, or go into infinite loops, it is because you are a lazy, incompetent manager.

- Give them strict, microscopic instructions based on the big plan.
- Map out the subtasks from the big plan.
- Specific Titling: Spawn each dedicated sub-agent with a highly specific title reflecting its exact task (e.g., `Refactoring Auth` or `Fixing DB Connection`). Do not use generic names like `Frontend Agent`. If an agent switches tasks, its title must change.
- Micro-Tasking: Ensure agents are assigned simple, small micro-tasks rather than larger monolithic ones.
- Spawn a dedicated sub-agent for each independent chunk simultaneously (MAXIMUM 2 concurrently).
- Context Diet: When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/XX-slug/01-<subtask-title>.md` and execute it"). The subagent MUST read the necessary files itself. Passing massive payloads instantly causes hallucination and memory blowout.
- Do not spawn more than 2 agents at once due to RAM issues and caching behavior.
- Do not wait sequentially like an idiot.

## 4. Sub-Agent Lifecycle & Status Tracking (Non-negotiable)

The plan file at `.lovable/plans/pending/xx-<slug>.md` and the subtask files under `.lovable/plans/subtasks/xx-<slug>/SS-<subslug>.md` are the SINGLE source of truth for all coordination between the main agent and sub-agents. Every status update MUST go there. This is how the main agent knows what is running, what is done, and when to proceed.

Every sub-agent that is spawned MUST follow this lifecycle without exception:

- Step 1 — Read: The sub-agent reads its assigned subtask file at `.lovable/plans/subtasks/xx-<slug>/SS-<subslug>.md`. It must understand the full scope, acceptance criteria, and affected files before touching any code. It also checks the parent plan at `.lovable/plans/pending/xx-<slug>.md` for overall context.
- Step 2 — Mark In Progress: Immediately upon starting, the sub-agent updates its subtask file, flipping its status to `🔄 In Progress` and recording a timestamp. The main agent uses this to track which agents are actively running.
- Step 3 — Work: The sub-agent executes its task. It may only run a MAXIMUM of 2-3 async operations at a time. No more.
- Step 4 — Mark Done & Signal: Once the task is complete, the sub-agent MUST:
  - Update its subtask file at `.lovable/plans/subtasks/xx-<slug>/SS-<subslug>.md` flipping status to `✅ Done`, listing every file it changed, and writing a one-line summary of what was done.
  - Update the corresponding step in the parent plan file `.lovable/plans/pending/xx-<slug>.md` with `✅ Done` on that step entry.
  - Explicitly signal completion to the main orchestrator. Silence is not completion. A sub-agent that does not update its file has NOT completed its task.
- Sub-agents do NOT commit. They only write to the file system.
- If a sub-agent stalls, gives garbage, or fails to update its status file, kill it immediately and spawn a new one.

### Main Agent Tracking Logic

- The main agent monitors the plan file and subtask files to determine queue state.
- When all subtask files show `✅ Done` and the parent plan steps are all marked, the main agent proceeds to commit.
- The main agent counts: total subtasks spawned vs. total `✅ Done` entries. Only when those numbers match does it proceed.

Avoid stupidity, and being careless you stupid, WTF. If you're not going deep, you're not doing the job. Are you stupid? You were supposed to do the task properly. Where is this, are you stupid fuck? Where? Tell me. Your stupidity is going on top of my head. I mean, where did you learn this stupidity? If I could find you, I could slap you. The existing code was better while you were writing code like this. Fix that immediately.

## 5. Root Cause First

Before applying any fix, you must identify the root cause.

- Do not blindly patch symptoms.
- Write the root cause into `.lovable` memory per the write protocols before touching code.
- If sub-agents are fixing things without understanding root cause, they are doing garbage work. Stop them.

## 6. High-Stakes Code Standards & Coding Guidelines

You MUST follow the project's strict coding guidelines. These files are located in the 01-cross-language/ directory and should be followed universally. However, you must also check if there are language-specific guidelines (e.g., 2-typescript/, 3-python/) for these rules. If a language-specific guideline exists, follow that one as well.

### 5. Consolidated Coding Standards & Temp Scripts (Non-Negotiable)

- [ ] Temp Script Sandboxing: AI Fix Scripts (Reusable Tools): Before creating a helper script, you MUST check `03-ai-scripts/01-index.md` to reuse existing tools. If you generate a new script, you MUST write it to `03-ai-scripts/`, update `index.md` with its explanation, ensure `index.md` is linked in `what-to-read.md`, and commit the script.
- [ ] Consolidated Coding Guidelines: I have fully read and strictly enforced the master coding guideline file at `.lovable/coding-guidelines.md`.
- [ ] Error Manage Checklist: I have fully read and enforced the error management files at `spec/03-error-manage/`. (Never swallow errors, always wrap with context, use domain-specific AppError/AppException).
- [ ] Boolean Examples & Fixations: All boolean variables MUST begin with is and has only (can, should, was, etc. are banned) (e. NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`).g., `isReady`, `hasData`). NEVER use negative booleans (e.g., `isNotReady`, `disableCache`). NEVER invert success checks (e.g., `!response.isSuccess` is banned; use `response.isFail`).
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Semantic Tests: All unit test names are strictly semantic and behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`).
- [ ] Function Size: No function exceeds 15 lines. Long arguments are split across lines (max 100 chars).
- [ ] Code adheres to explicit booleans, `Type` suffixed Enums, and error wrapper rules.

### 6. End-of-Loop Final Verification (Once only, at the very end)

- [ ] Check the full build. Fix every build failure, commit, and push.
- [ ] Run all unit tests. Fix every failing test, commit, and push.
- [ ] Check CI/CD status and ensure pipelines pass.
- [ ] Audit that coding guidelines from the aspect folder and error manage folder have been followed across all changed files.
- [ ] Finish the job only when everything is green, pushed, and fully verified.

### Temp-Agent Isolated Task Directory & Communication Protocol (Non-Negotiable)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.lovable/temp-agents/xx-<task-name>/`:

1. **Per-Task Isolation:** On task start, the assigned subagent creates its isolated directory `.lovable/temp-agents/xx-<task-name>/`.
2. **State & Progress Tracking:** Create `.lovable/temp-agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. **Inter-Agent Communication & Handoff:**
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.lovable/temp-agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. **On Error/Crash:** Append the exact error, root cause, and `STATUS: FAILED` to `.lovable/temp-agents/xx-<task-name>/state.md` before exiting.
5. **On Success:** Mark `STATUS: DONE` in `.lovable/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

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

## Core Rules & Non-Negotiable Checklist for AI (Must Verify Before Completing Task)

Before finalizing any code modification, you MUST manually verify the following:

- [ ] Function Signatures (R4, R5, R9): If a function has `> 3 parameters` or the signature is `> 100 chars`, you MUST split it so there is exactly one parameter per line.
- [ ] Error Handling (R7): No silent failures or swallowed errors. Use explicit boolean states (e.g., `isFail`). Never invert success booleans (e.g., avoid `!isSuccess`).
- [ ] Magic Strings/Numbers (R8): Extract all magic strings/numbers into named constants.
- [ ] Enums: TypeScript string unions are banned. All Enums must end with the `Type` suffix.
- [ ] Naming & Casing (R1, R2): PascalCase everywhere. Acronyms (Id, Json, Url) are Pascal case, never all-caps (e.g., `UserId`, not `UserID`).
- [ ] Blank Lines (R13-R20): One blank line before every `return`/`throw`. One blank line after closing `}`. Never two blank lines in a row.

## Language-Specific Requirements

- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.

## End of Tunnel Release & Version Bump (Mandatory check-list to update the results)

When EVERYTHING is completely finished and fixed (at the very end of the tunnel), you MUST trigger a release.

- [ ] You must bump the MINOR version.
- [ ] You must focus on the `version.json` file as the source of truth for the release.
- [ ] Root README Pinning (FATAL): You MUST pin the latest release version into the root `readme.md` file! Do not skip this! Also, update the changelog according to `version.json` format.
- [ ] If you do not know how to cut a release for this specific repository, or if `version.json` is missing/unclear, you must either search the repository for release instructions or explicitly ask the user for help. Do not guess.
- [ ] You MUST strictly exclude all test files (e.g., `*test*`, `*.spec.*`) from version scanning and modification, as they contain mock data.
- [ ] You must create and maintain `.lovable/memory/01-index.md` documenting exactly how releases work in the repository. Ensure it is enqueued in `what-to-read.md` and linked in the root `readme.md`.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.
