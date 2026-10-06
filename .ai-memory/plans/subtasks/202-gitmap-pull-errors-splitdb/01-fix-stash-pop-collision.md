# Subtask: Fix Stash Pop Untracked File Collision

## Objective
Implement a recovery strategy for `git stash pop` untracked file collisions in the `gitmap fix` workflow by automatically dropping the stash if it fails specifically due to the untracked file collision.

## Instructions
1. Open `cli/cmd/fix_execute.go`.
2. Add a new helper function `recoverStashCollision(step gitutil.RemediationStep, output string) bool`:
   - Check if the command includes a `stash pop` operation.
   - Check if `strings.Contains(strings.ToLower(output), "could not restore untracked files from stash")` AND `strings.Contains(strings.ToLower(output), "already exists, no checkout")`.
   - Extract the `repoPath` from `step.Args` (typically `step.Args[1]` if `step.Args[0] == "-C"`).
   - If a collision is detected and the path is extracted, execute `exec.Command("git", "-C", repoPath, "stash", "drop")`. Run it using `.Run()`.
   - If `exec.Command` succeeds, return `true`, else return `false`.
3. Locate `handleStepFailure` in `cli/cmd/fix_execute.go`.
4. Just below the `isBenignCommitClean` and `isBenignStashClean` checks, add a new conditional block:
   - Call `recoverStashCollision(step, outStr)`.
   - If it returns `true`, print a formatted recovery message (using `constants.ColorYellow`) such as `recovered (dropped stashed untracked collision)` and return `nil`.
5. Save the file and verify formatting.

## Execution Rules
- Strict relative paths only.
- Do not build or test during this step.
- Focus exclusively on the code modification.
