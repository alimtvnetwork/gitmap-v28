# Instruction (must follow): Execute Coding Guidelines — Database & Schema Rules

Trigger Keywords & Aliases: `cg-schema`, `cg-execute schema`, `audit schema`, `fix database guidelines`, `enforce schema rules`

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute database schema, entity models, and JSON serialization compliance across the entire repository by decomposing violations into subtasks, verifying/creating schema linters, and running a continuous N-step self-loop until 100% green without a single failure.

- [ ] /goal First N/2 steps will be given for spec writing for AI as given, deep codebase scanning across all SQL schemas, migrations, struct models, and JSON encoders, listing all relevant spec files with why and how, creating the Antigravity skill, and breaking down into microscopic subtasks for N/2 steps.
- [ ] /goal Second N/2 steps will be given to execute the created subtasks, refactoring tables, columns, primary keys, and JSON tags to adhere to schema guidelines, running the schema linter, and verifying all local CI gates exit with code 0.
- [ ] /learn Ingest `.lovable/memory/00-index.md`, `.lovable/strictly-avoid.md`, `spec/02-coding-guidelines/`, and `.lovable/coding-guidelines/coding-guidelines.md` before taking action.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec in .lovable/plans/pending/, Subtasks in .lovable/plans/subtasks/, Skill Creation, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Autonomous Execution, Schema Refactoring, Linter Verification, Local CI Runner Verification)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, you must check if this prompt is already installed as a native Antigravity Skill.

1. Check if `.agents/skills/cg-data-schema/skill.md` exists in the workspace. If it does NOT exist, you MUST create it now.
2. Extract the core instructions of this schema prompt and save it into that `skill.md` using standard YAML frontmatter:
   ```yaml
   ---
   name: cg-data-schema
   description: >-
     Autonomously audits, refactors, and validates repository-wide database schemas, PascalCase tables, camelCase columns, {TableName}Id primary keys, and Mermaid ERDs.
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

- **What to write:** Break down the parent task into a detailed architectural plan, complete schema violation inventory, target Mermaid ERDs, and embedded schema guidelines.
- **Where to save it:** Save this master plan into `.lovable/plans/pending/XX-schema-guidelines-audit.md`.
- **Create a Task-Specific Rule Set:** Write down 3-5 custom rules or constraints unique to this task inside the spec file.
- **Subtasks:** Create detailed subtask files inside `.lovable/plans/subtasks/XX-schema-guidelines/`.

---

## 3. Authoritative Spec Files Checklist (Non-Negotiable Action Items)

You MUST read, follow, and mechanically verify every single specification file below before and during execution:

- [ ] **`spec/02-coding-guidelines/05-data-schema-rules.md`** (or Section in `.lovable/coding-guidelines/coding-guidelines.md`)
  - **Why:** Establishes table, column, and primary key naming conventions.
  - **How:** Tables and entities MUST be `PascalCase`; columns and struct fields MUST be `camelCase`; JSON keys in API models MUST be `PascalCase`.
- [ ] **`spec/02-coding-guidelines/06-primary-key-rules.md`**
  - **Why:** Ensures uniform, deterministic primary key contracts across SQLite models.
  - **How:** All primary keys MUST be named `{TableName}Id` (integer auto-increment). Raw `id` or UUID primary keys are prohibited.
- [ ] **`spec/02-coding-guidelines/07-text-columns-and-nulls.md`**
  - **Why:** Standardizes nullable documentation fields.
  - **How:** Reference tables use `Description TEXT NULL`; transactional tables use `Notes TEXT NULL` and `Comments TEXT NULL`. All nullable without empty string defaults.
- [ ] **`spec/02-coding-guidelines/08-mermaid-erd-requirement.md`**
  - **Why:** Architectural visibility into entity relationships.
  - **How:** Every database-altering plan and PR MUST include a valid, complete Mermaid ERD diagram.

---

## 4. Mandatory Linter & CI/CD Connection Checklist

Code standards must be mechanically enforced by automated linters. You MUST verify or create the linter and connect it to CI:

- [ ] **Linter Script Identification:** Check if `linter-scripts/check-database-schema.py` exists in the repository.
- [ ] **Auto-Create Linter if Missing:** If missing, create `linter-scripts/check-database-schema.py` that scans for:
  1. Non-PascalCase table names.
  2. Non-camelCase column definitions.
  3. Primary keys not matching `{TableName}Id`.
  4. Missing Mermaid ERDs in schema modification plans.
- [ ] **Local Linter Command:** Execute and verify the linter locally:
  ```bash
  python linter-scripts/check-database-schema.py
  ```
- [ ] **CI/CD Local Runner Connection:** Register the linter script inside `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary:
  ```python
  JOBS["lint:schema"] = ["python", "linter-scripts/check-database-schema.py"]
  ```
- [ ] **GitHub Actions Workflow Connection:** Verify that `.github/workflows/ci.yml` contains a dedicated step running `python linter-scripts/check-database-schema.py`.

---

## 5. Phase 2: Autonomous Subtask Execution Loop (Steps PHASE_1_STEPS+1 to N)

> [!IMPORTANT]
> **AUTONOMOUS EXECUTION MANDATE — DO NOT STOP.**
> Sequentially execute each subtask until all database schema checks pass 100% green.

```text
STEP = 0
WHILE (STEP < PHASE_2_STEPS):
    STEP += 1

    1. Read the next subtask from .lovable/plans/subtasks/XX-schema-guidelines/
    2. Apply surgical refactoring (update SQL schemas, model structs, primary keys, JSON tags).
    3. Run the schema linter:
          python linter-scripts/check-database-schema.py
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

- slug: cg-data-schema
- priority: medium
- status: active
