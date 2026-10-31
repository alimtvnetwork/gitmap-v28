# Subtask 05: Cobra Argument Completion for Short Aliases & Release Install Hints

- **Task Code:** `TASK-05`
- **Slug:** `gitmap-installation-suggestion-issue - Task 05`
- **Parent Plan:** [365-gitmap-installation-suggestion-issue.md](../pending/365-gitmap-installation-suggestion-issue.md)
- **Status:** `pending`

---

## Scope & Deliverables
1. Update `cli/cmd/root_cobra_completion.go` to register argument completion for alias `gm`.
2. Update `cli/release/releaseinstallhint.go` to allow `repoName == "gitmap"` without version suffix.
3. Update `cli/cmdinstall/install_prompt.go` to print actionable suggestions when `gitmap` is already installed.
