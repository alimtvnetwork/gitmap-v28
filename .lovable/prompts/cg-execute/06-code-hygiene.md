# Instruction (must follow): Execute Coding Guidelines — Code Hygiene & File Caps

Trigger Keywords & Aliases: `cg-hygiene`, `cg-execute hygiene`, `audit hygiene`, `fix code hygiene`, `enforce file caps`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute repository code hygiene, file size caps, definition separation, and artifact protection compliance across the entire repository by decomposing violations into subtasks, verifying/creating hygiene linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all source files for file size caps (>300 lines), struct caps (>120 lines), inline definitions, committed artifacts, and versioning drift, listing all relevant spec files with why and how, creating the Antigravity skill, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, modularizing oversized files, extracting inline definitions into dedicated files, sanitizing gitignore, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, Modularization, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-code-hygiene/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this code hygiene prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-code-hygiene
   description: >-
     Autonomously audits, refactors, and validates repository file size caps (<=300 lines), struct caps (<=120 lines), dedicated definitions, and artifact protection.
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

- **What to write:** Break down the parent task into a detailed architectural plan, complete file cap violation inventory, code review guides, and embedded hygiene guidelines.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-code-hygiene-audit.md`.
- **Create a Task-Specific Rule Set:** Write down 3-5 custom rules or constraints unique to this task inside the spec file.
- **Subtasks:** Create detailed subtask files inside `.lovable/plans/subtasks/XX-code-hygiene/`.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/13-file-size-caps.md`** (or Section in `.lovable/coding-guidelines/coding-guidelines.md`)
  - **Why:** Prevents unwieldy monolithic source files.
  - **How:** Any file must NOT exceed 300 lines maximum; any class or struct must NOT exceed 120 lines maximum. Split oversized files into focused modules.
- [ ] **`spec/02-coding-guidelines/14-dedicated-definition-files.md`**
  - **Why:** Centralizes type, enum, and constant discovery.
  - **How:** Types, enums, constants, and interfaces must live in dedicated files (e.g. `types.go`, `constants.go`), never defined inline next to the first use.
- [ ] **`spec/02-coding-guidelines/15-artifact-protection.md`**
  - **Why:** Prevents repository bloat, security leaks, and build pollution.
  - **How:** Never commit generated code, compiled binaries (`.exe`, `.dll`, `.so`), cache files (`__pycache__`, `*.pyc`), or test reports. Ensure `.gitignore` ignores all build outputs.
- [ ] **`spec/02-coding-guidelines/16-version-source-of-truth.md`**
  - **Why:** Single point of version management across all languages.
  - **How:** Root `version.json` is the sole source of truth; never hardcode independent version strings across files.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-code-hygiene.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If missing, create `linter-scripts/check-code-hygiene.py` that scans for:
  1. Any source file exceeding 300 lines.
  2. Any struct/class file exceeding 120 lines.
  3. Inline enum or type definitions outside dedicated `types` files.
  4. Tracked binary or cache files matching `.gitignore` patterns.
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-code-hygiene.py
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:hygiene"] = ["python", "linter-scripts/check-code-hygiene.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running `python linter-scripts/check-code-hygiene.py`.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask until all code hygiene and file cap checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-code-hygiene/
    2. Apply surgical refactoring (split oversized files, extract type definitions, sanitize gitignore).
    3. Run the code hygiene linter:
          python linter-scripts/check-code-hygiene.py
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

- slug: cg-code-hygiene
- priority: medium
- status: active
