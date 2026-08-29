# Master Plan: Repository-Wide Boolean Conventions, Naming, Enums & Conditional Flattening (17-boolean-and-naming-audit.md)

> **Status:** PENDING  
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)  
> **Target:** Eliminate all explicit boolean comparisons (`== true`, `== false`), inverted success checks (`!isSuccess`), single-line if blocks, nested `if` statements (depth > 1), and missing `Type` enum suffixes across the entire repository.

---

## 1. Task-Specific Rule Set

1. **Implicit Booleans Only:** NEVER compare booleans with literals (`== true`, `== false`, `=== true`, `=== false`). Evaluate them directly: `if isReady { ... }` or `if !isReady { ... }`.
2. **Zero Inverted Success Checks:** Never invert positive success variables (`!response.isSuccess`). Use explicit failure variables (`response.isFail`, `isError`).
3. **Zero Nested `if` Blocks (Nesting Depth $\le 1$):** Flatten nested `if`s using guard clauses, early returns, and extracted helper functions. Keep happy paths un-indented.
4. **Enum `Type` Suffix:** All enum definitions must end with `Type` suffix (e.g. `UserRoleType`, `SyncStatusType`).
5. **Anti-Compression Integrity:** Never collapse `if/else` statements onto a single line to cheat line caps.

---

## 2. Subtask Breakdown

- **Subtask 01**: [`01-explicit-boolean-comparisons-and-inversions.md`](../subtasks/17-boolean-and-naming/01-explicit-boolean-comparisons-and-inversions.md)
- **Subtask 02**: [`02-flatten-nested-ifs-in-cmd.md`](../subtasks/17-boolean-and-naming/02-flatten-nested-ifs-in-cmd.md)
- **Subtask 03**: [`03-flatten-nested-ifs-in-packages.md`](../subtasks/17-boolean-and-naming/03-flatten-nested-ifs-in-packages.md)
- **Subtask 04**: [`04-enum-type-suffix-and-naming.md`](../subtasks/17-boolean-and-naming/04-enum-type-suffix-and-naming.md)
- **Subtask 05**: [`05-ci-linter-integration-and-verification.md`](../subtasks/17-boolean-and-naming/05-ci-linter-integration-and-verification.md)
