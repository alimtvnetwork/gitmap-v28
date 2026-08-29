# Instruction (must follow): Execute Coding Guidelines — Booleans, Naming & Enums

Trigger Keywords & Aliases: `cg-boolean`, `cg-execute boolean`, `audit boolean`, `fix boolean guidelines`, `enforce naming guidelines`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute boolean and identifier naming compliance across the entire repository by decomposing violations into subtasks, verifying/creating naming linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all active files for non-standard boolean prefixes, negative flags, generic variables, and enum suffixes, listing all relevant spec files with why and how, creating the Antigravity skill, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, refactoring booleans, variables, and enums to adhere strictly to semantic standards, running the naming linter, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action.
- [ ] /learn `.lovable/coding-guidelines/coding-guidelines.md` and it is must and /goal apply the guidelines in coding every aspect.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, Boolean/Enum Refactoring, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-boolean-naming/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this naming prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-boolean-naming
   description: >-
     Autonomously audits, refactors, and validates repository-wide boolean prefixes, semantic variable naming, and enum Type suffixes using custom linters.
   ---
   ```
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents strict, microscopic instructions.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet: When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/XX-boolean-naming/01-task.md` and execute it"). The subagent MUST read the necessary files itself.

---

## 2. Phase 1: Write the Implementation Spec & Subtasks FIRST (Steps 1 to PHASE_1_STEPS)

Before doing anything else, you MUST write a highly detailed execution spec.

- **What to write:** Break down the parent task into a detailed architectural plan, complete naming violation inventory, code review guides, and embedded naming guidelines.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-boolean-naming-audit.md`. Do not hallucinate folders.
- **Create a Task-Specific Rule Set:** Before executing, analyze the specific task domain and explicitly write down 3-5 custom rules or constraints unique to this task inside the spec file.
- **Subtasks:** You MUST break the plan down and create detailed subtask files inside `.lovable/plans/subtasks/XX-boolean-naming/`.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/01-boolean-naming.md`** (or Section in `.lovable/coding-guidelines/coding-guidelines.md`)
  - **Why:** Establishes mandatory positive boolean prefixes (`is`, `has`, `can`, `should`, `was`, `will`, `did`, `must`).
  - **How:** Invert all negative flags (`isNotReady` -> `isReady`), eliminate `flag`/`check`/`bool` suffixes, and ban explicit `== true` comparisons.
- [ ] **`spec/02-coding-guidelines/02-anti-garbage-naming.md`**
  - **Why:** Eliminates vague, generic, and unmaintainable identifiers.
  - **How:** Scan and replace all instances of `temp`, `data`, `obj`, `item`, `comp_100`, `helper1` with domain-specific descriptive names.
- [ ] **`spec/02-coding-guidelines/03-enum-conventions.md`**
  - **Why:** Cross-language type safety and enum distinction.
  - **How:** Ensure all enum types end with the `Type` suffix (e.g. `UserRoleType`, `SeverityType`) and reside in dedicated definition files.
- [ ] **`spec/02-coding-guidelines/04-no-boolean-flags.md`**
  - **Why:** Prevents ambiguous boolean parameters in function signatures.
  - **How:** Split boolean flag parameters into distinct functions (e.g. `renderExpanded()` and `renderCollapsed()`).

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-boolean-naming.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If missing, create `linter-scripts/check-boolean-naming.py` that scans for:
  1. Booleans lacking allowed prefixes (`is`, `has`, `can`, `should`, etc.).
  2. Negative boolean identifiers (`isNot`, `disable`, `hasNo`).
  3. Generic variable names (`temp`, `data`, `obj`, `comp_100`).
  4. Enums lacking the `Type` suffix.
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-boolean-naming.py
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:booleans"] = ["python", "linter-scripts/check-boolean-naming.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running `python linter-scripts/check-boolean-naming.py`.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask, applying surgical refactoring until all boolean and naming checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-boolean-naming/
    2. Apply surgical refactoring (rename booleans, remove garbage names, add Type suffixes).
    3. Run the naming linter:
          python linter-scripts/check-boolean-naming.py
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

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `.lovable/ai-fix-scripts/index.md` before writing new scripts.
- [ ] **Native File Manipulator:** Natively use `python .lovable/ai-fix-scripts/01-file-manipulator.py <command>`.
- [ ] **Go Generate Sync:** Run `go generate ./...` in modified Go directories.
- [ ] **Index Documentation:** Keep `.lovable/ai-fix-scripts/index.md` updated with collapsible details tags.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Boolean naming prefixes strictly verified (`is`, `has`, `can`, `should`).
- [ ] Zero generic garbage variable names (`temp`, `data`, `obj`, `comp_100`).
- [ ] All enums end with `Type` suffix.
- [ ] Local CI runner exits 0.

---

## Metadata

- slug: cg-boolean-naming
- priority: high
- status: active
