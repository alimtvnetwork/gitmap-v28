# Migration & Verification Spec: 319 - Package Constants & Enums DRY Restructuring

## 1. Multi-Agent Wave Deployment Strategy ($A = 3$, $H = 2$, Capacity = 6)
When the user triggers execution via `@[01-prompts/15-cg-execute/06-constants-and-enums.md]`, execution proceeds across the 5 target packages in 2 tightly coordinated waves:

### Wave 1: Core Engine Packages (3 Parallel Subagents)
- **Agent 1 (Primary Owner: `cli/cmdpull`):**
  - Files Owned: `cli/cmdpull/enums/*`, `cli/cmdpull/constants/*`, `cli/cmdpull/pull_dispatch.go`, `cli/cmdpull/pull_summary_single.go`, `cli/cmdpull/pull_summary_changes.go`.
  - Action: Instantiate `PullModeType`, `PullStatusType`, constants files, replace all raw string comparisons.
- **Agent 2 (Primary Owner: `cli/cloner` & `cli/cmdclone`):**
  - Files Owned: `cli/cloner/enums/*`, `cli/cloner/constants/*`, `cli/cloner/cloner.go`, `cli/cmdclone/enums/*`, `cli/cmdclone/constants/*`, `cli/cmdclone/cmd_clone.go`.
  - Action: Instantiate `CloneProtocolType`, `CloneDepthType`, `CloneStrategyType`, replace git CLI parameter literals.
- **Agent 3 (Primary Owner: `cli/store`):**
  - Files Owned: `cli/store/enums/*`, `cli/store/constants/*`, `cli/store/store.go`, `cli/store/db_sqlite.go`.
  - Action: Instantiate `StoreTableType`, `DbOpType`, pragmas, replace raw SQL string literals.

### Wave 2: Rollout Package 5 & Cross-Package Verification (3 Parallel Subagents)
- **Agent 1 (Primary Owner: `cli/cmdcommit`):**
  - Files Owned: `cli/cmdcommit/enums/*`, `cli/cmdcommit/constants/*`, `cli/cmdcommit/commit.go`, `cli/cmdcommit/types.go`.
  - Action: Instantiate `CommitVerbType`, `CommitFormatType`, module prefix templates, Conventional Commit validation.
- **Agent 2 (Quality & Guidelines Auditor):**
  - Files Owned: Read-only cross-audit of all 5 packages.
  - Action: Enforce $\le 15$ line function limits, positive booleans (`is*`, `has*`), zero nested `if` statements, blank lines around returns.
- **Agent 3 (Linter & AST Gate Verification):**
  - Files Owned: Read-only verification scripts executor.
  - Action: Execute FastGate, enum linter, and format validation.

---

## 2. Disjoint File Boundary Matrix
| Package | Exclusive Enum Path | Exclusive Constants Path | Modified Call-Sites |
| :--- | :--- | :--- | :--- |
| `cli/cmdpull` | `cli/cmdpull/enums/*.go` | `cli/cmdpull/constants/*.go` | `cli/cmdpull/pull_dispatch.go`, `cli/cmdpull/pull_cwd.go` |
| `cli/cloner` | `cli/cloner/enums/*.go` | `cli/cloner/constants/*.go` | `cli/cloner/cloner.go`, `cli/cloner/transport.go` |
| `cli/cmdclone` | `cli/cmdclone/enums/*.go` | `cli/cmdclone/constants/*.go` | `cli/cmdclone/clone_batch.go`, `cli/cmdclone/options.go` |
| `cli/store` | `cli/store/enums/*.go` | `cli/store/constants/*.go` | `cli/store/store.go`, `cli/store/repo_query.go` |
| `cli/cmdcommit` | `cli/cmdcommit/enums/*.go` | `cli/cmdcommit/constants/*.go` | `cli/cmdcommit/commit.go`, `cli/cmdcommit/formatter.go` |

---

## 3. Verification & Quality Gates
1. **AST & Syntax Gate:**
   - Execute FastGate AST validation to ensure 0 syntax errors and clean symbol resolution:
     ```powershell
     python 03-ai-scripts/50-fastgate.py
     ```
2. **Enum Coding Guideline Linter:**
   - Verify that all newly created and refactored enums conform to the `*Type` suffix and have `IsValid()` methods:
     ```powershell
     python linter-scripts/check-enum-guidelines.py
     ```
3. **Forbidden Strings & Magic Literal Scan:**
   - Confirm elimination of raw magic literals across target package boundaries:
     ```powershell
     python linter-scripts/check-forbidden-strings.py
     ```
4. **Formatting Invariant:**
   - Ensure 100% compliance with `gofmt`:
     ```powershell
     gofmt -w cli/cmdpull/ cli/cloner/ cli/cmdclone/ cli/store/ cli/cmdcommit/
     ```

---

## 4. Final Atomic Commit Invariant
- Following multi-agent completion and 100% green verification, the Lead Orchestrator commits using GitMap:
  ```powershell
  gitmap cpf "constants - centralize enums and constants across 5 packages"
  ```
- No raw `git commit` commands.
- No colons in commit arguments (colon is rendered by GitMap).
- Clean working tree maintained.
