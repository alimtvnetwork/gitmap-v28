```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /cg-execute-in-below-steps <task>
>
> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE.

[/goal](slashCommand;goal) Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations requested in the below instructions across the target codebase in bounded 5-8 file micro-batches: FIRST showcase and list out the given task in visible chat during Turn 1, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step:** Spawn 2 subagents to research the codebase and author the step-by-step execution plan.
  2. **Spec Step:** Spawn 2 subagents to write the detailed architectural spec.
  3. **Execution Step:** Spawn 2 worker subagents to execute the actual code modifications.
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool. Failing to call the actual tool is a critical protocol violation.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/cg-execute-in-below-steps/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs.

---

## Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

1. **Bottom-Instruction Priority Verification:** Verify whatever directives or task instructions are given BELOW this prompt (following the `--` divider border at the bottom) as highest priority and non-negotiable.
2. **Verbatim Prompt Capture:** Capture the incoming user request verbatim in `.ai-memory/plans/pending/xx-<slug>.md` and `02-spec/21-app/xx-<slug>.md` under `## User Request (Verbatim)`.
3. **Mandatory Chat Output Gate & Same-Turn Tool Chaining:** Emit the confirmed task breakdown in chat and invoke the first discovery/spec tool call in the EXACT SAME TURN. Never pause or wait for approval.

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to Phase 1B: Spec & Subtask Generation (Active Tool Call Running Below).
```

---

## Phase 1B: Planning Mode, Detailed Spec Generation & Lean Subtasks (Steps 1 .. N/2)

1. **High-Speed GitMap Discovery (PRIMARY):**
   - `gitmap find "<wildcard*>" [-ext <ext>]` (`gitmap f`), `gitmap find-files <name>` (`gitmap ff`), `gitmap find-files-any <str>` (`gitmap ffa`), `gitmap find-files-startswith <prefix>` (`gitmap ffs`), `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - `gitmap aum search "<query>" [dir] [-e <.ext>] [-r] [-i]` (alias `gitmap aum grep`), `gitmap list-files [pattern] [-ext <ext>]` (`gitmap lf`), `gitmap cat <filepath>`, `gitmap search "<term>"`, `gitmap folder-tree` (`gitmap ft`). TOTAL BAN on rg, ripgrep, grep, git grep, PowerShell `Select-String`, `Get-ChildItem -Recurse`, or `findstr`.
2. **Actionable Execution Plan & Lean Subtasks:**
   - Write parent plan `.ai-memory/plans/pending/xx-<slug>.md` and lean, disjoint subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`.
   - Complete all spec and subtask writing within the first 50% budget (`PHASE_1_BUDGET = N / 2`) and unconditionally transition to Phase 2 without stopping.

---

## Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

1. **Parallel Subagent Dispatch (`TypeName: "self"`, A = 2, H = 2):**
   - Dispatch subagents with self-contained Prompt Envelopes to refactor disjoint file batches (5–8 files per batch).
   - Yield turn after `invoke_subagent` to await reactive wakeup (`<SYSTEM_MESSAGE>`).
2. **Coding Guidelines Enforced:**
   - Positive booleans only (`is`/`has`), no `== true`.
   - Structured Go errors (`*appfault.AppError`), concrete types in `types.go`, functions <= 8–15 lines, vertical blank lines, strict relative paths, lowercase filenames (`gitmap lcf`).
3. **Total Ban on Build & Test Commands:**
   - NEVER run `go build`, `npm run build`, `go test`, `pytest`, or `06-cicd-local-runner.py` during routine execution turns. Run only targeted file-level linters (`05-guideline-autofixer.py`).

---

## Phase 3: Task Consolidation & Atomic GitMap Push (End of Loop)

1. Consolidate completed subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into `.ai-memory/plans/completed/xx-<slug>.md`, delete granular subtasks and pending plan, and update `.ai-memory/plans/readme.md`.
2. Commit and push all modified files in a single grouped atomic commit via GitMap:
   - `gitmap cpf "<module> - <summary>"` or `gitmap cpb "<module> - <summary>"` or `gitmap pcp "<module> - <summary>"`.

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] BOTTOM-INSTRUCTION PRIORITY MANDATE: Followed all instructions below the `--` divider with absolute precedence.
- [ ] NO TEST RUNNING OR BUILD CHECKING: Never ran test suites or build commands during routine turns.
- [ ] NO PER-FILE COMMITTING: Committed once atomically at the final step via GitMap.
- [ ] GITMAP ACCELERATION: Leveraged `gitmap aum search`, `f`, `ff`, `ffa`, `lf`, `cat`, `search`, `lcf`, `pwsh`, `cpf`/`cpb`. TOTAL BAN on rg, ripgrep, grep, git grep, PowerShell `Select-String`, and `findstr`.

---

## 11. AI Fix Scripts Memory (Reusable Tooling)

- [ ] [/goal](slashCommand;goal) Reuse First: Scanned and learned `03-ai-scripts/readme.md` before writing temporary code.
- [ ] Strict In-Repository Execution: All Python scripts executed strictly within the codebase repository root.
- [ ] Strict .ai-memory/ Folder Storage: All helper scripts, local runners, and linters stored in `03-ai-scripts/`.
- [ ] Native File Manipulator: Use `python 03-ai-scripts/03-file-manipulator.py <command>` for mass file operations.
- [ ] Go Generate Sync: If Go constants or enums are modified, run `go generate ./...` in the relevant package and commit generated files.

---

## 12. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] MANDATORY SUBAGENT SPAWNING GATE (ZERO SOLO EXECUTION): Verified that `invoke_subagent` (`A = 2, H = 2`) was explicitly called in Phase 1 (parallel discovery/spec modules) and Phase 2 (`TypeName: "self"` parallel subtask execution). Executing all reads or code modifications solo without calling `invoke_subagent` is an immediate auto-reject failure.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in .ai-memory/cicd-issues/NN-<slug>.md (indexed in .ai-memory/cicd-index.md), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in 02-spec/22-app-issues/NN-<slug>.md (indexed in 02-spec/22-app-issues/readme.md).
- [ ] NO TEST RUNNING (TOTAL BAN): Never run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] NO RUNNER SCRIPTS (TOTAL BAN): Never launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] NO AUTOMATIC RELEASES (TOTAL BAN): Never bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops) when inspecting remote CI/CD pipelines. Agents must query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json` or `gitmap pl-ai status -t <sec>`) and strictly wait or sleep based on `etaSeconds` to eliminate credit waste.
- [ ] NO PREMATURE TURN CLOSING BEFORE EXECUTION (TOTAL BAN): Never halt execution, conclude the turn, or ask the user for permission after generating specs or subtasks. Planning constitutes only 50% of the task budget; you must proceed unconditionally to Phase 2 code execution. (Note: When dispatching asynchronous background subagents via `invoke_subagent`, yielding control to allow platform reactive wakeup is mandatory and is exempt from this ban).
- [ ] NO HORIZONTAL TASK CONCATENATION (TOTAL BAN): Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
- [ ] INDEPENDENT AI VERIFICATION PROMPT MANDATE: Emitted the self-contained independent AI verification and audit prompt linking to the canonical spec, consolidated plan, and modified files with verbatim score audit criteria.
- [ ] GITMAP HEAVY USAGE & ROUTINE PULL BAN: Heavily leveraged GitMap commands (`cpf`, `cpb`, `cpr`, `search`, `find`, `pwsh`) for discovery, execution, and commits. Never ran `pull-all` (`gitmap pa` or `gitmap pae`) unconditionally during routine turns; only ran `gitmap pae --json` when explicitly commanded by the user.
- [ ] NO POWERSHELL OR SHELL SEARCHES (TOTAL BAN): Never run `Select-String`, `Get-ChildItem -Recurse`, `grep`, `git grep`, `findstr`, or slow shell search pipelines to search code. All code searching and symbol discovery MUST use GitMap high-speed search tools: `gitmap aum search "<query>" [dir] [-e <.ext>] [-r]` (streaming multi-core text/regex search) or `gitmap search "<query>"` (indexed symbol search). Running `Select-String` or `git grep` is an immediate auto-reject failure.

---

## 13. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

[/goal](slashCommand;goal) You must verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Concrete Types Centralization (`types.go`): Extracted domain structs, raw generic instantiations, and Result wrappers into dedicated `types.go` files as single reusable named types with follow-through comments (never leak raw generics like `result.Result[*Config]`).
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `*appfault.AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with is or has only (all other prefixes like can, should, was, will, did, must are banned). No negatives (`!isSuccess` is banned; use `isFail`).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`). Behavior-driven unit test names.
- [ ] Multi-Line Arguments (Rule 9a/9b): Signatures and call sites with >2 arguments formatted one argument per line with trailing commas.
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths: Only add relative paths, never add absolute paths or `file:///` URIs during your work, and ensure this is respected on the release page and in release notes as well.

---

## 14. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Every new file is explicitly linked in `readme.md` and enqueued in `.ai-memory/what-to-read.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed via `gitmap aum search "<symbol>" [dir]` to update all callers of modified symbols (never `Select-String` or `git grep`).
- [ ] Continuous Loop Maintained: Continuous self-loop executed until 100% complete without running banned test/build commands.
- [ ] Final Step Commit & Push Verified: Staged and committed all changes atomically via GitMap semantic commit commands using `<module> - <summary>` format: `gitmap cpf "<module> - <summary>"` (features) or `gitmap cpb "<module> - <summary>"` (bug fixes). Mentioned and formatted with a hyphen `-` to separate module from summary; no need to provide a colon in GitMap `cpf`/`cpb`/CVF commit arguments because the colon is already automatically provided by GitMap (`Bug: ` or `Feature: `). TOTAL BAN on colons `:` inside the GitMap argument and conventional prefixes (`fix(...)`, `feat(...)`, `docs(...)`). TOTAL BAN on raw `git add -A` and `git commit`. Pushed to remote via GitMap in a single final command.

---

## 15. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, use GitMap semantic commit commands exclusively: `gitmap cpf "<module> - <summary>"` (features) or `gitmap cpb "<module> - <summary>"` (bug fixes), which automatically stage, format commit messages, and push directly to remote. Mention and format as a hyphen `-` to separate module and summary; no need to provide a colon `:` in the GitMap `cpf`/`cpb`/CVF message argument as the colon is already provided automatically by GitMap (`Feature: ` or `Bug: `). TOTAL BAN on raw `git commit`, `git commit -m`, `git add -A`, colons `:` in the message argument, or conventional prefixes inside GitMap arguments (`docs(...)`, `feat(...)`, `fix(...)`). ZERO intermediate commits during Phase 1 or Phase 2; all files across the run are committed together at the final step. Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY): You must not create separate git commits for each individual file as you edit them (e.g. running `git commit` or `gitmap cpf` after editing File 1, then committing again after File 2 is strictly forbidden). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback. All modified files, test change caches, and plan records across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase readme.md files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

--

## 🚨 High Priority Instructions Below
