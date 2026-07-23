# 20 — revert Command

## Purpose

`gitmap revert <version>` checks out a specific release version tag and rebuilds/deploys the binary using the same handoff mechanism as `gitmap update`.

## Command Signature

```
gitmap revert <version>
```

Where `<version>` is a release tag (e.g. `v2.9.0`, `2.9.0`).

## Behavior

### Step 1 — Validate

1. Require exactly one positional argument (the version). Exit 1 with usage if missing.
2. Normalize the version with `release.NormalizeVersion()` (auto-prefix `v`).
3. Verify the tag exists locally via `release.TagExistsLocally()`. Exit 1 if not found.

### Step 2 — Checkout

4. Save the current branch name via `release.CurrentBranchName()`.
5. Run `git checkout <tag>` (detached HEAD) in the repo at `constants.RepoPath`.

### Step 3 — Handoff (same as update)

6. Create a handoff copy of the current binary (same as `update`: same-dir first, fallback to `%TEMP%`).
7. Launch the handoff copy with a hidden `revert-runner` command.
8. The runner generates a temporary PowerShell script that runs `.\run.ps1` (no `-Update` flag — just build + deploy).
9. The runner executes the script, piping stdout/stderr to the terminal.

### Step 4 — Post-revert

10. Print confirmation: `✓ Reverted to gitmap <version>`.
11. Run `update-cleanup` to remove handoff artifacts.

## Error Handling

- If `git checkout` fails, print the error and exit 1 (do not proceed to handoff).
- If the build/deploy script fails, exit with the script's exit code.
- On any failure after checkout, the repo remains on the checked-out tag (user must manually return).

## Implementation Files

| File                          | Responsibility                              |
|-------------------------------|---------------------------------------------|
| `cmd/revert.go`               | Command handler, validation, checkout       |
| `cmd/revertscript.go`         | Handoff copy, script generation, execution  |
| `constants/constants_cli.go`  | `CmdRevert`, help text, error messages      |

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return. No magic strings. No switch statements.
