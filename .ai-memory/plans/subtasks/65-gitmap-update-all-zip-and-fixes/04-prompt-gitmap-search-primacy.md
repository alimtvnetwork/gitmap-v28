# Subtask 65.4: Canonical Prompts GitMap Search Primacy & Sub-Node Diagnostics

> **Parent Plan:** [65-gitmap-update-all-zip-and-fixes.md](../../pending/65-gitmap-update-all-zip-and-fixes.md)  
> **Tracking Spec:** [02-component-spec.md](../../../../02-spec/21-app/65-gitmap-update-all-zip-and-fixes/02-component-spec.md)  
> **Status:** Pending  
> **Primary File Targets:** `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`  

---

## 1. Objective

Update the V6 continuous execution workflow prompt and its matching Antigravity skill:
1. Enforce an absolute ban on raw unindexed search utilities:
   - `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`, `Get-ChildItem -Recurse`.
2. Mandate GitMap high-speed search tools with explicit syntax examples:
   - `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
   - `gitmap search "<symbol>"`
   - `gitmap lf [dir]`
   - `gitmap cat <file>`
3. Add the Package Path Disambiguation Rule: when packages or modules have identical names in different parts of the repository, always refer to them using distinct relative repository paths.
4. Add the Failure Sub-Node Diagnostic Rule: when an operation fails, provide structured sub-node options presenting up to two distinct remediation paths.

---

## 2. Implementation Scope

### 2.1 Prompt Modification (`01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`)
- **Search Primacy Section:**
  - Update rule item (Search Primacy) in the preflight and Phase 2 guidelines.
  - Expand the ban list to explicitly mention `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.
  - Provide concrete invocation examples for `gitmap aum search` and `gitmap search`.
- **Package Path Disambiguation Rule:**
  - Inject explicit directive requiring agents to qualify repeated packages by their full relative path (`cli/cmdpull/pull.go` vs `cli/cmdignore/fix_ignore.go`) to prevent ambiguity.
- **Failure Sub-Node Diagnostics:**
  - Inject guideline instructing agents to branch failure diagnosis into sub-nodes with two alternative solutions (e.g. Option A: Rebase / Option B: Discard).
- **Checklist Ingestion:**
  - Update the turn-end checklist item to enforce GitMap search usage and verify zero usage of banned search tools.

### 2.2 Skill Synchronization (`.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`)
- Update the embedded skill prompt in `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` in complete lockstep with `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`.
- Ensure YAML frontmatter is preserved.

---

## 3. Coding Guidelines & Constraints
- Keep markdown formatting clean, semantic, and well-structured.
- Maintain consistency between `01-prompts/` and `.agents/skills/`.
- Ensure all command examples reflect actual tested GitMap CLI syntax.

---

## 4. Verification Steps
- Diff the prompt and skill files to ensure 100% parity of the execution rules.
- Verify that searching for banned tokens (`ripgrep`, `rg`, `Select-String`) in the prompt rules clearly marks them as forbidden.
