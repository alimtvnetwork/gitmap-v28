# Subtask 331.2 — Update Conflict Reconciliation and Release Scripts

> **Parent Plan:** `.ai-memory/plans/pending/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation.md`  
> **Target Package:** `03-ai-scripts/`  
> **Status:** Pending  

---

## 1. Objectives
1. In `03-ai-scripts/47-git-reconcile-and-resolve-conflict.py`:
   - Replace raw `git commit -m "fix: auto-resolve merge conflicts"` or similar fallbacks with `gitmap cbf "GIT - auto-resolve merge conflicts"`.
   - If `gitmap` is not in PATH or during bootstrapping, fallback to `git commit -m "Bug: GIT - auto-resolve merge conflicts"`.
2. In `03-ai-scripts/37-bump-version.py`:
   - Verify that release commit generation uses `gitmap cpr "RELEASE - vX.Y.Z <scope>"`.
3. In `01-prompts/09-commit-and-multi-agent-code-fix/08-git-reconcile-and-resolve-conflict.md`:
   - Update Step 6 instructions to use `gitmap cbf "GIT - auto-resolve merge conflicts"`.
