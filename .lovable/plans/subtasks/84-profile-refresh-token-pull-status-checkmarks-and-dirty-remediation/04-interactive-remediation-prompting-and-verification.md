# Subtask 04: Interactive Remediation Prompting & Verification

## Objective
Implement interactive user prompting in `gitmap pull` when dirty repositories exist (asking whether the user wants to remediate or display instructions, with non-interactive flags), and run complete quality gates.

## Files to Touch
- `gitmap/cmd/pull.go`
- `gitmap/cmd/remediation_box.go`
- `gitmap/cmd/pull_test.go`

## Detailed Implementation Steps
1. In `gitmap/cmd/pull.go`:
   - Add flags to `pullOptions`:
     - `autoFix bool` (`--fix`)
     - `yes bool` (`--yes`, `-y`)
     - `noFix bool` (`--no-fix`)
   - When `len(remItems) > 0`:
     - If `opts.noFix`: skip prompt, print remediation box with commands.
     - If `opts.yes` or `opts.autoFix`: automatically apply default recipe (stash).
     - Otherwise, if interactive terminal: prompt user:
       `Remediate dirty repository(ies)? [1] Stash all, [2] Commit WIP, [n] Skip/No [default]: `
       Execute choice or skip cleanly.
2. Run unit tests across `cmd` and `gitutil`.
3. Move completed plan to `.lovable/plans/completed/`.
