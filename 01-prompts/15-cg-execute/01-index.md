# Coding Guideline Execution Suite (`cg-execute`) — Index & Catalog (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously orchestrate and execute coding standard compliance across the repository by directly scanning files, generating verifiable audit specs, actively refactoring source code, and verifying with automated linters until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase to inventory all architectural violations and anti-patterns.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.lovable/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.lovable/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/01-index.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `spec/02-coding-guidelines/` for domain-specific architectural specifications.
15. [ ] /learn Ingest `spec/03-error-manage/` for error handling architectures and AppError.
16. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan, Spec with Violation Ledger in .lovable/plans/pending/, Subtasks, Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Active Code Refactoring, Disk Verification, Linter Verification, Local CI Runner Verification, Plan Completion)
```

---

## Prompts Catalog & Execution Order

Prompts are sequenced according to priority. Error management, control-flow flattening, and core semantics precede presentation styling:

| Sequence | File | Title | Trigger Keywords | Focus Area | Linter Hook (`linter-scripts/`) |
|:---:|---|---|---|---|---|
| **01** | [`01-index.md`](./01-index.md) | Prompts Catalog & Registry | `cg-index`, `cg-execute index` | Master overview, lifecycle architecture, prompt routing | — |
| **02** | [`02-error-management.md`](./02-error-management.md) | Error Management & Architecture | `cg-error`, `cg-execute error`, `audit error` | `AppError` context wrappers, zero dual-handling (leafs return `error`, callers handle exits), typed exit enums, specialized exit helpers, universal envelopes | `check-error-management.py`, `check-mws-error-codes.py` |
| **03** | [`03-nested-if-and-guard-clauses.md`](./03-nested-if-and-guard-clauses.md) | Nested `if` Elimination & Guard Clauses | `cg-nested-if`, `cg-execute nested-if`, `audit nested if` | Zero tolerance for nested if statements (depth > 1), condition inversion, early guard returns, validation helper extraction | `validate-guidelines.py`, `05-guideline-autofixer.py` |
| **04** | [`04-booleans-and-complex-conditions.md`](./04-booleans-and-complex-conditions.md) | Booleans, Negatives & Complex Conditions | `cg-boolean`, `cg-execute boolean`, `audit boolean` | Affirmative prefixes (`is`, `has`), zero `== true`, positive framing (ban `!isSuccess` / `isNot*`), elimination of mixed polarity (`if a && !b`) | `check-boolean-guidelines.py`, `05-guideline-autofixer.py` |
| **05** | [`05-naming-conventions-and-boolean-prefixes.md`](./05-naming-conventions-and-boolean-prefixes.md) | Variable & Boolean Naming Conventions, Anti-`ok` & Inverted Guards | `cg-naming`, `cg-execute naming`, `fix ok boolean`, `fix boolean naming` | Strict `is`/`has` prefixes, total ban on bare `ok` in type assertions / map lookups, positive boolean declarations with inverted `if` guard returns, lowercase filenames | `check-enum-and-boolean.mjs`, `08-naming-autofixer.py` |
| **06** | [`06-constants-and-enums.md`](./06-constants-and-enums.md) | Constants, Enums & Magic Literal Elimination | `cg-enums`, `cg-constants`, `cg-execute enums` | Mandatory `*Type` enum suffix, zero magic string/integer literals, centralized constants packages, typed enum function arguments | `check-enum-guidelines.py`, `05-guideline-autofixer.py` |
| **07** | [`07-data-and-schema.md`](./07-data-and-schema.md) | Database & Data Schema Rules | `cg-schema`, `cg-execute schema`, `audit schema` | PascalCase tables, camelCase columns, `{TableName}Id` primary/foreign keys, Mermaid ERDs, SQLite/ORM explicit relations, zero nested if statements | `validate-guidelines.py` |
| **08** | [`08-react-frontend-guidelines.md`](./08-react-frontend-guidelines.md) | React & Frontend Architecture | `cg-react`, `cg-execute react`, `audit react` | Component size caps (<= 100 lines, recommended <= 80), functions (<= 8–15 lines), zero nested if statements, `useEffect` minimization, immutable state, named hook objects | `check-enum-and-boolean.mjs` |
| **09** | [`09-code-hygiene.md`](./09-code-hygiene.md) | Code Hygiene, File Caps & Parameter Reduction | `cg-hygiene`, `cg-execute hygiene`, `audit hygiene` | File caps (<= 100 lines coding max, recommended <= 80), functions (<= 8–15 lines), parameter reduction via specialized helpers, struct caps (<= 120 lines) | `check-file-sizes.py`, `13-file-size-guard.py` |
| **10** | [`10-style-guidelines.md`](./10-style-guidelines.md) | Coding Style, Formatting & Line-Gaps | `cg-style`, `cg-execute style`, `audit style` | Return New Line concept (R13-R16: blank line before `return`/`throw`, blank line after `}`), zero nested if statements, <= 8–15 line function caps, MD022/MD032 spacing | `check-function-lengths.py`, `04-newline-fixer.py` |
| **11** | [`11-testing-and-coverage.md`](./11-testing-and-coverage.md) | Integration, E2E & Branch Test Coverage | `cg-test`, `cg-execute test`, `audit tests`, `write e2e tests` | Three-part semantic naming (`TestUnit_Scenario_Outcome`), <= 8-line function decomposition, zero nested if statements, positive/negative/error branch coverage | `validate-guidelines.py`, `06-cicd-local-runner.py` |
| **12** | [`12-relative-paths.md`](./12-relative-paths.md) | Strict Relative Git Paths & Absolute Path Elimination | `cg-relative-paths`, `cg-execute relative-paths`, `fix absolute paths`, `fix file paths` | Total ban on absolute paths (`/absolute/path/to/...`, `/home/...`) and `file:///` URIs, automatic conversion to relative Git paths | `check-relative-paths.py`, `07-relative-path-fixer.py` |
| **13** | [`13-cli-commands-and-help.md`](./13-cli-commands-and-help.md) | CLI Commands, Help Text Parity & Help UI Architecture | `cg-cli`, `cg-help`, `cg-execute cli`, `audit cli commands` | 100% command/subcommand registration in `--help`, flag descriptions, terminal usage examples, structured Help UI layout | `09-cli-help-auditor.py` |
| **14** | [`14-function-signatures-and-return-types.md`](./14-function-signatures-and-return-types.md) | Function Naming, Single Return Types & Result Envelopes | `cg-functions`, `cg-signatures`, `cg-return-types`, `cg-execute functions` | Semantic verb/predicate prefixes, single `Result[T]` return envelope, universal `*AppError` wrapping, zero generic `(T, error)` | `check-function-lengths.py`, `check-error-management.py` |
| **15** | [`15-typescript-guidelines-and-types.md`](./15-typescript-guidelines-and-types.md) | TypeScript Strict Typing, Discriminated Unions & Architecture | `cg-typescript`, `cg-ts`, `cg-execute ts`, `audit typescript` | Total ban on `any`, Discriminated Unions for state, strongly-typed `Result<T>` envelopes, `as const` object enums, exhaustive `assertNever` matching | `tsc --noEmit`, `check-enum-and-boolean.mjs` |
| **16** | [`16-multi-language-enums-and-traits.md`](./16-multi-language-enums-and-traits.md) | Multi-Language Enums, Traits & Pattern Matching | `cg-enums-traits`, `cg-enums`, `cg-execute enums`, `audit enums` | PHP 8.1+ Backed Enums + `HasEnumHelpers` trait, Rust ADT Enums + exhaustive `match`, Go custom enums + stringers, `*Type` suffixes | `check-enum-guidelines.py`, `05-guideline-autofixer.py` |
| **17** | [`17-terminal-ui-and-cli-styling.md`](./17-terminal-ui-and-cli-styling.md) | Terminal UI, CLI Styling, Lipgloss & Animation Architecture | `cg-terminal-ui`, `cg-cli-style`, `cg-lipgloss`, `audit terminal ui` | Bright bold 9X ANSI palette, Catppuccin pastel cycling, 2-column width caps (max width 26), intent banners, clone spinners, version footers | `09-cli-help-auditor.py` |
| **18** | [`18-function-argument-reduction-and-params.md`](./18-function-argument-reduction-and-params.md) | Argument Reduction, Parameter Structs & Return Architecture | `cg-argument-reduction`, `cg-params`, `cg-struct-params`, `cg-execute params` | Parameter structs (`*Params`) for >2–3 args, value-based passing, affirmative boolean fields (`is`/`has`), mandatory `*apperror.AppError` returns (zero void in Go), framework error conversion | `check-function-lengths.py`, `check-error-management.py` |

---

## Canonical Sizing Tier & Formatting Rules

- **Standard File Size:** Max 100 coding lines per file (recommended <= 80 lines).
- **Functions:** Target <= 8 lines of body logic; hard cap of <= 15 lines.
- **React Components:** Recommended <= 80 lines; standard max <= 100 lines.
- **Nested `if` Statements:** Zero tolerance (must be flattened with guard clauses).
- **NO Line-Compression Cheating:** Never collapse `if/else` onto a single line or delete blank lines to fit under line caps. Reduce size by decomposing into separate files.

---

## Anti-Hallucination & Verifiable Execution Protocol (Zero-Drift Mandate)

> [!CAUTION]
> **TOTAL BAN ON GHOST DIFFS, CODE TRUNCATION, AND FABRICATED PASSES:**
> AI agents frequently hallucinate progress by outputting chat diffs without modifying disk, omitting violations from memory, or inserting lazy stubs (`// ... rest of code ...`). Every step must be mechanically grounded in disk reality.

### 1. The Exhaustive Violation Ledger (Phase 1 Mandate)

- Phase 1 audit plans (`.lovable/plans/pending/XX-*-audit.md`) MUST include a markdown table tracking every single violation:
  `| Violation Id | File Path | Line Number | Exact Snippet | Planned Fix | Status (PENDING/DONE) |`
- Never group files into vague summaries like "various files in pkg/". Every file and line must be explicitly numbered.

### 2. The Disk Reality Verification Gate (Phase 2 Mandate)

- Before claiming a file was refactored, the agent MUST run `git status --porcelain` or `git diff --stat`.
- If a file claimed as modified does NOT show as modified on disk in git, the turn is an immediate hallucination failure.

### 3. Absolute Ban on Code Truncation & Stubs

- NEVER write placeholder comments like `// ... existing code ...`, `/* unchanged */`, `TODO`, `FIXME`, or `[N]`.
- All file edits MUST be 100% complete and functionally working.

### 4. Grounded Command Veracity

- NEVER claim a linter or test suite passed without executing the actual command via tool calls and capturing `exit code 0`. Fabricating test outputs results in an immediate strike.

### 5. Spec Citation & Relative Git Path Requirement (TOTAL BAN on Absolute Paths / `file:///` URIs)

- Every refactoring action, plan file (`.lovable/plans/pending/`), subtask (`.lovable/plans/subtasks/`), and memory log must cite the authoritative specification path from `spec/` or `.lovable/` using **strictly relative paths from the git root**.
- **TOTAL BAN:** NEVER write absolute filesystem paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///absolute/path/to/...`, `file:///absolute/path/to/...`) into any committed or created files.
  - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
  - ❌ **BAD:** `Target: file:///absolute/path/to/gitmap/cmd/login.go`
  - ✅ **GOOD:** `[SSH Commands](spec/13-generic-cli/01-index.md)`
  - ✅ **GOOD:** `Target: gitmap/cmd/login.go`

---

## Standardized N-Step Self-Loop Architecture

Every prompt in this suite operates using a strict two-phase loop budget:

### Phase 1: Scan, Spec & Subtasks (Steps 1 to N/2)

1. **Memory Ingestion:** Ingest `.lovable/coding-guidelines.md`, `.lovable/strictly-avoid.md`, and recent issues in `.lovable/memory/issues/`.
2. **High-Speed Violation Scan:** Run `python 03-ai-scripts/11-fast-file-scanner.py --check` and `python 03-ai-scripts/12-fast-cached-grep.py "<pattern>"` to detect AST violations across the codebase in milliseconds.
3. **Master Spec Creation:** Write `.lovable/plans/pending/xx-<slug>-audit.md` capturing the full violation ledger, affected files, line numbers, and acceptance criteria.
4. **Subtask Decomposition:** Break down the master plan into granular subtasks in `.lovable/plans/subtasks/xx-<slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
5. **Linter Hook Verification:** Check if the automated linter script exists in `linter-scripts/`. If missing, generate the linter script and connect it to `03-ai-scripts/06-cicd-local-runner.py` and CI/CD pipelines.

### Phase 2: Autonomous Code Refactoring & Verification (Steps N/2+1 to N)

1. **Autofixer-First Execution:** Run deterministic AST autofixers first on target files (`05-guideline-autofixer.py`, `08-naming-autofixer.py`, `04-newline-fixer.py`, `07-relative-path-fixer.py`) to automatically resolve 80-90% of mechanical violations.
2. **Cognitive Refactoring:** Agent performs surgical architectural refactoring on the remaining complex logic (<= 8–15 line functions, single return types, `*AppError` envelopes).
3. **Linter & CI Verification:** Execute the section's dedicated linter and run `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` ensuring all 19 gates pass 100% green (`exit 0`).
4. **Automated Plan Consolidation:** Run `python 03-ai-scripts/20-plan-consolidator.py` to archive completed subtasks and update `.lovable/plans/01-index.md`.
5. **Stage & Commit:** Group changes into clean commits (e.g. `refactor(guidelines): enforce <section> rules`).

---

## Mandatory Linter & CI/CD Integration Contract

Every prompt in this suite enforces that code standards must be mechanically verified by automated tooling:

1. **Specific Linter Script:** Each prompt names the exact relative path in `linter-scripts/` responsible for validating its rules.
2. **Auto-Creation Mandate:** If the linter script does not exist, the executing agent is strictly commanded to create it using Python/Node.js/Go.
3. **Local Execution:** The prompt provides the exact local command to execute the linter with zero external overhead.
4. **CI/CD Integration:** The prompt provides the exact configuration snippet to wire the linter script into `.github/workflows/ci.yml` and register it inside `06-cicd-local-runner.py` under the `JOBS` dictionary.

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.lovable/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.lovable/test-inventory.json` so associated tests are known for future release verification.

## STRICT AVOIDANCE: Anti-Compression & Formatting Integrity (No Cheating)

> [!CAUTION]
> **TOTAL BAN ON LINE-COMPRESSION CHEATING:**
> When enforcing file size (<= 100 coding lines) and function size (8–15 lines), AI agents frequently attempt to "cheat" the line counter by destroying formatting. This is strictly forbidden and results in immediate rejection.

- **NO Single-Line If/Else:** NEVER collapse `if/else`, return statements, or blocks into a single line (e.g. `if (x) return y;` or `if (x) { y(); }` are strictly forbidden). Every statement requires its own line and curly braces.
- **NO Deleting Required Blank Lines (R13-R16):** NEVER delete blank lines before `return`/`throw` or after closing `}` to artificially reduce file size.
- **NO Stripping Types or Comments:** NEVER remove TypeScript types, docstrings, or clean indentation to cram code into fewer lines.
- **Mandatory Solution:** The ONLY acceptable way to satisfy line limits is **legitimate modular decomposition** — extracting helper functions into separate files and breaking large components into child components.

---

## Metadata

- slug: cg-execute-index
- version: 2.0.0
- status: active

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.lovable/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.lovable/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.lovable/` (`.lovable/plans/`, `.lovable/01-index.md`, `.lovable/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.lovable/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.lovable/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.lovable/plans/subtasks/` to `.lovable/plans/completed/` and update `.lovable/plans/01-index.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.lovable/plan.md` and `.lovable/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.
