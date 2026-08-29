# Coding Guideline Execution Prompts Registry (`cg-execute`)

> **Location:** `.lovable/prompts/cg-execute/`  
> **Purpose:** Standardized, sequenced N-step autonomous prompts for auditing, refactoring, and enforcing every section of the coding guidelines across any codebase.

---

## 1. Prompt Registry

| Sequence | File | Title | Trigger Keywords | Focus Area |
|:---:|---|---|---|---|
| **01** | [`01-style-guidelines.md`](./01-style-guidelines.md) | Style, Formatting & Line-Gaps | `cg-style`, `cg-execute style` | Function length (<=15 lines), nested `if` flattening, line-gap rules, MD022/MD032 markdown spacing |
| **02** | [`02-error-management.md`](./02-error-management.md) | Error Management & Architecture | `cg-error`, `cg-execute error` | `AppError` wrappers, centralized `HandleError` dispatcher, buffer flushing, zero bare panics/exits |
| **03** | [`03-boolean-and-naming.md`](./03-boolean-and-naming.md) | Booleans, Naming & Enums | `cg-boolean`, `cg-execute boolean` | Positive boolean prefixes (`is`, `has`, `can`, `should`), semantic variable naming, `*Type` enum suffix |
| **04** | [`04-data-and-schema.md`](./04-data-and-schema.md) | Database & Schema Rules | `cg-schema`, `cg-execute schema` | PascalCase tables, camelCase columns, `{TableName}Id` keys, Mermaid ERDs, SQLite conventions |
| **05** | [`05-react-frontend-guidelines.md`](./05-react-frontend-guidelines.md) | React & Frontend Architecture | `cg-react`, `cg-execute react` | Component size (<=100 lines), `useEffect` minimization, immutability, hook named objects (no tuples) |
| **06** | [`06-code-hygiene-and-ci-linter.md`](./06-code-hygiene-and-ci-linter.md) | Code Hygiene & Master CI Linters | `cg-hygiene`, `cg-execute hygiene` | File caps (<=300 lines), dedicated definition files, artifact protection, end-to-end CI/CD linter suites |

---

## 2. Standard N-Step Execution Lifecycle

Every prompt in this suite executes using the standardized N-step autonomous loop:

```
┌────────────────────────────────────────────────────────┐
│                   Phase 1: Steps 1 .. N/2              │
│               (Planning, Audit & Task Decomposition)   │
│                                                        │
│  1. Scan codebase for section-specific violations.     │
│  2. Write master plan: .lovable/plans/pending/XX-...md │
│  3. Decompose into: .lovable/plans/subtasks/XX-.../    │
│  4. Verify/Create section linter in scripts/ & CI/CD.  │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│                  Phase 2: Steps N/2+1 .. N             │
│            (Autonomous Execution & Verification)       │
│                                                        │
│  1. Sequentially execute subtasks with strict rules.   │
│  2. Run automated test suites & section linters.       │
│  3. Ensure 100% clean CI/CD local runner pass (exit 0).│
│  4. Move completed plans to .lovable/plans/completed/  │
└────────────────────────────────────────────────────────┘
```

---

## 3. Mandatory CI/CD Linter Hook

If an automated linter check for a specific guideline category is missing from the repository, the executing agent is **mandated** to create the linter script (e.g. in `scripts/` or `linters/`) and connect it directly to the CI/CD local runner and pipeline before completing the run.
