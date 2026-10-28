# Architecture Specification: GitMap Skills and Features Update

**Sequence ID:** 337  
**Document:** `01-architecture-spec.md`  
**Task Slug:** `gitmap-skills-and-features-update`  
**Scope:** CLI Help & Examples, PowerShell Execution Hardening, Database Search Exclusivity, Skills Hygiene & Sync  

---

## 1. Executive Summary & Problem Statement

This specification defines surgical architectural fixes and feature enhancements across GitMap:
1. **PowerShell Argument Parsing & Script Execution Hardening**: In `cli/cmd/cmd_powershell.go`, scripts passed to `gitmap pwsh -c` or `gitmap ps -c` were invoked directly via `-Command strings.Join(cleanArgs, " ")`. On Windows, complex commands with variable assignments, semicolons, and pipes could result in PowerShell misinterpreting syntax tokens (e.g. `=: The term '=' is not recognized as a name of a cmdlet`). This is hardened to wrap script arguments safely into PowerShell scriptblocks or properly escaped `-Command` invocations.
2. **Help Documentation & Examples for Move/Merge (`merge-right` and `merge-left`)**: In `cli/helpdoc/examples_merge.go`, curated examples existed for `merge write` / `mr` (pull request creation), but lacked curated `HelpExample` entries for `merge-right` and `merge-left`. This left `gitmap help example merge-right` and `gitmap help example merge-left` rendering fallback placeholders.
3. **Coding Guidelines & Prompts Flag Examples**: `gitmap cg` and `gitmap install-prompts` lacked detailed flag examples in `cli/cmdcg/cg_help.go` and `cli/helpdoc/examples_dev.go`. Users running `gitmap help example cg` received "no help example yet", and `gitmap cg --help` lacked concrete flag usage examples.
4. **Database Search Strict Exclusivity**: Prohibit ad-hoc Python SQLite querying scripts in skill documentation, enforcing native compiled `gitmap db-find`, `gitmap db-read`, and `gitmap db-listing` exclusively.
5. **Skills Inventory & Hygiene Alignment**: Synchronize skills between `gitmap` and `coding-guidelines`, deprecating and archiving obsolete prompt wrappers (including OS clean and legacy plan tools) to streamline agent context.

---

## 2. Component Design & Modifications

### 2.1 PowerShell Command Execution (`cli/cmd/cmd_powershell.go`)
- Inspect arguments in `buildPowerShellCmd`:
  - When `-Command` or `-c` is passed with subsequent tokens, coalesce into a single unified script argument or scriptblock `& { <script> }`.
  - On Windows, ensure `-NoProfile` is always applied, and invoke `-Command` with a properly quoted block so that variable assignments (`$a = ...`) are not fragmented by the OS process argument parser.

### 2.2 Curated Merge Examples (`cli/helpdoc/examples_merge.go`)
- Register curated examples for:
  - `merge-right` (and aliases `mrg`)
  - `merge-left` (and alias `ml`)
  - `merge-both` (and alias `mb`)
- Provide realistic invocations with expected output showing diffing, conflict resolution, and commit/push behavior.

### 2.3 Curated CG & Prompt Install Examples (`cli/helpdoc/examples_dev.go` & `cli/cmdcg/cg_help.go`)
- In `cli/cmdcg/cg_help.go`:
  - Add concrete command and flag invocation examples in `buildCGUsageLines()`:
    - `gitmap cg install [targets...] [--all] [--dry-run]`
    - `gitmap cg install-prompts [targets...] [--all] [--dry-run]`
    - `gitmap cg sync --all [--exclude <csv>]`
    - `gitmap cg install-version-json --version 1.0.0`
- In `cli/helpdoc/examples_dev.go`:
  - Add curated registry entries for `cg`, `coding-guidelines`, `install-prompts`, and `cg install`.

### 2.4 Skills Documentation & Database Exclusivity
- In `.agents/skills/gitmap/skill.md`:
  - Reinforce Section 9 (Total Ban on Python SQLite Scripts) with explicit guidance on `gitmap db-find` pattern matching and `gitmap db-read`.
  - Remove deprecated skills from the repository's `.agents/skills/` and `.cursor/skills/` following the coding-guidelines reduction matrix.
