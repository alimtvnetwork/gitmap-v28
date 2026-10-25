# Specification: Git Commit Formatter Engine & Script Remediation

> **Module:** `cli/cmdcommit`, `03-ai-scripts/`  
> **Parent Program:** `331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation`  
> **File:** `02-spec/21-app/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation/01-formatter-and-script-spec.md`  
> **Status:** Draft / Ready for Future Agent Execution  

---

## 1. Scope & Component Changes

### 1.1 In `cli/cmdcommit/commit_formatter.go`
- **Current Problem:**
  `resolveHyphenModuleAndItem` checks `isRedundantModule(prefix, upperMod)`. When `upperMod == "RELEASE"` and `prefix == "Release: "`, it assumes `"RELEASE"` is redundant and attempts to split the remaining item for another nested module. When no nested module exists, it strips the module entirely.
- **Solution:**
  For release commits (`gitmap cpr`), define the canonical module format as:
  ```text
  Release: RELEASE - vX.Y.Z <scope>
  ```
  Allow `"RELEASE"` or `"VERSION"` as explicit valid modules for `cpr` rather than discarding them.
- **Unit Test Coverage:**
  Add table test cases in `cli/cmdcommit/commit_formatter_test.go`:
  - Input: `gitmap cpr "RELEASE - v6.615.0 context menu"` -> `Release: RELEASE - v6.615.0 context menu`
  - Input: `gitmap cpr "VERSION - v6.615.0 context menu"` -> `Release: RELEASE - v6.615.0 context menu` (canonicalized)
  - Input: `gitmap cpr "v6.615.0 context menu"` -> `Release: RELEASE - v6.615.0 context menu` (automatically injects `RELEASE - `)

### 1.2 In `03-ai-scripts/47-git-reconcile-and-resolve-conflict.py`
- **Current Problem:**
  Line 89-91 executes raw `git commit -m "Merge remote-tracking branch ... (mechanically resolved conflicts)"` or fallback commits like `fix: auto-resolve merge conflicts`.
- **Solution:**
  Replace with structured commit command:
  ```python
  # Attempt gitmap cbf first, fallback to standard git commit formatted correctly:
  commit_msg = "Bug: GIT - auto-resolve merge conflicts"
  ```
  Ensure all automated reconciliation commits conform strictly to GitMap's `"Fix: MODULE - summary"` or `"Bug: MODULE - summary"`.

### 1.3 In `01-prompts/17-release-management/02-minor-bump.md` & `37-bump-version.py`
- Ensure the instruction explicitly commands:
  ```text
  gitmap cpr "RELEASE - vX.Y.0 <scope>"
  ```
  and guarantees that `gitmap cpr` emits `Release: RELEASE - vX.Y.0 <scope>`.
