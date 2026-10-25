# Subtask 331.3 — Synchronize Coding Guidelines, Agent Rules, and Skills

> **Parent Plan:** `.ai-memory/plans/pending/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation.md`  
> **Target Files:** `.agents/rules/gitmap_guidelines.md`, `.agents/skills/`  
> **Status:** Pending  

---

## 1. Objectives
1. Update `.agents/rules/gitmap_guidelines.md`:
   - Add explicit definitions for all 5 commit verbs:
     - `gitmap cpf "MODULE - summary"` -> `Feature: MODULE - summary`
     - `gitmap cbf "MODULE - summary"` -> `Bug: MODULE - summary`
     - `gitmap cfr "MODULE - summary"` -> `Refactor: MODULE - summary`
     - `gitmap cpc "MODULE - summary"` -> `Chore: MODULE - summary`
     - `gitmap cpr "RELEASE - vX.Y.Z summary"` -> `Release: RELEASE - vX.Y.Z summary`
   - Add an absolute ban on raw `git commit -m "..."`.
2. Update `.agents/skills/git-reconcile-and-resolve-conflict/skill.md`:
   - Enforce using `gitmap cbf "GIT - auto-resolve merge conflicts"` instead of raw `git commit`.
3. Update `.agents/skills/release-management/skill.md` and `minor-bump/skill.md`:
   - Ensure the release message template is strictly `"RELEASE - vX.Y.0 <scope>"`.
