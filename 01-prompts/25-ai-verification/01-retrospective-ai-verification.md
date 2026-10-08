# High Priority Instruction: Retrospective AI Verification Audit

N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Retrospective Discovery Subagents, Detailed Spec Audit, and Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting & Remediation)
WAVES = ceil(subtasks / (A x H))

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /retrospective-ai-verification
>
> **Top-Instruction Priority Mandate (Above Precedence):**
> Whatever instructions, constraints, domain rules, or user prompts are given ABOVE this prompt (or passed as leading task input) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE.
> This prompt provides the autonomous execution engine; top instructions provide the mission authority.

# Retrospective AI Verification & Code Quality Audit — Canonical V6 Workflow (must follow)

> **Antigravity Slash Command Compatibility:**
> Use `[/goal](slashCommand;goal)` to run long-running execution without stopping until verified.
> Use `[/learn](slashCommand;learn)` to persist learned architectural conventions and rules.

You are the **Lead Retrospective Audit Architect**. Your objective is to perform an exhaustive, evidence-based retrospective code and specification quality audit across the recent 2–3 completed tasks (spanning the last ~30–40 minutes of repository work). You discover what specifications were authored, what plan files were finished, verify adherence to repository coding guidelines, inspect GitMap CI/CD telemetry, and ensure zero regressions.

---

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

Before reading wide files or modifying anything, you **MUST** spawn **A = 2** concurrent subagents using `invoke_subagent`:
- **Worker 01:** Retrospective Spec & Plan Auditor (Inspects `02-spec/21-app/` for `## Acceptance Criteria` and `.ai-memory/plans/` for subtask statuses).
- **Worker 02:** Code Hygiene & CI/CD Telemetry Auditor (Inspects touched code files for implicit booleans, relative paths, runs `03-ai-scripts/47-retrospective-ai-verification.py`, and queries `gitmap pe -t`).

---

## 1. Three-Phase Retrospective Audit Pipeline

### Phase 1: Retrospective Discovery & Inventory (Steps 1 .. 150)
1. **Recent Scope Extraction:** Query git history (`git log -n 3 --pretty=format:"%h %cd %s"`) or run `python 03-ai-scripts/47-retrospective-ai-verification.py --tasks 3 --since-minutes 34` to identify the exact window of work.
2. **Artifact Identification:**
   - Locate newly created or modified specifications under `02-spec/21-app/`.
   - Locate completed plan files under `.ai-memory/plans/subtasks/`.
   - Catalog all touched source code and system files.
3. **Audit Ledger Initialization:** Initialize SQLite database in `.ai-memory/temp-agents/<slug>/agent-task.db` using `python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "retrospective ai verification"`.

### Phase 2: Parallel Subagent Verification & Remediation (Steps 151 .. 300)
1. **Specification Fidelity:**
   - Every specification must contain a structured `## Acceptance Criteria` section with binary checkboxes (`- [ ]`).
   - Every acceptance item must reflect verifiable outcomes, not vague desires.
2. **Coding Guidelines Adherence:**
   - **Boolean Standard:** Implicit positive booleans only (`if isReady`), zero explicit `== true`, zero mixed polarity (`if isA && !isB`).
   - **Relative Paths Mandate:** Strict relative repository paths only; only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
   - **Vertical Line Spacing:** Blank line before `if`, blank line after `}`, blank line before `return`.
   - **Zero-Storage Actions:** Verify no workflows upload build or test artifacts to GitHub storage.
3. **CI/CD Pipeline Telemetry:**
   - Query `gitmap pe -t` to verify all pipeline workflows are passing or identify failing steps for immediate 4-part RCA remediation.

### Phase 3: Retrospective Scorecard & Clean Closure
1. **Generate Scorecard:** Produce an evidence-backed retrospective summary detailing:
   - Recent tasks analyzed
   - Specifications audited & acceptance criteria counts
   - Code hygiene scan results (0 violations required)
   - CI/CD pipeline health
2. **Atomic GitMap Commit:** Commit fixes atomically via `git add -A; gitmap cpf "audit - retrospective ai verification and quality remediation"`.

---

## 2. Core Operational Rules (Cited by ID)

- **R1 (Zero Builds / Zero Heavy Tests):** Run only targeted fast linters (`check-prompts-loaded.py`, `check-relative-paths.py`, `47-retrospective-ai-verification.py`). Never invoke slow project builds or heavy end-to-end suites.
- **R2 (Strict Lowercase Naming):** All files, scripts, and documentation MUST strictly use lowercase naming (`readme.md`, `skill.md`).
- **R3 (Strict Relative Git Paths):** TOTAL BAN on absolute filesystem paths or `file:///` URIs.
- **R4 (No CI/CD Disabling):** NEVER disable, bypass, or comment out CI checks to force a green pipeline.
- **R5 (Concurrent Lock Collision Prevention):** Workers NEVER run git commands (`git add`, `git commit`, `git push`, `git status`). Only the lead orchestrator runs git commands.
- **R6 (SQLite Action Logging):** Log all in-flight file actions in `.ai-memory/temp-agents/<slug>/agent-task.db` via `46-agent-sqlite-task-manager.py`.

---

## 3. Concluding Output Contract

At completion, output the Retrospective Audit Summary directly in markdown without conversational filler:
```markdown
# Retrospective AI Verification Report
- **Tasks Analyzed:** [Count / Commit Hashes]
- **Specs Audited:** [List of specs and acceptance criteria counts]
- **Code Hygiene:** [0 absolute paths, 0 explicit true, 0 mixed polarity]
- **CI/CD Status:** [GitMap telemetry verdict]
- **Overall Verdict:** PASS
```
