# Coding Guideline Execution Prompts Registry (`cg-execute`)

> **Location:** `.lovable/prompts/cg-execute/`
> **Purpose:** Standardized, sequenced N-step autonomous prompts for auditing, refactoring, creating linters, and enforcing every section of the coding guidelines across any repository.

---

## 1. Prompt Registry & Execution Sequence

| Sequence | File | Title | Trigger Keywords | Focus Area | Linter Script |
|:---:|---|---|---|---|---|
| **01** | [`01-index.md`](./01-index.md) | Prompts Catalog & Registry | — | Central navigation, lifecycle map, and execution documentation | — |
| **02** | [`02-error-management.md`](./02-error-management.md) | Error Management & Architecture | `cg-error`, `cg-execute error`, `audit error` | `AppError` wrappers, centralized `HandleError` dispatcher, buffer flushing, universal response envelopes, zero bare panics/exits | `linter-scripts/check-error-management.py` |
| **03** | [`03-boolean-and-naming.md`](./03-boolean-and-naming.md) | Booleans, Naming & Enums | `cg-boolean`, `cg-execute boolean`, `audit boolean` | Positive boolean prefixes (`is`, `has`, `can`, `should`), semantic variable naming, `*Type` enum suffix, no tuple returns | `linter-scripts/check-boolean-naming.py` |
| **04** | [`04-data-and-schema.md`](./04-data-and-schema.md) | Database & Schema Rules | `cg-schema`, `cg-execute schema`, `audit schema` | PascalCase tables, camelCase columns, `{TableName}Id` keys, Mermaid ERDs, SQLite conventions, enum join tables | `linter-scripts/check-database-schema.py` |
| **05** | [`05-react-frontend-guidelines.md`](./05-react-frontend-guidelines.md) | React & Frontend Architecture | `cg-react`, `cg-execute react`, `audit react` | Component size ($\le$ 100 lines), `useEffect` minimization, immutability, hook named objects (no tuples), `types.ts` separation | `linter-scripts/check-react-components.py` |
| **06** | [`06-code-hygiene.md`](./06-code-hygiene.md) | Code Hygiene & Repository Cleanliness | `cg-hygiene`, `cg-execute hygiene`, `audit hygiene` | File caps ($\le$ 300 lines), dedicated definition files, artifact protection (`.gitignore`), Single Source of Truth versioning | `linter-scripts/check-code-hygiene.py` |
| **07** | [`07-style-guidelines.md`](./07-style-guidelines.md) | Style, Formatting & Line-Gaps | `cg-style`, `cg-execute style`, `audit style` | Function length ($\le$ 15 lines), nested `if` flattening, line-gap rules (single blank line before returns/after closing braces), MD022/MD032 markdown spacing | `linter-scripts/check-style-guidelines.py` |

---

## 2. Standard N-Step Execution Lifecycle

Every prompt in this suite executes using the standardized N-step autonomous loop:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        Phase 1: Steps 1 .. N/2                         │
│                 (Planning, Audit & Task Decomposition)                 │
│                                                                        │
│  1. Scan codebase for section-specific violations.                     │
│  2. Ingest relevant spec files with Why and How verification.          │
│  3. Bootstrap Antigravity skill in .agents/skills/cg-.../skill.md      │
│  4. Write master plan: .lovable/plans/pending/XX-...-audit.md          │
│  5. Decompose into: .lovable/plans/subtasks/XX-.../                    │
│  6. Verify/Create section linter in linter-scripts/ and wire to CI/CD. │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                        Phase 2: Steps N/2+1 .. N                       │
│                  (Autonomous Execution & Verification)                 │
│                                                                        │
│  1. Sequentially execute subtasks with strict coding guidelines.       │
│  2. Run section linter: python linter-scripts/check-*.py               │
│  3. Run guideline autofixer: python .lovable/ai-fix-scripts/02-*.py    │
│  4. Run CI local runner: python .lovable/ai-fix-scripts/03-*.py        │
│  5. Ensure 100% clean CI pass (exit code 0).                           │
│  6. Move completed plans to .lovable/plans/completed/                  │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Mandatory Linter & CI/CD Connection Protocol

If an automated linter for a specific guideline category is missing from the repository:
1. The executing agent is **mandated** to create the linter script in `linter-scripts/check-<category>.py`.
2. The agent must test it locally via `python linter-scripts/check-<category>.py`.
3. The agent must connect it to the CI/CD local runner in `.lovable/ai-fix-scripts/03-cicd-local-runner.py` under the `JOBS` dictionary.
4. The agent must connect it to the GitHub Actions workflow in `.github/workflows/ci.yml`.
