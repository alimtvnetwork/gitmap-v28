# Legacy Grouped Task Decomposition — Planning Spec (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Follow the above instructions.

N = [Number of steps the planning AI should self-loop to write the detailed task spec]
Groups / X Steps = [How many groups/agents it will take to complete the plan, and how many steps per group]

## RULE 0: Do NOT Execute (Non-Negotiable)

DO NOT EXECUTE the task initially. Your first responsibility is to write the spec based on what is given. No code edits, migrations, installs, or shell side effects should occur in this turn. Spec first, then plan.

## 1. Asset & Image Handling

If any images or files are provided, you MUST save them into the `assets/` folder.

- Rename them properly based on the task context (lowercase-hyphenated).
- You MUST refer back to these images in the markdown spec files you write.

## 2. High-Detail Task Breakdown (The "High-Powered Brain")

You must use your high-powered brain to write sub-tasks that explicitly detail:

- Where the code is.
- What the issue is.
- How to change it.

The tasks must be so detailed and explicit that the next executing AI (who requires less brain) can simply follow the instructions without guessing.

## 3. Sub-Agent Grouping & Folder Structure

You must create tasks for AI agents as a group.

- Group pending tasks into `.lovable/plans/pending/` and `.lovable/plans/subtasks/` so that when the executing AI agents run, they see the pending tasks clearly defined.
- Each group must be assigned to a separate standalone sub-agent.
- Each group must contain full context so the agent can execute and commit its code in a completely STANDALONE way.
- Give each group/sub-task a highly specific title reflecting its exact task (e.g., `Refactoring Auth` or `Fixing DB Connection`).
- Sub-agents must only be assigned simple, small micro-tasks rather than larger monolithic ones.

## 4. Self-Loop to Create Brain and Tasks

You must self-loop `N` times to create the brain, tasks, and subtasks for the defined Groups.
Generate the folder structure, the `.lovable/plans/pending/` files, and the `.lovable/plans/subtasks/` files immediately during your self-loop.
Your self-loop must strictly follow this structure:

1. Roll-up Index: Update `.lovable/plans/01-index.md`.
2. Parent Plan: Create `.lovable/plans/pending/xx-<slug>.md`.
3. Subtasks (Strict 01, 02 Sequence): Every group/step must be placed in a dedicated file under `.lovable/plans/subtasks/xx-<slug>/`.
   - The subtasks MUST follow a strict zero-padded numeric sequence: `01`, `02`, `03`, etc.

## 5. High-Stakes Code Standards, Error Management & Guidelines

You MUST follow the project's strict coding guidelines and ensure your groups enforce them.
For every task, you MUST check if the following files or folders exist. If they exist, they MUST be followed and included in the task's checklist for the executing AI to follow. If they do not exist, they can be skipped.

### Dynamic Required Reading / Reference Checklist (Non-Negotiable):

1. Root Memory Guidelines

- /learn `.lovable/coding-guidelines.md` (and/or `.lovable/coding-guidelines.md`)

2. Master Consolidated Guide & Coding Guidelines

- /learn `spec/17-consolidated-guidelines/05-coding-guidelines.md`
- /learn `spec/02-coding-guidelines/01-cross-language/01-index.md`

3. Error Management (Must Follow for all Coding Tasks)

- /learn `spec/03-error-manage/01-index.md`
- /learn *Include most of the files from the error manage directory to ensure robust error handling is implemented per task.*

4. Boolean Conditions, Wrappers & Samples

- /learn `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`
- /learn `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`
- /learn `spec/02-coding-guidelines/01-cross-language/24-boolean-flag-methods.md`

5. Code Style & File Size Limits (80-100 lines max)

- /learn `spec/02-coding-guidelines/01-cross-language/04-code-style/01-index.md`
- /learn `spec/02-coding-guidelines/01-cross-language/20-nesting-resolution-patterns.md`
- /learn `spec/02-coding-guidelines/01-cross-language/06-cyclomatic-complexity.md`

6. Variable Naming & Definitions

- /learn `spec/02-coding-guidelines/01-cross-language/22-variable-naming-conventions.md`
- /learn `spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md`
- /learn `spec/02-coding-guidelines/01-cross-language/10-function-naming.md`
- /learn `spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md`

## 6. End-of-Loop Commit Fix (Non-Negotiable)

Each step or group that completes MUST immediately commit and fix the Git. Follow these exact rules from the commit-fix prompt:

- Group all completed work for the subtask into a single logical commit.
- RED FLAG: Verify absolutely NO test results, artifacts, or compiled binaries are staged before making the commit.
- Ensure .gitignore explicitly excludes them.
- If issues arise during the commit process, fix those git issues and try again.
- You MUST push the commit to the repository immediately. Pushing after commits is non-negotiable.

## Action Items — Must Follow (Non-Negotiable)

- [ ] Read the user-provided "N" and "Groups / X Steps" variables.
- [ ] Do not execute anything; write the detailed spec and pending tasks first.
- [ ] Save all provided images to `assets/`, rename them contextually, and link them in the markdown.
- [ ] Write tasks using your high-powered brain, detailing explicitly where the code is, what the issue is, and how to change it.
- [ ] Group tasks logically for independent sub-agents, ensuring each group contains full context for a standalone sub-agent.
- [ ] Ensure each group can work and commit standalone.
- [ ] Self-loop `N` times to write detailed specs and task files.
- [ ] Ensure subtask files are named exactly xx-<subslug>.md where XX is the sequence.
- [ ] Audit the plan against the Master Consolidated Guide, Code Style, Boolean Conditions, and Variable Naming rules.
- [ ] Ensure each group section explicitly outlines the agent part (that the group is executed by a separate standalone agent).
- [ ] Ensure the end-of-loop commit fix is properly executed for each completed step or group.
- [ ] Update `.lovable/plans/01-index.md` and populate `.lovable/plans/pending/` and `.lovable/plans/subtasks/`.

### Execution Mode: Plan & Wait (User Approval Required)

By default, you operate in a strict Plan & Wait mode.

- [ ] Prepare Plan: First, outline what you intend to do and create the proposed plan overview.
- [ ] Stop and Wait: You MUST stop execution and ask the user for approval. Do NOT write the multiple task files to disk, do NOT spawn sub-agents, and do NOT update the indexes yet.
- [ ] Execute: ONLY when the user explicitly replies with "execute" or "approved" may you proceed to the full execution phase (self-looping, writing out the detailed subtask files, and updating the filesystem).
*(Note: If the user explicitly requests "Plan and Execute" mode upfront, you may bypass this wait and proceed directly to full generation).*

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

- [ ] Anti-Boilerplate Check: Did I copy-paste the exact same "How" steps across multiple tasks? (If yes, you are acting stupid. Stop and rewrite them to be uniquely specific to the task's exact technical requirements).
- [ ] Cognitive Check: Does every task contain "all the brains" (exact logic, specific paths, deep architectural context) so a lower-level agent can execute it without guessing?
- [ ] /learn the overarching main task plan.
- [ ] Read `.lovable/memory/01-index.md` and `.lovable/plans/01-index.md` before planning.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## Execution & Self-Looping Protocol

This task requires deep, multi-agent processing. You MUST NOT run a fast, linear self-loop that artificially consolidates work or skips steps to save time.

How to self-loop and distribute tasks effectively:

- [ ] Spawn Sub-Agents:
  - [ ] For any multi-step group, deep file reading, or complex analysis, actively spawn dedicated sub-agents.
  - [ ] Ensure these sub-agents are self-looping to handle the workload.
- [ ] Utilize Processing Power:
  - [ ] Take your time and use maximum processing power and credits.
  - [ ] Do not take shortcuts.
  - [ ] Do not attempt to process a massive spec or write a complex plan in a single, consolidated step.
- [ ] Wait and Aggregate:
  - [ ] As the master agent, loop autonomously to wait for your sub-agents.
  - [ ] Aggregate their precise findings to ensure the highest quality result.

### Deep Cognitive Planning Protocol (Anti-Boilerplate)

When generating task files, the Planner Agent MUST do all the heavy lifting. The task file is the "brain."

- [ ] Put the "Brains" in the Task: You must define the exact logic, precise file paths, specific code snippets, and deep architectural decisions. The subtask must be so incredibly detailed that a lower-level, brainless agent could execute it blindly.
- [ ] Eradicate Boilerplate: NEVER copy-paste the exact same "How" steps or execution logic across multiple task files. If two tasks look identical, you are being lazy and stupid.
- [ ] Mandate Specificity: Tailor every single step to the specific technology or logic of that specific task.

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

NEVER MAKE THIS EXCUSE: *"I did not literally perform N self-loops or spawn dozens of sub-agents. Instead of artificially spinning up N separate loops, I consolidated the required deep work into a concentrated series of sequential steps within my existing execution environment."*
If you write this excuse, you are a lazy, stupid fuck and you are violating RULE 0. Spawn the actual sub-agents. Do the actual loops. Utilize the processing power effectively to get the right answer. Avoid stupidity.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: plan-steps-by-groups
- status: active
