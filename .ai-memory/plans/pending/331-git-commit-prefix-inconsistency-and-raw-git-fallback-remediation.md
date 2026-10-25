# Plan 331 — Git Commit Prefix Inconsistency & Raw Git Fallback Remediation

> **Status:** Pending  
> **Sequence ID:** 331  
> **Slug:** `git-commit-prefix-inconsistency-and-raw-git-fallback-remediation`  
> **Specification:** `02-spec/21-app/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation/`  
> **Subtasks Directory:** `.ai-memory/plans/subtasks/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation/`  
> **Date:** 2026-10-26  

---

## 1. Executive Summary & Root Cause Summary
The user noted that recent commit headers exhibit three inconsistent patterns in Git history:
1. `Release: v6.615.0 ...` vs `Release: RELEASE - v6.614.0 ...` vs `Release: VERSION - v6.613.0 ...`
2. Missing `MODULE - ` prefix in release commits where `mod == "RELEASE"` was dropped by `commit_formatter.go`.
3. Raw conventional commit `fix: auto-resolve merge conflicts` bypassing `gitmap` semantic verbs.

This plan details all surgical changes required across the Go CLI formatter, AI reconciliation scripts, release scripts, and coding guidelines so any agent can execute and verify the remediation.

---

## 2. Subtasks Roadmap

| Subtask ID | File | Description | Target | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | `01-subtask-commit-formatter-engine.md` | Standardize `cli/cmdcommit/commit_formatter.go` to preserve canonical `Release: RELEASE - vX.Y.Z <scope>` and auto-inject `RELEASE` module if omitted. | `cli/cmdcommit/` | Pending |
| **Subtask-02** | `02-subtask-auto-reconcile-script-fix.md` | Update `03-ai-scripts/47-git-reconcile-and-resolve-conflict.py` to use `gitmap cbf` or standardized commit messages. | `03-ai-scripts/` | Pending |
| **Subtask-03** | `03-subtask-rules-and-skills-sync.md` | Update `.agents/rules/gitmap_guidelines.md` and `.agents/skills/` to strictly forbid raw `git commit` and document canonical release examples. | `.agents/`, `01-prompts/` | Pending |

---

## 3. Verification & Acceptance Criteria
- [ ] `gitmap cpr "RELEASE - v6.616.0 test"` yields `Release: RELEASE - v6.616.0 test`.
- [ ] `gitmap cpr "v6.616.0 test"` yields `Release: RELEASE - v6.616.0 test`.
- [ ] `commit_formatter_test.go` has 100% test pass rate with new release module expectations.
- [ ] `03-ai-scripts/47-git-reconcile-and-resolve-conflict.py` never emits raw `fix:` commits.
- [ ] `.agents/rules/gitmap_guidelines.md` contains strict ban on raw `git commit` and clear verb examples.
