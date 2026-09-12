---
name: execute-parent-task
description: Autonomously orchestrate and execute parent tasks by decomposing them into subtasks and running continuous N-step self-loops with strict coding guidelines and error management.
---

# Instruction (must follow): Execute Parent Task (N-Step Continuous Loop & Multi-Agent)

/goal Autonomously orchestrate and execute the parent task by decomposing it into subtasks and running a continuous N-step self-loop until completion without a single failure.

```text
N = 130
```

N = is the number of steps that the agents will perform

- [ ] /goal First N/2 steps will be given for spec writing for AI as given and then breaking down into parts as instruction as subtasks for N/2 steps.
- [ ] /goal second N/2 steps will be given execute the created tasks with following coding guidelines and error manage properly.
- [ ] /learn Ingest `.lovable/memory/01-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, and `spec/03-error-manage/`, `.lovable/coding-guidelines.md` before taking action and also create agent rules in the repo if required to or missing from rules set of agent memory.
- [ ] /learn `.lovable/coding-guidelines.md` and it is must and /goal apply the guidelines in coding every aspect.

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents clear, lean instructions focused on the domain task itself rather than writing massive generic prompts.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet & Task Focus: When spawning a subagent, provide clear, lean instructions that focus on the actual domain task itself rather than writing massive generic meta-prompts. DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction. The subagent MUST read the necessary files itself.

## 2. Phase 1: Write the Implementation Spec & Subtasks FIRST

Before doing anything else, you MUST write a highly detailed execution spec.

- What to write: Break down the parent task into a detailed architectural plan, code review guides, and embedded coding guidelines.
- Where to save it: Save this master plan into `.lovable/plans/pending/xx-<slug>.md`. Do not hallucinate folders.
- **Strict Relative Git Paths Mandate (TOTAL BAN on Absolute Paths / `file:///` URIs):**
  When breaking down tasks into subtasks (`.lovable/plans/subtasks/xx-<slug>/*.md`), writing specs, updating memory logs (`.lovable/memory/issues/`), or citing guidelines, ALL file paths and markdown links MUST be **strictly relative to the git repository root**. NEVER write absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`) into any committed or created files.
  - ❌ **BAD (Absolute path / file URI):**
    `- [SSH Commands](file:///absolute/path/to/...)`
    `- [Target File](file:///absolute/path/to/cmd/main.go)`
  - ✅ **GOOD (Strict relative Git path):**
    `- [SSH Commands](spec/13-generic-cli/01-index.md)`
    `- Target File: cmd/main.go`
- **Create a Task-Specific Rule Set:** Before executing, analyze the specific task domain and explicitly write down 3-5 custom rules or constraints unique to this task inside the spec file. This prevents domain-specific regressions and forces sub-agents to follow exact architectures.
- Subtasks: You MUST break the plan down into a lean set of focused subtask files inside `.lovable/plans/subtasks/xx-<slug>/`. Every subtask file must contain actionable instructions focused on the domain task itself (Task Focus Over Meta-Prompting) with strictly relative Git paths.

## 3. Non-Negotiable Core Rules (Auto-Reject on Violation)

1. Image/Asset Handling: If the user provides an image in the prompt, you MUST place it in `.lovable/assets/<category>/xx-<slug>.<ext>`. NEVER place images in random root directories.

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/01-index.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.lovable/plans/01-index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] Staged files sanitized of artifact zips and temporary scratch files.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.lovable/coding-guidelines.md`.
- [ ] /learn and apply as a /goal  `.lovable/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.
- [ ] Error Manage Checklist: I have fully read and enforced the error management files at `spec/03-error-manage/`. I understand which files to follow (architecture, response envelopes) and how to follow them (never swallow errors, always wrap with context).
- [ ] Boolean Examples & Fixations: All boolean variables MUST begin with is and has only (can, should, was, etc. are banned) (e. NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`).g., `isReady`, `hasData`). NEVER use negative booleans (e.g., `isNotReady`, `disableCache`). NEVER invert success checks (e.g., `!response.isSuccess` is banned; use `response.isFail`).
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Semantic Tests: All unit test names are strictly semantic and behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`). `TestHandleComp100` is an immediate failure.
- [ ] Function Size: No function exceeds 15 lines. Long arguments are split across lines (max 100 chars).
- [ ] Error Handling (AppError): Errors use domain-specific `AppError` or custom `AppException` (for C#/OOP), not generic base `Error`.
- [ ] Code adheres to explicit booleans, `Type` suffixed Enums, and error wrapper rules.
- [ ] Formatting & Acronyms: Spacing rules are strictly followed. Acronyms are strictly PascalCase (`SwapIpWindows` not `SwapIPWindows`).
- [ ] Temp & Failure Folder Isolation: All temporary test files, caches, and runner artifacts are isolated within `.lovable/temp/`. Never create `.tmp/` at root. Failed tests/gates write to `.lovable/temp/failures/`; passing tests remain completely silent and produce zero disk files.
- [ ] Runner In-Flight ETA Wait Protocol: When executing or checking background runners/tests (with runner heartbeats every 25 seconds), agents MUST sleep/wait for 1 minute (60s) each time, or dynamically sleep for the remaining ETA duration read from `.lovable/temp/runner-eta.json` (or based on previous total approximate delay) instead of busy-polling.
- [ ] Centralized Test Inventory & Incremental Caching: All unit tests are cataloged in `.lovable/test-inventory.json` with strictly repository-relative paths (`target_file`, `test_file`). First run profiles all tests for baseline timings; subsequent runs execute incrementally only when target code or test files change. Slow test threshold defaults to `4.0s` (`GITMAP_SLOW_TEST_THRESHOLD`).
- [ ] Dual-Queue Worker Pools: Slow tests run in a dedicated 4-worker pool running at most 2 tests at a time per batch. Fast tests run in a 4-worker pool running at most 4 tests at a time, pulling in chunks of 100 tests from the test inventory queue until all are complete.
- [ ] Fast-forward commits created and pushed without rewriting published git history.
- [ ] Continuous loop maintained; only pausing to ask for "continue" on critical unrecoverable failures.

## 4. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal  You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: I have fully read and strictly enforced every file in `spec/02-coding-guidelines/` and `.lovable/coding-guidelines.md`.
- [ ] Error Management: I have read and enforced `spec/03-error-manage/`. I used `AppError`/`AppException` and did not swallow errors.
- [ ] Boolean Conventions: All booleans begin with is or has ONLY (all other prefixes like can, should, was, will, did, must are banned) (e.g., `isFail`, `hasData`). NO negatives (`!isSuccess` is banned, use `isFail`).
- [ ] Semantic Naming: Absolutely NO generic garbage names (`temp`, `data`, `obj`, `comp_100`). All unit tests are behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`).
- [ ] Formatting: Signatures > 3 parameters or > 100 chars are split to one parameter per line. Newlines around every Markdown header (MD022) and lists are surrounded by blank lines (MD032).
- [ ] Acronyms & Magic Strings: Acronyms are PascalCase (`UserId` not `UserID`). Magic strings/numbers are extracted to constants.
- [ ] /learn the section as a /goal (AI Fix Scripts Memory)[#AI Fix Scripts Memory]
- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

## 4. AI Fix Scripts Memory

- [ ] /goal All helper scripts were written to `03-ai-scripts/`, documented in its `03-ai-scripts/01-index.md`, and committed to Git for reuse.
- [ ] `03-ai-scripts/01-index.md` Should contain details tag, sequence of the script naming and why it is there and what it should do.
- [ ] Should first scan /learn all from that index.md file what is exist and can be reused.

## 5. Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] Echo Back the Spec: I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] Pre-Commit Diff Proof: I have executed `git status` or `git diff --stat` and verified that the files I claim to have modified are actually listed as modified in the terminal output before committing.
- [ ] No Placeholder Search: I ran a regex search for `TODO` and `\[.*\]` in my modified files and confirmed I left zero placeholders behind. I actually wrote the implementation.
- [ ] Index Sync Deadman Switch: I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.lovable/what-to-read.md`. I did not leave any orphaned files.
- [ ] Blast Radius Acknowledgment: Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## Task Consolidation & File Reduction (End of Loop)

> **CRITICAL:** To reduce markdown file count and bloat, you MUST consolidate subtasks when a parent task is 100% complete.

When all subtasks for a parent task (`.lovable/plans/pending/xx-<slug>.md`) are finished, execute this final cleanup step before ending the run:
1. Combine all the completed granular subtasks from `.lovable/plans/subtasks/xx-<slug>/*.md` into a single consolidated file at `.lovable/plans/completed/xx-<slug>.md`.
2. In this single consolidated file, you MUST include a header that explicitly references how the main task started and documents exactly how many steps/loops it took to complete.
3. Delete the original granular `.md` files in `.lovable/plans/subtasks/xx-<slug>/` so that only the single consolidated file remains.
4. Delete the original parent plan `.lovable/plans/pending/xx-<slug>.md`.
5. Update `.lovable/plans/01-index.md` to point to the newly consolidated completed file.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
