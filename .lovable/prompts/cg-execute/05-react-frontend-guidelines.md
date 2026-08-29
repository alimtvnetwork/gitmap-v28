# Instruction (must follow): Execute Coding Guidelines — React & Frontend Architecture

Trigger Keywords & Aliases: `cg-react`, `cg-execute react`, `audit react`, `fix frontend guidelines`, `enforce react rules`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute React frontend architecture compliance across the entire repository by decomposing violations into subtasks, verifying/creating component linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all React components, custom hooks, and state managers (`src/**/*.tsx`, `src/**/*.ts`), listing all relevant spec files with why and how, creating the Antigravity skill, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, modularizing oversized components (<= 100 lines), fixing `useEffect` guards, eliminating state mutations, removing tuple returns, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, React Refactoring, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-react-frontend/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this React prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-react-frontend
   description: >-
     Autonomously audits, refactors, and validates repository-wide React component line caps (<=100 lines), useEffect minimization, immutability, and named hook returns.
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

- **What to write:** Break down the parent task into a detailed architectural plan, complete component violation inventory, code review guides, and embedded React guidelines.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-react-guidelines-audit.md`.
- **Create a Task-Specific Rule Set:** Write down 3-5 custom rules or constraints unique to this task inside the spec file.
- **Subtasks:** Create detailed subtask files inside `.lovable/plans/subtasks/XX-react-guidelines/`.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/09-react-component-caps.md`** (or Section in `.lovable/coding-guidelines/coding-guidelines.md`)
  - **Why:** Prevents bloated, unmaintainable frontend components.
  - **How:** Every React component file (`.tsx`) MUST be $\le$ 100 lines. Extract child components, custom hooks, and helpers into dedicated files.
- [ ] **`spec/02-coding-guidelines/10-react-use-effect-rules.md`**
  - **Why:** Eliminates synchronization bugs and infinite render loops.
  - **How:** Minimize `useEffect` count (default 0); extract guards into positive named booleans (`isReadyToSync`); always return cleanup functions for acquired resources.
- [ ] **`spec/02-coding-guidelines/11-immutable-state-and-render.md`**
  - **Why:** Ensures reliable React render reconciliation and eliminates silent state loss.
  - **How:** Never mutate state in-place (`.push`, `.splice`, in-place object writes); use spread, `.map`, `.filter`, or `structuredClone`; render lists using expressions, no raw `for` loops in JSX.
- [ ] **`spec/02-coding-guidelines/12-no-public-tuples.md`**
  - **Why:** Self-documenting, explicit public shapes across hooks and component props.
  - **How:** Custom hooks must return named object types (e.g. `{ user, isLoading, error }`), NEVER raw tuples like `[User, boolean, Error]`. Component props must live in dedicated `types.ts`.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-react-components.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If missing, create `linter-scripts/check-react-components.py` that scans for:
  1. `.tsx` files exceeding 100 lines.
  2. In-place mutations on state variables.
  3. Custom hooks returning bare tuples.
  4. Raw `for` loops inside React render methods.
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-react-components.py
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:react"] = ["python", "linter-scripts/check-react-components.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running `python linter-scripts/check-react-components.py`.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask until all React frontend checks and build gates pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-react-guidelines/
    2. Apply surgical refactoring (split components <= 100 lines, extract hooks, enforce immutability).
    3. Run the React linter:
          python linter-scripts/check-react-components.py
    4. Run the frontend build:
          npm run build (or vite build)
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

- slug: cg-react-frontend
- priority: medium
- status: active
