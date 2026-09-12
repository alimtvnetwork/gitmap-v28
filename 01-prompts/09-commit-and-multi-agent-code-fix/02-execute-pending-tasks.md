# Multi-Agent Task Dispatcher & Subtask Execution — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously orchestrate and execute ALL pending tasks in a continuous N-step self-loop until the entire queue is completely resolved without a single failure.

## 1. Initial State: Audit, Casing Check & Re-Sequence Pending Tasks

Before starting execution, you must ensure the queue of pending tasks is properly ordered, named, and structured into Execution Waves.

- Verify root readme is strictly lowercase `readme.md`.
- Read the `.lovable/plans/pending/` directory and `.lovable/plans/01-index.md`.
- Count exactly how many pending tasks exist.
- Naming Correction: Check if the pending task files are correctly sequenced with a 2-digit numerical prefix (e.g., `01-<slug>.md`, `02-<slug>.md`).
- If naming is incorrect or missing prefixes, fix it immediately. Rename the files to follow sequential `01-`, `02-` format and update `.lovable/plans/01-index.md` to match in the same operation.

- **Create a Task-Specific Rule Set:** Before executing, analyze the specific task domain and explicitly write down 3-5 custom rules or constraints unique to this task inside the spec file. This prevents domain-specific regressions and forces sub-agents to follow exact architectures.
- Execution Waves: Group tasks into Execution Waves (Wave 1: Schemas/DB/wrappers; Wave 2: Business services; Wave 3: UI & docs).

## 2. Uninterrupted Autonomous Execution (Self-Looping & Locking Matrix)

You are the sole orchestrator. Your job is to complete ALL pending tasks without stopping.

- Make a Great Plan: Analyze all pending tasks and devise a comprehensive execution plan. Tasks exceeding 7 steps must be decomposed into `plans/subtasks/xx-<slug>/`.
- Do NOT Ask Questions: Do not stop to ask the user for permission. Do not stop to ask clarifying questions.
- File Collision Locking Matrix (`active-locks.json`): Register active target files in `.lovable/01-index.md` so parallel tasks touch completely disjoint files.
- Self-Loop: Self-loop continuously until every single pending task in the queue is verifiably completed.
- 3-Strike Rollback: If an agent fails unit tests or builds 3 consecutive times, automatically rollback dirty working tree (`git checkout -- <files>`), log failure context to `.lovable/plan.md`, and advance to the next disjoint task.

## 3. High-Stakes Code Standards & Root Cause Analysis

While executing the pending tasks, you must adhere strictly to the project's code standards and root cause protocols:

- Root Cause First: Find the root cause of every problem before applying any fix. Record the root cause into `.lovable/` memory before touching code.
- Error Management: All caught errors must be explicitly logged following guidelines in `spec/03-error-manage/`. Use the established query wrapper that automatically logs failures.
- No Magic Values: Do not introduce any magic strings or magic numbers anywhere unless explicitly for the logger.
- Enums Over Unions: Replace TypeScript string union types (e.g., `"pass" | "fail"`) with Enums.
- Enum Naming: Every Enum name must end with the `Type` suffix (e.g., `StatusType`). Enum values must use `PascalCase` (e.g., `ActiveState`).
- Boolean State Checks: Always use explicit boolean state checks like `response.isFail`. NEVER use inverted success booleans like `!response.isSuccess`.
- DRY Code: Reuse constants; never duplicate them.

## 4. Sub-Agent Orchestration

To speed up the work, you may spawn sub-agents to handle independent chunks of the pending tasks.

- Spawn a maximum of 2 to 3 sub-agents concurrently to avoid RAM and caching issues.
- Give sub-agents highly specific titles (e.g., `Refactoring Auth Service`, `Fixing DB Connection`).
- Enforce the sub-agent lifecycle: they must read their subtask file, mark it `🔄 In Progress`, do the work, and mark it `✅ Done` with a summary of changed files.
- Sub-agents only write to the file system; they NEVER commit to Git.

## 5. End-of-Loop Final Verification, Artifact Sanitizer & Push

Once ALL pending tasks have been completed and marked `✅ Done`:

- Final Verification: Check full build, run all local unit tests, and check CI/CD status. Fix any build failures or failing tests immediately.
- Artifact Sanitizer & Git History Guard: Ensure no zip archives, temporary scratch files, test data, or binaries are committed. Never rewrite published Git history (no force push, no rebase, no squash) to protect Lovable editor synchronization.
- Commit: Group all completed work into a single logical Git commit with a clear, descriptive message summarizing executed tasks.
- Push: Push the commit to the remote GitHub repository. Pushing after the commit is non-negotiable.

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

- [ ] Self-loop continuously until every pending task is completed; do not stop until the queue is completely empty.
- [ ] Audited `.lovable/plans/pending/` and re-sequenced task filenames to `01-`, `02-`, etc., if incorrectly named.
- [ ] Grouped tasks into Execution Waves and checked `.lovable/01-index.md` for file collisions.
- [ ] Executed autonomously via continuous self-looping without stopping to ask user questions.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back with `git checkout` and logged to `last-failure.md`.
- [ ] Followed all high-stakes code standards (Enums with `Type` suffix, `PascalCase` values, explicit `isFail` checks, no magic strings, DRY code).
- [ ] Root causes identified and logged in `.lovable/` before code was patched.
- [ ] Sub-agents followed strict lifecycle with specific titles and never exceeded 3 concurrent instances.
- [ ] End-of-loop verification passed: build is green, unit tests pass.
- [ ] Staged files sanitized of artifact zip bundles, temporary scripts, and test data.
- [ ] Fast-forward commit created and pushed without rewriting Git history.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

---

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Metadata

- slug: execute-pending-tasks
- status: active

## 5. Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] **Echo Back the Spec:** I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] **Pre-Commit Diff Proof:** I have executed `git status` or `git diff --stat` and verified that the files I claim to have modified are actually listed as modified in the terminal output before committing.
- [ ] **No Placeholder Search:** I ran a regex search for `TODO` and `\[.*\]` in my modified files and confirmed I left zero placeholders behind. I actually wrote the implementation.
- [ ] **Index Sync Deadman Switch:** I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.lovable/what-to-read.md`. I did not leave any orphaned files.
- [ ] **Blast Radius Acknowledgment:** Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.
