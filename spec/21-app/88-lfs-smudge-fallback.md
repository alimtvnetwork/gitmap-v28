# Git LFS Smudge Fallback

## Problem

When cloning a repository that contains dangling Git LFS pointers (objects that were never uploaded to the LFS server), `git clone` successfully fetches the Git objects but fails during the checkout phase because the `git-lfs filter-process` (the smudge filter) encounters a 404 error from the LFS API.

This leaves the user with a broken working tree ("Clone succeeded, but checkout failed").

## Solution

When `gitmap` executes a clone operation and detects the specific error signature (`smudge filter lfs failed`), it will:
1. Parse the standard error to extract the offending file path.
2. Prompt the user interactively: asking if they want to automatically drop the broken LFS pointer to fix the clone.
3. If confirmed, it executes the remediation sequence:
   - Sets `GIT_LFS_SKIP_SMUDGE=1` for subsequent commands in the destination directory.
   - Runs `git restore --source=HEAD :/` to populate the working tree (skipping smudge).
   - Runs `git rm --cached "<offending-file>" -q` to drop the file from the index.
   - Removes the physical pointer file via `os.Remove`.
   - Commits the removal with `git commit -m "chore(lfs): remove pointer for missing LFS object"`.
   - Pushes the fix back to `origin`.

## Architecture Implementation

- **Detection**: `clonefrom/execute.go` in `runGitClone` captures `CombinedOutput`. If it fails, a new function `detectLFSSmudgeError(output)` parses the output.
- **Prompting & Fixing**: To maintain separation of concerns, `runGitClone` will bubble a specific typed error or struct `LFSSmudgeError` up to `clonefrom.Execute`. Wait, since `clonefrom.Execute` currently returns a flat `Result{Status, Detail}`, we can enrich `Result` or we can just handle the prompt and fix inside `runGitClone` via a package-level injected `ConfirmCallback` (defaulting to stdin/stderr prompt) so that all consumers (`gitmap clone`, `gitmap cfr`, `gitmap reclone`) inherit the fix automatically without needing to wire up the prompt logic in 5 different commands.
- **Helper**: The sequence of git commands to execute the fix will be implemented as `runLFSFallbackFix(dest string, file string) error`.
