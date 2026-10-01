# Subtask 05: CPAR and FIA Commands

**Objective**: Implement `gitmap commit-push-all-repos (cpar)` and `gitmap fix-ignore-all (fia)` commands per Spec 197.

## Requirements
- `gitmap cpar [-y] [--review (r)] [--commit-only (co)]`: Scans all repos for pending changes, prompts for review/batch commit.
- `gitmap fix-ignore-all (fia) [-y]` and `gitmap fix-ignores-all-ssh (fias) [-y]`: Sweeps all repos for ignore consistency, dedupes ignores, and applies fixes. fias uses the PAS formula (SSH delegation).
- `gitmap see (c)` subcommands: `c commit pending`, `c gitignore issues (ig)`, `c errors ssh (ses)`
- Add commands to `cli/cmdcpar/`, `cli/cmdfixgit/` or `cli/cmdignore/`.

## Constraints
- Bounding Box: `cli/cmdcpar/*.go`, `cli/cmdignore/*fix*.go`, `cli/cmdsee/*.go`
- Coding Rules: Positive booleans only, `*appfault.AppError`, functions <= 15 lines. No builds/tests.
