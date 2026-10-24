# 00 — Program Overview: GitMap Command and Coding Guideline Synchronization

> **Program ID:** 318  
> **Slug:** `gitmap-command-and-coding-guideline-synchronization`  
> **Status:** Active / In-Flight  
> **Canonical Spec Root:** `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/`  
> **Master Plan:** `.ai-memory/plans/pending/318-gitmap-command-and-coding-guideline-synchronization.md`

---

## 1. User Request (Verbatim)

```text
On these that this is actually ran by cursor. Why it ran this code? That's the first observation. Second is that why we shouldn't have the alternative commands in `gitmap`. Even building the Rust tools, things like that, installing and building, we should have `gitmap` command and these commands as an alternative and things like that should be updated in the `gitmap` skills in the coding guideline and sync with all the other repositories. Do you understand? Can you please follow through and make sure in the code base there is no magic string other than logs. Every magic string should be replaced with proper constant or enum. If it is from a same group, then it should be coming from the enum by reading the enum guideline from the coding guidelines inside the spec. Can you please follow through and do that? Also make sure that in the execute prompt in V6, it should also refer to the coding guideline error manage, and some of the principles should be embedded in the execute parent task in inner steps in the coding guideline. Make sure of that. And do not write, let's say, commits like this third bracket, update the coding guideline for the execute prompt in V6. It should be done using the `gitmap` CFR CBF, and it should be in uppercase, and in bracket, there should be no third brackets for the commits. So that needs to be corrected everywhere with an example. Do you understand the task? Can you follow through?

## task: GitMap Command and Coding Guideline Synchronization
## slug: gitmap-command-and-coding-guideline-synchronization

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-gitmap-command-and-coding-guideline-synchronization/ and enqueue plan task in .ai-memory/plans/pending/XX-gitmap-command-and-coding-guideline-synchronization.md (subtasks in .ai-memory/plans/subtasks/XX-gitmap-command-and-coding-guideline-synchronization/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Ensure all alternative commands are updated in `gitmap` skills and synchronized across all repositories.
6. Replace all magic strings in the codebase with proper constants or enums, following the enum guideline from the coding guidelines inside the spec.
7. Update the execute prompt in V6 to refer to the coding guideline error management and embed relevant principles in the execute parent task.
8. Correct commit message format to use `gitmap` CFR CBF, in uppercase, without third brackets.

## Must follow and spawn agent using
@[.agents/skills/execute-parent-task-with-n-steps-v6]
```

---

## 2. Strategic Objectives & Architectural Pillars

### Pillar 1: Root Cause & Systems Intent
- **Cursor Ad-Hoc Scripts:** AI agents running inside Cursor, Muse, or Antigravity resorted to writing raw Python snippets (`gitmap py -c "..."`) and manual shell snippets (`taskkill /PID`, `$env:LIBCLANG_PATH=...`, `cargo test`) because first-class GitMap alternative commands for Rust, Cargo, Tauri, process management, and hardware/disk telemetry were not documented in agent skills.
- **GitMap Alternative Tooling:** GitMap must document official commands (`gitmap cargo ...`, `gitmap crg ...`, `gitmap aum locate llvm`, `gitmap env set ...`, `gitmap ports ...`, `gitmap os user kill ...`, `gitmap disk ...`, `gitmap os storage ...`) in `.agents/skills/gitmap/skill.md`.
- **Magic String Elimination:** Eliminate raw string comparisons across the codebase in favor of strongly-typed PascalCase enums ending in `Type` adhering strictly to `02-spec/02-coding-guidelines/04-constants-and-enums/`.
- **Execute V6 Error Management Integration:** Canonical prompt `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` and skill `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` must formally embed Leaf Error Sovereignty, no `os.Exit` in leaf methods, `*appfault.AppError` wrappers, typed exit codes, and database audit capture.
- **Third Bracket Elimination:** Banish square brackets (`[` `]`) in commit headers (`feat([module]):`). Standardize on GitMap semantic hyphen format with uppercase module names (`gitmap cpf "MODULE - summary"`, `gitmap cpb "MODULE - summary"`, `gitmap cfr "MODULE - summary"`).

### Pillar 2: Scope Partition
- **In-Scope:**
  - `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/` (modular specs 00 to 03).
  - `.ai-memory/plans/pending/318-gitmap-command-and-coding-guideline-synchronization.md` and subtasks.
  - `.agents/skills/gitmap/skill.md` (alternative command recipes).
  - `cli/cmdcommit/commit_formatter.go` (bracket removal and hyphen preservation).
  - `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/03-cli-commands-and-commit-format-spec.md` (revert bracket rule).
  - `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md` (ban third brackets).
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` & `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` (error management & commit format).
  - Target `cli/` modules for magic string to enum migration.
- **Excluded:**
  - Non-owned third-party Git remotes (`common-linux-installer`, `alim.karim.profile`, `kubernetes-training`).
  - Full local compilation (`go build`, `cargo build`) or test suites per Rule R1.

---

## 3. Modular Specification Structure

| Document | Title | Ownership |
| :--- | :--- | :--- |
| `00-overview.md` | Program Overview & Requirements Ingestion | Lead Orchestrator |
| `01-alternative-commands-and-skills-spec.md` | GitMap Alternative Commands, Recipes & Skills | Authoring Subagent 1 |
| `02-commit-format-and-bracket-elimination-spec.md` | Commit Format Standardization & Third Bracket Removal | Authoring Subagent 2 |
| `03-magic-strings-and-error-management-spec.md` | Magic Strings to Enums & V6 Error Management Integration | Authoring Subagent 3 |
