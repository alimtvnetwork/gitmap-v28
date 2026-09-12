# Inventory Audit of Pending Tasks — Read-Only Proposal (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Perform a strictly read-only scan of the entire repository, `spec/`, and `.lovable/` directory to compile a comprehensive, deduplicated inventory of every pending task, subtask, unresolved issue, and open requirement structured into Execution Waves.

CRITICAL CONSTRAINT: This prompt is strictly for inventorying, structuring, and sequencing pending work. It MUST NOT execute code modifications, build changes, or launch the execution loop. Batch execution is handled by dedicated execution prompts.

---

## Folder Structure Reference (What and Where to Read)

To ensure zero blind spots, the AI must systematically inspect the following authoritative folder structure:

```text
.lovable/
  memory/
    index.md                                    # Master memory index
    what-to-read.md                             # Authoritative reading order
    specs/                                      # Verbatim user directives & requirements
    learned/                                    # Institutional knowledge & learnings
    workflow/                                   # Workflow state markers
  plans/
    index.md                                    # Roll-up index of all plans (pending & completed)
    pending/xx-<slug>.md                        # Active, pending parent execution plans
    completed/xx-<slug>.md                      # Historical completed plans (reference only)
    subtasks/xx-<slug>/SS-<subslug>.md          # Granular subtasks linked to parent plans
  issues/                                       # General pending bug reports
  pending-issues/                               # Active issues queue
  cicd-issues/                                  # CI/CD specific failures & blockers
  ambiguous-questions/
    01-new-ambiguity/xx-<slug>.md               # Blocking & non-blocking open questions
    02-ambiguity-resolved/xx-<slug>.md          # Binding resolved decisions
  suggestions.md                                # Suggestions tracker
  prompts.md                                    # Canonical prompt registry
spec/
  01-spec-authoring-guide/                      # Spec authoring conventions & required files
  02-coding-guidelines/                         # Coding standards & zero-tolerance rules
  03-error-manage/                              # Error management & logging philosophy
  04-database-conventions/                      # Database schema & query rules
  21-app/                                       # App specification & feature requirements
  (note: numbers in spec/<NN>-<slug>/ can switch between projects; dynamically traverse all nested .md files)
readme.md                                       # Root repository guide (strictly lowercase)
```

---

## Step 1: Deep Inspection & Deduplication Protocol

1. Scan Every Pending Source:
   - Open and read `.lovable/plans/01-index.md`, `.lovable/plans/pending/`, and all `.lovable/plans/subtasks/` files with `Status:` not `completed`.
   - Open and read all files in `.lovable/issues/`, `.lovable/pending-issues/`, and `.lovable/cicd-issues/`.
   - Open and read all open questions in `.lovable/ambiguous-questions/01-new-ambiguity/`.
   - Open and read unfulfilled directives in `.lovable/memory/` and `spec/21-app/`.
   - Open and read active suggestions in `.lovable/suggestions.md`.
2. Deduplicate Across Sources:
   - If a feature is referenced across a spec, a plan, and an issue, consolidate it into ONE primary task with cross-references to all origin files.
3. Step-Count Rubric & Decomposition Alert:
   - Trivial change (1 file, single edit): 1 step.
   - Small change (1-2 files, 1 verification step): 2-3 steps.
   - Standard task (multi-file, logic + UI/backend, test): 4-7 steps.
   - Cross-cutting task (schema + API + UI + full tests): 8-15 steps.
   - Automatic Subtask Decomposition Alert: When a pending task exceeds 7 steps, flag it with `[DECOMPOSITION REQUIRED]` to split it into `.lovable/plans/subtasks/xx-<slug>/` before entering the execution queue.
4. Ambiguity Impact Severity Scoring:
   - High Blast Radius: Blocks multiple core plans or schemas.
   - Medium Blast Radius: Blocks a single isolated feature.
   - Low Blast Radius: Non-blocking cosmetic or detail refinement.

---

## Step 2: Refined Output Format (Structured by Execution Waves)

Present the inventory to the user in this exact markdown structure:

```markdown

# Pending Tasks Inventory

## Executive Summary

- Total Pending Tasks: [N]
- Blocking Ambiguities: [K]
- Pending Plans: [P]
- Pending Issues & CI/CD: [I]
- Unimplemented Spec Scope: [U]
- Sources Scanned: [.lovable/plans/pending/, .lovable/plans/subtasks/, .lovable/issues/, .lovable/cicd-issues/, .lovable/ambiguous-questions/01-new-ambiguity/, spec/]

---

## Execution Waves

### Wave 1: Independent Foundations (DB Schemas, Wrappers, Core Models)

- Can run immediately in parallel across 3 subagents (disjoint files).

#### 1.1 [Task Title]

- Source File(s): `[path/to/file.md]`
- Type: `Plan` | `Issue` | `CI/CD` | `Spec-Scope`
- Status: `Pending` | `In-Progress`
- Estimated Steps: `[N] steps`
- Dependencies: `None`
- Outcome / Intent: [Clear 1-2 sentence description of what "done" looks like]

### Wave 2: Business Logic, Services & Endpoints

- Requires Wave 1 foundations to complete before execution.

#### 2.1 [Task Title]

- Source File(s): `[path/to/file.md]`
- Dependencies: `Wave 1 Task #[X]`
- Estimated Steps: `[N] steps`
- Outcome / Intent: [Description]

### Wave 3: UI Components, Views & Documentation

- Requires Wave 2 business services.

#### 3.1 [Task Title]

- Source File(s): `[path/to/file.md]`
- Dependencies: `Wave 2 Task #[Y]`
- Estimated Steps: `[N] steps`
- Outcome / Intent: [Description]

---

## Blocking Ambiguities (Ranked by Blast Radius)

- [HIGH/MED/LOW] `[slug]`: [Question summary] — *Blocks: Task #[X], Task #[Y]*
```

---

## Step 3: Follow-Up Execution Proposal (Strictly Positive Question Framing)

At the very end of the inventory report, the AI must ask the user whether to trigger execution using the standard 3-agent parallel loop:

```text
Would you like to start the continuous self-loop to execute and resolve these pending tasks in parallel using a maximum of 3 concurrent subagents? (Yes / No)
```

*(If the user answers Yes, the next turn will invoke `execute-robust-loop` or `execute-batched-loop` to begin execution).*

---

## Checklist: What This Instruction MUST Do and MUST NOT Do

### What to Do (Mandatory):

- [ ] Read all `.lovable/` pending folders, subtask files, issue trackers, and spec requirements in full.
- [ ] Deduplicate tasks appearing across multiple plan, spec, or issue files.
- [ ] Calculate concrete step counts based on actual file contents using the rubric.
- [ ] Flag tasks exceeding 7 steps for subtask decomposition.
- [ ] Rank open ambiguities by severity and blast radius.
- [ ] Sequence pending tasks into logical Execution Waves (Wave 1, Wave 2, Wave 3).
- [ ] Deliver the clean executive summary and prioritized task list.
- [ ] Present the 3-agent parallel execution proposal with strictly positive framing.

### What NOT to Do (Banned / Auto-Reject):

- [ ] DO NOT execute any code modifications, installs, or migrations during this turn.
- [ ] DO NOT launch sub-agents or start the execution loop inside this prompt.
- [ ] DO NOT cherry-pick or omit any pending tasks (missing even one task is a failure).
- [ ] DO NOT invent step counts without reading the target files.
- [ ] DO NOT use negative question phrasing or double negatives.

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

## Before Writing Code

Read and follow spec folders `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `spec/04-database-conventions/` before writing any code. Error management must be followed. Code must be DRY.

---

## Metadata

- slug: inventory-pending-tasks
- status: active
