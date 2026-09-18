# Post-Release Auto-Commit

The `release`, `release-branch`, `release-self`, and `release-pending` commands feature a post-workflow auto-commit behavior:

1. If only files in the `.gitmap/release/` directory (or legacy `.release/` directory) are modified or deleted after returning to the original branch, the system automatically commits and pushes silently. This classification ensures that directory migration deletions (tracked by Git) are treated as part of the release, allowing for a silent flow.
2. If changes exist elsewhere, the system prompts the user to auto-commit all changes. On decline, it commits only `.gitmap/release/` files.
3. Push failures due to remote movement trigger an automated `git pull --rebase` recovery.

## Auto-Confirm (`-y` / `--yes`)

All release commands accept `-y` (short) or `--yes` (long) to skip the interactive commit prompt. When set, the auto-commit step behaves as if the user answered "yes" — all changed files (not just `.gitmap/release/`) are committed and pushed without user input. This enables fully non-interactive release workflows:

```bash
gitmap release v2.55.0 -y
gitmap r v2.55.0 -y          # alias
gitmap release-branch release/v2.55.0 -y
gitmap release-pending -y
gitmap release-self v2.55.0 -y
```

The `-y` flag is stored as `Yes bool` in `release.Options` and passed through to `AutoCommit(version, dryRun, yes)`. The `promptAndCommit` function checks the `yes` parameter before reading from stdin.

## Bypass

The auto-commit step can be bypassed entirely via `--no-commit`. When both `--no-commit` and `-y` are set, `--no-commit` takes precedence (no commit is attempted).
