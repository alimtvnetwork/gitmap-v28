# Master Plan: 318 — GitMap Command and Coding Guideline Synchronization

> **Status:** Completed  
> **Sequence ID:** 318  
> **Slug:** `gitmap-command-and-coding-guideline-synchronization`  
> **Canonical Spec:** `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/`  
> **Subtask Root:** `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/`

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
@[\1](file;.agents/skills/\1)
```

---

## 2. 5-Pillar Architectural Mandate

### Pillar 1: Strategic Intent & Systems Reasoning
Provide a comprehensive unification across GitMap CLI commands, agent skills, coding standards, and commit conventions:
1. **Alternative Command Recipes:** Prevent AI tools (such as Cursor, Muse, Antigravity) from executing ad-hoc Python snippets or raw shell commands by documenting official, first-class GitMap alternative commands for Rust, Cargo, Tauri, process termination, and disk metrics in `.agents/skills/gitmap/skill.md`.
2. **Commit Hygiene & Third Bracket Removal:** Banish square brackets (`[` `]`) in commit scopes (`feat([module]):`). Enforce the canonical GitMap semantic hyphen format with uppercase module tags (`gitmap cpf "MODULE - summary"`, `gitmap cpb "MODULE - summary"`, `gitmap cfr "MODULE - summary"`).
3. **Magic String Elimination:** Modernize comparison logic in `cli/cmd/`, `cli/constants/` by migrating raw string literals to PascalCase enums ending in `Type`.
4. **V6 Execute Prompt Error Management:** Update `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` and `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` with Leaf Error Sovereignty, zero `os.Exit` in leaf methods, `*appfault.AppError` wrappers, typed exit codes, and database audit capture.

### Pillar 2: Explicit Scope Boundaries
- **In-Scope:**
  - `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/` (00 to 03).
  - `.ai-memory/plans/pending/318-gitmap-command-and-coding-guideline-synchronization.md`.
  - `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/` (subtasks 01 to 03).
  - `.agents/skills/gitmap/skill.md` (alternative command recipes).
  - `cli/cmdcommit/commit_formatter.go` & `commit_formatter_test.go` (bracket removal).
  - `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/03-cli-commands-and-commit-format-spec.md` (revert bracket rule).
  - `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md` (ban third brackets).
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` & `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`.
  - Key magic strings in `cli/constants/os_targets.go`, `cli/cmd/cmd_explorer.go`, `cli/cmd/cmd_whoami.go`, `cli/cmd/alias.go`, `cli/cmd/async_cmd.go`.
- **Excluded:**
  - Non-owned third-party Git remotes (`common-linux-installer`, `alim.karim.profile`, `kubernetes-training`).
  - Full local compilation (`go build`, `cargo build`) or test suites per Rule R1.

### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
- **Strict Search Ban:** No `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`. Only GitMap search commands.
- **Strict Relative Git Paths:** Zero drive letters or `file:///` URIs in any markdown or commit text.
- **Worker Git Ban:** Subagents never run Git commands (`git add`, `git commit`, `git push`, etc.).
- **Zero Builds & Tests (Rule R1):** Never run global build or test suites; fast file-scoped linters only.
- **Commit Formatting Rule:** GitMap hyphen format without colons in the message argument; zero square brackets (`[` `]`) in commit headers.

### Pillar 4: Blast Radius & Dependency Topology
- Subsystems touched: `cli/cmdcommit/`, `cli/constants/`, `.agents/skills/`, `01-prompts/14-execute/`, `02-spec/`.
- File locks avoided by disjoint subagent assignments.

### Pillar 5: Multi-Agent Concurrency Topology & Wave Orchestration
- $A = 3$ subagents, $H = 2$ hands (Capacity = 6 concurrent operations).
- Stage 1: 3 `research` subagents (Completed).
- Stage 2: 3 `self` authoring subagents for modular specs & subtasks (Completed).
- Stage 3: 3 `self` execution subagents for code, skills, and prompts (Completed).

---

## 3. Subtask Decomposition & Tracking

| Subtask ID | Title | Spec Ref | Status |
| :--- | :--- | :--- | :--- |
| **Subtask-01** | GitMap Alternative Commands & Skills Enhancement | `01-alternative-commands-and-skills-spec.md` | Completed |
| **Subtask-02** | Commit Format Standardization & Third Bracket Removal | `02-commit-format-and-bracket-elimination-spec.md` | Completed |
| **Subtask-03** | Magic String Enum Refactoring & V6 Error Management Integration | `03-magic-strings-and-error-management-spec.md` | Completed |

---

## 4. Verification Evidence Contract
- FastGate pre-commit linter (`python 03-ai-scripts/50-fastgate.py`): 100% PASS (< 1.5s).
- Path hygiene verification (`python linter-scripts/check-relative-paths.py`): PASS across all repository files.
- Atomic GitMap commit with uppercase module and hyphen format without colons (`gitmap cpf "SPEC - add program 318 gitmap command and coding guideline sync"`).
