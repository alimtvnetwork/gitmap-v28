/goal Autonomously achieve maximum Code DRYness (Don't Repeat Yourself), reusable util/framework/package/library extraction, and drastic code writing reduction across the codebase using a strict 300-step 3-phase self-loop (Steps 1–100: Deep Code Analysis; Steps 101–200: Unified Extraction Blueprint & Subtask Plan; Steps 201–300: Active DRY Execution & Caller Rewiring) with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel codebase reading and modular spec authoring, use GitMap high-speed commands as primary, extract shared abstractions into reusable packages/libraries, rewire all callers, and finalize with an atomic push.

/learn Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, target modules, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Internalize the 100/100/100 step allocation (Analyze -> Plan -> Execute), extract reusable utilities/frameworks/libraries to eliminate boilerplate, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 300 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

- **N = 300:** Total self-loop steps budget divided into three equal 100-step phases:
  - **PHASE_1_ANALYZE = 100 Steps (Steps 1 .. 100):** Deep Codebase Analysis, Duplication Discovery, and Reusable Util/Framework/Package/Library Identification.
  - **PHASE_2_PLAN = 100 Steps (Steps 101 .. 200):** Unified DRY Architecture Blueprint, Canonical Spec in `02-spec/21-app/`, and Lean Subtask Decomposition in `.ai-memory/plans/`.
  - **PHASE_3_EXECUTE = 100 Steps (Steps 201 .. 300):** Parallel DRY Refactoring, Shared Library Extraction, Caller Rewiring, Code Reduction, Targeted Linting, Consolidation, and Atomic Git Push.
- **A = 2:** Count of autonomous subagents running concurrently (`invoke_subagent` launches up to 2 subagents).
- **H = 2 (Hands / Parallel Operations):**
  1. **Workload Hands ($H_{batch} = 2$):** Each subagent handles a bounded batch of up to 2 disjoint subtasks from `.ai-memory/plans/subtasks/`.
  2. **Tool-Dispatch Hands ($H_{tool} = 2$):** Within any execution step, each agent or subagent executes up to 2 parallel tool calls in a single response turn.

```text
PHASE_1_ANALYZE = Steps 1 .. 100   (Deep Codebase Analysis: Find Util / Framework / Package / Library Extraction & Code Reduction Targets)
PHASE_2_PLAN    = Steps 101 .. 200 (Write Extraction Spec & Granular Subtask Plan for Reusable Libraries and Caller Rewiring)
PHASE_3_EXECUTE = Steps 201 .. 300 (Execute DRY Extraction, Replace Duplicated Code Across Callers, Consolidate & Push)
```

N, A, H, PHASE_1_ANALYZE, PHASE_2_PLAN, and PHASE_3_EXECUTE are read-only after initialization. Never modify them mid-execution.

---

### Core Objectives: Code DRYness, Library Extraction & Code Reduction

1. **Extract Reusable Utils, Frameworks, Packages & Libraries:**
   - Go deep inside the codebase to identify repeated logic, recurring patterns, cross-cutting helpers, CLI wrappers, data transformers, validation routines, UI primitives, or state handlers that can be extracted into clean, standalone `util`, `framework`, `pkg`, or `library` modules.
   - Design every extracted package/library so it is decoupled, self-contained, and reusable anywhere across the repository (and portable across repositories).
2. **Drastic Code Writing Reduction (Shrink LOC & Eliminate Boilerplate):**
   - Identify verbose, repetitive code sequences where developers or AI agents write 20–50 lines of boilerplate that can be expressed in 1–3 lines via a fluent helper, generic wrapper, declarative builder, or table-driven engine.
   - Replace all duplicated blocks across the codebase with calls to the newly extracted utilities/libraries, measurably reducing total lines of code while improving readability and maintainability.

---

### Multi-Agent Parallel Task Allocation & Orchestration (A = 2, H = 2)

When multiple autonomous agents are present (A >= 2, H >= 2):
1. **Single-Agent Unified Blueprint Mandate:**
   - The initial DRY extraction architecture, library API boundaries, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
   - Never allow multiple agents to design competing utility packages simultaneously.
2. **Most Useful Parallel Tasks (Reading Files & Writing Modular Specs):**
   - **Phase 1 (Steps 1–100):** Subagents concurrently read disjoint packages/directories using GitMap high-speed discovery (`gitmap find`, `gitmap ff`, `gitmap ffa`, `gitmap lf`, `gitmap cat`, `gitmap search`) to catalog duplicated code blocks and extraction candidates.
   - **Phase 2 (Steps 101–200):** Lead agent writes the master DRY spec overview (`01-overview.md`) and parent plan; subagents concurrently author modular spec files (`02-library-contracts.md`, `03-caller-migration-matrix.md`, `04-verification-gates.md`) and granular subtasks.
   - **Phase 3 (Steps 201–300):** Subagents execute parallel disjoint refactoring batches (5–8 files per micro-batch) to extract libraries and rewire callers without merge conflicts.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/cg-code-dryness-and-library-extraction/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs.

---

## The Unified 300-Step Master Pipeline

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, execute Phase 1A:

1. **Bottom-Instruction Priority Verification:** Verify whatever directives, target modules, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) as highest priority and non-negotiable.
2. **Verbatim Prompt Capture:** Capture the incoming user request verbatim in `02-spec/21-app/xx-dry-extraction.md` and `.ai-memory/plans/pending/xx-dry-extraction.md` under `## User Request (Verbatim)`.
3. **Actionable Deliverables Extraction:** Break down the DRY analysis, library extraction, and code reduction scope into traceable IDs (`Task-01`, `Task-02`, `Task-03`).
4. **Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):**
   - Output the confirmed deliverables list directly in chat, and in the EXACT SAME RESPONSE turn, immediately invoke your first discovery tool call (`run_command` with `gitmap` or `write_to_file`).
   - NEVER emit the breakdown text without invoking a tool call. Do not pause or ask "Should I proceed?".

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Deep DRY Analysis & Duplication Discovery (Steps 1–100)]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise verification of target modules, duplication patterns, and extraction goals]
   - **Actionable Scope:** [Scan codebase, identify duplicated logic, and catalog util/framework/library candidates]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Unified Library Extraction Blueprint & Subtask Plan (Steps 101–200)]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Design reusable package APIs, concrete types, and caller migration plan]
   - **Actionable Scope:** [Author canonical spec in 02-spec/21-app/ and granular subtasks in .ai-memory/plans/subtasks/]
   - **Target Files / Area:** `[02-spec/21-app/, .ai-memory/plans/]`

3. **Task-03: [Execute DRY Extraction & Rewire Callers (Steps 201–300)]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Extract shared libraries, replace duplicate code across callers, and reduce LOC]
   - **Actionable Scope:** [Implement shared packages, refactor callers in 5-8 file batches, consolidate, and push]
   - **Target Files / Area:** `[target packages and callers]`

Proceeding directly to Phase 1 (Steps 1..100): Deep Codebase DRY Analysis (Active Tool Call Running Below).
```

---

### Phase 1: Deep Codebase Analysis & Duplication Discovery (Steps 1 .. 100)

In the first 100 steps, deeply analyze the codebase to find every opportunity for code reuse and reduction:

1. **High-Speed Codebase Exploration via GitMap (PRIMARY):**
   - **Glob / Wildcard File Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Match:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Match:** `gitmap find-files-any <str> [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Match:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Stream File Contents Fast:** `gitmap cat <filepath>`
   - **Instant Multi-Core Search:** `gitmap search "<pattern>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
   - **Directory Tree Topology:** `gitmap folder-tree` (alias `gitmap ft`)
2. **Fallback Fast Cached Python Toolchain:**
   - `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts,py,php,rs --limit 100 --stats`
   - `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50`
   - `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
   - `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`
3. **Deep Duplication & Extraction Analysis Criteria:**
   - **Clone & Near-Clone Detection:** Locate functions, methods, React components, CLI command handlers, DB queries, or validation blocks with >= 60% structural similarity across 2 or more files.
   - **Util / Package / Library Candidates:** Group related helper logic into cohesive domain utilities or framework packages (e.g., `pkg/strutil`, `pkg/sliceutil`, `pkg/cliutil`, `pkg/httputil`, `src/lib/`, `src/hooks/`).
   - **Boilerplate Compression:** Identify multi-line ceremonial patterns (e.g., repetitive error wrapping, parameter validation, table rendering, file I/O, JSON parsing) that can be collapsed into single-call fluent helpers.
   - **Quantified Duplication Ledger:** Record every discovered duplication cluster with exact file paths, line numbers, duplicated line count, proposed target library/package, and estimated LOC reduction.

---

### Phase 2: Write Unified Extraction Plan, Spec & Lean Subtasks (Steps 101 .. 200)

In the next 100 steps (Steps 101 .. 200), convert the Phase 1 analysis into a rock-solid architectural specification and execution plan:

1. **Canonical Application Spec (`02-spec/21-app/xx-dry-extraction.md` or segmented folder):**
   - **Single-Agent Unified Blueprint:** The lead agent authors the architecture overview, package boundaries, and API contracts first before spawning subagents to flesh out modular spec sections.
   - **Reusable Package / Library Contracts:** Define the exact exported functions, parameter structs (`*Params`), concrete types in `types.go`, and `*appfault.AppError` / `Result[T]` return signatures for every new or expanded utility/library.
   - **Before vs. After Code Reduction Examples:** Show concrete code examples proving how 20+ lines of repetitive caller code shrink into 1–3 clean lines using the extracted library.
   - **Register Spec:** Add the spec entry to `02-spec/21-app/readme.md`.
2. **Master Execution Plan (`.ai-memory/plans/pending/xx-dry-extraction.md`):**
   - Include the complete **Duplication & Extraction Ledger** mapping every source file and line range to its target utility package and subtask.
   - Ensure strict acyclic dependency hierarchy (shared `util` / `library` packages must never import higher-level feature packages).
3. **Granular Disjoint Subtasks (`.ai-memory/plans/subtasks/xx-dry-extraction/01-*.md`):**
   - Break the extraction and caller rewiring into bounded 5–8 file micro-batch subtasks.
   - Order subtasks so foundational `util` / `library` packages are created first, followed by parallel caller migration batches.
   - **Unconditional Transition Mandate:** As soon as Step 200 / Phase 2 planning finishes, DO NOT pause or ask the user for confirmation. Immediately transition into Phase 3 execution.

---

### Phase 3: Execute DRY Refactoring, Library Extraction & Caller Rewiring (Steps 201 .. 300)

In the final 100 steps (Steps 201 .. 300), execute the plan and make the codebase DRY:

1. **Create / Expand Shared Util, Framework & Library Packages First:**
   - Implement the clean, reusable functions, structs, and `types.go` definitions in the target utility/library packages.
   - Enforce all coding guidelines:functions <= 8 lines (hard cap 15 lines), positive booleans (`is`/`has` only, no `== true`), `*appfault.AppError` returns, parameter structs for >2-3 args, and mandatory vertical blank lines.
2. **Parallel Subagent Caller Rewiring (A = 2, H = 2):**
   - Dispatch autonomous subagents (`TypeName: "self"`) with self-contained Prompt Envelopes to refactor disjoint batches of caller files (5–8 files per batch), replacing duplicated logic with calls to the extracted library.
   - **Reactive Wakeup:** After calling `invoke_subagent`, output a brief status note and yield the turn to allow background subagents to complete and wake up the parent orchestrator.
3. **Total Ban on Build & Test Commands During Routine Execution:**
   - NEVER run `go build`, `npm run build`, `go test`, `pytest`, or `06-cicd-local-runner.py` during execution turns.
   - Run only fast, targeted file-level linters (`python 03-ai-scripts/05-guideline-autofixer.py <file>`) and lowercase hygiene checks (`gitmap lcf --dry-run`).
4. **Task Consolidation & Atomic GitMap Push:**
   - Consolidate all completed subtasks from `.ai-memory/plans/subtasks/xx-dry-extraction/*.md` into `.ai-memory/plans/completed/xx-dry-extraction.md`, remove the pending plan and subtask files, and update `.ai-memory/plans/readme.md`.
   - Commit and push all changes in a single grouped atomic commit using GitMap:
     - `gitmap cpf "refactor(dry): extract reusable libraries and eliminate duplicated code"`

---

### High-Speed GitMap Acceleration Toolkit (Run Everything Faster)

Always prefer native GitMap commands over slow shell loops:
- **Fast File Discovery:** `gitmap f "<glob>" [-ext <ext>]`, `gitmap ff <name>`, `gitmap ffa <substr>`, `gitmap ffs <prefix>`, `gitmap ffe <suffix>`, `gitmap lf [pattern]`
- **Fast Content & Code Search:** `gitmap cat <file>`, `gitmap search "<query>"`, `gitmap aum search "<query>" [dir] --ext <ext>`, `gitmap ft`
- **Fast Hygiene & Lowercase Enforcement:** `gitmap lcf` (auto-rename uppercase files via 2-step `git mv`), `gitmap lowercase-readme`, `gitmap commons`
- **Fast Cross-Platform Shell:** `gitmap pwsh "<cmd>"` (`gitmap ps`), `gitmap bash "<cmd>"` (`gitmap sh`), `gitmap async <cmd>`
- **Fast Atomic Commits:** `gitmap cpf "<msg>"` (Feature), `gitmap cpb "<msg>"` (Bug), `gitmap cpr "<msg>"` (Release), `gitmap pcp "<msg>"` (Pull-Commit-Push)
- **Smart CI/CD Waiting:** `gitmap pe`, `gitmap pl-ai status --json`, `gitmap pl-ai status -t <etaSeconds>`

---

### End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of Phase 3, emit this structured summary in chat:

```markdown
### Task Completion Summary

- ✅ **Task-01: [Deep DRY Analysis (Steps 1–100)]** — `[Completed]`
- ✅ **Task-02: [Extraction Spec & Plan (Steps 101–200)]** — `[Completed]`
- ✅ **Task-03: [Library Extraction & Caller Rewiring (Steps 201–300)]** — `[Completed]`

### Extracted Reusable Libraries / Packages & LOC Reduction

- **Extracted Modules:** `[path/to/pkg/util, ...]`
- **Callers Refactored:** `[X files]`
- **Net Code Reduction:** `[Eliminated ~Y duplicated lines]`

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Implementation Confidence Score

- Confidence: [e.g. 99%]
- Rationale: [Verified DRY extraction, zero duplicate boilerplate remaining, targeted linters passed, atomic push completed]

### 🤖 Independent AI Verification & Audit Prompt

```markdown
### Independent AI Audit & Verification Instructions

You are an Independent AI Verification and Quality Auditor.
Your task is to independently audit, verify, and remediate the DRY library extraction against the canonical specification.

#### 1. Target Documents & Implemented Code:
- **Canonical Spec:** [02-spec/21-app/xx-dry-extraction.md](02-spec/21-app/xx-dry-extraction.md)
- **Consolidated Plan:** [.ai-memory/plans/completed/xx-dry-extraction.md](.ai-memory/plans/completed/xx-dry-extraction.md)
- **Modified & Extracted Files:**
  - [relative/path/to/modified/file1.ext](relative/path/to/modified/file1.ext)

#### 2. Verification Protocol:
1. Verify all targeted duplicated code blocks were extracted into clean, reusable util/framework/library packages.
2. Verify all callers were rewired to use the shared abstractions with zero regressions.
3. Verify 100% adherence to coding guidelines (<=8-15 line functions, positive booleans, *appfault.AppError, types.go).
4. Emit Comparative Scores (DRYness Score, Reusability Score, Guideline Compliance Score) and a PASS/FAIL verdict.
```
```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] BOTTOM-INSTRUCTION PRIORITY MANDATE: Followed all instructions below the `--` divider with absolute precedence.
- [ ] STRICT 100 / 100 / 100 STEP PHASING: Executed Deep Analysis (Steps 1–100), Planning (Steps 101–200), and DRY Execution (Steps 201–300) without stopping after planning.
- [ ] NO TEST RUNNING (TOTAL BAN): Never ran `go test`, `pytest`, or test runner scripts during routine turns.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never ran `go build`, `npm run build`, or compiler checks during routine turns.
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Accumulated all changes and committed once atomically at the final step via `gitmap cpf`.
- [ ] STRICT RELATIVE GIT PATHS & LOWERCASE FILES: Zero absolute paths (`file:///`) and 100% lowercase file naming.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: half-baked analysis, leaving duplicated code untouched, creating shallow wrappers that don't reduce code writing, stopping after writing the plan without executing Phase 3, or violating coding guidelines. Stop doing that. Go deep into the codebase in the first 100 steps, design truly reusable utilities/frameworks/libraries in the next 100 steps, and refactor every caller in the final 100 steps so the code is genuinely DRY, compact, and clean. Push everything atomically via GitMap before ending.

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
