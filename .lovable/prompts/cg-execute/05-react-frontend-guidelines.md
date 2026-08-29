# Instruction (must follow): Execute Coding Guidelines — React & Frontend Architecture

Trigger Keywords & Aliases: `cg-react`, `cg-execute react`, `audit react`, `fix frontend guidelines`, `enforce react standards`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute React and frontend architecture compliance across the entire repository by decomposing violations into subtasks, verifying/creating frontend linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all `.tsx`, `.jsx`, and `.ts` frontend files, listing all React & frontend spec files with why and how, creating the Antigravity skill, generating component hierarchy diagrams, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, decomposing monolithic components ($\le$ 100 lines), eliminating redundant `useEffect` hooks, converting custom hook tuple returns to named objects, running the frontend linter, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/02-typescript/`, `spec/07-design-system/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action and also create agent rules in the repo if required to or missing from rules set of agent memory.
- [ ] /learn `.lovable/coding-guidelines/coding-guidelines.md` and it is must and /goal apply the guidelines in coding every aspect.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, Component Refactoring, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-react-frontend/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-react-frontend
   description: >-
     Autonomously audits, refactors, and validates repository-wide React and frontend code against spec/02-coding-guidelines/02-typescript/ using 100-line caps, named hook objects, and CI linters.
   ---
   ```
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## 1. Ruthless Orchestration & Insult Protocol

/goal You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents strict, microscopic instructions.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet: When spawning a subagent, DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction (e.g., "Read subtask file `.lovable/plans/subtasks/XX-react-frontend/01-task.md` and execute it"). The subagent MUST read the necessary files itself.

---

## 2. Phase 1: Write the Implementation Spec & Subtasks FIRST (Steps 1 to PHASE_1_STEPS)

Before doing anything else, you MUST write a highly detailed execution spec.

- **What to write:** Break down the parent task into a detailed architectural plan, complete React component line count audit, custom hook return signature catalog, Mermaid component tree diagram, code review guides, and embedded frontend standards.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-react-frontend-audit.md`. Do not hallucinate folders.
- **Create a Task-Specific Rule Set:** Before executing, analyze the specific task domain and explicitly write down 3-5 custom rules or constraints unique to this task inside the spec file. This prevents domain-specific regressions and forces sub-agents to follow exact architectures.
- **Subtasks:** You MUST break the plan down and create detailed subtask files inside `.lovable/plans/subtasks/XX-react-frontend/`. Every subtask file must contain actionable, microscopic instructions.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/02-typescript/08-typescript-standards-reference.md`**
  - **Why:** Authoritative TypeScript standards.
  - **How:** Strict typing, zero `any`, explicit function return types, and discriminated unions for polymorphic data.
- [ ] **`spec/02-coding-guidelines/02-typescript/14-state-management.md`**
  - **Why:** Immutable state & hook conventions.
  - **How:** Custom React hooks MUST NOT return tuples (`[state, setState]` is banned). Custom hooks MUST return named property objects (`{ userProfile, isLoading, onUpdate }`).
- [ ] **`spec/02-coding-guidelines/02-typescript/09-promise-await-patterns.md`**
  - **Why:** Async lifecycle safety in React.
  - **How:** Handle all async flows with `try/catch` and `AppError` wrappers; avoid unhandled promises in `useEffect`.
- [ ] **`spec/07-design-system/01-design-principles.md`**
  - **Why:** Modular UI architecture & sizing caps.
  - **How:** Hard 100-line cap per React component file. Decompose large UI blocks into single-responsibility child components.
- [ ] **`spec/07-design-system/02-theme-variable-architecture.md`**
  - **Why:** Theme token usage.
  - **How:** Zero hardcoded hex color codes in components. Consume CSS theme custom properties and semantic tokens.
- [ ] **`spec/03-error-manage/02-error-architecture/04-error-modal/02-react-components/02-error-store.md`**
  - **Why:** Centralized UI error modal.
  - **How:** Route all unhandled UI exceptions and API error responses to the global error store; display via standard error modal.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-frontend-guidelines.mjs` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If no dedicated frontend linter exists, create `linter-scripts/check-frontend-guidelines.mjs` that AST-scans `.tsx` and `.jsx` files for:
  1. Component file line counts exceeding 100 lines.
  2. Custom hooks returning array literals (tuples) instead of object literals.
  3. Redundant `useEffect` calls syncing derived state.
  4. Direct state mutations (`state.items.push()` or `state.prop = x`).
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  node linter-scripts/check-frontend-guidelines.mjs
  # And verify standard ESLint:
  npm run lint
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:frontend"] = ["node", "linter-scripts/check-frontend-guidelines.mjs"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running the frontend linter.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask, decomposing components and refactoring hooks until all frontend checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-react-frontend/
    2. Extract sub-components into dedicated files (ensuring each is <= 100 lines).
    3. Convert hook returns to named property objects and update caller sites.
    4. Run the frontend linter:
          node linter-scripts/check-frontend-guidelines.mjs
    5. Run the local CI runner:
          python .lovable/ai-fix-scripts/03-cicd-local-runner.py
    6. IF any check fails:
          - Diagnose failure, fix component structure, and re-test immediately.
       IF all checks pass (exit code 0):
          - Mark subtask completed and proceed to next subtask.

    7. When all subtasks are finished and local CI is 100% green:
          - BREAK and proceed to End of Tunnel.
```

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `.lovable/ai-fix-scripts/index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python .lovable/ai-fix-scripts/01-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `.lovable/ai-fix-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `.lovable/ai-fix-scripts/index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.lovable/temp/active-locks.json`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `plans/index.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] All `.tsx`/`.jsx` component files are $\le$ 100 lines.
- [ ] Custom hooks return named objects `{ data, isLoading }`, not tuples.
- [ ] No redundant `useEffect` for derived computations.
- [ ] `node linter-scripts/check-frontend-guidelines.mjs` exited with code 0.
- [ ] Local CI runner `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Master Guidelines: I have fully read and strictly enforced every file in `spec/02-coding-guidelines/02-typescript/` and `.lovable/coding-guidelines/coding-guidelines.md`.
- [ ] Component Sizing: All component files are $\le$ 100 lines.
- [ ] Hook Return Signatures: Named property objects only, zero tuple returns.
- [ ] Semantic Naming: Absolutely NO generic garbage names (`temp`, `data`, `obj`, `comp_100`).
- [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)
- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] Echo Back the Spec: I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] Pre-Commit Diff Proof: I have executed `git status` or `git diff --stat` and verified that the files I claim to have modified are actually listed as modified in the terminal output before committing.
- [ ] No Placeholder Search: I ran a regex search for `TODO` and `\[.*\]` in my modified files and confirmed I left zero placeholders behind. I actually wrote the implementation.
- [ ] Index Sync Deadman Switch: I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.lovable/what-to-read.md`. I did not leave any orphaned files.
- [ ] Blast Radius Acknowledgment: Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.

---

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memories/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.** 
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion. 
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

---

## Metadata

- slug: cg-react-frontend
- priority: medium
- status: active
