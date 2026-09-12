# LFS Smudge Fallback

## Overview

Gitmap intercepts Git LFS smudge filter failures (404s when pointers point to missing objects on the LFS server) during cloning. Instead of leaving the user with a broken working tree, it interactively prompts to auto-fix the problem.

## Implementation Details

1. `clonefrom.detectLFSSmudgeError` parses `git clone` stderr.
2. `clonefrom.confirmYesNo` asks for permission to fix.
3. `clonefrom.executeLFSFix` implements the PowerShell script sequence in Go using `exec.Command` injected with `GIT_LFS_SKIP_SMUDGE=1`.

## File Locations

- `cli/clonefrom/execute.go` (Integration point inside `runGitClone`)
- `cli/clonefrom/execute_lfs_fix.go` (Regex detection and Git execution logic)
- `cli/clonefrom/prompt.go` (Isolated confirmation prompt)
- `.lovable/memory/specs/02-lfs-smudge-rca.md` (Original bug report and RCA)
