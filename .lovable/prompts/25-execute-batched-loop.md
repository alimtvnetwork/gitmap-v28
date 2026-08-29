# Execute Batched Loop (3 Agents, Chunked Commits)

- slug: execute-batched-loop
- status: active
- trigger: Run if again if said, go, continue or next

## Prompt

# Execute Batched Loop (3 Agents, Chunked Commits)

## Goal

Your objective is to execute pending tasks from `.lovable/plans/pending/` using a strictly batched multi-agent loop. You must use exactly 3 sub-agents, assign exactly 3 tasks per agent, handle crashes via the `.lovable/temp/` folder, and push chunked commits to Git without failure. At the end of every loop, you must explicitly list task statistics in your output window.

## Phase 1: Pre-Flight & Gitignore Enforcement (Non-Negotiable)

1. **Check Git Status:** The working tree must be clean. 
2. **Enforce Gitignore:** You MUST verify that `.lovable/temp/` is explicitly added to `.gitignore`. If it is not, update `.gitignore` immediately. This folder is for crash identification and must never be committed. This is non-negotiable.
3. **Garbage Collection:** Wipe any orphaned state files in `.lovable/temp/` from previous runs.

## Phase 2: Allocation & Execution (Strict 3x3 Rule)

1. **Strict Limits:** You must spawn exactly 3 sub-agents to run in parallel. Never more, never less (unless fewer than 3 task chunks remain).
2. **Chunking Micro-Tasks:** Each agent must be assigned exactly a chunk of 3 *simple, small* micro-tasks to complete sequentially in their own context. Do not assign monolithic tasks. 
3. **Disjoint Assignment:** Ensure that the tasks assigned to the 3 parallel agents touch different files or components to prevent Git merge conflicts.
4. **Temp Folder Logging & Specific Titling (Mandatory):** Before an agent starts, you must:
   - Spawn the sub-agent with a highly specific title reflecting its exact task (e.g., `Refactoring Auth` or `Fixing DB Connection`). Do not use generic names like `Frontend Agent`. If an agent switches chunks, its title must change.
   - Log its assigned chunk of tasks to `.lovable/temp/XX-agent-state.md`.
5. **Crash Identification & Recovery:** 
   - Every crash needs to be identified using the `.lovable/temp/` folder. 
   - If an agent fails or crashes, read its state from the temp folder to know exactly what chunk it was doing and how it crashed. 
   - Reason about the crash, log the failure to `.lovable/memory/last-failure.md` if it cannot be recovered, and restart a new agent to pick up the dropped chunk.

## Phase 3: Code Quality (Non-Negotiable)

While executing tasks, you and your agents MUST adhere to these strict coding guidelines:
- Follow guidelines in `spec/02-coding-guidelines/` and `spec/03-error-manage/`.
- **No Magic Strings/Numbers:** Do not introduce any magic strings or numbers anywhere unless it is explicitly for the logger (and state that in the typing).
- **TypeScript Enums:** Never use string union types (e.g., `"pass" | "fail"`). You must use Enums. Every single Enum must end with the suffix `Type` (e.g., `StatusType`).
- **Explicit Booleans:** Always use explicit boolean state checks (e.g., `response.isFail`). Never invert success booleans (e.g., `!response.isSuccess`).
- **DRY Code:** Code must be DRY. Reuse constants and wrappers.

## Phase 4: Chunked Delivery & File Moving

When a chunk of tasks is completed by the agents, you MUST do the following before starting the next loop iteration:
1. **Move Files:** `mv` the completed task files from `.lovable/plans/pending/` to `.lovable/plans/completed/`.
2. **Status Flip:** Open the moved files and change `Status: pending` to `Status: completed`.
3. **Update Indexes:** Update `.lovable/plans/index.md` to reflect the new file locations.
4. **Commit & Push (Non-Negotiable):** When a chunk of tasks is done, group the changes, run tests (no live API calls), commit the code with a clear descriptive message, and **push to git without a failure**. This must happen at the end of every chunk iteration.

## Phase 5: Output Window Stats (Mandatory)

Every time you return a response to the user or complete a loop iteration, you must explicitly output the following statistics in your output window. This is very, very important:
- **Tasks Done (This Chunk):** [Number of tasks completed]
- **Total Completed:** [Total number of tasks in `.lovable/plans/completed/`]
- **Total Pending:** [Number of tasks remaining in `.lovable/plans/pending/`]
- **Remaining Tasks List:** [List the specific filenames/slugs of the tasks remaining]

## Pre-Reply / Loop Checklist (Must verify every loop iteration)

- [ ] `.gitignore` was updated to exclude `.lovable/temp/`.
- [ ] Strictly 3 agents were spawned, each assigned a chunk of exactly 3 tasks.
- [ ] Pre-flight state written to `.lovable/temp/` for every agent to identify crashes.
- [ ] Crashes were handled via temp folder state recovery.
- [ ] No end-to-end API tests were executed.
- [ ] Completed task files were `mv`'d and `index.md` was updated.
- [ ] The completed chunk was committed and pushed to git successfully without failure.
- [ ] The output window explicitly lists "Done", "Pending", and remaining task names.
