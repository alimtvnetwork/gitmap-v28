/goal Autonomously orchestrate and apply concrete, surgical refactoring fixes for all coding guideline violations requested in the below instructions across the target codebase in bounded 5-8 file micro-batches with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and AST inspection, use GitMap AUM as primary, establish a single-agent blueprint during Phase 1 planning (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

/learn Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec with Violation Ledger in .ai-memory/plans/pending/, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Active Code Refactoring, Targeted Linters, Consolidation, Atomic Push)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

### Multi-Agent Parallel Task Allocation & Orchestration (A = 2, H = 2)

1. **Single-Agent Unified Blueprint Mandate:** The initial plan, violation scoping, and lookahead roadmap MUST be authored by a single lead agent first as a unified blueprint before delegating work to subagents.
2. **Parallel Subagent Tasks (A = 2, H = 2):**
   - **Reading Files:** Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting AST violations in parallel using GitMap high-speed discovery (`gitmap find`, `ff`, `ffa`, `ffs`, `ffe`, `lf`, `cat`, `search`, `ft`).
   - **Writing Modular Specs & Violation Ledgers:** Authoring modular spec sections and populating the violation ledger in parallel adhering to the lead agent's blueprint.
   - **Disjoint Refactoring:** Subagents execute parallel disjoint refactoring tasks across non-overlapping files (5–8 files per micro-batch) and run targeted file-level linters (`exit 0`).

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
   - `gitmap list-files [pattern] [-ext <ext>]` (`gitmap lf`), `gitmap cat <filepath>`, `gitmap search "<term>"`, `gitmap aum search "<query>" [dir] --ext <ext>`, `gitmap folder-tree` (`gitmap ft`)
2. **Actionable Execution Plan & Lean Subtasks:**
   - Write parent plan `.ai-memory/plans/pending/xx-<slug>.md` and lean, disjoint subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`.
   - Complete all spec and subtask writing within the first 50% budget (`PHASE_1_STEPS = N / 2`) and unconditionally transition to Phase 2 without stopping.

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
   - `gitmap cpf "<summary>"` or `gitmap cpb "<summary>"` or `gitmap pcp "<summary>"`.

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] BOTTOM-INSTRUCTION PRIORITY MANDATE: Followed all instructions below the `--` divider with absolute precedence.
- [ ] NO TEST RUNNING OR BUILD CHECKING: Never ran test suites or build commands during routine turns.
- [ ] NO PER-FILE COMMITTING: Committed once atomically at the final step via GitMap.
- [ ] GITMAP ACCELERATION: Leveraged `gitmap f`, `ff`, `ffa`, `lf`, `cat`, `search`, `lcf`, `pwsh`, `cpf`/`cpb`.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: skipping guidelines, leaving placeholders, or stopping after planning. Read the codebase deeply, refactor every target file cleanly, and push everything atomically before ending.

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
