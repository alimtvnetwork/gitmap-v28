# Specification: Coding Guidelines, Skills, and Rule Enforcement

> **Module:** `.agents/rules/`, `.agents/skills/`, `02-spec/02-coding-guidelines/`  
> **Parent Program:** `331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation`  
> **File:** `02-spec/21-app/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation/02-guidelines-and-skills-spec.md`  
> **Status:** Draft / Ready for Future Agent Execution  

---

## 1. Updates Required in `.agents/rules/gitmap_guidelines.md`

Add explicit rules and positive examples for all semantic commit verbs:

```markdown
### 4. GitMap Command Discipline & Commit Standards
- **Mandatory Pre-Flight Pull**: Always execute `gitmap pull` (fallback `git pull`) before creating or modifying code.
- **Dedicated Semantic Commit Verbs**:
  - Features:   `gitmap cpf "MODULE - summary"` -> creates `Feature: MODULE - summary`
  - Bug Fixes:  `gitmap cbf "MODULE - summary"` -> creates `Bug: MODULE - summary`
  - Refactors:  `gitmap cfr "MODULE - summary"` -> creates `Refactor: MODULE - summary`
  - Chores:     `gitmap cpc "MODULE - summary"` -> creates `Chore: MODULE - summary`
  - Releases:   `gitmap cpr "RELEASE - vX.Y.Z summary"` -> creates `Release: RELEASE - vX.Y.Z summary`
- **TOTAL BAN ON RAW GIT COMMITS**: Autonomous agents must NEVER execute raw `git commit -m "fix: ..."` or `git commit -m "feat: ..."`. Every commit MUST use `gitmap` verbs (`cpc`, `cpr`, `cbf`, `cfr`, `cpf`).
- **TOTAL BAN ON THIRD BRACKETS**: Never use `[module]:` or `type([module]):`.
```

---

## 2. Updates Required in Skills

1. **`.agents/skills/git-reconcile-and-resolve-conflict/skill.md`**:
   - Update Step 6 (Merge Commit Creation):
     Change from raw `git commit -m "Merge ..."` or `fix: ...` to:
     ```bash
     gitmap cbf "GIT - auto-resolve merge conflicts"
     ```
2. **`.agents/skills/release-management/skill.md` & `minor-bump/skill.md`**:
   - Emphasize canonical release commit syntax:
     ```bash
     gitmap cpr "RELEASE - vX.Y.0 <scope>"
     ```
3. **`.agents/skills/gitmap/skill.md`**:
   - Expand the commit command reference table to include explicit outputs for each verb.
