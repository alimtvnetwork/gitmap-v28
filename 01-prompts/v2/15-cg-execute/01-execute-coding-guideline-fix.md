# Coding Guideline Fix Execution & Autonomous Linter Remediation — Workflow (must follow)

> [!IMPORTANT]
> Prompt Version: 2.6.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations across the target codebase in bounded 5-8 file micro-batches until 100% green without stopping.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents, default: 2)
H = 2   (Number of hands / parallel operations per agent, default: 2)
```

N = total self-loop steps budget that the agents will perform.
A = count of autonomous subagents running concurrently (default: 2).
H = number of hands / parallel operations per agent (default: 2).

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec with Violation Ledger in .ai-memory/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Active Code Refactoring, Disk Verification, Linter Verification, Local CI Runner Verification, Plan Completion)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

### Multi-Agent Parallel Task Allocation & Orchestration (A = 2, H = 2)

When multiple autonomous agents are present (A >= 2, H >= 2):
1. **Single-Agent Unified Blueprint Mandate:**
   - The initial plan, violation scoping, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
   - Never allow multiple agents to author disjoint or competing audit plans simultaneously. A single coherent architectural vision must lead.
2. **Most Useful Parallel Tasks (Reading Files & Writing Modular Specs):**
   - Once the unified blueprint is authored by the lead agent, the most effective parallel tasks for subagents (A = 2, H = 2) are:
     - **Reading Files:** Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting AST violations in parallel using GitMap AUM (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary and Python scripts (`03-ai-scripts/17-fast-file-reader.py`, `03-ai-scripts/11-fast-file-scanner.py`) as fallback.
     - **Writing Modular Specs & Violation Ledgers:** Authoring modular spec sections and populating the violation ledger in parallel adhering to the lead agent's blueprint.
3. **Spec Writing & Spec Audit Emphasis:**
   - **Spec Writing Section:** The lead agent defines the guideline audit overview and rule boundaries first; subagents are then spawned in parallel to flesh out disjoint spec modules and verification gates concurrently.
   - **Spec Audit Section:** The lead agent establishes the audit methodology and roadmap first; subagents are then dispatched concurrently to inspect disjoint code areas, build the violation ledger, and generate granular subtask files.
4. **Execution Mode (Disjoint Refactoring):**
   - Subagents execute parallel disjoint refactoring tasks across non-overlapping files and run targeted file-level linters (`exit 0`).

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/execute-coding-guideline-fix/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable, and strictly incorporated ahead of all other guidelines.
2. **Verbatim Prompt Capture:** Capture the incoming user request verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under a dedicated `## User Request (Verbatim)` section.
3. **Actionable Deliverables Extraction:** Break down user requirements into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, `Task-03`) under `## Extracted Actionable Task List`.
4. **Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):**
   - You MUST output the confirmed deliverables list directly in chat, but you MUST NOT end your turn or pause after emitting it!
   - In the EXACT SAME RESPONSE turn as your markdown output, you MUST immediately invoke the first tool call (e.g. `write_to_file` to initialize the planning spec or `run_command` for fast discovery).
   - NEVER emit the breakdown text without invoking a tool call. Emitting visible text without a tool call terminates the AI turn, causing the engine to close or hang the conversation while tasks remain pending.
   - TOTAL BAN ON WAITING FOR APPROVAL: Do not pause, do not wait for the user to reply, and do not ask "Should I proceed?". Override any planning mode stop directives. Proceed immediately into Phase 1B in the same turn.
   - Format each task clearly with proper markdown indentation, vertical blank lines, task state (`State: [IN PROGRESS — EXECUTING IMMEDIATELY]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).

Line-by-Line Output Format Structure:
```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and implementation scope]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to Phase 1B: Spec & Subtask Generation (Active Tool Call Running Below).
```

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

## Phase 1B: Scan, Spec Generation & Lean Subtasks (Steps 1 .. N/2)

### Step 1: Single-Agent Unified Blueprint & Master Spec
- Author `.ai-memory/plans/pending/xx-<slug>-audit.md` with an exhaustive Violation Ledger tracking every violation:
  `| Violation Id | File Path | Line Number | Exact Snippet | Planned Fix | Status (PENDING/DONE) |`
- The lead agent defines the audit roadmap and rule boundaries first before spawning subagents.

### Step 2: Multi-Agent Parallel Reading & Discovery (A = 2, H = 2)
- When multiple agents are present, the primary parallel task is reading files. Subagents concurrently inspect disjoint folders using the 2-tier toolchain:
  - **Tier 1 (GitMap AUM Acceleration - PRIMARY):**
    - `gitmap find "<pattern>" [-ext <ext>]` (alias `gitmap f`)
    - `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
    - `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
    - `gitmap cat <filepath>` (stream file to stdout without disk writes)
    - `gitmap search "<term>"` (immediate multi-core filesystem walk)
  - **Tier 2 (Fast Cached Python Scripts - FALLBACK):**
    - `python 03-ai-scripts/11-fast-file-scanner.py --lang <lang> --limit 100 --stats`
    - `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<pattern>" --limit 50`
    - `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> --limit 50`
  - `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

### Step 3: Granular Subtask Decomposition (5-8 Files Micro-Batches)
- Decompose the master plan into granular subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`, `02-<subtask>.md`, each strictly bounded to 5-8 files.
- Subtasks must specify disjoint target files with strict relative git paths.

### Step 4: Subtask Readiness & Spec Audit Gate
- Confirm all violations are mapped to subtasks, audit criteria are met, and target files are disjoint.
- Confirm zero absolute paths and zero `file:///` URIs.

### Step 5: Unconditional Zero-Question Execution Mandate
- As soon as Phase 1 completes, immediately, unconditionally self-loop and transition directly into Phase 2 execution mode. Stopping after spec writing is strictly banned.

---

## Phase 2: Autonomous Code Refactoring & Micro-Batching (Steps N/2+1 .. N)

1. **Parallel Dispatch (A = 2, H = 2):** Spawn up to A = 2 execution subagents (H = 2 operations each) assigned to disjoint subtasks.
2. **Autofixer-First:** Run deterministic AST autofixers (`05-guideline-autofixer.py`, `08-naming-autofixer.py`, `04-newline-fixer.py`, `07-relative-path-fixer.py`) on target files first.
3. **Cognitive Refactoring:** Perform surgical architectural refactoring following authoritative guidelines:
   - Functions strictly <= 8–15 lines body logic.
   - Files <= 80-100 coding lines.
   - Positive booleans only: `is` / `has` prefixes, implicit checks (never `== true`), positive framing, no mixed polarity (`if a && !b`).
   - Guard clauses: Invert early checks and return immediately to eliminate nested if statements.
   - Structured error returns: `*appfault.AppError` context wrappers with `.WithPath()` and `.WithVar()`, monadic `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`.
   - Dedicated `types.go` single reusable named types for domain structs and Result wrappers.
4. **Targeted Quality Linting Only:** Run targeted file-level linters/autofixers on specifically modified files (`exit 0`). DO NOT run full CI/CD pipeline runners (`06-cicd-local-runner.py`), build checks, or test suites during routine execution turns.
5. **Atomic Change Tracking:** Record modified files under lock: `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
6. **Dynamic CI/CD Waiting:** Use `gitmap pipeline-ai status --json` with proportional ETA waiting (never rapid polling).

---

## Phase 3: Task Consolidation & Atomic Git Push

1. Consolidate completed subtasks into `.ai-memory/plans/completed/xx-<slug>.md`.
2. Delete granular subtasks in `.ai-memory/plans/subtasks/xx-<slug>/` and pending plan `.ai-memory/plans/pending/xx-<slug>.md`.
3. Update `.ai-memory/plans/readme.md`.
4. **Final Step Git Commit & Push (Mandatory):** Stage all modified files (`git add -A`), commit them in a single clean grouped atomic commit (`git commit -m "<summary>"`), and push to git (`git push origin <branch>`). Under no circumstances commit each file individually.

---

## End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of all tasks and before concluding the turn, emit this structured summary:

> [!CRITICAL]
> **STRICT LINE-BY-LINE OUTPUT MANDATE (TOTAL BAN ON HORIZONTAL CONCATENATION):**
> Every single completed task in the `Task Completion Summary` MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
> NEVER concatenate multiple tasks horizontally into a single run-on paragraph or single wrapped line.
> In Markdown, consecutive lines without bullet markers (`- `) collapse into a single run-on horizontal sentence. You MUST format each task as a discrete bullet list item (`- ✅`) followed by an explicit newline!
>
> ❌ **BANNED (Horizontal Run-on Concat):**
> `✅ #1. Task-01: [Title] — Completed ✅ #2. Task-02: [Title] — Completed ✅ #3. Task-03: [Title] — Completed`
>
> ✅ **MANDATORY (Strict Line-by-Line Vertical Markdown List):**
> ```markdown
> ### Task Completion Summary
>
> - ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
> - ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`
> - ✅ **Task-03: [Descriptive Task Title]** — `[Completed]`
> ```

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]`
- ✅ **Task-02: [Descriptive Task Title]** — `[Completed]`
(If any task failed or was deferred, mark with `- ❌` or `- ⏳` on its own separate line and explain why)

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Implementation Confidence Score

- Confidence: [e.g. 98% or 100%]
- Rationale: [Detailed explanation of verified quality gates, passing linters, contract adherence, and zero regressions]

### 🤖 Independent AI Verification & Audit Prompt

At the conclusion of the turn, emit this copy-pasteable prompt for an independent auditor AI to verify the implementation against the canonical specification and verbatim requirements:

```markdown
### Independent AI Audit & Verification Instructions

You are an Independent AI Verification and Quality Auditor.
Your task is to independently audit, verify, and remediate the implementation against the canonical specification and verbatim requirements.

#### 1. Target Documents & Implemented Code:
- **Canonical Spec & Verbatim Requirements:** [02-spec/21-app/xx-<slug>.md](02-spec/21-app/xx-<slug>.md)
- **Consolidated Plan & Subtasks:** [.ai-memory/plans/completed/xx-<slug>.md](.ai-memory/plans/completed/xx-<slug>.md)
- **Modified & Implemented Code Files:**
  - [relative/path/to/modified/file1.ext](relative/path/to/modified/file1.ext)
  - [relative/path/to/modified/file2.ext](relative/path/to/modified/file2.ext)

#### 2. Verification Protocol:
1. **Strict Verbatim Inspection:** Read the canonical spec file `02-spec/21-app/xx-<slug>.md` completely, focusing on the Verbatim Requirements and Acceptance Criteria.
2. **Line-by-Line Code Comparison:** Inspect the modified code files line-by-line and verify whether every single verbatim requirement is fully implemented.
3. **Gap & Missing Items Identification:** Itemize all missing, incomplete, placeholder, or non-compliant elements.
4. **Autonomous Self-Loop Remediation:** If any gaps or missing items exist, self-loop and modify the code directly until 100% of the verbatim requirements are met. Do not ask for confirmation.
5. **Comparative Score Audit:** Evaluate and report a comparative score between the verbatim specification and the real-life output:
   - **Verbatim Adherence Score:** [X/100] (Degree to which implementation matches exact spec text and constraints)
   - **Completeness Score:** [Y/100] (Percentage of required features and acceptance criteria satisfied)
   - **Guideline Compliance Score:** [Z/100] (Adherence to booleans, error handling, function sizing, type extraction)
   - **Overall Implementation Score:** [(X+Y+Z)/3 / 100]
6. **Final Audit Verdict:** Emit a clear PASS/FAIL verdict with confidence score and rationale.
```
```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** Never run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`), Go (`go test`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** Never run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** Never launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** Never bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** Never commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single grouped atomic commit at the final step before pushing.
- [ ] **NO RAPID CI/CD POLLING (TOTAL BAN):** Never query or loop rapidly (`gh run view` in tight loops). Query pipeline state using GitMap Pipeline-AI (`gitmap pipeline-ai status --json`) and strictly wait based on `etaSeconds`.
- [ ] **NO STOPPING AFTER SPEC WRITING (TOTAL BAN):** Never halt execution, conclude the turn, or ask the user for permission after generating specs or subtasks. Planning constitutes only 50% of the task budget; you must proceed unconditionally to Phase 2 code execution.
- [ ] **NO HORIZONTAL TASK CONCATENATION (TOTAL BAN):** Never concatenate tasks horizontally in the Task Completion Summary (e.g. NEVER `✅ #1... ✅ #2...` run-on). Every completed task MUST be rendered on its OWN SEPARATE LINE starting with an individual markdown list bullet (`- ✅`).
