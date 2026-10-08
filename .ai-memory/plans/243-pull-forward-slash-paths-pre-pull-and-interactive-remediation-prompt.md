# Plan: Forward-Slash Path Hygiene, Pre-Pull Remediation & Interactive Fix Prompt

- **Plan Slug:** `243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt`
- **Related Spec:** `02-spec/21-app/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/`
- **Status:** COMPLETED

---

## Subtasks Breakdown

1. **Subtask 01:** [01-forward-slash-path-normalization-everywhere.md](subtasks/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/01-forward-slash-path-normalization-everywhere.md)
   - Eliminate escaped Windows backslashes (`D:\\work\\gitmap`) across all git commands, hints, and output.
   - Use `filepath.ToSlash(...)` for repository roots and file paths in `cli/cmdpull/pull_remediation_hint.go` and `cli/cmdpull/pull_efficient_render.go`.

2. **Subtask 02:** [02-pull-before-changes-in-remediation-workflows.md](subtasks/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/02-pull-before-changes-in-remediation-workflows.md)
   - Ensure that remediation execution pulls latest remote changes before applying stashes or commits.

3. **Subtask 03:** [03-interactive-remediation-prompt-all-single-skip.md](subtasks/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/03-interactive-remediation-prompt-all-single-skip.md)
   - Lift the prompt suppression check in `handlePullRemediationForRecords` for `pull-all` / `pa`.
   - Implement interactive prompt options: `[a/1] Fix all`, `[s/2] Fix single / select`, `[k/q] Skip / Exit`.

4. **Subtask 04:** [04-testing-quality-gates-and-release-ceremony.md](subtasks/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/04-testing-quality-gates-and-release-ceremony.md)
   - Run AST linters (nested ifs, boolean guidelines, legacy refs, gofmt, golangci-lint).
   - Bump minor version to `v6.516.0`, commit with hyphen format, and verify CI/CD green.
