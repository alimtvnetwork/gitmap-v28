# Instruction (must follow): Coding Guideline Execution — Data, Database & Schema Rules

Trigger Keywords & Aliases: `cg-schema`, `cg-execute schema`, `fix database schema`, `execute data guidelines`

```text
N = 100
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to reviewing all database schemas, models, SQL migrations, and JSON structures, writing the master execution plan to `.lovable/plans/pending/XX-schema-guidelines-audit.md`, and decomposing it into subtasks in `.lovable/plans/subtasks/XX-schema-guidelines/`.
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to executing those subtasks in an autonomous self-loop until all database tables, columns, and serialization models strictly adhere to the schema guidelines.
- [ ] /goal **Linter Mandate**: If an automated schema validation or migration integrity checker does not exist in CI/CD, you MUST create a database schema linter and connect it directly to the CI/CD local runner and workflows.
- [ ] /learn Ingest `.lovable/coding-guidelines/coding-guidelines.md`, `.lovable/strictly-avoid.md`, and `.lovable/memory/00-index.md` before touching any code.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## 1. Data & Schema Non-Negotiable Checklist

You MUST audit and strictly enforce every rule below across the entire codebase:

### A. Naming & Case Conventions
1. **Entities & Tables**: Strictly `PascalCase` (e.g. `UserAccount`, `TransactionRecord`, `ProjectProfile`).
2. **Fields & Columns**: Strictly `camelCase` (e.g. `createdAt`, `userRole`, `isArchived`).
3. **JSON Keys**: Strictly `PascalCase` in API and serialization envelopes (e.g. `{"UserId": 1, "CreatedAt": "..."}`).
4. **Primary Keys**: Integer auto-increment, named `{TableName}Id` (e.g. `UserAccountId`, `TransactionRecordId`). No bare `id` or raw UUID primary keys.

### B. Relational Integrity & Columns
1. **Type / Category / Status Columns**: Never use free-form string columns for statuses. Use a 1-N or N-M relation with a registered enum or join table.
2. **Standard Text Columns**:
   - Reference/Entity tables: `Description TEXT NULL`.
   - Transactional tables: `Notes TEXT NULL` and `Comments TEXT NULL`.
   - All standard text columns must be nullable without empty string defaults.
3. **Foreign Keys & Indices**: Explicit foreign key constraints and indexed lookup columns. SQLite is the default database engine.

### C. Visual Documentation (Mermaid ERD)
1. **Mandatory Mermaid ERD**: Any pull request, plan, or commit that creates or modifies database tables MUST include a complete Mermaid ERD diagram documenting tables, keys, and relationships.

---

## 2. Phase 1: Planning, Audit & Subtask Decomposition (Steps 1 .. N/2)

1. **Schema Audit**: Inspect all migration scripts, SQLite table creations, model struct definitions, and JSON tags.
2. **Master Plan**: Write a detailed execution plan to `.lovable/plans/pending/XX-schema-guidelines-audit.md`. Include the current and target Mermaid ERDs.
3. **Subtask Files**: Decompose into subtask files in `.lovable/plans/subtasks/XX-schema-guidelines/` (e.g. `01-task.md`, `02-task.md`, ...).
4. **Linter Connection**: Implement a schema linter (e.g. `python scripts/lint-schema.py`) to verify column casing, primary key formats, and foreign keys in CI/CD.

---

## 3. Phase 2: Autonomous Execution Loop (Steps N/2+1 .. N)

1. **Apply Migrations**: Update table definitions, struct fields, and SQL queries to use `{TableName}Id`, `camelCase` columns, and `PascalCase` entities.
2. **Validate JSON Serializers**: Ensure JSON tags match PascalCase contract requirements and contract tests pass.
3. **Execute Suite**: Run database integration tests and the CI local runner.
4. **Update Status**: Mark completed tasks as `DONE`, move completed plans to `.lovable/plans/completed/`, and update `.lovable/plans/index.md`.

---

## 4. Pre-Commit Verification Checklist

- [ ] All table names are PascalCase and all column names are camelCase.
- [ ] Primary keys are `{TableName}Id` integers.
- [ ] JSON keys in serializers use PascalCase.
- [ ] Mermaid ERD is documented in the plan/spec.
- [ ] Schema linter is integrated into CI/CD and passes with exit 0.
