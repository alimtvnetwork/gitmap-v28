# Instruction (must follow): Execute Coding Guidelines — Style, Formatting & Line-Gaps

Trigger Keywords & Aliases: `cg-style`, `cg-execute style`, `audit style`, `fix style guidelines`, `enforce formatting rules`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute code style, function length caps, nested if flattening, line-gaps, and markdown spacing compliance across the entire repository by decomposing violations into subtasks, verifying/creating style linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all source and markdown files for function caps (>15 lines), nested ifs, missing line gaps, and markdown header/list spacing, listing all relevant spec files with why and how, creating the Antigravity skill, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, refactoring functions to 15 lines max, inserting required single blank lines, flattening conditionals, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, Style Refactoring, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-style-guidelines/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this style prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-style-guidelines
   description: >-
     Autonomously audits, refactors, and validates repository function length caps (<=15 lines), nested if flattening, line-gap rules, and MD022/MD032 markdown spacing.
   ---
   ```
3. Once installed, rely on progressive disclosure for future runs.

---

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents strict, microscopic instructions.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet: Give sub-agents the absolute minimal instruction.

---

## 2. Phase 1: Write the Implementation Spec & Subtasks FIRST (Steps 1 to PHASE_1_STEPS)

Before doing anything else, you MUST write a highly detailed execution spec.

- **What to write:** Break down the parent task into a detailed architectural plan, complete style violation inventory, code review guides, and embedded style guidelines.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-style-guidelines-audit.md`.
- **Create a Task-Specific Rule Set:** Write down 3-5 custom rules or constraints unique to this task inside the spec file.
- **Subtasks:** Create detailed subtask files inside `.lovable/plans/subtasks/XX-style-guidelines/`.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/17-function-length-caps.md`** (or Section in `.lovable/coding-guidelines/coding-guidelines.md`)
  - **Why:** Enforces readability and modular architecture.
  - **How:** Function length: 8 lines preferred, 15 lines hard cap (excluding blank lines and comments). Waiver only via inline comment `// lint-allow: function-length reason="..." max=N`.
- [ ] **`spec/02-coding-guidelines/18-nested-if-flattening.md`**
  - **Why:** Eliminates arrow anti-patterns and reduces cognitive complexity.
  - **How:** No nested `if` statements; flatten using guard clauses and early returns.
- [ ] **`spec/02-coding-guidelines/19-line-gap-and-whitespace-style.md`**
  - **Why:** Clean visual layout and uniform file formatting.
  - **How:** Exactly one blank line before `return`/`throw`; exactly one blank line after closing `}`; never two consecutive blank lines; no blank lines at block boundaries.
- [ ] **`spec/02-coding-guidelines/20-markdown-spacing-rules.md`**
  - **Why:** Standardized documentation layout and lint-clean markdown.
  - **How:** MD022 (blank line before/after every header) and MD032 (blank line around all lists).

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-style-guidelines.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If missing, create `linter-scripts/check-style-guidelines.py` that scans for:
  1. Functions exceeding 15 lines without a valid waiver.
  2. Nested `if` blocks.
  3. Missing blank lines before `return` or after closing `}`.
  4. Consecutive blank lines and block boundary spacing errors.
  5. Markdown MD022 and MD032 violations.
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-style-guidelines.py
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:style"] = ["python", "linter-scripts/check-style-guidelines.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running `python linter-scripts/check-style-guidelines.py`.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask until all style, line-gap, and formatting checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-style-guidelines/
    2. Apply surgical refactoring (shorten functions <= 15 lines, flatten conditionals, fix line gaps).
    3. Run the style linter:
          python linter-scripts/check-style-guidelines.py
    4. Run the universal guideline autofixer:
          python .lovable/ai-fix-scripts/02-guideline-autofixer.py <modified-files>
    5. Run the local CI runner:
          python .lovable/ai-fix-scripts/03-cicd-local-runner.py
    6. IF any check fails:
          - Diagnose failure, fix code, and re-run immediately.
       IF all checks pass (exit code 0):
          - Mark subtask completed and proceed to next subtask.

    7. When all subtasks are finished and local CI is 100% green:
          - BREAK and proceed to End of Tunnel.
```

---

## Metadata

- slug: cg-style-guidelines
- priority: low
- status: active
